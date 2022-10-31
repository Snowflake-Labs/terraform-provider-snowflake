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
			Comment:               String("test account"),
			DefaultSecondaryRoles: StringSlice([]string{"ALL", "READ", "WRITE"}),
			Disabled:              Bool(false),
		},
	}
	return ts.client.Users.Create(context.Background(), options)
}

func (ts *testSuite) TestListUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	users, err := ts.client.Users.List(context.Background(), UserListOptions{Pattern: "USER%"})
	ts.NoError(err)
	ts.Equal(1, len(users))

	ts.NoError(ts.client.Users.Delete(context.Background(), user.Name))
}

func (ts *testSuite) TestReadUser() {
	user, err := ts.createUser()
	ts.NoError(err)

	entity, err := ts.client.Users.Read(context.Background(), user.Name)
	ts.NoError(err)
	ts.Equal(entity.Name, user.Name)

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
