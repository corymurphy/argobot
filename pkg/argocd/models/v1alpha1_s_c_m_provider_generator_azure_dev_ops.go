// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1alpha1SCMProviderGeneratorAzureDevOps SCMProviderGeneratorAzureDevOps defines connection info specific to Azure DevOps.
//
// swagger:model v1alpha1SCMProviderGeneratorAzureDevOps
type V1alpha1SCMProviderGeneratorAzureDevOps struct {

	// access token ref
	AccessTokenRef *V1alpha1SecretRef `json:"accessTokenRef,omitempty"`

	// Scan all branches instead of just the default branch.
	AllBranches bool `json:"allBranches,omitempty"`

	// The URL to Azure DevOps. If blank, use https://dev.azure.com.
	API string `json:"api,omitempty"`

	// Azure Devops organization. Required. E.g. "my-organization".
	Organization string `json:"organization,omitempty"`

	// Azure Devops team project. Required. E.g. "my-team".
	TeamProject string `json:"teamProject,omitempty"`
}

// Validate validates this v1alpha1 s c m provider generator azure dev ops
func (m *V1alpha1SCMProviderGeneratorAzureDevOps) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateAccessTokenRef(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SCMProviderGeneratorAzureDevOps) validateAccessTokenRef(formats strfmt.Registry) error {
	if swag.IsZero(m.AccessTokenRef) { // not required
		return nil
	}

	if m.AccessTokenRef != nil {
		if err := m.AccessTokenRef.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("accessTokenRef")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("accessTokenRef")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1alpha1 s c m provider generator azure dev ops based on the context it is used
func (m *V1alpha1SCMProviderGeneratorAzureDevOps) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateAccessTokenRef(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SCMProviderGeneratorAzureDevOps) contextValidateAccessTokenRef(ctx context.Context, formats strfmt.Registry) error {

	if m.AccessTokenRef != nil {

		if swag.IsZero(m.AccessTokenRef) { // not required
			return nil
		}

		if err := m.AccessTokenRef.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("accessTokenRef")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("accessTokenRef")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1SCMProviderGeneratorAzureDevOps) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1SCMProviderGeneratorAzureDevOps) UnmarshalBinary(b []byte) error {
	var res V1alpha1SCMProviderGeneratorAzureDevOps
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
