// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"strconv"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1alpha1SCMProviderGeneratorAWSCodeCommit SCMProviderGeneratorAWSCodeCommit defines connection info specific to AWS CodeCommit.
//
// swagger:model v1alpha1SCMProviderGeneratorAWSCodeCommit
type V1alpha1SCMProviderGeneratorAWSCodeCommit struct {

	// Scan all branches instead of just the default branch.
	AllBranches bool `json:"allBranches,omitempty"`

	// Region provides the AWS region to discover repos.
	// if not provided, AppSet controller will infer the current region from environment.
	Region string `json:"region,omitempty"`

	// Role provides the AWS IAM role to assume, for cross-account repo discovery
	// if not provided, AppSet controller will use its pod/node identity to discover.
	Role string `json:"role,omitempty"`

	// TagFilters provides the tag filter(s) for repo discovery
	TagFilters []*V1alpha1TagFilter `json:"tagFilters"`
}

// Validate validates this v1alpha1 s c m provider generator a w s code commit
func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTagFilters(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) validateTagFilters(formats strfmt.Registry) error {
	if swag.IsZero(m.TagFilters) { // not required
		return nil
	}

	for i := 0; i < len(m.TagFilters); i++ {
		if swag.IsZero(m.TagFilters[i]) { // not required
			continue
		}

		if m.TagFilters[i] != nil {
			if err := m.TagFilters[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tagFilters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tagFilters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1alpha1 s c m provider generator a w s code commit based on the context it is used
func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTagFilters(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) contextValidateTagFilters(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.TagFilters); i++ {

		if m.TagFilters[i] != nil {

			if swag.IsZero(m.TagFilters[i]) { // not required
				return nil
			}

			if err := m.TagFilters[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tagFilters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tagFilters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1SCMProviderGeneratorAWSCodeCommit) UnmarshalBinary(b []byte) error {
	var res V1alpha1SCMProviderGeneratorAWSCodeCommit
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}