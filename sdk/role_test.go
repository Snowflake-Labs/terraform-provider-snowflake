package sdk

import "context"

func (ts *testSuite) createRole() (*Role, error) {
	options := RoleCreateOptions{
		Name: "ROLE_TEST",
		RoleProperties: &RoleProperties{
			Comment: String("test account"),
		},
	}
	return ts.client.Roles.Create(context.Background(), options)
}

func (ts *testSuite) TestListRole() {
	role, err := ts.createRole()
	ts.NoError(err)

	roles, err := ts.client.Roles.List(context.Background(), RoleListOptions{Pattern: "ROLE%"})
	ts.NoError(err)
	ts.Equal(1, len(roles))

	ts.NoError(ts.client.Roles.Delete(context.Background(), role.Name))
}

func (ts *testSuite) TestReadRole() {
	role, err := ts.createRole()
	ts.NoError(err)

	entity, err := ts.client.Roles.Read(context.Background(), role.Name)
	ts.NoError(err)
	ts.Equal(role.Name, entity.Name)

	ts.NoError(ts.client.Roles.Delete(context.Background(), role.Name))
}

func (ts *testSuite) TestCreateRole() {
	role, err := ts.createRole()
	ts.NoError(err)
	ts.NoError(ts.client.Roles.Delete(context.Background(), role.Name))
}

func (ts *testSuite) TestUpdateRole() {
	role, err := ts.createRole()
	ts.NoError(err)

	options := RoleUpdateOptions{
		RoleProperties: &RoleProperties{
			Comment: String("updated comment"),
		},
	}
	afterUpdate, err := ts.client.Roles.Update(context.Background(), role.Name, options)
	ts.NoError(err)
	ts.Equal(*options.RoleProperties.Comment, afterUpdate.Comment)

	ts.NoError(ts.client.Roles.Delete(context.Background(), role.Name))
}
