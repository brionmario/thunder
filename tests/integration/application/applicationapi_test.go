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

package application

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"
)

const (
	testServerURL = "https://localhost:8095"
)

var (
	testApp = Application{
		Name:                      "Test App",
		Description:               "Test application for API testing",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://testapp.example.com",
		LogoURL:                   "https://testapp.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:                "test_app_client",
					ClientSecret:            "test_app_secret",
					RedirectURIs:            []string{"http://localhost/testapp/callback"},
					GrantTypes:              []string{"authorization_code", "client_credentials"},
					ResponseTypes:           []string{"code"},
					TokenEndpointAuthMethod: "client_secret_basic",
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}

	appToCreate = Application{
		Name:                      "App To Create",
		Description:               "Application to create for API testing",
		IsRegistrationFlowEnabled: true,
		URL:                       "https://apptocreate.example.com",
		LogoURL:                   "https://apptocreate.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:                "app_to_create_client",
					ClientSecret:            "app_to_create_secret",
					RedirectURIs:            []string{"http://localhost/apptocreate/callback"},
					GrantTypes:              []string{"authorization_code", "client_credentials"},
					ResponseTypes:           []string{"code"},
					TokenEndpointAuthMethod: "client_secret_basic",
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}

	appToUpdate = Application{
		Name:                      "Updated App",
		Description:               "Updated Description",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://appToUpdate.example.com",
		LogoURL:                   "https://appToUpdate.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:                "updated_client_id",
					ClientSecret:            "updated_secret",
					RedirectURIs:            []string{"http://localhost/callback2"},
					GrantTypes:              []string{"authorization_code"},
					ResponseTypes:           []string{"code"},
					TokenEndpointAuthMethod: "client_secret_basic",
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}
)

var (
	testAppID       string
	testAppInstance Application
)

type ApplicationAPITestSuite struct {
	suite.Suite
}

func TestApplicationAPITestSuite(t *testing.T) {

	suite.Run(t, new(ApplicationAPITestSuite))
}

// SetupSuite creates test applications for the test suite
func (ts *ApplicationAPITestSuite) SetupSuite() {
	// Create test application
	app1ID, err := createApplication(testApp)
	if err != nil {
		ts.T().Fatalf("Failed to create test application during setup: %v", err)
	}
	testAppID = app1ID

	// Build the test app structure for validations
	testAppInstance = testApp
	testAppInstance.ID = testAppID
	if len(testAppInstance.InboundAuthConfig) > 0 && testAppInstance.InboundAuthConfig[0].OAuthAppConfig != nil {
		testAppInstance.ClientID = testAppInstance.InboundAuthConfig[0].OAuthAppConfig.ClientID
	}
}

// TearDownSuite cleans up test applications
func (ts *ApplicationAPITestSuite) TearDownSuite() {
	// Delete the test application
	if testAppID != "" {
		err := deleteApplication(testAppID)
		if err != nil {
			ts.T().Logf("Failed to delete test application during teardown: %v", err)
		}
	}
}

// Test application listing
func (ts *ApplicationAPITestSuite) TestApplicationListing() {

	req, err := http.NewRequest("GET", testServerURL+"/applications", nil)
	if err != nil {
		ts.T().Fatalf("Failed to create request: %v", err)
	}

	// Configure the HTTP client to skip TLS verification
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // Skip certificate verification
		},
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Validate the response
	if resp.StatusCode != http.StatusOK {
		ts.T().Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	// Parse the response body
	var appList ApplicationList
	err = json.NewDecoder(resp.Body).Decode(&appList)
	if err != nil {
		ts.T().Fatalf("Failed to parse response body: %v", err)
	}

	totalResults := appList.TotalResults
	if totalResults == 0 {
		ts.T().Fatalf("Response does not contain a valid total results count")
	}

	appCount := appList.Count
	if appCount == 0 {
		ts.T().Fatalf("Response does not contain a valid application count")
	}

	applicationListLength := len(appList.Applications)
	if applicationListLength == 0 {
		ts.T().Fatalf("Response does not contain any applications")
	}

	// Verify that the test application is present in the list
	testApps := []Application{testAppInstance}
	for _, expectedApp := range testApps {
		found := false
		for _, app := range appList.Applications {
			if app.ID == expectedApp.ID &&
				app.Name == expectedApp.Name &&
				app.Description == expectedApp.Description &&
				app.ClientID == expectedApp.ClientID {
				found = true
				break
			}
		}
		if !found {
			ts.T().Fatalf("Test application not found in list: %+v", expectedApp)
		}
	}
}

