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

// V1alpha1ComparedTo ComparedTo contains application source and target which was used for resources comparison
//
// swagger:model v1alpha1ComparedTo
type V1alpha1ComparedTo struct {

	// destination
	Destination *V1alpha1ApplicationDestination `json:"destination,omitempty"`

	// IgnoreDifferences is a reference to the application's ignored differences used for comparison
	IgnoreDifferences []*V1alpha1ResourceIgnoreDifferences `json:"ignoreDifferences"`

	// source
	Source *V1alpha1ApplicationSource `json:"source,omitempty"`

	// Sources is a reference to the application's multiple sources used for comparison
	Sources []*V1alpha1ApplicationSource `json:"sources"`
}

// Validate validates this v1alpha1 compared to
func (m *V1alpha1ComparedTo) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateDestination(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateIgnoreDifferences(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSource(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSources(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1ComparedTo) validateDestination(formats strfmt.Registry) error {
	if swag.IsZero(m.Destination) { // not required
		return nil
	}

	if m.Destination != nil {
		if err := m.Destination.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destination")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destination")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1ComparedTo) validateIgnoreDifferences(formats strfmt.Registry) error {
	if swag.IsZero(m.IgnoreDifferences) { // not required
		return nil
	}

	for i := 0; i < len(m.IgnoreDifferences); i++ {
		if swag.IsZero(m.IgnoreDifferences[i]) { // not required
			continue
		}

		if m.IgnoreDifferences[i] != nil {
			if err := m.IgnoreDifferences[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ignoreDifferences" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("ignoreDifferences" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1alpha1ComparedTo) validateSource(formats strfmt.Registry) error {
	if swag.IsZero(m.Source) { // not required
		return nil
	}

	if m.Source != nil {
		if err := m.Source.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("source")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("source")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1ComparedTo) validateSources(formats strfmt.Registry) error {
	if swag.IsZero(m.Sources) { // not required
		return nil
	}

	for i := 0; i < len(m.Sources); i++ {
		if swag.IsZero(m.Sources[i]) { // not required
			continue
		}

		if m.Sources[i] != nil {
			if err := m.Sources[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("sources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("sources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// ContextValidate validate this v1alpha1 compared to based on the context it is used
func (m *V1alpha1ComparedTo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateDestination(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateIgnoreDifferences(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSource(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSources(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1alpha1ComparedTo) contextValidateDestination(ctx context.Context, formats strfmt.Registry) error {

	if m.Destination != nil {

		if swag.IsZero(m.Destination) { // not required
			return nil
		}

		if err := m.Destination.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("destination")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("destination")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1ComparedTo) contextValidateIgnoreDifferences(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.IgnoreDifferences); i++ {

		if m.IgnoreDifferences[i] != nil {

			if swag.IsZero(m.IgnoreDifferences[i]) { // not required
				return nil
			}

			if err := m.IgnoreDifferences[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("ignoreDifferences" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("ignoreDifferences" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *V1alpha1ComparedTo) contextValidateSource(ctx context.Context, formats strfmt.Registry) error {

	if m.Source != nil {

		if swag.IsZero(m.Source) { // not required
			return nil
		}

		if err := m.Source.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("source")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("source")
			}
			return err
		}
	}

	return nil
}

func (m *V1alpha1ComparedTo) contextValidateSources(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Sources); i++ {

		if m.Sources[i] != nil {

			if swag.IsZero(m.Sources[i]) { // not required
				return nil
			}

			if err := m.Sources[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("sources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("sources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1alpha1ComparedTo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1alpha1ComparedTo) UnmarshalBinary(b []byte) error {
	var res V1alpha1ComparedTo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}