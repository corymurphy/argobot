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

// V1alpha1DuckTypeGenerator DuckType defines a generator to match against clusters registered with ArgoCD.
//
// swagger:model v1alpha1DuckTypeGenerator
type V1alpha1DuckTypeGenerator struct {

	// ConfigMapRef is a ConfigMap with the duck type definitions needed to retrieve the data
	//              this includes apiVersion(group/version), kind, matchKey and validation settings
	// Name is the resource name of the kind, group and version, defined in the ConfigMapRef
	// RequeueAfterSeconds is how long before the duckType will be rechecked for a change
	ConfigMapRef string `json:"configMapRef,omitempty"`

	// label selector
	LabelSelector *V1LabelSelector `json:"labelSelector,omitempty"`

	// name
	Name string `json:"name,omitempty"`

	// requeue after seconds
	RequeueAfterSeconds int64 `json:"requeueAfterSeconds,omitempty"`

	// template
	Template *V1alpha1ApplicationSetTemplate `json:"template,omitempty"`

	// Values contains key/value pairs which are passed directly as parameters to the template
	Values map[string]string `json:"values,omitempty"`
}

// Validate validates this v1alpha1 duck type generator
func (m *V1alpha1DuckTypeGenerator) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateLabelSelector(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTemplate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1DuckTypeGenerator) validateLabelSelector(formats strfmt.Registry) error {
	if swag.IsZero(m.LabelSelector) { // not required
		return nil
	}

	if m.LabelSelector != nil {
		if err := m.LabelSelector.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("labelSelector")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("labelSelector")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1DuckTypeGenerator) validateTemplate(formats strfmt.Registry) error {
	if swag.IsZero(m.Template) { // not required
		return nil
	}

	if m.Template != nil {
		if err := m.Template.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("template")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("template")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1alpha1 duck type generator based on the context it is used
func (m *V1alpha1DuckTypeGenerator) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateLabelSelector(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTemplate(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1DuckTypeGenerator) contextValidateLabelSelector(ctx context.Context, formats strfmt.Registry) error {

	if m.LabelSelector != nil {

		if swag.IsZero(m.LabelSelector) { // not required
			return nil
		}

		if err := m.LabelSelector.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("labelSelector")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("labelSelector")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1DuckTypeGenerator) contextValidateTemplate(ctx context.Context, formats strfmt.Registry) error {

	if m.Template != nil {

		if swag.IsZero(m.Template) { // not required
			return nil
		}

		if err := m.Template.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("template")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("template")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1DuckTypeGenerator) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1DuckTypeGenerator) UnmarshalBinary(b []byte) error {
	var res V1alpha1DuckTypeGenerator
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
