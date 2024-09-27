package unit

import (
	"io/fs"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
	"github.com/gruntwork-io/terratest/modules/k8s"
)

func TestTemplateBase(t *testing.T) {
	options := &helm.Options{
		KubectlOptions: k8s.NewKubectlOptions("", "", "phuquoc"),
		SetValues: map[string]string{
			"chartVersion": "0.0.1",
		},
		// ValuesFiles: []string{testCase.ValuesPath},
	}
	actual := helm.RenderTemplate(t, options, "../../charts/argobot", "test", nil)

	if *update {
		t.Error("here")
		err := os.WriteFile("../testdata/template_expected.yaml", []byte(actual), fs.FileMode(0644))
		if err != nil {
			t.Fatal(err)
		}
	}

	expected, _ := os.ReadFile("../testdata/template_expected.yaml")

	if string(expected) != actual {
		// TODO: This is just temporary, i'm working on a better experience in a future pr
		// in the meantime, use your own diff tool for debugging
		os.WriteFile("../../.debug/actual.yaml", []byte(actual), 0644)
		os.WriteFile("../../.debug/expected.yaml", []byte(expected), 0644)
		t.Fatalf(`
rendered output does not match

-------- expected ---------
%s
---------------------------

-------- actual  ----------
%s
---------------------------

`, expected, actual)
	}
}
