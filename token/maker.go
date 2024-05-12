package token

import "time"

// Interface for managing tokens
type Maker interface {
	//Create a new token for username and specific duration
	CreateToken(user_id int64, email string, role string, duration time.Duration) (string, *Payload, error)

	//Check if the token is valid or not
	VerifyToken(token string) (*Payload, error)
}
