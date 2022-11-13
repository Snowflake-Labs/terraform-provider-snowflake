package sdk

import (
	"context"
)

func (ts *testSuite) createUser() (*User, error) {
	options := UserCreateOptions{
		Name:     "USER_TEST",
		Password: String("Test123456"),
		UserProperties: &UserProperties{
			FirstName:             String("John"),
			LastName:              String("Hi"),
			Comment:               String("test user"),
			DefaultSecondaryRoles: StringSlice([]string{"ALL", "READ", "WRITE"}),
			Disabled:              Bool(false),
		},
	}
	return ts.client.Users.Create(context.Background(), options)
}

func (ts *testSuite) TestListUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	limit := 1
	users, err := ts.client.Users.List(context.Background(), UserListOptions{
		Pattern: "USER%",
		Limit:   Int(1),
	})
	ts.NoError(err)
	ts.Equal(limit, len(users))

	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}

func (ts *testSuite) TestReadUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	entity, err := ts.client.Users.Read(context.Background(), user.Name)
	ts.NoError(err)
	ts.Equal(user.Name, entity.Name)

	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}

func (ts *testSuite) TestCreateUser() {
	user, err := ts.createUser()
	ts.NoError(err)
	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}

func (ts *testSuite) TestUpdateUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	options := UserUpdateOptions{
		UserProperties: &UserProperties{
			Email:     String("test@gmail.com"),
			FirstName: String("Krebs"),
			LastName:  String("Great"),
		},
	}
	afterUpdate, err := ts.client.Users.Update(context.Background(), user.Name, options)
	ts.NoError(err)
	ts.Equal(*options.UserProperties.Email, afterUpdate.Email)
	ts.Equal(*options.UserProperties.FirstName, afterUpdate.FirstName)
	ts.Equal(*options.UserProperties.LastName, afterUpdate.LastName)

	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}

func (ts *testSuite) TestRenameUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	newUser := "NEW_USER_TEST"
	ts.NoError(ts.client.Users.Rename(context.Background(), user.Name, newUser))
	ts.NoError(ts.client.Users.Delete(context.Background(), newUser))
}

func (ts *testSuite) TestResetPassword() {
	user, err := ts.createUser()
	ts.NoError(err)

	result, err := ts.client.Users.ResetPassword(context.Background(), user.Name)
	ts.NoError(err)
	ts.NotEmpty(result.Status)
	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}
