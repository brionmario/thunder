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

	dbmodel "github.com/asgardeo/thunder/internal/system/database/model"
	"github.com/asgardeo/thunder/tests/mocks/database/clientmock"
	"github.com/asgardeo/thunder/tests/mocks/database/modelmock"
	"github.com/asgardeo/thunder/tests/mocks/database/providermock"
)

type GroupStoreTestSuite struct {
	suite.Suite
}

func TestGroupStoreTestSuite(t *testing.T) {
	suite.Run(t, new(GroupStoreTestSuite))
}

const queryGroupExistsID = "GRQ-GROUP_MGT-15"

type validateGroupIDsSetupFn func(
	*providermock.DBProviderInterfaceMock,
	*clientmock.DBClientInterfaceMock,
)

type validateGroupIDsOverrideFn func(*bool) func()

type validateGroupIDsPostAssertFn func(
	*testing.T,
	*providermock.DBProviderInterfaceMock,
	*clientmock.DBClientInterfaceMock,
	bool,
)

func assertBuilderErrorPostconditions(
	t *testing.T,
	providerMock *providermock.DBProviderInterfaceMock,
	dbClientMock *clientmock.DBClientInterfaceMock,
	builderCalled bool,
) {
	require.True(t, builderCalled)
	dbClientMock.AssertNotCalled(t, "Query", mock.Anything, mock.Anything)
}

func assertEmptyInputPostconditions(
	t *testing.T,
	providerMock *providermock.DBProviderInterfaceMock,
	dbClientMock *clientmock.DBClientInterfaceMock,
	builderCalled bool,
) {
	require.False(t, builderCalled)
	providerMock.AssertNotCalled(t, "GetDBClient", mock.Anything)
	dbClientMock.AssertNotCalled(t, "Query", mock.Anything, mock.Anything)
}

type stubSQLResult struct {
	rows int64
}

