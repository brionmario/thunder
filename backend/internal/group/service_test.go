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
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	oupkg "github.com/asgardeo/thunder/internal/ou"
	"github.com/asgardeo/thunder/internal/system/error/serviceerror"
	"github.com/asgardeo/thunder/tests/mocks/oumock"
	"github.com/asgardeo/thunder/tests/mocks/usermock"
)

type GroupServiceTestSuite struct {
	suite.Suite
}

func TestGroupServiceTestSuite(t *testing.T) {
	suite.Run(t, new(GroupServiceTestSuite))
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupList_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupListCount").
		Return(3, nil).
		Once()
	storeMock.On("GetGroupList", 2, 1).
		Return([]GroupBasicDAO{
			{ID: "g1", Name: "group-1", Description: "desc-1", OrganizationUnitID: "ou-1"},
			{ID: "g2", Name: "group-2", Description: "desc-2", OrganizationUnitID: "ou-2"},
		}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupList(2, 1)

	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 3, response.TotalResults)
	require.Equal(t, 2, response.Count)
	require.Equal(t, 2, response.StartIndex)
	require.Len(t, response.Groups, 2)
	require.Equal(t, "group-1", response.Groups[0].Name)
	require.Len(t, response.Links, 3)
	require.Equal(t, "/groups?offset=0&limit=2", response.Links[0].Href)
	require.Equal(t, "first", response.Links[0].Rel)
	require.Equal(t, "/groups?offset=0&limit=2", response.Links[1].Href)
	require.Equal(t, "prev", response.Links[1].Rel)
	require.Equal(t, "/groups?offset=2&limit=2", response.Links[2].Href)
	require.Equal(t, "last", response.Links[2].Rel)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupList_InvalidPagination() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupList(0, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidLimit, *err)
	storeMock.AssertNotCalled(t, "GetGroupListCount")
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupList_StoreError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupListCount").
		Return(0, errors.New("count failure")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupList(5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
	storeMock.AssertNotCalled(t, "GetGroupList", mock.Anything, mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupList_ListError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupListCount").
		Return(2, nil).
		Once()
	storeMock.On("GetGroupList", 5, 0).
		Return(nil, errors.New("list failure")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupList(5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}
func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupsByOrganizationUnitCount", "ou-123").
		Return(4, nil).
		Once()
	storeMock.On("GetGroupsByOrganizationUnit", "ou-123", 2, 0).
		Return([]GroupBasicDAO{
			{ID: "g1", Name: "group-1", OrganizationUnitID: "ou-123"},
			{ID: "g2", Name: "group-2", OrganizationUnitID: "ou-123"},
		}, nil).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{ID: "ou-123"}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 2, 0)

	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 4, response.TotalResults)
	require.Equal(t, 2, response.Count)
	require.Len(t, response.Groups, 2)
	require.Equal(t, "group-1", response.Groups[0].Name)
	require.Equal(t, 1, response.StartIndex)
	require.Len(t, response.Links, 2)
	require.Equal(t, "next", response.Links[0].Rel)
	require.Equal(t, "/groups?offset=2&limit=2", response.Links[0].Href)
	require.Equal(t, "last", response.Links[1].Rel)
	require.Equal(t, "/groups?offset=2&limit=2", response.Links[1].Href)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_InvalidPath() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("  ", 10, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
	ouServiceMock.AssertNotCalled(t, "GetOrganizationUnitByPath", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_OUNotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{}, &oupkg.ErrorOrganizationUnitNotFound).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 10, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
	storeMock.AssertNotCalled(t, "GetGroupsByOrganizationUnitCount", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_OUServiceError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	expectedErr := &serviceerror.ServiceError{
		Code: "OU-5000",
		Type: serviceerror.ServerErrorType,
	}
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{}, expectedErr).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 5, 0)

	require.Nil(t, response)
	require.Equal(t, expectedErr, err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_InvalidPagination() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{ID: "ou-1"}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 0, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidLimit, *err)
	storeMock.AssertNotCalled(t, "GetGroupsByOrganizationUnitCount", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_CountError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupsByOrganizationUnitCount", "ou-123").
		Return(0, errors.New("count fail")).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{ID: "ou-123"}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
	storeMock.AssertNotCalled(t, "GetGroupsByOrganizationUnit", mock.Anything, mock.Anything, mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupsByPath_ListError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroupsByOrganizationUnitCount", "ou-123").
		Return(1, nil).
		Once()
	storeMock.On("GetGroupsByOrganizationUnit", "ou-123", 5, 0).
		Return(nil, errors.New("list fail")).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root/child").
		Return(oupkg.OrganizationUnit{ID: "ou-123"}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
		ouService:  ouServiceMock,
	}

	response, err := service.GetGroupsByPath("root/child", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("CheckGroupNameConflictForCreate", "engineering", "ou-001").
		Return(nil).
		Once()

	storeMock.On("ValidateGroupIDs", []string{"grp-002"}).
		Return([]string{}, nil).
		Once()

	storeMock.On("CreateGroup", mock.MatchedBy(func(group GroupDAO) bool {
		return group.Name == "engineering" &&
			group.OrganizationUnitID == "ou-001" &&
			len(group.Members) == 2
	})).
		Return(nil).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{ID: "ou-001"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", []string{"usr-001"}).
		Return([]string{}, nil).
		Once()

	request := CreateGroupRequest{
		Name:               "engineering",
		Description:        "Engineers",
		OrganizationUnitID: "ou-001",
		Members: []Member{
			{ID: "usr-001", Type: MemberTypeUser},
			{ID: "grp-002", Type: MemberTypeGroup},
		},
	}

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroup(request)

	require.Nil(t, err)
	require.NotNil(t, group)
	require.Equal(t, "engineering", group.Name)
	require.Equal(t, "ou-001", group.OrganizationUnitID)
	require.Len(t, group.Members, 2)
	require.NotEmpty(t, group.ID)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_InvalidOU() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-unknown").
		Return(oupkg.OrganizationUnit{}, &oupkg.ErrorOrganizationUnitNotFound).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", []string{}).
		Return([]string{}, nil).
		Maybe()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	request := CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-unknown",
	}

	group, err := service.CreateGroup(request)

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidOUID, *err)
	storeMock.AssertNotCalled(t, "CheckGroupNameConflictForCreate", mock.Anything, mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_InvalidUserIDs() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("CheckGroupNameConflictForCreate", "engineering", "ou-001").
		Return(nil).
		Maybe()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Maybe()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{ID: "ou-001"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", []string{"usr-invalid"}).
		Return([]string{"usr-invalid"}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	request := CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-001",
		Members: []Member{
			{ID: "usr-invalid", Type: MemberTypeUser},
		},
	}

	group, err := service.CreateGroup(request)

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidUserMemberID, *err)
	storeMock.AssertNotCalled(t, "CreateGroup", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_NameConflict() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("CheckGroupNameConflictForCreate", "engineering", "ou-001").
		Return(ErrGroupNameConflict).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{ID: "ou-001"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroup(CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-001",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNameConflict, *err)
	storeMock.AssertNotCalled(t, "CreateGroup", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_CheckConflictError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("CheckGroupNameConflictForCreate", "engineering", "ou-001").
		Return(errors.New("db failure")).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{ID: "ou-001"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroup(CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-001",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
	storeMock.AssertNotCalled(t, "CreateGroup", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_CreateError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("CheckGroupNameConflictForCreate", "engineering", "ou-001").
		Return(nil).
		Once()
	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()
	storeMock.On("CreateGroup", mock.Anything).
		Return(errors.New("create fail")).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{ID: "ou-001"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroup(CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-001",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroup_InternalOUError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-001").
		Return(oupkg.OrganizationUnit{}, &serviceerror.ServiceError{
			Code: "OU-5000",
			Type: serviceerror.ServerErrorType,
		}).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroup(CreateGroupRequest{
		Name:               "engineering",
		OrganizationUnitID: "ou-001",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroupByPath_InvalidPath() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	userServiceMock := usermock.NewUserServiceInterfaceMock(t)

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroupByPath(" ", CreateGroupByPathRequest{Name: "n"})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
	ouServiceMock.AssertNotCalled(t, "GetOrganizationUnitByPath", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroupByPath_OUServiceError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	expectedErr := &serviceerror.ServiceError{
		Code: "OU-5000",
		Type: serviceerror.ServerErrorType,
	}
	ouServiceMock.On("GetOrganizationUnitByPath", "root").
		Return(oupkg.OrganizationUnit{}, expectedErr).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroupByPath("root", CreateGroupByPathRequest{Name: "n"})

	require.Nil(t, group)
	require.Equal(t, expectedErr, err)
}

func (suite *GroupServiceTestSuite) TestGroupService_CreateGroupByPath_OUNotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnitByPath", "root").
		Return(oupkg.OrganizationUnit{}, &oupkg.ErrorOrganizationUnitNotFound).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.CreateGroupByPath("root", CreateGroupByPathRequest{Name: "n"})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroup_MissingID() {
	t := suite.T()
	service := &groupService{}

	group, err := service.GetGroup("")

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorMissingGroupID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroup_InternalError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, errors.New("db error")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	group, err := service.GetGroup("grp-001")

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroup_NotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-404").
		Return(GroupDAO{}, ErrGroupNotFound).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	group, err := service.GetGroup("grp-404")

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_MissingID() {
	t := suite.T()
	service := &groupService{}

	group, err := service.UpdateGroup("", UpdateGroupRequest{})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorMissingGroupID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_InvalidRequest() {
	t := suite.T()
	service := &groupService{}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "old-name",
			Description:        "legacy",
			OrganizationUnitID: "ou-old",
		}, nil).
		Once()

	storeMock.On("CheckGroupNameConflictForUpdate", "new-name", "ou-new", "grp-001").
		Return(nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	storeMock.On("UpdateGroup", mock.MatchedBy(func(group GroupDAO) bool {
		return group.ID == "grp-001" &&
			group.Name == "new-name" &&
			group.OrganizationUnitID == "ou-new"
	})).
		Return(nil).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-new").
		Return(oupkg.OrganizationUnit{ID: "ou-new"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "new-name",
		OrganizationUnitID: "ou-new",
	})

	require.Nil(t, err)
	require.NotNil(t, group)
	require.Equal(t, "new-name", group.Name)
	require.Equal(t, "ou-new", group.OrganizationUnitID)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_NameConflict() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "old-name",
			OrganizationUnitID: "ou-old",
		}, nil).
		Once()

	storeMock.On("CheckGroupNameConflictForUpdate", "new-name", "ou-new", "grp-001").
		Return(ErrGroupNameConflict).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-new").
		Return(oupkg.OrganizationUnit{ID: "ou-new"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "new-name",
		OrganizationUnitID: "ou-new",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNameConflict, *err)
	storeMock.AssertNotCalled(t, "UpdateGroup", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_GroupNotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, ErrGroupNotFound).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_GetGroupError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, errors.New("db error")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_ValidateOUError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "name",
			OrganizationUnitID: "ou-old",
		}, nil).
		Once()

	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-new").
		Return(oupkg.OrganizationUnit{}, &oupkg.ErrorOrganizationUnitNotFound).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Maybe()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Maybe()

	service := &groupService{
		groupStore:  storeMock,
		ouService:   ouServiceMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou-new",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidOUID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_InvalidUserIDs() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "name",
			OrganizationUnitID: "ou",
		}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", []string{"usr-1"}).
		Return([]string{"usr-1"}, nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Maybe()

	service := &groupService{
		groupStore:  storeMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "usr-1", Type: MemberTypeUser},
		},
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidUserMemberID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_ValidateGroupIDsError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "name",
			OrganizationUnitID: "ou",
		}, nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return(nil, errors.New("validate fail")).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "grp-2", Type: MemberTypeGroup},
		},
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_InvalidGroupIDs() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "name",
			OrganizationUnitID: "ou",
		}, nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{"grp-2"}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "grp-2", Type: MemberTypeGroup},
		},
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidGroupMemberID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_CheckConflictError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "old-name",
			OrganizationUnitID: "ou",
		}, nil).
		Once()

	storeMock.On("CheckGroupNameConflictForUpdate", "new", "ou", "grp-001").
		Return(errors.New("db error")).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "new",
		OrganizationUnitID: "ou",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_UpdateGroup_UpdateError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{
			ID:                 "grp-001",
			Name:               "old-name",
			OrganizationUnitID: "ou",
		}, nil).
		Once()

	storeMock.On("CheckGroupNameConflictForUpdate", "new", "ou", "grp-001").
		Return(nil).
		Once()

	storeMock.On("ValidateGroupIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	storeMock.On("UpdateGroup", mock.Anything).
		Return(errors.New("update fail")).
		Once()

	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", mock.Anything).
		Return([]string{}, nil).
		Once()

	service := &groupService{
		groupStore:  storeMock,
		userService: userServiceMock,
	}

	group, err := service.UpdateGroup("grp-001", UpdateGroupRequest{
		Name:               "new",
		OrganizationUnitID: "ou",
	})

	require.Nil(t, group)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_DeleteGroup_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{ID: "grp-001"}, nil).
		Once()
	storeMock.On("DeleteGroup", "grp-001").
		Return(nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.DeleteGroup("grp-001")

	require.Nil(t, err)
}

func (suite *GroupServiceTestSuite) TestGroupService_DeleteGroup_MissingID() {
	t := suite.T()
	service := &groupService{}

	err := service.DeleteGroup("")

	require.NotNil(t, err)
	require.Equal(t, ErrorMissingGroupID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_DeleteGroup_GetGroupError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, errors.New("db error")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.DeleteGroup("grp-001")

	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
	storeMock.AssertNotCalled(t, "DeleteGroup", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_DeleteGroup_DeleteError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{ID: "grp-001"}, nil).
		Once()
	storeMock.On("DeleteGroup", "grp-001").
		Return(errors.New("delete fail")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.DeleteGroup("grp-001")

	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_DeleteGroup_NotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, ErrGroupNotFound).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.DeleteGroup("grp-001")

	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
	storeMock.AssertNotCalled(t, "DeleteGroup", "grp-001")
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_Success() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{ID: "grp-001"}, nil).
		Once()
	storeMock.On("GetGroupMemberCount", "grp-001").
		Return(3, nil).
		Once()
	storeMock.On("GetGroupMembers", "grp-001", 2, 0).
		Return([]Member{
			{ID: "usr-001", Type: MemberTypeUser},
			{ID: "grp-002", Type: MemberTypeGroup},
		}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupMembers("grp-001", 2, 0)

	require.Nil(t, err)
	require.NotNil(t, response)
	require.Equal(t, 3, response.TotalResults)
	require.Equal(t, 2, response.Count)
	require.Equal(t, 1, response.StartIndex)
	require.Len(t, response.Members, 2)
	require.Len(t, response.Links, 2)
	require.Equal(t, "next", response.Links[0].Rel)
	require.Equal(t, "/groups/grp-001/members?offset=2&limit=2", response.Links[0].Href)
	require.Equal(t, "last", response.Links[1].Rel)
	require.Equal(t, "/groups/grp-001/members?offset=2&limit=2", response.Links[1].Href)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_GroupNotFound() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, ErrGroupNotFound).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupMembers("grp-001", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorGroupNotFound, *err)
	storeMock.AssertNotCalled(t, "GetGroupMemberCount", mock.Anything)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_InvalidPagination() {
	t := suite.T()
	service := &groupService{}

	response, err := service.GetGroupMembers("grp-001", 0, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidLimit, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_MissingID() {
	t := suite.T()
	service := &groupService{}

	response, err := service.GetGroupMembers("", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorMissingGroupID, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_GetGroupError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{}, errors.New("db error")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupMembers("grp-001", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_CountError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{ID: "grp-001"}, nil).
		Once()
	storeMock.On("GetGroupMemberCount", "grp-001").
		Return(0, errors.New("count fail")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupMembers("grp-001", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestGroupService_GetGroupMembers_ListError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("GetGroup", "grp-001").
		Return(GroupDAO{ID: "grp-001"}, nil).
		Once()
	storeMock.On("GetGroupMemberCount", "grp-001").
		Return(1, nil).
		Once()
	storeMock.On("GetGroupMembers", "grp-001", 5, 0).
		Return(nil, errors.New("list fail")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	response, err := service.GetGroupMembers("grp-001", 5, 0)

	require.Nil(t, response)
	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestValidatePaginationParams() {
	t := suite.T()
	err := validatePaginationParams(0, 0)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidLimit, *err)

	err = validatePaginationParams(5, -1)
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidOffset, *err)

	err = validatePaginationParams(5, 0)
	require.Nil(t, err)
}

func (suite *GroupServiceTestSuite) TestBuildPaginationLinks() {
	t := suite.T()
	links := buildPaginationLinks("/groups", 5, 5, 18)

	require.Len(t, links, 4)
	require.Equal(t, "first", links[0].Rel)
	require.Equal(t, "/groups?offset=0&limit=5", links[0].Href)
	require.Equal(t, "prev", links[1].Rel)
	require.Equal(t, "/groups?offset=0&limit=5", links[1].Href)
	require.Equal(t, "next", links[2].Rel)
	require.Equal(t, "/groups?offset=10&limit=5", links[2].Href)
	require.Equal(t, "last", links[3].Rel)
	require.Equal(t, "/groups?offset=15&limit=5", links[3].Href)
}

func (suite *GroupServiceTestSuite) TestValidateCreateGroupRequestMissingFields() {
	t := suite.T()
	service := &groupService{}

	err := service.validateCreateGroupRequest(CreateGroupRequest{})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)

	err = service.validateCreateGroupRequest(CreateGroupRequest{
		Name: "name",
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
}

func (suite *GroupServiceTestSuite) TestValidateCreateGroupRequestInvalidMembers() {
	t := suite.T()
	service := &groupService{}

	err := service.validateCreateGroupRequest(CreateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "id", Type: "invalid"},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)

	err = service.validateCreateGroupRequest(CreateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "", Type: MemberTypeUser},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
}

func (suite *GroupServiceTestSuite) TestValidateUpdateGroupRequestMissingFields() {
	t := suite.T()
	service := &groupService{}

	err := service.validateUpdateGroupRequest(UpdateGroupRequest{})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)

	err = service.validateUpdateGroupRequest(UpdateGroupRequest{
		Name: "name",
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
}

func (suite *GroupServiceTestSuite) TestValidateUpdateGroupRequestInvalidMembers() {
	t := suite.T()
	service := &groupService{}

	err := service.validateUpdateGroupRequest(UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "id", Type: "invalid"},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)

	err = service.validateUpdateGroupRequest(UpdateGroupRequest{
		Name:               "name",
		OrganizationUnitID: "ou",
		Members: []Member{
			{ID: "", Type: MemberTypeGroup},
		},
	})
	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidRequestFormat, *err)
}

func (suite *GroupServiceTestSuite) TestValidateOUHandlesInternalError() {
	t := suite.T()
	ouServiceMock := oumock.NewOrganizationUnitServiceInterfaceMock(t)
	ouServiceMock.On("GetOrganizationUnit", "ou-1").
		Return(oupkg.OrganizationUnit{}, &serviceerror.ServiceError{
			Code: "OU-5000",
			Type: serviceerror.ServerErrorType,
		}).
		Once()

	service := &groupService{
		ouService: ouServiceMock,
	}

	err := service.validateOU("ou-1")

	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestValidateAndProcessHandlePath() {
	t := suite.T()
	service := &groupService{}

	testCases := []struct {
		name        string
		handlePath  string
		expectError bool
	}{
		{
			name:        "empty string",
			handlePath:  "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			handlePath:  "   ",
			expectError: true,
		},
		{
			name:        "only slashes",
			handlePath:  "///",
			expectError: true,
		},
		{
			name:        "double slash between handles",
			handlePath:  "root//child",
			expectError: true,
		},
		{
			name:        "single slash",
			handlePath:  "/",
			expectError: true,
		},
		{
			name:        "valid handles",
			handlePath:  "root/child",
			expectError: false,
		},
		{
			name:        "valid handles with surrounding whitespace and slashes",
			handlePath:  "  /root/child/  ",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := service.validateAndProcessHandlePath(tc.handlePath)
			if tc.expectError {
				require.NotNil(t, err)
				require.Equal(t, ErrorInvalidRequestFormat, *err)
				return
			}

			require.Nil(t, err)
		})
	}
}

func (suite *GroupServiceTestSuite) TestValidateUserIDsHandlesServiceError() {
	t := suite.T()
	userServiceMock := usermock.NewUserServiceInterfaceMock(t)
	userServiceMock.On("ValidateUserIDs", []string{"usr-001"}).
		Return([]string{}, &serviceerror.ServiceError{
			Code: "USR-5000",
			Type: serviceerror.ServerErrorType,
		}).
		Once()

	service := &groupService{
		userService: userServiceMock,
	}

	err := service.validateUserIDs([]string{"usr-001"})

	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}

func (suite *GroupServiceTestSuite) TestValidateGroupIDsHandlesInvalidIDs() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("ValidateGroupIDs", []string{"grp-001"}).
		Return([]string{"grp-001"}, nil).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.ValidateGroupIDs([]string{"grp-001"})

	require.NotNil(t, err)
	require.Equal(t, ErrorInvalidGroupMemberID, *err)
}

func (suite *GroupServiceTestSuite) TestValidateGroupIDsHandlesStoreError() {
	t := suite.T()
	storeMock := newGroupStoreInterfaceMock(t)
	storeMock.On("ValidateGroupIDs", []string{"grp-001"}).
		Return(nil, errors.New("db error")).
		Once()

	service := &groupService{
		groupStore: storeMock,
	}

	err := service.ValidateGroupIDs([]string{"grp-001"})

	require.NotNil(t, err)
	require.Equal(t, ErrorInternalServerError, *err)
}
