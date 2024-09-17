package argocd

type Client interface {
	Plan(app string, sha string) (string, error)
}
