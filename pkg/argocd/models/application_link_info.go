// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ApplicationLinkInfo application link info
//
// swagger:model applicationLinkInfo
type ApplicationLinkInfo struct {

	// description
	Description string `json:"description,omitempty"`

	// icon class
	IconClass string `json:"iconClass,omitempty"`

	// title
	Title string `json:"title,omitempty"`

	// url
	URL string `json:"url,omitempty"`
}

// Validate validates this application link info
func (m *ApplicationLinkInfo) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this application link info based on context it is used
func (m *ApplicationLinkInfo) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationLinkInfo) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationLinkInfo) UnmarshalBinary(b []byte) error {
	var res ApplicationLinkInfo
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
