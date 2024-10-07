package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/argoproj/argo-cd/v2/controller"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/settings"
	argoappv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	apiclient "github.com/argoproj/argo-cd/v2/reposerver/apiclient"
	"github.com/argoproj/argo-cd/v2/util/argo"
	argodiff "github.com/argoproj/argo-cd/v2/util/argo/diff"
	"github.com/argoproj/argo-cd/v2/util/argo/normalizers"
	"github.com/argoproj/gitops-engine/pkg/sync/hook"
	"github.com/argoproj/gitops-engine/pkg/sync/ignore"
	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	models "github.com/corymurphy/argobot/pkg/argocd/models"
	"github.com/corymurphy/argobot/pkg/logging"
	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func Plan(ctx context.Context, argoSettings *settings.Settings, app *argoappv1.Application, resources *application.ManagedResourcesResponse, target *apiclient.ManifestResponse, revision string, log logging.SimpleLogging) (string, bool, error) {
	var plan string
	hasDiff := false

	liveObjs, err := LiveObjects(resources.Items)
	if err != nil {
		return plan, hasDiff, err
	}

	items := make([]objKeyLiveTarget, 0)

	var unstructureds []*unstructured.Unstructured

	for _, mfst := range target.Manifests {
		obj, err := models.UnmarshalToUnstructured(mfst)
		if err != nil {
			return plan, hasDiff, nil
		}
		unstructureds = append(unstructureds, obj)
	}

	groupedObjs, err := groupObjsByKey(unstructureds, liveObjs, app.Spec.Destination.Namespace)
	if err != nil {
		return plan, hasDiff, err
	}

	items, err = groupObjsForDiff(resources, groupedObjs, items, argoSettings, app.InstanceName(argoSettings.ControllerNamespace), app.Spec.Destination.Namespace)
	if err != nil {
		return plan, hasDiff, err
	}

	for _, item := range items {
		if item.target != nil && hook.IsHook(item.target) || item.live != nil && hook.IsHook(item.live) {
			continue
		}
		overrides := make(map[string]argoappv1.ResourceOverride)
		for k := range argoSettings.ResourceOverrides {
			val := argoSettings.ResourceOverrides[k]
			overrides[k] = *val
		}

		// TODO remove hardcoded IgnoreAggregatedRoles and retrieve the
		// compareOptions in the protobuf
		ignoreAggregatedRoles := false
		ignoreNormalizerOpts := normalizers.IgnoreNormalizerOpts{
			JQExecutionTimeout: normalizers.DefaultJQExecutionTimeout,
		}
		diffConfig, err := argodiff.NewDiffConfigBuilder().
			WithDiffSettings(app.Spec.IgnoreDifferences, overrides, ignoreAggregatedRoles, ignoreNormalizerOpts).
			WithTracking(argoSettings.AppLabelKey, argoSettings.TrackingMethod).
			WithNoCache().
			// WithLogger(logutils.NewLogrusLogger(logutils.NewWithCurrentConfig())).
			Build()
		if err != nil {
			return plan, hasDiff, err
		}
		diffRes, err := argodiff.StateDiff(item.live, item.target, diffConfig)
		if err != nil {
			return plan, hasDiff, err
		}

		if diffRes.Modified || item.target == nil || item.live == nil {
			plan = plan + fmt.Sprintf("\n===== %s/%s %s/%s ======\n", item.key.Group, item.key.Kind, item.key.Namespace, item.key.Name)
			var live *unstructured.Unstructured
			var target *unstructured.Unstructured
			if item.target != nil && item.live != nil {
				target = &unstructured.Unstructured{}
				live = item.live
				err = json.Unmarshal(diffRes.PredictedLive, target)
				if err != nil {
					return plan, hasDiff, err
				}
			} else {
				live = item.live
				target = item.target
			}
			if !hasDiff {
				hasDiff = true
			}
			out, err := diff(item.key.Name, live, target)
			if err != nil {
				return out, hasDiff, err
			}
			plan = plan + out
		}
	}
	return plan, hasDiff, nil
}