func (s stubSQLResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (s stubSQLResult) RowsAffected() (int64, error) {
	return s.rows, nil
}

type errSQLResult struct {
	err error
}

func (e errSQLResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (e errSQLResult) RowsAffected() (int64, error) {
	return 0, e.err
}

// testExecRollbackError is a helper function to test rollback errors during database operations.
func testExecRollbackError(t *testing.T, query string, operation func(*groupStore, GroupDAO) error) {
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(nil, errors.New("exec failed")).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := operation(store, group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListCountSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupListCount).
		Return([]map[string]interface{}{{"total": int64(7)}}, nil).
		Once()

	count, err := store.GetGroupListCount()

	require.NoError(t, err)
	require.Equal(t, 7, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListCountClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("no client")).
		Once()

	count, err := store.GetGroupListCount()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
	require.Equal(t, 0, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListCountQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupListCount).
		Return(nil, errors.New("boom")).
		Once()

	count, err := store.GetGroupListCount()

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute count query")
	require.Equal(t, 0, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	rows := []map[string]interface{}{
		{
			"group_id":    "g1",
			"name":        "Group 1",
			"description": "Desc 1",
			"ou_id":       "ou-1",
		},
		{
			"group_id":    "g2",
			"name":        "Group 2",
			"description": "Desc 2",
			"ou_id":       "ou-2",
		},
	}

	dbClientMock.
		On("Query", QueryGetGroupList, 5, 0).
		Return(rows, nil).
		Once()

	groups, err := store.GetGroupList(5, 0)

	require.NoError(t, err)
	require.Len(t, groups, 2)
	require.Equal(t, "g1", groups[0].ID)
	require.Equal(t, "Group 2", groups[1].Name)
	require.Equal(t, "ou-2", groups[1].OrganizationUnitID)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListProviderError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("boom")).
		Once()

	groups, err := store.GetGroupList(1, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
	require.Nil(t, groups)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupList, 1, 0).
		Return(nil, errors.New("query fail")).
		Once()

	groups, err := store.GetGroupList(1, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute group list query")
	require.Nil(t, groups)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupListInvalidRow() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupList, 1, 0).
		Return([]map[string]interface{}{
			{
				"group_id": "g1",
				"name":     "Group 1",
				// description missing to trigger error
				"ou_id": "ou-1",
			},
		}, nil).
		Once()

	groups, err := store.GetGroupList(1, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to build group from result row")
	require.Nil(t, groups)
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		Description:        "Core team",
		OrganizationUnitID: "ou-1",
		Members: []Member{
			{ID: "user-1", Type: MemberTypeUser},
		},
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryCreateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "user-1").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Commit").
		Return(nil).
		Once()

	err := store.CreateGroup(group)

	require.NoError(t, err)
	txMock.AssertNotCalled(t, "Rollback")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupInsertError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		Description:        "Core team",
		OrganizationUnitID: "ou-1",
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryCreateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(nil, errors.New("insert failed")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.CreateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute query")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client error")).
		Once()

	err := store.CreateGroup(GroupDAO{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupInsertRollbackError() {
	t := suite.T()
	testExecRollbackError(t, QueryCreateGroup.Query, func(store *groupStore, group GroupDAO) error {
		return store.CreateGroup(group)
	})
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupAddMemberRollbackError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001", Members: []Member{{ID: "usr-1", Type: MemberTypeUser}}}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryCreateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "usr-1").
		Return(nil, errors.New("member fail")).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.CreateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupCommitError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001", Members: []Member{{ID: "usr-1", Type: MemberTypeUser}}}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryCreateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "usr-1").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Commit").
		Return(errors.New("commit fail")).
		Once()

	err := store.CreateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to commit transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupBeginTxError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(nil, errors.New("begin fail")).
		Once()

	err := store.CreateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to begin transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CreateGroupAddMemberError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		OrganizationUnitID: "ou-1",
		Members: []Member{
			{ID: "user-1", Type: MemberTypeUser},
		},
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryCreateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "user-1").
		Return(nil, errors.New("member fail")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.CreateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to add member to group")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupByID, "grp-001").
		Return([]map[string]interface{}{
			{
				"group_id":    "grp-001",
				"name":        "Engineering",
				"description": "Core team",
				"ou_id":       "ou-1",
			},
		}, nil).
		Once()

	group, err := store.GetGroup("grp-001")

	require.NoError(t, err)
	require.Equal(t, "Engineering", group.Name)
	require.Equal(t, "ou-1", group.OrganizationUnitID)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	_, err := store.GetGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupBuildError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupByID, "grp-001").
		Return([]map[string]interface{}{{"name": "group"}}, nil).
		Once()

	_, err := store.GetGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to parse group_id")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupByID, "grp-001").
		Return(nil, errors.New("query fail")).
		Once()

	_, err := store.GetGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute query")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupUnexpectedResults() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupByID, "grp-001").
		Return([]map[string]interface{}{
			{"group_id": "grp-001"},
			{"group_id": "grp-002"},
		}, nil).
		Once()

	_, err := store.GetGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected number of results")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupNotFound() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupByID, "grp-404").
		Return([]map[string]interface{}{}, nil).
		Once()

	group, err := store.GetGroup("grp-404")

	require.ErrorIs(t, err, ErrGroupNotFound)
	require.Empty(t, group.ID)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMembersSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMembers, "grp-001", 2, 0).
		Return([]map[string]interface{}{
			{"member_id": "usr-1", "member_type": "user"},
			{"member_id": "grp-2", "member_type": "group"},
		}, nil).
		Once()

	members, err := store.GetGroupMembers("grp-001", 2, 0)

	require.NoError(t, err)
	require.Len(t, members, 2)
	require.Equal(t, MemberTypeUser, members[0].Type)
	require.Equal(t, "grp-2", members[1].ID)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMembersQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMembers, "grp-001", 2, 0).
		Return(nil, errors.New("query failed")).
		Once()

	members, err := store.GetGroupMembers("grp-001", 2, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get group members")
	require.Nil(t, members)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMembersDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	_, err := store.GetGroupMembers("grp-001", 1, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMemberCountSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMemberCount, "grp-001").
		Return([]map[string]interface{}{
			{"total": int64(3)},
		}, nil).
		Once()

	count, err := store.GetGroupMemberCount("grp-001")

	require.NoError(t, err)
	require.Equal(t, 3, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMemberCountQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMemberCount, "grp-001").
		Return(nil, errors.New("query fail")).
		Once()

	_, err := store.GetGroupMemberCount("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get group member count")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMemberCountDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	_, err := store.GetGroupMemberCount("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMemberCountEmptyResult() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMemberCount, "grp-001").
		Return([]map[string]interface{}{}, nil).
		Once()

	count, err := store.GetGroupMemberCount("grp-001")

	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupMemberCountInvalidFormat() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupMemberCount, "grp-001").
		Return([]map[string]interface{}{{"total": "invalid"}}, nil).
		Once()

	count, err := store.GetGroupMemberCount("grp-001")

	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupRowsAffectedZero() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		Description:        "Core",
		OrganizationUnitID: "ou-1",
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 0}, nil).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.UpdateGroup(group)

	require.ErrorIs(t, err, ErrGroupNotFound)
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	err := store.UpdateGroup(GroupDAO{})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupExecRollbackError() {
	t := suite.T()
	testExecRollbackError(t, QueryUpdateGroup.Query, func(store *groupStore, group GroupDAO) error {
		return store.UpdateGroup(group)
	})
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupRowsAffectedError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(errSQLResult{err: errors.New("rows fail")}, nil).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupRowsAffectedRollbackError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 0}, nil).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupUpdateMembersRollbackError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001", Members: []Member{{ID: "usr-1", Type: MemberTypeUser}}}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, group.ID).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "usr-1").
		Return(nil, errors.New("member fail")).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupCommitError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, group.ID).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, mock.Anything, mock.Anything).
		Return(stubSQLResult{rows: 1}, nil).
		Maybe()

	txMock.
		On("Commit").
		Return(errors.New("commit fail")).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to commit transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupBeginTxError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(nil, errors.New("begin fail")).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to begin transaction")
}
func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupExecError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{ID: "grp-001"}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(nil, errors.New("exec fail")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute query")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupDeleteMembersError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		OrganizationUnitID: "ou-1",
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, group.ID).
		Return(nil, errors.New("delete fail")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.UpdateGroup(group)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete existing group member assignments")
}

