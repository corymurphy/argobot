package argocd

import (
	"context"
	"fmt"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/corymurphy/argobot/pkg/logging"
)

var LockKey = "corymurphy.io/argobot/lockedby"

type Locker struct {
	ArgoClient *ApplicationsClient
	Log        logging.SimpleLogging
}

func NewLocker(client *ApplicationsClient, log logging.SimpleLogging) *Locker {
	return &Locker{
		ArgoClient: client,
		Log:        log,
	}
}
func (l *Locker) Lock(ctx context.Context, name string, by string) error {
	app, err := l.ArgoClient.Get(name)
	if err != nil {
		return fmt.Errorf("failed to get application %s for locking", name)
	}
	if app == nil {
		return fmt.Errorf("application %s not found", name)
	}

	hasKey := false
	for i, info := range app.Spec.Info {
		if info.Name == LockKey {
			app.Spec.Info[i].Value = by
			l.Log.Debug("locking application %s by %s", name, by)
			hasKey = true
			break
		}
	}

	if !hasKey {
		app.Spec.Info = append(app.Spec.Info, v1alpha1.Info{
			Name:  LockKey,
			Value: by,
		})
		l.Log.Debug("adding lock key for application %s by %s", name, by)
	}
	if err := l.ArgoClient.PutApplicationSpec(name, &app.Spec); err != nil {
		return fmt.Errorf("failed to update application %s with lock information: %w", name, err)
	}
	l.Log.Debug("application %s locked by %s", name, by)

	return nil
}

func (l *Locker) Unlock(ctx context.Context, name string) error {
	app, err := l.ArgoClient.Get(name)
	if err != nil {
		return fmt.Errorf("failed to get application %s for unlocking", name)
	}
	if app == nil {
		return fmt.Errorf("application %s not found", name)
	}

	for i, info := range app.Spec.Info {
		if info.Name == LockKey {
			app.Spec.Info = append(app.Spec.Info[:i], app.Spec.Info[i+1:]...)
			l.Log.Debug("unlocking application %s", name)
			break
		}
	}

	if err := l.ArgoClient.PutApplicationSpec(name, &app.Spec); err != nil {
		return fmt.Errorf("failed to update application %s with unlock information: %w", name, err)
	}
	l.Log.Debug("application %s unlocked", name)

	return nil
}
