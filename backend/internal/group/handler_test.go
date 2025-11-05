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

package group

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/asgardeo/thunder/internal/system/config"
	serverconst "github.com/asgardeo/thunder/internal/system/constants"
	"github.com/asgardeo/thunder/internal/system/error/apierror"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/internal/system/log"
)

type flakyResponseWriter struct {
	*httptest.ResponseRecorder
	failNext bool
}

type GroupHandlerTestSuite struct {
	suite.Suite
}

func TestGroupHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GroupHandlerTestSuite))
}

func (suite *GroupHandlerTestSuite) SetupTest() {
	config.ResetThunderRuntime()

	err := config.InitializeThunderRuntime("", &config.Config{})
	suite.Require().NoError(err)
}

func (suite *GroupHandlerTestSuite) TearDownTest() {
	config.ResetThunderRuntime()
}

func (suite *GroupHandlerTestSuite) ensureRuntime() {
	config.ResetThunderRuntime()

	err := config.InitializeThunderRuntime("", &config.Config{})
	suite.Require().NoError(err)
}

func newFlakyResponseWriter() *flakyResponseWriter {
	return &flakyResponseWriter{
		ResponseRecorder: httptest.NewRecorder(),
		failNext:         true,
	}
}

func (w *flakyResponseWriter) Write(b []byte) (int, error) {
	if w.failNext {
		w.failNext = false
		return 0, errors.New("write failure")
	}
	return w.ResponseRecorder.Write(b)
}

