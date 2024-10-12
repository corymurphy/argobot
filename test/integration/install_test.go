package integration

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestTemplateBase(t *testing.T) {
	releaseName := "integration"
	helmChartPath, _ := filepath.Abs("../../charts/argobot")
	argoCdChartPath, _ := filepath.Abs("../../charts/argocd")
	kubectlOptions := k8s.NewKubectlOptions(*kubeContext, "", *kubeNamespace)

	tag := utils.InsecureRandom(8)
	image := fmt.Sprintf("ghcr.io/corymurphy/containers/argobot:%s", tag)
	buildOptions := &docker.BuildOptions{
		Tags: []string{image},
	}
	docker.Build(t, "../../", buildOptions)

	shell.RunCommand(t, shell.Command{
		Command: "kind",
		Args: []string{
			"load",
			"docker-image",
			image,
		},
	})

	k8s.CreateNamespaceE(t, kubectlOptions, *kubeNamespace)
	k8s.KubectlApply(t, kubectlOptions, "./testdata/secret.yaml")

	argocd := &helm.Options{
		KubectlOptions: kubectlOptions,
		ValuesFiles:    []string{argoCdChartPath + "/values.yaml"},
		SetValues:      map[string]string{"argo-cd.crds.install": *createCrds, "apps.create": "false"},
		ExtraArgs: map[string][]string{
			"upgrade": {"--timeout", "60s", "--install", "--wait-for-jobs", "--wait", "--create-namespace", "--namespace", *kubeNamespace},
		},
	}

	helm.Upgrade(t, argocd, argoCdChartPath, "argocd1")
	argocd.SetValues["apps.create"] = "true"
	helm.Upgrade(t, argocd, argoCdChartPath, "argocd1")
	defer helm.Delete(t, argocd, "argocd1", true)

	options := &helm.Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"image.tag": tag,
		},
		ExtraArgs: map[string][]string{
			"upgrade": {"--timeout", "30s", "--install", "--wait-for-jobs", "--wait", "--create-namespace", "--namespace", *kubeNamespace},
		},
	}

	helm.Upgrade(t, options, helmChartPath, releaseName)
	defer helm.Delete(t, options, releaseName, true)

	services := k8s.ListServices(t, kubectlOptions, v1.ListOptions{LabelSelector: fmt.Sprintf("app.kubernetes.io/name=argobot,app.kubernetes.io/instance=%s", releaseName)})
	if len(services) < 1 {
		t.Fatalf("expected at least 1 service, found %d", len(services))
	}
	for _, service := range services {
		k8s.WaitUntilServiceAvailable(t, kubectlOptions, service.Name, 10, 1*time.Second)
	}
}
