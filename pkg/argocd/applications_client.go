package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	apiclient "github.com/argoproj/argo-cd/v2/reposerver/apiclient"
	// "github.com/argoproj/argo-cd/v2/cmpserver/apiclient"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/settings"
	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	// argp
	// models "github.com/corymurphy/argobot/pkg/argocd/models"
)

type ApplicationsClient struct {
	BaseUrl string
	Token   string
}

// TODO: refactor to be less repetitive

func (c *ApplicationsClient) Apply(app string, revision string) (string, error) {
	a := ApplicationSyncRequest{
		Name:     app,
		Revision: revision,
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/api/v1/applications/%s/sync", c.BaseUrl, app)
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))

	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to sync app %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *ApplicationsClient) List() (*argoappv1.ApplicationList, error) {
	var apps argoappv1.ApplicationList
	url := fmt.Sprintf("%s/api/v1/applications", c.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return &apps, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return &apps, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &apps, err
	}

	if err := json.Unmarshal(body, &apps); err != nil {
		return &apps, err
	}

	return &apps, nil
}

func (c *ApplicationsClient) Get(name string) (*argoappv1.Application, error) {
	var app *argoappv1.Application
	url := fmt.Sprintf("%s/api/v1/applications/%s", c.BaseUrl, name)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return app, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return app, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return app, err
	}

	if err := json.Unmarshal(body, &app); err != nil {
		return app, err
	}

	return app, nil
}

func (c *ApplicationsClient) GetManifest(name string, revision string) (*apiclient.ManifestResponse, error) {
	var app *apiclient.ManifestResponse

	// http://localhost:8081/api/v1/applications/helloworld/manifests?revision=47110b135dfe3e64e9199f66945532f378f05b4b

	url := fmt.Sprintf("%s/api/v1/applications/%s/manifests", c.BaseUrl, name)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return app, err
	}

	if revision != "" {
		query := req.URL.Query()
		query.Add("revision", revision)
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return app, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return app, err
	}

	if err := json.Unmarshal(body, &app); err != nil {
		return app, err
	}

	return app, nil
}

func (c *ApplicationsClient) ManagedResources(app string) (*application.ManagedResourcesResponse, error) {
	var resources *application.ManagedResourcesResponse
	url := fmt.Sprintf("%s/api/v1/applications/%s/managed-resources", c.BaseUrl, app)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return resources, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resources, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return resources, err
	}

	if err := json.Unmarshal(body, &resources); err != nil {
		return resources, err
	}

	return resources, nil
}

func (c *ApplicationsClient) GetSettings() (settings.Settings, error) {
	var settings settings.Settings
	url := fmt.Sprintf("%s/api/v1/settings", c.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return settings, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return settings, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return settings, err
	}

	if err := json.Unmarshal(body, &settings); err != nil {
		return settings, err
	}

	return settings, nil
}

func (c *ApplicationsClient) PutApplicationSpec(name string, spec *argoappv1.ApplicationSpec) error {
	data, err := json.Marshal(spec)
	if err != nil {
		return fmt.Errorf("failed to marshal application spec: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/applications/%s/spec", c.BaseUrl, name)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request to update application spec: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to update application spec: %s", body)
	}
	return nil
}
