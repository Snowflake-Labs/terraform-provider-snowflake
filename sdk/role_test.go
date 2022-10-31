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