func (suite *GroupStoreTestSuite) TestGroupStore_UpdateGroupSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}
	group := GroupDAO{
		ID:                 "grp-001",
		Name:               "Engineering",
		Description:        "Core",
		OrganizationUnitID: "ou-1",
		Members: []Member{
			{ID: "user-1", Type: MemberTypeUser},
		},
	}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryUpdateGroup.Query, group.ID, group.OrganizationUnitID, group.Name, group.Description).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, group.ID).
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, group.ID, MemberTypeUser, "user-1").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Commit").
		Return(nil).
		Once()

	err := store.UpdateGroup(group)

	require.NoError(t, err)
	txMock.AssertNotCalled(t, "Rollback")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupBeginTxError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(nil, errors.New("begin fail")).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to begin transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupMembersError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(nil, errors.New("delete fail")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete group members")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupExecError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(nil, errors.New("delete fail")).
		Once()

	txMock.
		On("Rollback").
		Return(nil).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to execute query")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupCommitError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Commit").
		Return(errors.New("commit fail")).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to commit transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupRowsAffectedError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(errSQLResult{err: errors.New("rows fail")}, nil).
		Once()

	txMock.
		On("Commit").
		Return(nil).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get rows affected")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Commit").
		Return(nil).
		Once()

	err := store.DeleteGroup("grp-001")

	require.NoError(t, err)
	txMock.AssertNotCalled(t, "Rollback")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupMembersRollbackError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(nil, errors.New("delete fail")).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupRollbackAfterDeleteError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(nil, errors.New("delete fail")).
		Once()

	txMock.
		On("Rollback").
		Return(errors.New("rollback fail")).
		Once()

	err := store.DeleteGroup("grp-001")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to rollback transaction")
}

