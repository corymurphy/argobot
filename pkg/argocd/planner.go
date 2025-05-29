package argocd

import (
	"context"

	"github.com/corymurphy/argobot/pkg/logging"
)

type Planner struct {
	ArgoClient *ApplicationsClient
	Log        logging.SimpleLogging
}

func NewPlanner(client *ApplicationsClient, log logging.SimpleLogging) *Planner {
	return &Planner{
		ArgoClient: client,
		Log:        log,
	}
}

func (p *Planner) Plan(ctx context.Context, name string, revision string) (string, bool, error) {
	var plan string
	var diff bool = false

	resources, err := p.ArgoClient.ManagedResources(name)

	if err != nil {
		return plan, diff, err
	}

	live, err := p.ArgoClient.Get(name)

	if err != nil {
		return plan, diff, err
	}

	// live.Spec.Info
	// p.ArgoClient.

	// for i, info := range live.Spec.Info {
	// 	if info.Name == "corymurphy.io/argobot/lockedby" {
	// 		info.Value =
	// 		break
	// 	}
	// }

	target, err := p.ArgoClient.GetManifest(name, revision)
	if err != nil {
		return plan, diff, err
	}

	settings, err := p.ArgoClient.GetSettings()
	if err != nil {
		return plan, diff, err
	}

	return Plan(ctx, &settings, live, resources, target, revision, p.Log)
}
