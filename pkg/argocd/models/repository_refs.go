// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RepositoryRefs A subset of the repository's named refs
//
// swagger:model repositoryRefs
type RepositoryRefs struct {

	// branches
	Branches []string `json:"branches"`

	// tags
	Tags []string `json:"tags"`
}

// Validate validates this repository refs
func (m *RepositoryRefs) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this repository refs based on context it is used
func (m *RepositoryRefs) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RepositoryRefs) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepositoryRefs) UnmarshalBinary(b []byte) error {
	var res RepositoryRefs
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