func mapStringToValues(m map[string]string) url.Values {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, v)
	}
	return values
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesOptionsGroups() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	registerRoutes(mux, newGroupHandler(nil))

	req := httptest.NewRequest(http.MethodOptions, "/groups", nil)
	resp := httptest.NewRecorder()

	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesGroupIDDispatch() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)
	registerRoutes(mux, handler)

	serviceMock.
		On("GetGroup", "grp-001").
		Return(&Group{ID: "grp-001"}, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesGroupMembersDispatch() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)
	registerRoutes(mux, handler)

	serviceMock.
		On("GetGroupMembers", "grp-001", serverconst.DefaultPageSize, 0).
		Return(&MemberListResponse{}, nil).
		Once()

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesGroupIDNotFoundPath() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	registerRoutes(mux, newGroupHandler(nil))

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/unknown", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNotFound, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesOptionsGroupID() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	registerRoutes(mux, newGroupHandler(nil))

	req := httptest.NewRequest(http.MethodOptions, "/groups/grp-001", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestRegisterRoutesOptionsTreePath() {
	t := suite.T()
	suite.ensureRuntime()
	mux := http.NewServeMux()
	registerRoutes(mux, newGroupHandler(nil))

	req := httptest.NewRequest(http.MethodOptions, "/groups/tree/root", nil)
	resp := httptest.NewRecorder()
	mux.ServeHTTP(resp, req)

	require.Equal(t, http.StatusNoContent, resp.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListRequestSuccess() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups?limit=3&offset=2", nil)
	rr := httptest.NewRecorder()

	expected := &GroupListResponse{
		TotalResults: 5,
		StartIndex:   3,
		Count:        2,
		Groups: []GroupBasic{
			{ID: "g1", Name: "group-1"},
			{ID: "g2", Name: "group-2"},
		},
	}

	serviceMock.
		On("GetGroupList", 3, 2).
		Return(expected, nil).
		Once()

	handler.HandleGroupListRequest(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, serverconst.ContentTypeJSON, rr.Header().Get(serverconst.ContentTypeHeaderName))

	var body GroupListResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, expected.TotalResults, body.TotalResults)
	require.Equal(t, expected.Count, body.Count)
	require.Len(t, body.Groups, 2)
	require.Equal(t, expected.Groups[0].Name, body.Groups[0].Name)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListRequestInvalidLimit() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups?limit=invalid", nil)
	rr := httptest.NewRecorder()

	handler.HandleGroupListRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, serverconst.ContentTypeJSON, rr.Header().Get(serverconst.ContentTypeHeaderName))

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidLimit.Code, body.Code)
	require.Equal(t, ErrorInvalidLimit.Error, body.Message)

	serviceMock.AssertNotCalled(t, "GetGroupList", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups", nil)
	writer := newFlakyResponseWriter()

	serviceMock.
		On("GetGroupList", serverconst.DefaultPageSize, 0).
		Return(&GroupListResponse{}, nil).
		Once()

	handler.HandleGroupListRequest(writer, req)

	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListRequestClientErrorEncodingFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups?limit=invalid", nil)
	writer := newFlakyResponseWriter()

	handler.HandleGroupListRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "GetGroupList", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListRequestServiceError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups", nil)
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroupList", serverconst.DefaultPageSize, 0).
		Return((*GroupListResponse)(nil), &ErrorInternalServerError).
		Once()

	handler.HandleGroupListRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListByPathRequestSuccess() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/ous/root/groups", nil)
	req.SetPathValue("path", "root")
	rr := httptest.NewRecorder()

	expected := &GroupListResponse{
		TotalResults: 1,
		StartIndex:   1,
		Count:        1,
		Groups: []GroupBasic{
			{ID: "g1", Name: "root-group"},
		},
	}

	serviceMock.
		On("GetGroupsByPath", "root", serverconst.DefaultPageSize, 0).
		Return(expected, nil).
		Once()

	handler.HandleGroupListByPathRequest(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var body GroupListResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, expected.TotalResults, body.TotalResults)
	require.Equal(t, expected.Groups[0].Name, body.Groups[0].Name)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListByPathRequestMissingPath() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/ous//groups", nil)
	rr := httptest.NewRecorder()

	handler.HandleGroupListByPathRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidRequestFormat.Code, body.Code)

	serviceMock.AssertNotCalled(t, "GetGroupsByPath", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListByPathRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/ous/root/groups", nil)
	req.SetPathValue("path", "root")
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroupsByPath", "root", serverconst.DefaultPageSize, 0).
		Return(nil, &ErrorInternalServerError).
		Once()

	handler.HandleGroupListByPathRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListByPathRequestPaginationError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/ous/root/groups?limit=invalid", nil)
	req.SetPathValue("path", "root")
	rr := httptest.NewRecorder()

	handler.HandleGroupListByPathRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidLimit.Code, body.Code)
	serviceMock.AssertNotCalled(t, "GetGroupsByPath", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupListByPathRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/ous/root/groups", nil)
	req.SetPathValue("path", "root")
	writer := newFlakyResponseWriter()

	serviceMock.
		On("GetGroupsByPath", "root", serverconst.DefaultPageSize, 0).
		Return(&GroupListResponse{}, nil).
		Once()

	handler.HandleGroupListByPathRequest(writer, req)

	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestInvalidJSON() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader("{invalid json"))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	handler.HandleGroupPostRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, serverconst.ContentTypeJSON, rr.Header().Get(serverconst.ContentTypeHeaderName))

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidRequestFormat.Code, body.Code)
	require.Contains(t, body.Description, "Failed to parse request body")

	serviceMock.AssertNotCalled(t, "CreateGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestSanitizesPayload() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{
		"name": "  Team <script> ",
		"description": " desc ",
		"organizationUnitId": " ou-001 ",
		"members": [
			{"id": " member-1 ", "type": "user"}
		]
	}`

	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(payload))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroup", mock.MatchedBy(func(request CreateGroupRequest) bool {
			return request.Name == "Team &lt;script&gt;" &&
				request.Description == "desc" &&
				request.OrganizationUnitID == "ou-001" &&
				len(request.Members) == 1 &&
				request.Members[0].ID == "member-1" &&
				request.Members[0].Type == MemberTypeUser
		})).
		Return(&Group{ID: "grp-001", Name: "Team &lt;script&gt;", OrganizationUnitID: "ou-001"}, nil).
		Once()

	handler.HandleGroupPostRequest(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)
	require.Equal(t, serverconst.ContentTypeJSON, rr.Header().Get(serverconst.ContentTypeHeaderName))

	var body Group
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, "grp-001", body.ID)
	require.Equal(t, "Team &lt;script&gt;", body.Name)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestServiceError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{"name":"group","organizationUnitId":"ou"}`
	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(payload))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroup", mock.AnythingOfType("group.CreateGroupRequest")).
		Return(nil, &ErrorGroupNameConflict).
		Once()

	handler.HandleGroupPostRequest(rr, req)

	require.Equal(t, http.StatusConflict, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorGroupNameConflict.Code, body.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{"name":"group","organizationUnitId":"ou"}`
	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(payload))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroup", mock.AnythingOfType("group.CreateGroupRequest")).
		Return((*Group)(nil), &ErrorInternalServerError).
		Once()

	handler.HandleGroupPostRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{"name":"team","organizationUnitId":"ou"}`
	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader(payload))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	serviceMock.
		On("CreateGroup", mock.MatchedBy(func(request CreateGroupRequest) bool {
			return request.Name == "team" && request.OrganizationUnitID == "ou"
		})).
		Return(&Group{ID: "grp-001", Name: "team"}, nil).
		Once()

	handler.HandleGroupPostRequest(writer, req)

	require.Equal(t, http.StatusCreated, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/groups", strings.NewReader("{"))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	handler.HandleGroupPostRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "CreateGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestSuccess() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{"name":"name"}`
	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader(payload))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroupByPath", "root", CreateGroupByPathRequest{Name: "name"}).
		Return(&Group{ID: "grp-001", Name: "name"}, nil).
		Once()

	handler.HandleGroupPostByPathRequest(rr, req)

	require.Equal(t, http.StatusCreated, rr.Code)

	var body Group
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, "grp-001", body.ID)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestInvalidJSON() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader("{"))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	handler.HandleGroupPostByPathRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidRequestFormat.Code, body.Code)

	serviceMock.AssertNotCalled(t, "CreateGroupByPath", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestInvalidPath() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous//groups", strings.NewReader(`{"name":"n"}`))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	handler.HandleGroupPostByPathRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	serviceMock.AssertNotCalled(t, "CreateGroupByPath", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestServiceError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader(`{"name":"n"}`))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroupByPath", "root", CreateGroupByPathRequest{Name: "n"}).
		Return(nil, &ErrorGroupNotFound).
		Once()

	handler.HandleGroupPostByPathRequest(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader(`{"name":"team"}`))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	serviceMock.
		On("CreateGroupByPath", "root", CreateGroupByPathRequest{Name: "team"}).
		Return(&Group{ID: "grp-001", Name: "team"}, nil).
		Once()

	handler.HandleGroupPostByPathRequest(writer, req)

	require.Equal(t, http.StatusCreated, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader("{"))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	handler.HandleGroupPostByPathRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "CreateGroupByPath", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPostByPathRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPost, "/ous/root/groups", strings.NewReader(`{"name":"n"}`))
	req.SetPathValue("path", "root")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("CreateGroupByPath", "root", CreateGroupByPathRequest{Name: "n"}).
		Return((*Group)(nil), &ErrorInternalServerError).
		Once()

	handler.HandleGroupPostByPathRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupGetRequestNotFound() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-404", nil)
	req.SetPathValue("id", "grp-404")
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroup", "grp-404").
		Return(nil, &ErrorGroupNotFound).
		Once()

	handler.HandleGroupGetRequest(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorGroupNotFound.Code, body.Code)
	require.Equal(t, ErrorGroupNotFound.Error, body.Message)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupGetRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroup", "grp-001").
		Return((*Group)(nil), &ErrorInternalServerError).
		Once()

	handler.HandleGroupGetRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupGetRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001", nil)
	req.SetPathValue("id", "grp-001")
	writer := newFlakyResponseWriter()

	serviceMock.
		On("GetGroup", "grp-001").
		Return(&Group{ID: "grp-001"}, nil).
		Once()

	handler.HandleGroupGetRequest(writer, req)

	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupGetRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/", nil)
	writer := newFlakyResponseWriter()

	handler.HandleGroupGetRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "GetGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupGetRequestMissingID() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/", nil)
	rr := httptest.NewRecorder()

	handler.HandleGroupGetRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorMissingGroupID.Code, body.Code)

	serviceMock.AssertNotCalled(t, "GetGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestSanitizesPayload() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	payload := `{
		"name":"  Updated <payload> ",
		"description": " desc ",
		"organizationUnitId": " ou-123 ",
		"members": [{"id":" member-1 ","type":"group"}]
	}`

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001", strings.NewReader(payload))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("UpdateGroup", "grp-001", mock.MatchedBy(func(request UpdateGroupRequest) bool {
			return request.Name == "Updated &lt;payload&gt;" &&
				request.Description == "desc" &&
				request.OrganizationUnitID == "ou-123" &&
				len(request.Members) == 1 &&
				request.Members[0].ID == "member-1" &&
				request.Members[0].Type == MemberTypeGroup
		})).
		Return(&Group{ID: "grp-001", Name: "Updated &lt;payload&gt;"}, nil).
		Once()

	handler.HandleGroupPutRequest(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)

	var body Group
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, "grp-001", body.ID)
	require.Equal(t, "Updated &lt;payload&gt;", body.Name)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestInvalidJSON() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001", strings.NewReader("{"))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	handler.HandleGroupPutRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorInvalidRequestFormat.Code, body.Code)

	serviceMock.AssertNotCalled(t, "UpdateGroup", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestDecodeErrorEncodingFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001", strings.NewReader("{"))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	handler.HandleGroupPutRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "UpdateGroup", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestMissingID() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/", strings.NewReader("{}"))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	handler.HandleGroupPutRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorMissingGroupID.Code, body.Code)

	serviceMock.AssertNotCalled(t, "UpdateGroup", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestServiceError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001",
		strings.NewReader(`{"name":"n","organizationUnitId":"ou"}`))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("UpdateGroup", "grp-001", mock.AnythingOfType("group.UpdateGroupRequest")).
		Return(nil, &ErrorGroupNotFound).
		Once()

	handler.HandleGroupPutRequest(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001",
		strings.NewReader(`{"name":"n","organizationUnitId":"ou"}`))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	rr := httptest.NewRecorder()

	serviceMock.
		On("UpdateGroup", "grp-001", mock.AnythingOfType("group.UpdateGroupRequest")).
		Return(nil, &ErrorInternalServerError).
		Once()

	handler.HandleGroupPutRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/grp-001",
		strings.NewReader(`{"name":"team","organizationUnitId":"ou"}`))
	req.SetPathValue("id", "grp-001")
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	serviceMock.
		On("UpdateGroup", "grp-001", mock.AnythingOfType("group.UpdateGroupRequest")).
		Return(&Group{ID: "grp-001"}, nil).
		Once()

	handler.HandleGroupPutRequest(writer, req)

	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupPutRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodPut, "/groups/", strings.NewReader("{}"))
	req.Header.Set(serverconst.ContentTypeHeaderName, serverconst.ContentTypeJSON)
	writer := newFlakyResponseWriter()

	handler.HandleGroupPutRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "UpdateGroup", mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupDeleteRequestMissingID() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodDelete, "/groups/", nil)
	rr := httptest.NewRecorder()

	handler.HandleGroupDeleteRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorMissingGroupID.Code, body.Code)

	serviceMock.AssertNotCalled(t, "DeleteGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupDeleteRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodDelete, "/groups/", nil)
	writer := newFlakyResponseWriter()

	handler.HandleGroupDeleteRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "DeleteGroup", mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupDeleteRequestConflict() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodDelete, "/groups/grp-001", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("DeleteGroup", "grp-001").
		Return(&ErrorCannotDeleteGroup).
		Once()

	handler.HandleGroupDeleteRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorCannotDeleteGroup.Code, body.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupDeleteRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodDelete, "/groups/grp-001", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("DeleteGroup", "grp-001").
		Return(&ErrorInternalServerError).
		Once()

	handler.HandleGroupDeleteRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupDeleteRequestSuccess() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodDelete, "/groups/grp-001", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("DeleteGroup", "grp-001").
		Return(nil).
		Once()

	handler.HandleGroupDeleteRequest(rr, req)

	require.Equal(t, http.StatusNoContent, rr.Code)
	require.Empty(t, rr.Body.String())
}
func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestSuccess() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members?limit=2&offset=1", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	expected := &MemberListResponse{
		TotalResults: 3,
		StartIndex:   2,
		Count:        2,
		Members: []Member{
			{ID: "usr-1", Type: MemberTypeUser},
			{ID: "grp-2", Type: MemberTypeGroup},
		},
	}

	serviceMock.
		On("GetGroupMembers", "grp-001", 2, 1).
		Return(expected, nil).
		Once()

	handler.HandleGroupMembersGetRequest(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, serverconst.ContentTypeJSON, rr.Header().Get(serverconst.ContentTypeHeaderName))

	var body MemberListResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, expected.TotalResults, body.TotalResults)
	require.Len(t, body.Members, 2)
	require.Equal(t, expected.Members[0].ID, body.Members[0].ID)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestInvalidLimit() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members?limit=NaN", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	handler.HandleGroupMembersGetRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	serviceMock.AssertNotCalled(t, "GetGroupMembers", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestServiceError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroupMembers", "grp-001", serverconst.DefaultPageSize, 0).
		Return(nil, &ErrorGroupNotFound).
		Once()

	handler.HandleGroupMembersGetRequest(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestEncodeError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members", nil)
	req.SetPathValue("id", "grp-001")
	writer := newFlakyResponseWriter()

	serviceMock.
		On("GetGroupMembers", "grp-001", serverconst.DefaultPageSize, 0).
		Return(&MemberListResponse{}, nil).
		Once()

	handler.HandleGroupMembersGetRequest(writer, req)

	require.Equal(t, http.StatusOK, writer.Code)
	require.Equal(t, "Failed to encode response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestErrorResponseEncodeFailure() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups//members", nil)
	writer := newFlakyResponseWriter()

	handler.HandleGroupMembersGetRequest(writer, req)

	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
	serviceMock.AssertNotCalled(t, "GetGroupMembers", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestParsePaginationParamsInvalidOffset() {
	t := suite.T()
	limit, offset, err := parsePaginationParams(mapStringToValues(map[string]string{
		"limit":  "10",
		"offset": "abc",
	}))

	require.Zero(t, limit)
	require.Zero(t, offset)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidOffset, *err)
}

func (suite *GroupHandlerTestSuite) TestExtractAndValidatePathEncodeFailure() {
	t := suite.T()
	writer := newFlakyResponseWriter()
	req := httptest.NewRequest(http.MethodGet, "/ous//groups", nil)
	logger := log.GetLogger().With(log.String("component", "test"))

	path, failed := extractAndValidatePath(writer, req, logger)

	require.True(t, failed)
	require.Equal(t, "", path)
	require.Equal(t, http.StatusBadRequest, writer.Code)
	require.Equal(t, "Failed to encode error response\n", writer.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestInternalError() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups/grp-001/members", nil)
	req.SetPathValue("id", "grp-001")
	rr := httptest.NewRecorder()

	serviceMock.
		On("GetGroupMembers", "grp-001", serverconst.DefaultPageSize, 0).
		Return((*MemberListResponse)(nil), &ErrorInternalServerError).
		Once()

	handler.HandleGroupMembersGetRequest(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleGroupMembersGetRequestMissingID() {
	t := suite.T()
	serviceMock := NewGroupServiceInterfaceMock(t)
	handler := newGroupHandler(serviceMock)

	req := httptest.NewRequest(http.MethodGet, "/groups//members", nil)
	rr := httptest.NewRecorder()

	handler.HandleGroupMembersGetRequest(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)

	serviceMock.AssertNotCalled(t, "GetGroupMembers", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupHandlerTestSuite) TestHandleErrorInternalServer() {
	t := suite.T()
	handler := newGroupHandler(nil)
	rr := httptest.NewRecorder()
	logger := log.GetLogger().With(log.String("component", "test"))

	handler.handleError(rr, logger, &serviceerror.ServiceError{
		Type:             serviceerror.ServerErrorType,
		Code:             "GRP-9999",
		Error:            "boom",
		ErrorDescription: "explosion",
	})

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Equal(t, "Internal server error\n", rr.Body.String())
}

func (suite *GroupHandlerTestSuite) TestHandleErrorClientError() {
	t := suite.T()
	handler := newGroupHandler(nil)
	rr := httptest.NewRecorder()
	logger := log.GetLogger().With(log.String("component", "test"))

	handler.handleError(rr, logger, &ErrorGroupNameConflict)

	require.Equal(t, http.StatusConflict, rr.Code)

	var body apierror.ErrorResponse
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &body))
	require.Equal(t, ErrorGroupNameConflict.Code, body.Code)
}