// Test application get by ID
func (ts *ApplicationAPITestSuite) TestApplicationGetByID() {
	// Create an application for get testing
	appID, err := createApplication(appToCreate)
	if err != nil {
		ts.T().Fatalf("Failed to create application for get test: %v", err)
	}
	defer func() {
		// Clean up the created application
		if err := deleteApplication(appID); err != nil {
			ts.T().Logf("Failed to delete application after get test: %v", err)
		}
	}()

	// Build the expected app structure for validation
	expectedApp := appToCreate
	expectedApp.ID = appID
	if len(expectedApp.InboundAuthConfig) > 0 && expectedApp.InboundAuthConfig[0].OAuthAppConfig != nil {
		expectedApp.ClientID = expectedApp.InboundAuthConfig[0].OAuthAppConfig.ClientID
	}
	
	retrieveAndValidateApplicationDetails(ts, expectedApp)
}

// Test application update
func (ts *ApplicationAPITestSuite) TestApplicationUpdate() {
	// Create an application for update testing
	appID, err := createApplication(appToCreate)
	if err != nil {
		ts.T().Fatalf("Failed to create application for update test: %v", err)
	}
	defer func() {
		// Clean up the created application
		if err := deleteApplication(appID); err != nil {
			ts.T().Logf("Failed to delete application after update test: %v", err)
		}
	}()

	// Add the ID to the application to update
	appToUpdateWithID := appToUpdate
	appToUpdateWithID.ID = appID

	appJSON, err := json.Marshal(appToUpdateWithID)
	if err != nil {
		ts.T().Fatalf("Failed to marshal appToUpdate: %v", err)
	}

	reqBody := bytes.NewReader(appJSON)
	req, err := http.NewRequest("PUT", testServerURL+"/applications/"+appID, reqBody)
	if err != nil {
		ts.T().Fatalf("Failed to create update request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send update request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	// For update operations, verify the response directly
	var updatedApp Application
	if err = json.NewDecoder(resp.Body).Decode(&updatedApp); err != nil {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Failed to decode update response: %v. Response: %s", err, string(responseBody))
	}

	// Client secret should be present in the update response
	if len(updatedApp.InboundAuthConfig) > 0 &&
		updatedApp.InboundAuthConfig[0].OAuthAppConfig != nil &&
		updatedApp.InboundAuthConfig[0].OAuthAppConfig.ClientSecret == "" {
		ts.T().Fatalf("Expected client secret in update response but got empty string")
	}

	// Now validate by getting the application (which should not have client secret)
	// Make sure client ID is properly set in the root level before validation
	if len(appToUpdateWithID.InboundAuthConfig) > 0 &&
		appToUpdateWithID.InboundAuthConfig[0].OAuthAppConfig != nil {
		appToUpdateWithID.ClientID = appToUpdateWithID.InboundAuthConfig[0].OAuthAppConfig.ClientID
	}

	retrieveAndValidateApplicationDetails(ts, appToUpdateWithID)
}

func retrieveAndValidateApplicationDetails(ts *ApplicationAPITestSuite, expectedApp Application) {

	req, err := http.NewRequest("GET", testServerURL+"/applications/"+expectedApp.ID, nil)
	if err != nil {
		ts.T().Fatalf("Failed to create request: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	// Check if the response Content-Type is application/json
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		ts.T().Fatalf("Expected Content-Type application/json, got %s", contentType)
	}

	var app Application
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &app)
	if err != nil {
		ts.T().Fatalf("Failed to parse response body: %v\nResponse body: %s", err, string(body))
	}

	// For GET operations, client secret should be empty in the response
	// Make sure expectedApp has client secret cleared for proper comparison
	appForComparison := expectedApp
	if len(appForComparison.InboundAuthConfig) > 0 && appForComparison.InboundAuthConfig[0].OAuthAppConfig != nil {
		// Make sure client ID is in root object
		appForComparison.ClientID = appForComparison.InboundAuthConfig[0].OAuthAppConfig.ClientID
		// Remove client secret for GET comparison
		appForComparison.InboundAuthConfig[0].OAuthAppConfig.ClientSecret = ""
	}

	// Ensure certificate is set in expected app if it's null
	if appForComparison.Certificate == nil {
		appForComparison.Certificate = &ApplicationCert{
			Type:  "NONE",
			Value: "",
		}
	}

	if !app.equals(appForComparison) {
		appJSON, _ := json.MarshalIndent(app, "", "  ")
		expectedJSON, _ := json.MarshalIndent(appForComparison, "", "  ")
		ts.T().Fatalf("Application mismatch:\nGot:\n%s\n\nExpected:\n%s", string(appJSON), string(expectedJSON))
	}
}

func createApplication(app Application) (string, error) {
	appJSON, err := json.Marshal(app)
	if err != nil {
		return "", fmt.Errorf("failed to marshal application: %w", err)
	}

	reqBody := bytes.NewReader(appJSON)
	req, err := http.NewRequest("POST", testServerURL+"/applications", reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		responseBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("expected status 201, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	// For create operations, directly parse the response to a full Application
	var createdApp Application
	err = json.NewDecoder(resp.Body).Decode(&createdApp)
	if err != nil {
		return "", fmt.Errorf("failed to parse response body: %w", err)
	}

	// Verify client secret is present in the create response
	if len(createdApp.InboundAuthConfig) > 0 &&
		createdApp.InboundAuthConfig[0].OAuthAppConfig != nil &&
		createdApp.InboundAuthConfig[0].OAuthAppConfig.ClientSecret == "" {
		return "", fmt.Errorf("expected client secret in create response but got empty string")
	}

	id := createdApp.ID
	if id == "" {
		return "", fmt.Errorf("response does not contain id")
	}
	return id, nil
}

func deleteApplication(appID string) error {
	req, err := http.NewRequest("DELETE", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		responseBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("expected status 204, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}
	return nil
}

// TestApplicationCreationWithDefaults tests that applications created without grant_types, response_types, or token_endpoint_auth_method get proper defaults
func (ts *ApplicationAPITestSuite) TestApplicationCreationWithDefaults() {
	appWithDefaults := Application{
		Name:                      "App With Defaults",
		Description:               "Application to test default values",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://defaults.example.com",
		LogoURL:                   "https://defaults.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:     "defaults_app_client",
					ClientSecret: "defaults_app_secret",
					RedirectURIs: []string{"http://localhost/defaults/callback"},
					// Intentionally omitting GrantTypes, ResponseTypes, and TokenEndpointAuthMethod
					PKCERequired: false,
					PublicClient: false,
				},
			},
		},
	}

	appID, err := createApplication(appWithDefaults)
	if err != nil {
		ts.T().Fatalf("Failed to create application: %v", err)
	}

	req, err := http.NewRequest("GET", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		ts.T().Fatalf("Failed to create GET request: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	var retrievedApp Application
	err = json.NewDecoder(resp.Body).Decode(&retrievedApp)
	if err != nil {
		ts.T().Fatalf("Failed to decode response: %v", err)
	}

	// Verify defaults were applied
	if len(retrievedApp.InboundAuthConfig) > 0 && retrievedApp.InboundAuthConfig[0].OAuthAppConfig != nil {
		oauthConfig := retrievedApp.InboundAuthConfig[0].OAuthAppConfig

		ts.Assert().Equal([]string{"authorization_code"}, oauthConfig.GrantTypes, "Default grant_types should be ['authorization_code']")
		ts.Assert().Equal([]string{"code"}, oauthConfig.ResponseTypes, "Default response_types should be ['code']")
		ts.Assert().Equal("client_secret_basic", oauthConfig.TokenEndpointAuthMethod, "Default token_endpoint_auth_method should be 'client_secret_basic'")
	}

	err = deleteApplication(appID)
	if err != nil {
		ts.T().Logf("Failed to delete test application: %v", err)
	}
}

// TestApplicationCreationWithInvalidTokenEndpointAuthMethod tests validation of invalid token_endpoint_auth_method values
func (ts *ApplicationAPITestSuite) TestApplicationCreationWithInvalidTokenEndpointAuthMethod() {
	appWithInvalidAuthMethod := Application{
		Name:                      "App With Invalid Auth Method",
		Description:               "Application to test invalid token endpoint auth method",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://invalid.example.com",
		LogoURL:                   "https://invalid.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:                "invalid_auth_app_client",
					ClientSecret:            "invalid_auth_app_secret",
					RedirectURIs:            []string{"http://localhost/invalid/callback"},
					GrantTypes:              []string{"authorization_code"},
					ResponseTypes:           []string{"code"},
					TokenEndpointAuthMethod: "invalid_auth_method", // Invalid value
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}

	_, err := createApplication(appWithInvalidAuthMethod)
	if err == nil {
		ts.T().Fatalf("Expected validation error for invalid token_endpoint_auth_method, but application was created successfully")
	}

	appWithEmptyAuthMethod := Application{
		Name:                      "App With Empty Auth Method",
		Description:               "Application to test empty token endpoint auth method",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://empty.example.com",
		LogoURL:                   "https://empty.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:                "empty_auth_app_client",
					ClientSecret:            "empty_auth_app_secret",
					RedirectURIs:            []string{"http://localhost/empty/callback"},
					GrantTypes:              []string{"authorization_code"},
					ResponseTypes:           []string{"code"},
					TokenEndpointAuthMethod: "",
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}

	appID, err := createApplication(appWithEmptyAuthMethod)
	if err != nil {
		ts.T().Fatalf("Failed to create application with empty token_endpoint_auth_method: %v", err)
	}

	req, err := http.NewRequest("GET", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		ts.T().Fatalf("Failed to create GET request: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	var retrievedApp Application
	err = json.NewDecoder(resp.Body).Decode(&retrievedApp)
	if err != nil {
		ts.T().Fatalf("Failed to decode response: %v", err)
	}

	if len(retrievedApp.InboundAuthConfig) > 0 && retrievedApp.InboundAuthConfig[0].OAuthAppConfig != nil {
		oauthConfig := retrievedApp.InboundAuthConfig[0].OAuthAppConfig
		ts.Assert().Equal("client_secret_basic", oauthConfig.TokenEndpointAuthMethod, "Empty token_endpoint_auth_method should get default 'client_secret_basic'")
	}

	err = deleteApplication(appID)
	if err != nil {
		ts.T().Logf("Failed to delete test application: %v", err)
	}
}

// TestApplicationCreationWithPartialDefaults tests applications with some fields missing (partial defaults)
func (ts *ApplicationAPITestSuite) TestApplicationCreationWithPartialDefaults() {
	appWithPartialDefaults := Application{
		Name:                      "App With Partial Defaults",
		Description:               "Application to test partial default values",
		IsRegistrationFlowEnabled: false,
		URL:                       "https://partial.example.com",
		LogoURL:                   "https://partial.example.com/logo.png",
		AuthFlowGraphID:           "auth_flow_config_basic",
		RegistrationFlowGraphID:   "registration_flow_config_basic",
		Certificate: &ApplicationCert{
			Type:  "NONE",
			Value: "",
		},
		InboundAuthConfig: []InboundAuthConfig{
			{
				Type: "oauth2",
				OAuthAppConfig: &OAuthAppConfig{
					ClientID:     "partial_app_client",
					ClientSecret: "partial_app_secret",
					RedirectURIs: []string{"http://localhost/partial/callback"},
					// GrantTypes missing - should get default
					ResponseTypes:           []string{"code"},     // Explicitly set
					TokenEndpointAuthMethod: "client_secret_post", // Explicitly set
					PKCERequired:            false,
					PublicClient:            false,
				},
			},
		},
	}

	appID, err := createApplication(appWithPartialDefaults)
	if err != nil {
		ts.T().Fatalf("Failed to create application: %v", err)
	}

	// Verify that defaults were applied by getting the application
	req, err := http.NewRequest("GET", testServerURL+"/applications/"+appID, nil)
	if err != nil {
		ts.T().Fatalf("Failed to create GET request: %v", err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		ts.T().Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		ts.T().Fatalf("Expected status 200, got %d. Response: %s", resp.StatusCode, string(responseBody))
	}

	var retrievedApp Application
	err = json.NewDecoder(resp.Body).Decode(&retrievedApp)
	if err != nil {
		ts.T().Fatalf("Failed to decode response: %v", err)
	}

	if len(retrievedApp.InboundAuthConfig) > 0 && retrievedApp.InboundAuthConfig[0].OAuthAppConfig != nil {
		oauthConfig := retrievedApp.InboundAuthConfig[0].OAuthAppConfig

		ts.Assert().Equal([]string{"authorization_code"}, oauthConfig.GrantTypes, "Missing grant_types should get default ['authorization_code']")
		ts.Assert().Equal([]string{"code"}, oauthConfig.ResponseTypes, "Explicitly set response_types should be preserved")
		ts.Assert().Equal("client_secret_post", oauthConfig.TokenEndpointAuthMethod, "Explicitly set token_endpoint_auth_method should be preserved")
	}

	err = deleteApplication(appID)
	if err != nil {
		ts.T().Logf("Failed to delete test application: %v", err)
	}
}
