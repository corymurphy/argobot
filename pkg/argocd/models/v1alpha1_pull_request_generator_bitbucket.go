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

// V1alpha1PullRequestGeneratorBitbucket PullRequestGeneratorBitbucket defines connection info specific to Bitbucket.
//
// swagger:model v1alpha1PullRequestGeneratorBitbucket
type V1alpha1PullRequestGeneratorBitbucket struct {

	// The Bitbucket REST API URL to talk to. If blank, uses https://api.bitbucket.org/2.0.
	API string `json:"api,omitempty"`

	// basic auth
	BasicAuth *V1alpha1BasicAuthBitbucketServer `json:"basicAuth,omitempty"`

	// bearer token
	BearerToken *V1alpha1BearerTokenBitbucketCloud `json:"bearerToken,omitempty"`

	// Workspace to scan. Required.
	Owner string `json:"owner,omitempty"`

	// Repo name to scan. Required.
	Repo string `json:"repo,omitempty"`
}

// Validate validates this v1alpha1 pull request generator bitbucket
func (m *V1alpha1PullRequestGeneratorBitbucket) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateBasicAuth(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateBearerToken(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1PullRequestGeneratorBitbucket) validateBasicAuth(formats strfmt.Registry) error {
	if swag.IsZero(m.BasicAuth) { // not required
		return nil
	}

	if m.BasicAuth != nil {
		if err := m.BasicAuth.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basicAuth")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("basicAuth")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1PullRequestGeneratorBitbucket) validateBearerToken(formats strfmt.Registry) error {
	if swag.IsZero(m.BearerToken) { // not required
		return nil
	}

	if m.BearerToken != nil {
		if err := m.BearerToken.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bearerToken")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bearerToken")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1alpha1 pull request generator bitbucket based on the context it is used
func (m *V1alpha1PullRequestGeneratorBitbucket) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateBasicAuth(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateBearerToken(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1PullRequestGeneratorBitbucket) contextValidateBasicAuth(ctx context.Context, formats strfmt.Registry) error {

	if m.BasicAuth != nil {

		if swag.IsZero(m.BasicAuth) { // not required
			return nil
		}

		if err := m.BasicAuth.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("basicAuth")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("basicAuth")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1PullRequestGeneratorBitbucket) contextValidateBearerToken(ctx context.Context, formats strfmt.Registry) error {

	if m.BearerToken != nil {

		if swag.IsZero(m.BearerToken) { // not required
			return nil
		}

		if err := m.BearerToken.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("bearerToken")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("bearerToken")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1PullRequestGeneratorBitbucket) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1PullRequestGeneratorBitbucket) UnmarshalBinary(b []byte) error {
	var res V1alpha1PullRequestGeneratorBitbucket
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
