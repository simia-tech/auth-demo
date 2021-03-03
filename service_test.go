package authdemo_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func TestAuthentication(t *testing.T) {
	e := setUpTestEnvironment(t)
	defer e.tearDown()

	testFn := func(clientID, clientSecret string, expectErr error, expectAccessToken bool) func(*testing.T) {
		return func(t *testing.T) {
			client := &clientcredentials.Config{
				ClientID:     clientID,
				ClientSecret: clientSecret,
				Scopes:       []string{"fosite"},
				TokenURL:     e.serviceBaseURL + "/token",
				AuthStyle:    oauth2.AuthStyleInHeader,
			}

			token, err := client.Token(oauth2.NoContext)
			if expectErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
			if expectAccessToken {
				assert.NotEmpty(t, token.AccessToken)
			}
		}
	}

	t.Run("UsingValidClientCredentials", testFn("auth-client", "foobar", nil, true))
	t.Run("UsingInvalidClientCredentials", testFn("auth-client", "invalid", &oauth2.RetrieveError{}, false))
}
