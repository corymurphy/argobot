package dev

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/corymurphy/argobot/pkg/utils"
	"github.com/gruntwork-io/terratest/modules/docker"
	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func Test_LocalDevelopmentInstall(t *testing.T) {
	// TODO: make these cli args
	releaseName := "argobot-dev"
	namespace := releaseName
	kubeContext := "kind-kind"
	helmChartPath, _ := filepath.Abs("../../charts/argobot")
	kubectlOptions := k8s.NewKubectlOptions(kubeContext, "", namespace)

	version := "0.25.0"

	tag := utils.InsecureRandom(8)
	image := fmt.Sprintf("ghcr.io/corymurphy/containers/argobot:%s", tag)
	buildOptions := &docker.BuildOptions{
		Tags: []string{image},
		Load: true,
	}
	docker.Build(t, "../../", buildOptions)
	// shell.RunCommand(t, shell.Command{
	// 	Command: "kind",
	// 	Args: []string{
	// 		"load",
	// 		"docker-image",
	// 		image,
	// 	},
	// })

	k8s.CreateNamespaceE(t, kubectlOptions, namespace)
	k8s.KubectlApply(t, kubectlOptions, "../../.secrets/secrets.yaml")

	options := &helm.Options{
		KubectlOptions: kubectlOptions,
		SetValues: map[string]string{
			"chartVersion":       strings.TrimSpace(version),
			"image.tag":          tag,
			"webServer.logLevel": "debug",
		},
		ExtraArgs: map[string][]string{
			"upgrade": {"--timeout", "15s", "--install", "--wait-for-jobs", "--wait", "--create-namespace", "--namespace", namespace},
		},
		ValuesFiles: []string{"../../.secrets/values.yaml"},
	}
	helm.Upgrade(t, options, helmChartPath, releaseName)
}