func diff(name string, live *unstructured.Unstructured, target *unstructured.Unstructured) (string, error) {
	var result string
	tempDir, err := os.MkdirTemp("", "argocd-diff")
	if err != nil {
		return result, err
	}
	targetFile := path.Join(tempDir, name)
	targetData := []byte("")
	if target != nil {
		targetData, err = yaml.Marshal(target)
		if err != nil {
			return result, err
		}
	}
	err = os.WriteFile(targetFile, targetData, 0o644)
	if err != nil {
		return result, err
	}
	liveFile := path.Join(tempDir, fmt.Sprintf("%s-live.yaml", name))
	liveData := []byte("")
	if live != nil {
		liveData, err = yaml.Marshal(live)
		if err != nil {
			return result, err
		}
	}
	err = os.WriteFile(liveFile, liveData, 0o644)
	if err != nil {
		return result, err
	}
	var args []string
	cmd := exec.Command("diff", append(args, liveFile, targetFile)...)
	if cmd.Err != nil {
		return result, cmd.Err
	}

	// diff is expected to return non-zero
	output, _ := cmd.CombinedOutput()
	return string(output), nil
}

func groupObjsByKey(localObs []*unstructured.Unstructured, liveObjs []*unstructured.Unstructured, appNamespace string) (map[kube.ResourceKey]*unstructured.Unstructured, error) {
	namespacedByGk := make(map[schema.GroupKind]bool)
	for i := range liveObjs {
		if liveObjs[i] != nil {
			key := kube.GetResourceKey(liveObjs[i])
			namespacedByGk[schema.GroupKind{Group: key.Group, Kind: key.Kind}] = key.Namespace != ""
		}
	}
	localObs, _, err := controller.DeduplicateTargetObjects(appNamespace, localObs, &resourceInfoProvider{namespacedByGk: namespacedByGk})
	if err != nil {
		return nil, err
	}
	objByKey := make(map[kube.ResourceKey]*unstructured.Unstructured)
	for i := range localObs {
		obj := localObs[i]
		if !(hook.IsHook(obj) || ignore.Ignore(obj)) {
			objByKey[kube.GetResourceKey(obj)] = obj
		}
	}
	return objByKey, nil
}

func groupObjsForDiff(resources *application.ManagedResourcesResponse, objs map[kube.ResourceKey]*unstructured.Unstructured, items []objKeyLiveTarget, argoSettings *settings.Settings, appName, namespace string) ([]objKeyLiveTarget, error) {
	resourceTracking := argo.NewResourceTracking()
	for _, res := range resources.Items {
		live := &unstructured.Unstructured{}
		err := json.Unmarshal([]byte(res.NormalizedLiveState), &live)
		// errors.CheckError(err)
		if err != nil {
			return items, err
		}

		key := kube.ResourceKey{Name: res.Name, Namespace: res.Namespace, Group: res.Group, Kind: res.Kind}
		if key.Kind == kube.SecretKind && key.Group == "" {
			// Don't bother comparing secrets, argo-cd doesn't have access to k8s secret data
			delete(objs, key)
			continue
		}
		if local, ok := objs[key]; ok || live != nil {
			if local != nil && !kube.IsCRD(local) {
				err = resourceTracking.SetAppInstance(local, argoSettings.AppLabelKey, appName, namespace, argoappv1.TrackingMethod(argoSettings.GetTrackingMethod()))
				if err != nil {
					return items, err
				}
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
	return items, nil
}

// TODO: this is just temporary while i build the proof of concept
type objKeyLiveTarget struct {
	key    kube.ResourceKey
	live   *unstructured.Unstructured
	target *unstructured.Unstructured
}

type resourceInfoProvider struct {
	namespacedByGk map[schema.GroupKind]bool
}

func (p *resourceInfoProvider) IsNamespaced(gk schema.GroupKind) (bool, error) {
	return p.namespacedByGk[gk], nil
}
