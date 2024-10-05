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

// ProjectSyncWindowsResponse project sync windows response
//
// swagger:model projectSyncWindowsResponse
type ProjectSyncWindowsResponse struct {

	// windows
	Windows []*V1alpha1SyncWindow `json:"windows"`
}

// Validate validates this project sync windows response
func (m *ProjectSyncWindowsResponse) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateWindows(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSyncWindowsResponse) validateWindows(formats strfmt.Registry) error {
	if swag.IsZero(m.Windows) { // not required
		return nil
	}

	for i := 0; i < len(m.Windows); i++ {
		if swag.IsZero(m.Windows[i]) { // not required
			continue
		}

		if m.Windows[i] != nil {
			if err := m.Windows[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("windows" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("windows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this project sync windows response based on the context it is used
func (m *ProjectSyncWindowsResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateWindows(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ProjectSyncWindowsResponse) contextValidateWindows(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Windows); i++ {

		if m.Windows[i] != nil {

			if swag.IsZero(m.Windows[i]) { // not required
				return nil
			}

			if err := m.Windows[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("windows" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("windows" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *ProjectSyncWindowsResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ProjectSyncWindowsResponse) UnmarshalBinary(b []byte) error {
	var res ProjectSyncWindowsResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}