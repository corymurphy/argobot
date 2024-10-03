// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// SessionGetUserInfoResponse The current user's userInfo info
//
// swagger:model sessionGetUserInfoResponse
type SessionGetUserInfoResponse struct {

	// groups
	Groups []string `json:"groups"`

	// iss
	Iss string `json:"iss,omitempty"`

	// logged in
	LoggedIn bool `json:"loggedIn,omitempty"`

	// username
	Username string `json:"username,omitempty"`
}

// Validate validates this session get user info response
func (m *SessionGetUserInfoResponse) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this session get user info response based on context it is used
func (m *SessionGetUserInfoResponse) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *SessionGetUserInfoResponse) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *SessionGetUserInfoResponse) UnmarshalBinary(b []byte) error {
	var res SessionGetUserInfoResponse
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
