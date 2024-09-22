package integration

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/shell"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func random(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func TestTemplateBase(t *testing.T) {
	releaseName := "integration"
	helmChartPath, _ := filepath.Abs("../../../charts/argobot")
	kubectlOptions := k8s.NewKubectlOptions(*kubeContext, "", *kubeNamespace)

	tag := random(8)
	image := fmt.Sprintf("ghcr.io/corymurphy/argobot:%s", tag)
	buildOptions := &docker.BuildOptions{
		Tags: []string{image},
	}
	docker.Build(t, "../../../", buildOptions)

	shell.RunCommand(t, shell.Command{
		Command: "kind",
		Args: []string{
			"load",
			"docker-image",
			image,
		},
	})

	options := &helm.Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"chartVersion": "0.0.1",
			"image.tag":    tag,
		},
		ExtraArgs: map[string][]string{
			"upgrade": {"--timeout", "15s", "--install", "--wait-for-jobs", "--wait", "--create-namespace", "--namespace", *kubeNamespace},
		},
		// ValuesFiles: []string{testCase.ValuesPath},
	}

	helm.Upgrade(t, options, helmChartPath, releaseName)
	defer helm.Delete(t, options, releaseName, true)

	// k8s.expo

	services := k8s.ListServices(t, kubectlOptions, v1.ListOptions{LabelSelector: fmt.Sprintf("app.kubernetes.io/name=argobot,app.kubernetes.io/instance=%s", releaseName)})
	if len(services) < 1 {
		t.Fatalf("expected at least 1 service, found %d", len(services))
	}
	for _, service := range services {
		// service.
		k8s.WaitUntilServiceAvailable(t, kubectlOptions, service.Name, 10, 1*time.Second)
	}
}
