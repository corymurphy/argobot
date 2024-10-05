package argocd

import (
	"context"
	"encoding/json"

	"github.com/argoproj/gitops-engine/pkg/sync/hook"
	"github.com/argoproj/gitops-engine/pkg/sync/ignore"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	"github.com/argoproj/pkg/errors"
	models "github.com/corymurphy/argobot/pkg/argocd/models"
	settings "google.golang.org/genproto/googleapis/cloud/securitycenter/settings/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// type DifferenceOption struct {
// 	local           string
// 	localRepoRoot   string
// 	revision        string
// 	cluster         *argoappv1.Cluster
// 	res             *repoapiclient.ManifestResponse
// 	serversideRes   *repoapiclient.ManifestResponse
// 	revisions       []string
// 	sourcePositions []int64
// }

func Plan(ctx context.Context, app *models.ApplicationApplicationResponse, resources *models.ApplicationManagedResourcesResponse, target *models.RepositoryManifestResponse, revision string) (string, error) {
	var plan string

	// liveObjs, err := LiveObjects(resources.Items)
	// if err != nil {
	// 	return plan, err
	// }

	// items := make([]objKeyLiveTarget, 0)

	// var unstructureds []*unstructured.Unstructured

	// for _, mfst := range target.Manifests {
	// 	obj, err := models.UnmarshalToUnstructured(mfst)
	// 	if err != nil {
	// 		return plan, nil
	// 	}
	// 	unstructureds = append(unstructureds, obj)
	// }
	// groupedObjs := groupObjsByKey(unstructureds, liveObjs, app.Spec.Destination.Namespace)
	// items = groupObjsForDiff(resources, groupedObjs, items, argoSettings, app.InstanceName(argoSettings.ControllerNamespace), app.Spec.Destination.Namespace)

	// state := resources.Items[0].LiveObject()
	// state.LiveState
	// liveObjs, err := cmdutil.LiveObjects(resources.Items)

	// get app -> app *argoappv1.Application
	// liveObjs, err := cmdutil.LiveObjects(resources.Items)
	// http://localhost:8081/api/v1/projects/{name}/detailed

	return plan, nil
}

func groupObjsByKey(localObs []*unstructured.Unstructured, liveObjs []*unstructured.Unstructured, appNamespace string) map[kube.ResourceKey]*unstructured.Unstructured {
	namespacedByGk := make(map[schema.GroupKind]bool)
	for i := range liveObjs {
		if liveObjs[i] != nil {
			key := kube.GetResourceKey(liveObjs[i])
			namespacedByGk[schema.GroupKind{Group: key.Group, Kind: key.Kind}] = key.Namespace != ""
		}
	}
	localObs, _, err := controller.DeduplicateTargetObjects(appNamespace, localObs, &resourceInfoProvider{namespacedByGk: namespacedByGk})
	errors.CheckError(err)
	objByKey := make(map[kube.ResourceKey]*unstructured.Unstructured)
	for i := range localObs {
		obj := localObs[i]
		if !(hook.IsHook(obj) || ignore.Ignore(obj)) {
			objByKey[kube.GetResourceKey(obj)] = obj
		}
	}
	return objByKey
}

func groupObjsForDiff(resources *application.ManagedResourcesResponse, objs map[kube.ResourceKey]*unstructured.Unstructured, items []objKeyLiveTarget, argoSettings *settings.Settings, appName, namespace string) []objKeyLiveTarget {
	resourceTracking := argo.NewResourceTracking()
	for _, res := range resources.Items {
		live := &unstructured.Unstructured{}
		err := json.Unmarshal([]byte(res.NormalizedLiveState), &live)
		errors.CheckError(err)

		key := kube.ResourceKey{Name: res.Name, Namespace: res.Namespace, Group: res.Group, Kind: res.Kind}
		if key.Kind == kube.SecretKind && key.Group == "" {
			// Don't bother comparing secrets, argo-cd doesn't have access to k8s secret data
			delete(objs, key)
			continue
		}
		if local, ok := objs[key]; ok || live != nil {
			if local != nil && !kube.IsCRD(local) {
				err = resourceTracking.SetAppInstance(local, argoSettings.AppLabelKey, appName, namespace, argoappv1.TrackingMethod(argoSettings.GetTrackingMethod()))
				errors.CheckError(err)
			}

			items = append(items, objKeyLiveTarget{key, live, local})
			delete(objs, key)
		}
	}
	for key, local := range objs {
		if key.Kind == kube.SecretKind && key.Group == "" {
			// Don't bother comparing secrets, argo-cd doesn't have access to k8s secret data
			delete(objs, key)
			continue
		}
		items = append(items, objKeyLiveTarget{key, nil, local})
	}
	return items
}

// TODO: this is just temporary while i build the proof of concept
type objKeyLiveTarget struct {
	key    ResourceKey
	live   *unstructured.Unstructured
	target *unstructured.Unstructured
}