func (suite *GroupStoreTestSuite) TestGroupStore_DeleteGroupRowsAffectedZero() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)
	txMock := modelmock.NewTxInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("BeginTx").
		Return(txMock, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryDeleteGroup.Query, "grp-001").
		Return(stubSQLResult{rows: 0}, nil).
		Once()

	txMock.
		On("Commit").
		Return(nil).
		Once()

	err := store.DeleteGroup("grp-001")

	require.NoError(t, err)
}
func (suite *GroupStoreTestSuite) TestGroupStore_ValidateGroupIDs() {
	t := suite.T()
	type testCase struct {
		name            string
		groupIDs        []string
		setup           validateGroupIDsSetupFn
		overrideBuilder validateGroupIDsOverrideFn
		wantInvalid     []string
		wantErr         string
		postAssert      validateGroupIDsPostAssertFn
	}

	queryMatcher := func() interface{} {
		return mock.MatchedBy(func(q dbmodel.DBQuery) bool { return q.ID == queryGroupExistsID })
	}

	testCases := []testCase{
		{
			name:     "returns missing IDs",
			groupIDs: []string{"grp-1", "grp-2"},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, dbClientMock *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(dbClientMock, nil).
					Once()

				dbClientMock.
					On("Query", queryMatcher(), "grp-1", "grp-2").
					Return([]map[string]interface{}{{"group_id": "grp-1"}}, nil).
					Once()
			},
			wantInvalid: []string{"grp-2"},
		},
		{
			name:     "preserves invalid order including empty IDs",
			groupIDs: []string{"grp-miss", "", "grp-hit"},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, dbClientMock *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(dbClientMock, nil).
					Once()

				dbClientMock.
					On("Query", queryMatcher(), "grp-miss", "", "grp-hit").
					Return([]map[string]interface{}{{"group_id": "grp-hit"}}, nil).
					Once()
			},
			wantInvalid: []string{"grp-miss", ""},
		},
		{
			name:     "query error",
			groupIDs: []string{"grp-1"},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, dbClientMock *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(dbClientMock, nil).
					Once()

				dbClientMock.
					On("Query", queryMatcher(), "grp-1").
					Return(nil, errors.New("query fail")).
					Once()
			},
			wantErr: "failed to execute query",
		},
		{
			name:     "builder error",
			groupIDs: []string{"grp-1"},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, dbClientMock *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(dbClientMock, nil).
					Once()
			},
			overrideBuilder: func(builderCalled *bool) func() {
				originalBuilder := buildBulkGroupExistsQueryFunc
				buildBulkGroupExistsQueryFunc = func(groupIDs []string) (dbmodel.DBQuery, []interface{}, error) {
					if builderCalled != nil {
						*builderCalled = true
					}
					return dbmodel.DBQuery{}, nil, errors.New("builder fail")
				}
				return func() { buildBulkGroupExistsQueryFunc = originalBuilder }
			},
			wantErr:    "failed to build bulk group exists query",
			postAssert: assertBuilderErrorPostconditions,
		},
		{
			name:     "db client error",
			groupIDs: []string{"grp-1"},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, _ *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(nil, errors.New("client fail")).
					Once()
			},
			wantErr: "failed to get database client",
		},
		{
			name:        "empty input returns immediately",
			groupIDs:    []string{},
			wantInvalid: []string{},
			overrideBuilder: func(builderCalled *bool) func() {
				originalBuilder := buildBulkGroupExistsQueryFunc
				buildBulkGroupExistsQueryFunc = func(groupIDs []string) (dbmodel.DBQuery, []interface{}, error) {
					if builderCalled != nil {
						*builderCalled = true
					}
					return originalBuilder(groupIDs)
				}
				return func() { buildBulkGroupExistsQueryFunc = originalBuilder }
			},
			postAssert: assertEmptyInputPostconditions,
		},
		{
			name:     "all empty values treated as invalid",
			groupIDs: []string{"", ""},
			setup: func(providerMock *providermock.DBProviderInterfaceMock, dbClientMock *clientmock.DBClientInterfaceMock) {
				providerMock.
					On("GetDBClient", "identity").
					Return(dbClientMock, nil).
					Once()

				dbClientMock.
					On("Query", queryMatcher(), "", "").
					Return([]map[string]interface{}{}, nil).
					Once()
			},
			wantInvalid: []string{"", ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			providerMock := providermock.NewDBProviderInterfaceMock(t)
			dbClientMock := clientmock.NewDBClientInterfaceMock(t)

			var builderCalled bool
			if tc.overrideBuilder != nil {
				restore := tc.overrideBuilder(&builderCalled)
				t.Cleanup(restore)
			}

			if tc.setup != nil {
				tc.setup(providerMock, dbClientMock)
			}

			store := &groupStore{dbProvider: providerMock}
			invalid, err := store.ValidateGroupIDs(tc.groupIDs)

			if tc.wantErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.wantErr)
				require.Nil(t, invalid)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantInvalid, invalid)
			}

			if tc.postAssert != nil {
				tc.postAssert(t, providerMock, dbClientMock, builderCalled)
			}

			providerMock.AssertExpectations(t)
			dbClientMock.AssertExpectations(t)
		})
	}
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitCountDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	_, err := store.GetGroupsByOrganizationUnitCount("ou-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitCountQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupsByOrganizationUnitCount, "ou-1").
		Return(nil, errors.New("query fail")).
		Once()

	_, err := store.GetGroupsByOrganizationUnitCount("ou-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get group count by organization unit")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitCountEmpty() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupsByOrganizationUnitCount, "ou-1").
		Return([]map[string]interface{}{}, nil).
		Once()

	count, err := store.GetGroupsByOrganizationUnitCount("ou-1")

	require.NoError(t, err)
	require.Equal(t, 0, count)
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitCountUnexpectedFormat() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupsByOrganizationUnitCount, "ou-1").
		Return([]map[string]interface{}{{"total": "not-number"}}, nil).
		Once()

	_, err := store.GetGroupsByOrganizationUnitCount("ou-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "unexpected response format")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	_, err := store.GetGroupsByOrganizationUnit("ou-1", 10, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitQueryError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupsByOrganizationUnit, "ou-1", 10, 0).
		Return(nil, errors.New("query fail")).
		Once()

	_, err := store.GetGroupsByOrganizationUnit("ou-1", 10, 0)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get groups by organization unit")
}

