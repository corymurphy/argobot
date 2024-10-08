package server_test

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"

	"github.com/corymurphy/argobot/pkg/argocd"
	"github.com/corymurphy/argobot/pkg/assert"
	"github.com/corymurphy/argobot/pkg/env"
	"github.com/corymurphy/argobot/pkg/github"
	"github.com/corymurphy/argobot/pkg/logging"
	"github.com/corymurphy/argobot/pkg/server"
	"github.com/rs/zerolog"
)

func Test_HealthCheck(t *testing.T) {
	s := server.NewServer(&env.Config{}, logging.NewLogger(logging.Silent), nil)
	req, _ := http.NewRequest("GET", "/health", bytes.NewBuffer(nil))
	w := httptest.NewRecorder()
	s.Health(w, req)
	assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	body, _ := io.ReadAll(w.Result().Body)
	assert.Equal(t, "application/json", w.Result().Header["Content-Type"][0])
	assert.Equal(t,
		`{
  "status": "ok"
}`, string(body))
}

func Test_PRCommentHandler(t *testing.T) {

	requests := map[string]bool{
		"/app/installations/345345345/access_tokens":     false,
		"/repos/atlas8518/argocd-data/issues/1/comments": false,
		"/repos/atlas8518/argocd-data/pulls/1":           false,
		"/api/v1/applications/testapp/sync":              false,
	}
	mu := sync.Mutex{}

	// ctx := context.Background()

	mockServer := mockServer(requests, &mu, t)
	defer mockServer.Close()

	config, err := env.ReadConfig("../testdata/argobot_config.yaml")
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(config.PrivateKeyFilePath)
	if err != nil {
		panic(err)
	}
	config.Github.V3APIURL = mockServer.URL
	config.Github.WebURL = mockServer.URL
	config.Github.App.PrivateKey = string(content)

	// logger := assert.NewTestLogger(t)
	logger := logging.NewLogger(logging.Silent)

	for _, testCase := range *NewServerTestCases(mockServer) {

		// TODO: fix this signature
		// s := server.NewServer(testCase.Config, logger, testCase.ApplyClient)
		s := server.NewServer(testCase.Config, logger, nil)

		w := httptest.NewRecorder()

		req := NewWebhookRequest(testCase.BodyPath, testCase.Config)

		s.PRCommentHandler.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	}

	for path, called := range requests {
		if !called {
			t.Errorf("%s was not called", path)
		}
	}
}

type ServerTestCase struct {
	BodyPath           string
	Config             *env.Config
	PlanClient         argocd.PlanClient
	ApplyClient        argocd.ApplyClient
	ExpectedStatusCode uint
	ExpectedRequests   map[string]bool
}

func NewServerTestCases(mockServer *httptest.Server) *[]ServerTestCase {

	config, err := env.ReadConfig("../testdata/argobot_config.yaml")
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(config.PrivateKeyFilePath)
	if err != nil {
		panic(err)
	}
	config.Github.V3APIURL = mockServer.URL
	config.Github.WebURL = mockServer.URL
	config.Github.App.PrivateKey = string(content)
	// config.ArgoConfig.Server = strings.ReplaceAll(mockServer.URL, "http://", "")

	return &[]ServerTestCase{
		{
			BodyPath:           "../testdata/comments/pullrequest_comment_bot_plan.json",
			Config:             config,
			ExpectedStatusCode: http.StatusOK,
		},
		{
			BodyPath: "../testdata/comments/pullrequest_comment_user_plan.json",
			Config:   config,
			// PlanClient:         NewMockPlanClient("../testdata/argocd_plan_diff"),

			ExpectedStatusCode: http.StatusOK,
		},
		{
			BodyPath: "../testdata/comments/pullrequest_comment_user_apply.json",
			Config:   config,
			ApplyClient: &argocd.ApplicationsClient{
				Token:   "123",
				BaseUrl: mockServer.URL,
			},
			ExpectedStatusCode: http.StatusOK,
		},
	}
}

func NewWebhookRequest(bodyPath string, config *env.Config) *http.Request {
	body, err := os.ReadFile(bodyPath)
	if err != nil {
		panic(err)
	}

	signatureSha1 := github.GenerateMAC(body, []byte(config.Github.App.WebhookSecret), sha1.New)
	signatureSha1Value := fmt.Sprintf("sha1=%s", hex.EncodeToString(signatureSha1))

	signatureSha256 := github.GenerateMAC(body, []byte(config.Github.App.WebhookSecret), sha256.New)
	signatureSha256Value := fmt.Sprintf("sha256=%s", hex.EncodeToString(signatureSha256))

	req, err := http.NewRequest("POST", "/api/github/hook", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-GitHub-Event", "issue_comment")
	req.Header.Add("X-GitHub-Hook-Installation-Target-ID", "345345345")
	req.Header.Add("X-GitHub-Hook-Installation-Target-Type", "integration")
	req.Header.Add("X-Hub-Signature", signatureSha1Value)
	req.Header.Add("X-Hub-Signature-256", signatureSha256Value)

	log := zerolog.New(os.Stdout)
	req = req.WithContext(log.WithContext(req.Context()))
	return req
}

func mockServer(requests map[string]bool, mu *sync.Mutex, t *testing.T) *httptest.Server {

	router := http.NewServeMux()
	router.HandleFunc("/app/installations/345345345/access_tokens", func(w http.ResponseWriter, r *http.Request) {
		response, _ := os.ReadFile("../testdata/github_access_token_response.json")
		mu.Lock()
		requests["/app/installations/345345345/access_tokens"] = true
		mu.Unlock()
		fmt.Fprint(w, string(response))
	})
	router.HandleFunc("/repos/atlas8518/argocd-data/issues/1/comments", func(w http.ResponseWriter, r *http.Request) {
		response, _ := os.ReadFile("../testdata/github_issues_comment_response.json")
		mu.Lock()
		requests["/repos/atlas8518/argocd-data/issues/1/comments"] = true
		actual, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		expected, err := os.ReadFile("../testdata/comments/pullrequest_comment_user_plan_diff.json")
		if err != nil {
			panic(err)
		}
		if string(expected) != string(actual) {
			t.Error("pull request comment plan diff did not match expected")
		}
		mu.Unlock()
		fmt.Fprint(w, string(response))
	})
	router.HandleFunc("/repos/atlas8518/argocd-data/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		response, _ := os.ReadFile("../testdata/github_get_pulls_response.json")
		mu.Lock()
		requests["/repos/atlas8518/argocd-data/pulls/1"] = true
		mu.Unlock()
		fmt.Fprint(w, string(response))
	})
	router.HandleFunc("/api/v1/applications/testapp/sync", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		requests["/api/v1/applications/testapp/sync"] = true
		mu.Unlock()
		fmt.Fprint(w, string(""))
	})
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t.Error(r)
	})
	return httptest.NewServer(router)
}
