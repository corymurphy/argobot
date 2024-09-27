package argocd

type ApplicationSyncRequest struct {
	Name     string `json:"name"`
	Revision string `json:"revision"`
	DryRun   bool   `json:"dryRun"`
}