func (suite *GroupStoreTestSuite) TestGroupStore_GetGroupsByOrganizationUnitSuccess() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(dbClientMock, nil).
		Once()

	dbClientMock.
		On("Query", QueryGetGroupsByOrganizationUnit, "ou-1", 10, 0).
		Return([]map[string]interface{}{
			{"group_id": "grp-1", "ou_id": "ou-1", "name": "g1", "description": "desc"},
		}, nil).
		Once()

	groups, err := store.GetGroupsByOrganizationUnit("ou-1", 10, 0)

	require.NoError(t, err)
	require.Len(t, groups, 1)
	require.Equal(t, "g1", groups[0].Name)
	require.Equal(t, "desc", groups[0].Description)
}

func (suite *GroupStoreTestSuite) TestCheckGroupNameConflictForCreateDetectsConflict() {
	t := suite.T()
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	dbClientMock.
		On("Query", QueryCheckGroupNameConflict, "engineering", "ou-1").
		Return([]map[string]interface{}{
			{"count": int64(1)},
		}, nil).
		Once()

	err := checkGroupNameConflictForCreate(dbClientMock, "engineering", "ou-1")

	require.ErrorIs(t, err, ErrGroupNameConflict)
}

func (suite *GroupStoreTestSuite) TestCheckGroupNameConflictForCreateQueryError() {
	t := suite.T()
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	dbClientMock.
		On("Query", QueryCheckGroupNameConflict, "engineering", "ou-1").
		Return(nil, errors.New("query fail")).
		Once()

	err := checkGroupNameConflictForCreate(dbClientMock, "engineering", "ou-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to check group name conflict")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CheckGroupNameConflictForCreateDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	err := store.CheckGroupNameConflictForCreate("engineering", "ou-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestCheckGroupNameConflictForUpdateSuccess() {
	t := suite.T()
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	dbClientMock.
		On("Query", QueryCheckGroupNameConflictForUpdate, "engineering", "ou-1", "grp-1").
		Return([]map[string]interface{}{
			{"count": int64(0)},
		}, nil).
		Once()

	err := checkGroupNameConflictForUpdate(dbClientMock, "engineering", "ou-1", "grp-1")

	require.NoError(t, err)
}

func (suite *GroupStoreTestSuite) TestCheckGroupNameConflictForUpdateDetectsConflict() {
	t := suite.T()
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	dbClientMock.
		On("Query", QueryCheckGroupNameConflictForUpdate, "engineering", "ou-1", "grp-1").
		Return([]map[string]interface{}{
			{"count": int64(1)},
		}, nil).
		Once()

	err := checkGroupNameConflictForUpdate(dbClientMock, "engineering", "ou-1", "grp-1")

	require.ErrorIs(t, err, ErrGroupNameConflict)
}

func (suite *GroupStoreTestSuite) TestCheckGroupNameConflictForUpdateQueryError() {
	t := suite.T()
	dbClientMock := clientmock.NewDBClientInterfaceMock(t)

	dbClientMock.
		On("Query", QueryCheckGroupNameConflictForUpdate, "engineering", "ou-1", "grp-1").
		Return(nil, errors.New("query fail")).
		Once()

	err := checkGroupNameConflictForUpdate(dbClientMock, "engineering", "ou-1", "grp-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to check group name conflict")
}

func (suite *GroupStoreTestSuite) TestGroupStore_CheckGroupNameConflictForUpdateDBClientError() {
	t := suite.T()
	providerMock := providermock.NewDBProviderInterfaceMock(t)
	store := &groupStore{dbProvider: providerMock}

	providerMock.
		On("GetDBClient", "identity").
		Return(nil, errors.New("client fail")).
		Once()

	err := store.CheckGroupNameConflictForUpdate("engineering", "ou-1", "grp-1")

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get database client")
}

func (suite *GroupStoreTestSuite) TestBuildGroupFromResultRowInvalidGroupID() {
	t := suite.T()
	_, err := buildGroupFromResultRow(map[string]interface{}{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "group_id")
}

func (suite *GroupStoreTestSuite) TestBuildGroupFromResultRowInvalidName() {
	t := suite.T()
	_, err := buildGroupFromResultRow(map[string]interface{}{
		"group_id": "grp-1",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "name")
}

func (suite *GroupStoreTestSuite) TestBuildGroupFromResultRowInvalidDescription() {
	t := suite.T()
	_, err := buildGroupFromResultRow(map[string]interface{}{
		"group_id": "grp-1",
		"name":     "group",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "description")
}

func (suite *GroupStoreTestSuite) TestBuildGroupFromResultRowInvalidOUID() {
	t := suite.T()
	_, err := buildGroupFromResultRow(map[string]interface{}{
		"group_id":    "grp-1",
		"name":        "group",
		"description": "desc",
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "ou_id")
}

func (suite *GroupStoreTestSuite) TestBuildBulkGroupExistsQueryEmpty() {
	t := suite.T()
	_, _, err := buildBulkGroupExistsQuery([]string{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "groupIDs list cannot be empty")
}

func (suite *GroupStoreTestSuite) TestAddMembersToGroupReturnsError() {
	t := suite.T()
	txMock := modelmock.NewTxInterfaceMock(t)

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, "grp-001", MemberTypeUser, "usr-1").
		Return(nil, errors.New("insert fail")).
		Once()

	err := addMembersToGroup(txMock, "grp-001", []Member{{ID: "usr-1", Type: MemberTypeUser}})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to add member to group")
}

func (suite *GroupStoreTestSuite) TestUpdateGroupMembersDeleteError() {
	t := suite.T()
	txMock := modelmock.NewTxInterfaceMock(t)

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(nil, errors.New("delete fail")).
		Once()

	err := updateGroupMembers(txMock, "grp-001", nil)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to delete existing group member assignments")
}

func (suite *GroupStoreTestSuite) TestUpdateGroupMembersAddError() {
	t := suite.T()
	txMock := modelmock.NewTxInterfaceMock(t)

	txMock.
		On("Exec", QueryDeleteGroupMembers.Query, "grp-001").
		Return(stubSQLResult{rows: 1}, nil).
		Once()

	txMock.
		On("Exec", QueryAddMemberToGroup.Query, "grp-001", MemberTypeUser, "usr-1").
		Return(nil, errors.New("member fail")).
		Once()

	err := updateGroupMembers(txMock, "grp-001", []Member{{ID: "usr-1", Type: MemberTypeUser}})

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to assign members to group")
}
