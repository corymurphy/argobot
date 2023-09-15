package argocd

type Client interface {
	Diff(app string, sha string) (string, error)
}
