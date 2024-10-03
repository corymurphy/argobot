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

// V1alpha1ApplicationSourceJsonnet ApplicationSourceJsonnet holds options specific to applications of type Jsonnet
//
// swagger:model v1alpha1ApplicationSourceJsonnet
type V1alpha1ApplicationSourceJsonnet struct {

	// ExtVars is a list of Jsonnet External Variables
	ExtVars []*V1alpha1JsonnetVar `json:"extVars"`

	// Additional library search dirs
	Libs []string `json:"libs"`

	// TLAS is a list of Jsonnet Top-level Arguments
	Tlas []*V1alpha1JsonnetVar `json:"tlas"`
}

// Validate validates this v1alpha1 application source jsonnet
func (m *V1alpha1ApplicationSourceJsonnet) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateExtVars(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTlas(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1ApplicationSourceJsonnet) validateExtVars(formats strfmt.Registry) error {
	if swag.IsZero(m.ExtVars) { // not required
		return nil
	}

	for i := 0; i < len(m.ExtVars); i++ {
		if swag.IsZero(m.ExtVars[i]) { // not required
			continue
		}

		if m.ExtVars[i] != nil {
			if err := m.ExtVars[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("extVars" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("extVars" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1alpha1ApplicationSourceJsonnet) validateTlas(formats strfmt.Registry) error {
	if swag.IsZero(m.Tlas) { // not required
		return nil
	}

	for i := 0; i < len(m.Tlas); i++ {
		if swag.IsZero(m.Tlas[i]) { // not required
			continue
		}

		if m.Tlas[i] != nil {
			if err := m.Tlas[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tlas" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tlas" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1alpha1 application source jsonnet based on the context it is used
func (m *V1alpha1ApplicationSourceJsonnet) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateExtVars(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateTlas(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1ApplicationSourceJsonnet) contextValidateExtVars(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.ExtVars); i++ {

		if m.ExtVars[i] != nil {

			if swag.IsZero(m.ExtVars[i]) { // not required
				return nil
			}

			if err := m.ExtVars[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("extVars" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("extVars" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1alpha1ApplicationSourceJsonnet) contextValidateTlas(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Tlas); i++ {

		if m.Tlas[i] != nil {

			if swag.IsZero(m.Tlas[i]) { // not required
				return nil
			}

			if err := m.Tlas[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("tlas" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("tlas" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1ApplicationSourceJsonnet) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1ApplicationSourceJsonnet) UnmarshalBinary(b []byte) error {
	var res V1alpha1ApplicationSourceJsonnet
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
