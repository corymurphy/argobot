package github_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/corymurphy/argobot/pkg/assert"
	"github.com/google/go-github/v72/github"
	"github.com/gregjones/httpcache"
	"github.com/palantir/go-githubapp/githubapp"
	"gopkg.in/yaml.v2"
)

func mockServer(requests map[string]bool, mu *sync.Mutex) *httptest.Server {

	router := http.NewServeMux()
	router.HandleFunc("/app/installations/123456/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		response, _ := os.ReadFile("../testdata/github_access_token_response.json")
		mu.Lock()
		requests["/app/installations/123456/access_tokens"] = true
		mu.Unlock()
		fmt.Fprint(w, string(response))
	})
	router.HandleFunc("/repos/atlas8518/argocd-data/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		response, _ := os.ReadFile("../testdata/github_issues_comment_response.json")
		mu.Lock()
		requests["/repos/atlas8518/argocd-data/issues/1/comments"] = true
		mu.Unlock()
		fmt.Fprint(w, string(response))
	})

	return httptest.NewServer(router)
}

func Test_TokenRequested(t *testing.T) {

	requests := map[string]bool{
		"/app/installations/123456/access_tokens":        false,
		"/repos/atlas8518/argocd-data/issues/1/comments": false,
	}
	mu := sync.Mutex{}

	ctx := context.Background()

	mockServer := mockServer(requests, &mu)
	defer mockServer.Close()

	var event github.IssueCommentEvent
	payload, _ := os.ReadFile("../testdata/comments/pullrequest_comment_user_plan.json")

	if err := json.Unmarshal(payload, &event); err != nil {
		t.Error(err)
	}

	serialized := fmt.Sprintf(`
web_url: %[1]s
v3_api_url: %[1]s
app:
  integration_id: 123456
  webhookSecret: fc1b391fa17718cfdbf9497ec9bfe59
`, mockServer.URL)

	var config githubapp.Config
	if err := yaml.Unmarshal([]byte(serialized), &config); err != nil {
		t.Error(err)
	}

	content, _ := os.ReadFile("../../pkg/testdata/github_app.pem")
	config.App.PrivateKey = string(content)

	cc, err := githubapp.NewDefaultCachingClientCreator(
		config,
		githubapp.WithClientUserAgent("argobot/0.1.0"),
		githubapp.WithClientTimeout(3*time.Second),
		githubapp.WithClientCaching(false, func() httpcache.Cache { return httpcache.NewMemoryCache() }),
	)

	if err != nil {
		t.Error(err)
	}

	client, err := cc.NewInstallationClient(config.App.IntegrationID)
	if err != nil {
		t.Error(err)
	}

	repo := event.GetRepo()
	prNum := event.GetIssue().GetNumber()
	repoOwner := repo.GetOwner().GetLogin()
	repoName := repo.GetName()
	body := "test"

	response := github.IssueComment{
		Body: &body,
	}

	if _, _, err = client.Issues.CreateComment(ctx, repoOwner, repoName, prNum, &response); err != nil {
		t.Error(err)
	}

	for path, requested := range requests {
		assert.True(t, requested, path)
	}
}
