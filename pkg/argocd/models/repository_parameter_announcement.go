// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// RepositoryParameterAnnouncement repository parameter announcement
//
// swagger:model repositoryParameterAnnouncement
type RepositoryParameterAnnouncement struct {

	// array is the default value of the parameter if the parameter is an array.
	Array []string `json:"array"`

	// collectionType is the type of value this parameter holds - either a single value (a string) or a collection
	// (array or map). If collectionType is set, only the field with that type will be used. If collectionType is not
	// set, `string` is the default. If collectionType is set to an invalid value, a validation error is thrown.
	CollectionType string `json:"collectionType,omitempty"`

	// itemType determines the primitive data type represented by the parameter. Parameters are always encoded as
	// strings, but this field lets them be interpreted as other primitive types.
	ItemType string `json:"itemType,omitempty"`

	// map is the default value of the parameter if the parameter is a map.
	Map map[string]string `json:"map,omitempty"`

	// name is the name identifying a parameter.
	Name string `json:"name,omitempty"`

	// required defines if this given parameter is mandatory.
	Required bool `json:"required,omitempty"`

	// string is the default value of the parameter if the parameter is a string.
	String string `json:"string,omitempty"`

	// title is a human-readable text of the parameter name.
	Title string `json:"title,omitempty"`

	// tooltip is a human-readable description of the parameter.
	Tooltip string `json:"tooltip,omitempty"`
}

// Validate validates this repository parameter announcement
func (m *RepositoryParameterAnnouncement) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this repository parameter announcement based on context it is used
func (m *RepositoryParameterAnnouncement) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *RepositoryParameterAnnouncement) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *RepositoryParameterAnnouncement) UnmarshalBinary(b []byte) error {
	var res RepositoryParameterAnnouncement
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}