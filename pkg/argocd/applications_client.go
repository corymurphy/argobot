package argocd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	models "github.com/corymurphy/argobot/pkg/argocd/models"
)

type ApplicationsClient struct {
	BaseUrl string
	Token   string
}

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

func (c *ApplicationsClient) List() (models.V1alpha1ApplicationList, error) {
	var apps models.V1alpha1ApplicationList
	url := fmt.Sprintf("%s/api/v1/applications", c.BaseUrl)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return apps, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return apps, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return apps, err
	}

	if err := json.Unmarshal(body, &apps); err != nil {
		return apps, err
	}

	return apps, nil
}
