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

// V1alpha1SyncStrategy SyncStrategy controls the manner in which a sync is performed
//
// swagger:model v1alpha1SyncStrategy
type V1alpha1SyncStrategy struct {

	// apply
	Apply *V1alpha1SyncStrategyApply `json:"apply,omitempty"`

	// hook
	Hook *V1alpha1SyncStrategyHook `json:"hook,omitempty"`
}

// Validate validates this v1alpha1 sync strategy
func (m *V1alpha1SyncStrategy) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateApply(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateHook(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SyncStrategy) validateApply(formats strfmt.Registry) error {
	if swag.IsZero(m.Apply) { // not required
		return nil
	}

	if m.Apply != nil {
		if err := m.Apply.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("apply")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("apply")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1SyncStrategy) validateHook(formats strfmt.Registry) error {
	if swag.IsZero(m.Hook) { // not required
		return nil
	}

	if m.Hook != nil {
		if err := m.Hook.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hook")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("hook")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this v1alpha1 sync strategy based on the context it is used
func (m *V1alpha1SyncStrategy) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateApply(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateHook(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1SyncStrategy) contextValidateApply(ctx context.Context, formats strfmt.Registry) error {

	if m.Apply != nil {

		if swag.IsZero(m.Apply) { // not required
			return nil
		}

		if err := m.Apply.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("apply")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("apply")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1SyncStrategy) contextValidateHook(ctx context.Context, formats strfmt.Registry) error {

	if m.Hook != nil {

		if swag.IsZero(m.Hook) { // not required
			return nil
		}

		if err := m.Hook.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("hook")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("hook")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1SyncStrategy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1SyncStrategy) UnmarshalBinary(b []byte) error {
	var res V1alpha1SyncStrategy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
