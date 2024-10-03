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

// V1alpha1KustomizeReplica v1alpha1 kustomize replica
//
// swagger:model v1alpha1KustomizeReplica
type V1alpha1KustomizeReplica struct {

	// count
	Count *IntstrIntOrString `json:"count,omitempty"`

	// Name of Deployment or StatefulSet
	Name string `json:"name,omitempty"`
}

// Validate validates this v1alpha1 kustomize replica
func (m *V1alpha1KustomizeReplica) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCount(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1KustomizeReplica) validateCount(formats strfmt.Registry) error {
	if swag.IsZero(m.Count) { // not required
		return nil
	}

	if m.Count != nil {
		if err := m.Count.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("count")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("count")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1alpha1 kustomize replica based on the context it is used
func (m *V1alpha1KustomizeReplica) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCount(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1KustomizeReplica) contextValidateCount(ctx context.Context, formats strfmt.Registry) error {

	if m.Count != nil {

		if swag.IsZero(m.Count) { // not required
			return nil
		}

		if err := m.Count.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("count")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("count")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1KustomizeReplica) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1KustomizeReplica) UnmarshalBinary(b []byte) error {
	var res V1alpha1KustomizeReplica
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
