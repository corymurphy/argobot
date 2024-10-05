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

// RepositoryHelmAppSpec HelmAppSpec contains helm app name  in source repo
//
// swagger:model repositoryHelmAppSpec
type RepositoryHelmAppSpec struct {

	// helm file parameters
	FileParameters []*V1alpha1HelmFileParameter `json:"fileParameters"`

	// name
	Name string `json:"name,omitempty"`

	// the output of `helm inspect values`
	Parameters []*V1alpha1HelmParameter `json:"parameters"`

	// value files
	ValueFiles []string `json:"valueFiles"`

	// the contents of values.yaml
	Values string `json:"values,omitempty"`
}

// Validate validates this repository helm app spec
func (m *RepositoryHelmAppSpec) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFileParameters(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateParameters(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RepositoryHelmAppSpec) validateFileParameters(formats strfmt.Registry) error {
	if swag.IsZero(m.FileParameters) { // not required
		return nil
	}

	for i := 0; i < len(m.FileParameters); i++ {
		if swag.IsZero(m.FileParameters[i]) { // not required
			continue
		}

		if m.FileParameters[i] != nil {
			if err := m.FileParameters[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("fileParameters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("fileParameters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *RepositoryHelmAppSpec) validateParameters(formats strfmt.Registry) error {
	if swag.IsZero(m.Parameters) { // not required
		return nil
	}

	for i := 0; i < len(m.Parameters); i++ {
		if swag.IsZero(m.Parameters[i]) { // not required
			continue
		}

		if m.Parameters[i] != nil {
			if err := m.Parameters[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("parameters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("parameters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this repository helm app spec based on the context it is used
func (m *RepositoryHelmAppSpec) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateFileParameters(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateParameters(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RepositoryHelmAppSpec) contextValidateFileParameters(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.FileParameters); i++ {

		if m.FileParameters[i] != nil {

			if swag.IsZero(m.FileParameters[i]) { // not required
				return nil
			}

			if err := m.FileParameters[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("fileParameters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("fileParameters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *RepositoryHelmAppSpec) contextValidateParameters(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Parameters); i++ {

		if m.Parameters[i] != nil {

			if swag.IsZero(m.Parameters[i]) { // not required
				return nil
			}

			if err := m.Parameters[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("parameters" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("parameters" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *RepositoryHelmAppSpec) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepositoryHelmAppSpec) UnmarshalBinary(b []byte) error {
	var res RepositoryHelmAppSpec
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}