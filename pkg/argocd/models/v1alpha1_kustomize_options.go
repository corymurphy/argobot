// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1alpha1KustomizeOptions KustomizeOptions are options for kustomize to use when building manifests
//
// swagger:model v1alpha1KustomizeOptions
type V1alpha1KustomizeOptions struct {

	// BinaryPath holds optional path to kustomize binary
	BinaryPath string `json:"binaryPath,omitempty"`

	// BuildOptions is a string of build parameters to use when calling `kustomize build`
	BuildOptions string `json:"buildOptions,omitempty"`
}

// Validate validates this v1alpha1 kustomize options
func (m *V1alpha1KustomizeOptions) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v1alpha1 kustomize options based on context it is used
func (m *V1alpha1KustomizeOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1KustomizeOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1KustomizeOptions) UnmarshalBinary(b []byte) error {
	var res V1alpha1KustomizeOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
