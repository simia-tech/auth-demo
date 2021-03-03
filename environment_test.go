package authdemo_test

import (
	"testing"

	authdemo "github.com/simia-tech/auth-demo"
	"github.com/stretchr/testify/require"
)

type environment struct {
	serviceBaseURL string
	tearDown       func()
}

func setUpTestEnvironment(tb testing.TB) *environment {
	s, err := authdemo.NewService("tcp", "localhost:0")
	require.NoError(tb, err)

	return &environment{
		serviceBaseURL: s.BaseURL(),
		tearDown: func() {
			require.NoError(tb, s.Close())
		},
	}
}
