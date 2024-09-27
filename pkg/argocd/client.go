package argocd

type Client interface {
	Plan(app string, sha string) (string, error)
	Apply(app string, revision string) (string, error)
}

type PlanClient interface {
	Plan(app string, sha string) (string, error)
}

type ApplyClient interface {
	Apply(app string, revision string) (string, error)
}
