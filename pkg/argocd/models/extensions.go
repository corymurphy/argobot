package models

import (
	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func (r V1alpha1ResourceDiff) LiveObject() (*unstructured.Unstructured, error) {
	return UnmarshalToUnstructured(r.LiveState)
}

func UnmarshalToUnstructured(resource string) (*unstructured.Unstructured, error) {
	if resource == "" || resource == "null" {
		return nil, nil
	}
	var obj unstructured.Unstructured
	err := json.Unmarshal([]byte(resource), &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}
