package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/laugart7272/inscrips/util"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func CreateRandomUser(t *testing.T) User {
	arg := CreateUserParams{
		Name:           util.RandomName(),
		LastName:       util.RandomName(),
		Email:          util.RandomEmail(),
		HashedPassword: util.RandomString(10),
		Phone:          util.RandomPhone(),
		UserType:       util.RandomName(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, arg.LastName, user.LastName)
	require.Equal(t, arg.Email, user.Email)

	require.NotZero(t, user.ID)
	require.NotZero(t, user.CreatedAt)

	return user
}

func TestCreateUser(t *testing.T) {
	CreateRandomUser(t)
}

func TestGetUser(t *testing.T) {
	User := CreateRandomUser(t)
	User1, err := testQueries.GetUser(context.Background(), User.ID)
	require.NoError(t, err)
	require.NotEmpty(t, User1)

	require.Equal(t, User.ID, User1.ID)
	require.Equal(t, User.Name, User1.Name)
	require.WithinDuration(t, User.CreatedAt.Time, User1.CreatedAt.Time, time.Second)
}

func TestUpdateUser(t *testing.T) {
	User := CreateRandomUser(t)

	arg := UpdateUserParams{
		ID:              User.ID,
		Name:            util.RandomName(),
		LastName:        util.RandomName(),
		Email:           util.RandomEmail(),
		HashedPassword:  util.RandomString(10),
		Phone:           util.RandomPhone(),
		UserType:        util.RandomName(),
		IsEmailVerified: true,
		Role:            util.ReViewerRole,
	}

	User1, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, User1)

	require.Equal(t, User.ID, User1.ID)
	// require.Equal(t, User.Balance, User1.Balance)
	require.WithinDuration(t, User.CreatedAt.Time, User1.CreatedAt.Time, time.Second)
}

func TestDeleteUser(t *testing.T) {
	User := CreateRandomUser(t)

	err := testQueries.DeleteUser(context.Background(), User.ID)
	require.NoError(t, err)

	User1, err := testQueries.GetUser(context.Background(), User.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, User1)
}

func TestListUsers(t *testing.T) {
	n := 10
	for i := 0; i < n; i++ {
		CreateRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	Users, err := testQueries.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, Users, 5)

	for _, User := range Users {
		require.NotEmpty(t, User)
	}
}
