package argocd

import (
	models "github.com/corymurphy/argobot/pkg/argocd/models"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func LiveObjects(resources []*models.V1alpha1ResourceDiff) ([]*unstructured.Unstructured, error) {
	objs := make([]*unstructured.Unstructured, len(resources))
	for i, resState := range resources {
		obj, err := resState.LiveObject()
		if err != nil {
			return nil, err
		}
		objs[i] = obj
	}
	return objs, nil
}

type ResourceKey struct {
	Group     string
	Kind      string
	Namespace string
	Name      string
}
