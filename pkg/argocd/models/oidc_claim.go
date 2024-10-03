// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// OidcClaim oidc claim
//
// swagger:model oidcClaim
type OidcClaim struct {

	// essential
	Essential bool `json:"essential,omitempty"`

	// value
	Value string `json:"value,omitempty"`

	// values
	Values []string `json:"values"`
}

// Validate validates this oidc claim
func (m *OidcClaim) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this oidc claim based on context it is used
func (m *OidcClaim) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *OidcClaim) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *OidcClaim) UnmarshalBinary(b []byte) error {
	var res OidcClaim
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
