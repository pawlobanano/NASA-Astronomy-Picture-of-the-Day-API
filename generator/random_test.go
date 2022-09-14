package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRandomString(t *testing.T) {
	// given
	var n uint8 = 10

	// when
	result := RandomString(n, AlphanumericLowerDash)

	// then
	require.NotEmpty(t, result)
	require.Len(t, result, 10)
}

func TestRandomStringUpperBoundary(t *testing.T) {
	// given
	var n uint8 = 255

	// when
	result := RandomString(n, AlphanumericLowerDash)

	// then
	require.NotEmpty(t, result)
	require.Len(t, result, 255)
}

func TestRandomStringLowerBoundary(t *testing.T) {
	// given
	var n uint8 = 0

	// when
	result := RandomString(n, AlphanumericLowerDash)

	// then
	require.Empty(t, result)
	require.Len(t, result, 0)
	require.Equal(t, "", result)
}

func TestRandomEmailContainsSuffix(t *testing.T) {
	// given
	prefix := "https://"
	suffix := ".com"

	// when
	result := RandomURL()

	// then
	require.NotEmpty(t, result)
	require.Contains(t, result, prefix)
	require.Contains(t, result, suffix)
	require.Len(t, result, 32)
}
