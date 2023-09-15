// temporary until this is used by the cli
// nolint,used,deadcode
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/corymurphy/argobot/pkg/auth"
	"github.com/corymurphy/argobot/pkg/env"
)

type APIClient struct {
	Host   string // example: api.github.com
	Config *env.Config
	Client *http.Client
}

func NewApiClient(config env.Config) *APIClient {
	return &APIClient{
		Config: &config,
		Host:   "api.github.com", // TODO: allow this to be configurable?
		Client: &http.Client{},
	}
}

func (a *APIClient) NewAccessToken(installationID string) (AccessToken, error) {

	jwt, err := auth.CreateJWT([]byte(a.Config.Github.App.PrivateKey))

	if err != nil {
		return AccessToken{}, err
	}

	url := fmt.Sprintf("https://%s/app/installations/%s/access_tokens", a.Host, installationID)

	req, err := http.NewRequest("POST", url, nil)

	if err != nil {
		return AccessToken{}, fmt.Errorf("error while creating access token request %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := a.Client.Do(req)

	if err != nil {
		return AccessToken{}, fmt.Errorf("error while requesting new access token: %s", err)
	}

	body, err := a.parseResponseBody(resp)

	if err != nil {
		return AccessToken{}, err
	}

	token := AccessToken{}

	if err := json.Unmarshal(body, &token); err != nil {
		return AccessToken{}, fmt.Errorf("unable to deserialize github response %s", fmt.Sprint(body))
	}

	return token, nil
}

func (a *APIClient) parseResponseBody(resp *http.Response) ([]byte, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte("no body returned"), fmt.Errorf("unable to parse response body %s", err)
	}
	return body, nil
}
