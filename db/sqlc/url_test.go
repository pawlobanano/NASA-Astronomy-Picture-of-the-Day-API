package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/pawlobanano/UGF3ZcWCIEdvZ29BcHBzIE5BU0E/generator"
	"github.com/stretchr/testify/require"
)

func TestCreateURL(t *testing.T) {
	// given
	arg := CreateURLParams{
		ID:  uuid.New(),
		URL: generator.RandomURL(),
	}

	// when
	createdURL, err := testQueries.CreateURL(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, createdURL)

	// then
	require.Equal(t, arg.ID, createdURL.ID)
	require.Equal(t, arg.URL, createdURL.URL)
	require.NotZero(t, createdURL.ID)
	require.NotZero(t, createdURL.CreatedAt)
}
