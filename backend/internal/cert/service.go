/*
 * Copyright (c) 2025, WSO2 LLC. (https://www.wso2.com).
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

// Package cert provides the implementation for managing certificates in the system.
package cert

import (
	"errors"

	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
	sysutils "github.com/asgardeo/thunder/internal/system/utils"
)

const loggerComponentName = "CertificateService"

// CertificateServiceInterface defines the methods for certificate service operations.
type CertificateServiceInterface interface {
	GetCertificateByID(id string) (*Certificate, *serviceerror.ServiceError)
	GetCertificateByReference(refType CertificateReferenceType, refID string) (
		*Certificate, *serviceerror.ServiceError)
	CreateCertificate(cert *Certificate) (*Certificate, *serviceerror.ServiceError)
	UpdateCertificateByID(id string, cert *Certificate) (*Certificate, *serviceerror.ServiceError)
	UpdateCertificateByReference(refType CertificateReferenceType, refID string, cert *Certificate) (
		*Certificate, *serviceerror.ServiceError)
	DeleteCertificateByID(id string) *serviceerror.ServiceError
	DeleteCertificateByReference(refType CertificateReferenceType, refID string) *serviceerror.ServiceError
}

// certificateService implements the CertificateServiceInterface for managing certificates.
type certificateService struct {
	Store certificateStoreInterface
}

// TODO: Need to remove once the authz service and token service are refactored to use DI.

// NewCertificateService creates a new instance of CertificateService.
func NewCertificateService() CertificateServiceInterface {
	return &certificateService{
		Store: newCachedBackedCertificateStore(),
	}
}

func newCertificateService(store certificateStoreInterface) CertificateServiceInterface {
	return &certificateService{
		Store: store,
	}
}

// GetCertificateByID retrieves a certificate by its ID.
func (s *certificateService) GetCertificateByID(id string) (*Certificate, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if id == "" {
		return nil, &ErrorInvalidCertificateID
	}

	certObj, err := s.Store.GetCertificateByID(id)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to get certificate by ID", log.String("id", id), log.Error(err))
		return nil, &ErrorInternalServerError
	}
	if certObj == nil {
		logger.Debug("Certificate not found for ID", log.String("id", id))
		return nil, &ErrorCertificateNotFound
	}

	return certObj, nil
}

// GetCertificateByReference retrieves a certificate by its reference type and ID.
func (s *certificateService) GetCertificateByReference(refType CertificateReferenceType,
	refID string) (*Certificate, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if !isValidReferenceType(refType) {
		return nil, &ErrorInvalidReferenceType
	}
	if refID == "" {
		return nil, &ErrorInvalidReferenceID
	}

	certObj, err := s.Store.GetCertificateByReference(refType, refID)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to get certificate by reference", log.String("refType", string(refType)),
			log.String("refID", refID), log.Error(err))
		return nil, &ErrorInternalServerError
	}
	if certObj == nil {
		logger.Debug("Certificate not found for reference", log.String("refType", string(refType)),
			log.String("refID", refID))
		return nil, &ErrorCertificateNotFound
	}

	return certObj, nil
}

// CreateCertificate creates a new certificate.
func (s *certificateService) CreateCertificate(cert *Certificate) (*Certificate,
	*serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if err := validateCertificateForCreation(cert); err != nil {
		return nil, err
	}

	// Check if a certificate with the same reference already exists
	existingCert, err := s.Store.GetCertificateByReference(cert.RefType, cert.RefID)
	if err != nil && !errors.Is(err, ErrCertificateNotFound) {
		logger.Error("Failed to check existing certificate", log.String("refType", string(cert.RefType)),
			log.String("refID", cert.RefID), log.Error(err))
		return nil, &ErrorInternalServerError
	}
	if existingCert != nil {
		return nil, &ErrorCertificateAlreadyExists
	}

	cert.ID = sysutils.GenerateUUID()
	err = s.Store.CreateCertificate(cert)
	if err != nil {
		logger.Error("Failed to create certificate", log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return cert, nil
}

// UpdateCertificateByID updates an existing certificate by its ID.
func (s *certificateService) UpdateCertificateByID(id string, cert *Certificate) (
	*Certificate, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if id == "" {
		return nil, &ErrorInvalidCertificateID
	}
	if err := validateCertificate(cert); err != nil {
		return nil, err
	}

	// Get the existing certificate to validate reference
	existingCert, err := s.Store.GetCertificateByID(id)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to get existing certificate", log.String("id", id), log.Error(err))
		return nil, &ErrorInternalServerError
	}
	if existingCert == nil {
		logger.Debug("Certificate not found for update", log.String("id", id))
		return nil, &ErrorCertificateNotFound
	}

	// Validate the reference is not changed
	if existingCert.RefType != cert.RefType || existingCert.RefID != cert.RefID {
		return nil, &ErrorReferenceUpdateIsNotAllowed
	}

	err = s.Store.UpdateCertificateByID(existingCert, cert)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to update certificate by ID", log.String("id", id), log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return cert, nil
}

// UpdateCertificateByReference updates an existing certificate by its reference type and ID.
func (s *certificateService) UpdateCertificateByReference(refType CertificateReferenceType,
	refID string, cert *Certificate) (*Certificate, *serviceerror.ServiceError) {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if !isValidReferenceType(refType) {
		return nil, &ErrorInvalidReferenceType
	}
	if refID == "" {
		return nil, &ErrorInvalidReferenceID
	}
	if err := validateCertificate(cert); err != nil {
		return nil, err
	}

	// Get the existing certificate to validate reference consistency
	existingCert, err := s.Store.GetCertificateByReference(refType, refID)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to get existing certificate", log.String("refType", string(refType)),
			log.String("refID", refID), log.Error(err))
		return nil, &ErrorInternalServerError
	}
	if existingCert == nil {
		logger.Debug("Certificate not found for update", log.String("refType", string(refType)),
			log.String("refID", refID))
		return nil, &ErrorCertificateNotFound
	}

	// Validate the reference is not changed
	if existingCert.RefType != cert.RefType || existingCert.RefID != cert.RefID {
		return nil, &ErrorReferenceUpdateIsNotAllowed
	}

	cert.ID = existingCert.ID
	err = s.Store.UpdateCertificateByReference(existingCert, cert)
	if err != nil {
		if errors.Is(err, ErrCertificateNotFound) {
			return nil, &ErrorCertificateNotFound
		}
		logger.Error("Failed to update certificate by reference", log.String("refType", string(refType)),
			log.String("refID", refID), log.Error(err))
		return nil, &ErrorInternalServerError
	}

	return cert, nil
}

// DeleteCertificateByID deletes a certificate by its ID.
func (s *certificateService) DeleteCertificateByID(id string) *serviceerror.ServiceError {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if id == "" {
		return &ErrorInvalidCertificateID
	}

	err := s.Store.DeleteCertificateByID(id)
	if err != nil {
		logger.Error("Failed to delete certificate by ID", log.String("id", id), log.Error(err))
		return &ErrorInternalServerError
	}

	return nil
}

// DeleteCertificateByReference deletes a certificate by its reference type and ID.
func (s *certificateService) DeleteCertificateByReference(refType CertificateReferenceType,
	refID string) *serviceerror.ServiceError {
	logger := log.GetLogger().With(log.String(log.LoggerKeyComponentName, loggerComponentName))

	if !isValidReferenceType(refType) {
		return &ErrorInvalidReferenceType
	}
	if refID == "" {
		return &ErrorInvalidReferenceID
	}

	err := s.Store.DeleteCertificateByReference(refType, refID)
	if err != nil {
		logger.Error("Failed to delete certificate by reference", log.String("refType", string(refType)),
			log.String("refID", refID), log.Error(err))
		return &ErrorInternalServerError
	}

	return nil
}

// isValidReferenceType checks if the provided reference type is valid.
func isValidReferenceType(refType CertificateReferenceType) bool {
	switch refType {
	case CertificateReferenceTypeApplication, CertificateReferenceTypeIDP:
		return true
	default:
		return false
	}
}

// isValidCertificateType checks if the provided certificate type is valid.
func isValidCertificateType(certType CertificateType) bool {
	switch certType {
	case CertificateTypeNone, CertificateTypeJWKS, CertificateTypeJWKSURI:
		return true
	default:
		return false
	}
}

// validateCertificate checks if the provided certificate is valid.
func validateCertificate(cert *Certificate) *serviceerror.ServiceError {
	if cert == nil {
		return &ErrorInvalidCertificateValue
	}
	if cert.ID == "" {
		return &ErrorInvalidCertificateID
	}
	if cert.RefID == "" {
		return &ErrorInvalidReferenceID
	}
	if !isValidReferenceType(cert.RefType) {
		return &ErrorInvalidReferenceType
	}
	if !isValidCertificateType(cert.Type) {
		return &ErrorInvalidCertificateType
	}
	if len(cert.Value) < 10 || len(cert.Value) > 4096 {
		return &ErrorInvalidCertificateValue
	}
	return nil
}

// validateCertificateForCreation checks if the provided certificate is valid for creation.
func validateCertificateForCreation(cert *Certificate) *serviceerror.ServiceError {
	if cert == nil {
		return &ErrorInvalidCertificateValue
	}
	if cert.RefID == "" {
		return &ErrorInvalidReferenceID
	}
	if !isValidReferenceType(cert.RefType) {
		return &ErrorInvalidReferenceType
	}
	if !isValidCertificateType(cert.Type) {
		return &ErrorInvalidCertificateType
	}
	if len(cert.Value) < 10 || len(cert.Value) > 4096 {
		return &ErrorInvalidCertificateValue
	}
	return nil
}
