package authdemo

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ory/fosite"
	"github.com/ory/fosite/compose"
	"github.com/ory/fosite/handler/oauth2"
	"github.com/ory/fosite/storage"
	"github.com/ory/fosite/token/hmac"
	"gopkg.in/square/go-jose.v2"
)

var hmacStrategy = &oauth2.HMACSHAStrategy{
	Enigma: &hmac.HMACStrategy{
		GlobalSecret: []byte("global-validation-key-0123456789"),
	},
	AccessTokenLifespan:   time.Hour,
	AuthorizeCodeLifespan: time.Minute,
}

// Service implements the authentication service.
type Service struct {
	listener net.Listener
	provider fosite.OAuth2Provider
}

// NewService returns an initialized service that listens to the provided network and address.
func NewService(network, address string) (*Service, error) {
	issuerKeyPair, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate key pair: %w", err)
	}

	store := storage.NewMemoryStore()
	store.Clients = map[string]fosite.Client{
		"auth-client": &fosite.DefaultClient{
			ID:         "auth-client",
			Secret:     []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
			GrantTypes: []string{"client_credentials"},
			Scopes:     []string{"fosite"},
		},
	}
	store.IssuerPublicKeys = map[string]storage.IssuerPublicKeys{
		"issuer@auth-demo.com": {
			Issuer: "issuer@auth-demo.com",
			KeysBySub: map[string]storage.SubjectPublicKeys{
				"auth-demo": {
					Subject: "auth-demo",
					Keys: map[string]storage.PublicKeyScopes{
						"123": {
							Key: &jose.JSONWebKey{
								Key:       issuerKeyPair.Public(),
								Algorithm: string(jose.RS256),
								Use:       "sig",
								KeyID:     "123",
							},
							Scopes: []string{"fosite"},
						},
					},
				},
			},
		},
	}

	l, err := net.Listen(network, address)
	if err != nil {
		return nil, fmt.Errorf("listen [%s %s]: %w", network, address, err)
	}
	log.Printf("opened http listener at %s %s", l.Addr().Network(), l.Addr().String())

	s := &Service{
		listener: l,
		provider: compose.Compose(
			&compose.Config{},
			store,
			hmacStrategy,
			nil,
			compose.OAuth2ClientCredentialsGrantFactory,
			compose.OAuth2TokenIntrospectionFactory,
		),
	}

	router := mux.NewRouter()
	router.HandleFunc("/token", s.token)

	go http.Serve(l, router)

	return s, nil
}

// Close tears down the service.
func (s *Service) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	return nil
}

// BaseURL returns the base url.
func (s *Service) BaseURL() string {
	return fmt.Sprintf("http://%s", s.listener.Addr())
}

func (s *Service) token(w http.ResponseWriter, r *http.Request) {
	ctx := fosite.NewContext()

	accessRequest, err := s.provider.NewAccessRequest(ctx, r, &oauth2.JWTSession{})
	if err != nil {
		s.provider.WriteAccessError(w, accessRequest, err)
		return
	}

	accessResponse, err := s.provider.NewAccessResponse(ctx, accessRequest)
	if err != nil {
		s.provider.WriteAccessError(w, accessRequest, err)
		return
	}

	s.provider.WriteAccessResponse(w, accessRequest, accessResponse)
}
