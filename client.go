package auth

import (
	"errors"
	"net/http"

	"github.com/mrehanabbasi/supabase-auth-go/endpoints"
)

var ErrInvalidProjectReference = errors.New("cannot create auth client: invalid project reference")

var _ Client = &client{}

type client struct {
	*endpoints.Client
}

type Config struct {
	BaseURL string
	APIKey  string
}

// Set up a new Auth client.
//
// projectReference: The project reference is the unique identifier for your
// Supabase project. It can be found in the Supabase dashboard under project
// settings as Reference ID.
//
// apiKey: The API key is used to authenticate requests to the Auth server.
// This should be your anon key.
//
// This function does not validate your project reference. Requests will fail
// if you pass in an invalid project reference.
func New(projectReference string, apiKey string) Client {
	return &client{
		Client: endpoints.New(projectReference, apiKey),
	}
}

// Set up a new Auth client with custom auth URL.
//
// cfg: cfg is the configuration struct which contains custom auth URL
// and API Key. The API key is used to authenticate requests to the Auth server.
// This should be your anon key.
//
// This function does not validate your auth URL. Requests will fail
// if you pass in an invalid auth URL.
func NewWithCustomAuthURL(cfg Config) Client {
	return &client{
		Client: endpoints.New("", cfg.APIKey).WithCustomAuthURL(cfg.BaseURL),
	}
}

// Set up a new Auth client with custom auth URL.
//
// cfg: cfg is the configuration struct which contains custom auth URL
// and API Key. The API key is used to authenticate requests to the Auth server.
// This should be your anon key.
//
// c: This is the stdlib *http.Client if you want to use your own HTTP client.
//
// This function does not validate your auth URL. Requests will fail
// if you pass in an invalid auth URL.
func NewWithCustomAuthURLAndHTTPClient(cfg Config, c *http.Client) Client {
	return &client{
		Client: endpoints.
			New("", cfg.APIKey).
			WithCustomAuthURL(cfg.BaseURL).
			WithClient(c),
	}
}

func (c client) WithCustomAuthURL(url string) Client {
	return &client{
		Client: c.Client.WithCustomAuthURL(url),
	}
}

func (c client) WithToken(token string) Client {
	return &client{
		Client: c.Client.WithToken(token),
	}
}

func (c client) WithClient(httpClient *http.Client) Client {
	return &client{
		Client: c.Client.WithClient(httpClient),
	}
}
