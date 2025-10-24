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

package authz

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type AuthorizeServiceTestSuite struct {
	suite.Suite
	service        AuthorizeServiceInterface
	mockAuthzStore *AuthorizationCodeStoreInterfaceMock
	testAuthzCode  AuthorizationCode
	testClientID   string
	testCode       string
}

func TestAuthorizeServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthorizeServiceTestSuite))
}

func (suite *AuthorizeServiceTestSuite) SetupTest() {
	suite.mockAuthzStore = NewAuthorizationCodeStoreInterfaceMock(suite.T())
	suite.service = newAuthorizeService(suite.mockAuthzStore)

	suite.testClientID = "test-client-id"
	suite.testCode = "test-auth-code"

	suite.testAuthzCode = AuthorizationCode{
		CodeID:           "test-code-id",
		Code:             suite.testCode,
		ClientID:         suite.testClientID,
		RedirectURI:      "https://client.example.com/callback",
		AuthorizedUserID: "test-user-id",
		TimeCreated:      time.Now().Add(-5 * time.Minute),
		ExpiryTime:       time.Now().Add(5 * time.Minute),
		Scopes:           "read write",
		State:            AuthCodeStateActive,
	}
}

func (suite *AuthorizeServiceTestSuite) TestNewAuthorizeService() {
	service := newAuthorizeService(suite.mockAuthzStore)
	assert.NotNil(suite.T(), service)
	assert.Implements(suite.T(), (*AuthorizeServiceInterface)(nil), service)
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_Success() {
	// Mock store to return valid authorization code
	suite.mockAuthzStore.On("GetAuthorizationCode", suite.testClientID, suite.testCode).
		Return(suite.testAuthzCode, nil)
	suite.mockAuthzStore.On("DeactivateAuthorizationCode", suite.testAuthzCode).
		Return(nil)

	result, err := suite.service.GetAuthorizationCodeDetails(suite.testClientID, suite.testCode)

	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), suite.testAuthzCode.Code, result.Code)
	assert.Equal(suite.T(), suite.testAuthzCode.ClientID, result.ClientID)
	assert.Equal(suite.T(), suite.testAuthzCode.AuthorizedUserID, result.AuthorizedUserID)

	suite.mockAuthzStore.AssertExpectations(suite.T())
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_StoreError() {
	// Mock store to return error
	suite.mockAuthzStore.On("GetAuthorizationCode", suite.testClientID, suite.testCode).
		Return(AuthorizationCode{}, errors.New("database error"))

	result, err := suite.service.GetAuthorizationCodeDetails(suite.testClientID, suite.testCode)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "invalid authorization code", err.Error())

	suite.mockAuthzStore.AssertExpectations(suite.T())
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_EmptyCode() {
	// Mock store to return authorization code with empty code string
	emptyAuthzCode := suite.testAuthzCode
	emptyAuthzCode.Code = ""

	suite.mockAuthzStore.On("GetAuthorizationCode", suite.testClientID, suite.testCode).
		Return(emptyAuthzCode, nil)

	result, err := suite.service.GetAuthorizationCodeDetails(suite.testClientID, suite.testCode)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "invalid authorization code", err.Error())

	suite.mockAuthzStore.AssertExpectations(suite.T())
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_DeactivationError() {
	// Mock store to return valid code but fail on deactivation
	suite.mockAuthzStore.On("GetAuthorizationCode", suite.testClientID, suite.testCode).
		Return(suite.testAuthzCode, nil)
	suite.mockAuthzStore.On("DeactivateAuthorizationCode", suite.testAuthzCode).
		Return(errors.New("deactivation failed"))

	result, err := suite.service.GetAuthorizationCodeDetails(suite.testClientID, suite.testCode)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed to invalidate authorization code", err.Error())

	suite.mockAuthzStore.AssertExpectations(suite.T())
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_EmptyClientID() {
	// Mock store to be called with empty client ID
	suite.mockAuthzStore.On("GetAuthorizationCode", "", suite.testCode).
		Return(AuthorizationCode{}, errors.New("invalid client"))

	result, err := suite.service.GetAuthorizationCodeDetails("", suite.testCode)

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "invalid authorization code", err.Error())

	suite.mockAuthzStore.AssertExpectations(suite.T())
}

func (suite *AuthorizeServiceTestSuite) TestGetAuthorizationCodeDetails_EmptyCodeString() {
	// Mock store to be called with empty code string
	suite.mockAuthzStore.On("GetAuthorizationCode", suite.testClientID, "").
		Return(AuthorizationCode{}, errors.New("invalid code"))

	result, err := suite.service.GetAuthorizationCodeDetails(suite.testClientID, "")

	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "invalid authorization code", err.Error())

	suite.mockAuthzStore.AssertExpectations(suite.T())
}
