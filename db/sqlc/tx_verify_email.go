package db

import (
	"context"
)

type VerifyEmailTxParams struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *Store) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})

		if err != nil {
			return err
		}

		//Get User
		user, err := q.GetUser(ctx, int64(result.VerifyEmail.ID))
		if err != nil {
			return err
		}

		user, err = q.UpdateUser(ctx, UpdateUserParams{
			ID:              user.ID,
			Name:            user.Name,
			LastName:        user.LastName,
			Email:           user.Email,
			Phone:           user.Phone,
			HashedPassword:  user.HashedPassword,
			UserType:        user.UserType,
			IsEmailVerified: true,
			Role:            user.Role,
			AvatarPath:      user.AvatarPath,
		})

		return err
	})

	return result, err
}
