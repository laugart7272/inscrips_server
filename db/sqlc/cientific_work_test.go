package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindInscription(t *testing.T) {
	found := false
	//Find if Exists
	arg_findinscription := FindInscriptionParams{
		UserID:  1,
		EventID: 1,
	}

	inscription, err := testQueries.FindInscription(context.Background(), arg_findinscription)
	if err != nil {
		if err == sql.ErrNoRows {
			// return nil, status.Errorf(codes.Internal, "inscription not foud: %s", err)
			//Creating Inscription
			found = false
		}
	} else {
		//Inscription Found !
		found = true
	}

	require.NoError(t, err)
	require.NotEmpty(t, inscription)

	require.True(t, found)

}
