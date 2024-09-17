package argocd

type MockClient struct {
	PlanResponse string
	Error        error
}

func (c *MockClient) Plan(app string, sha string) (string, error) {
	return c.PlanResponse, c.Error
}
