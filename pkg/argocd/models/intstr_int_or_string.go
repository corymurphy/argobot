// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// IntstrIntOrString IntOrString is a type that can hold an int32 or a string.  When used in
// JSON or YAML marshalling and unmarshalling, it produces or consumes the
// inner type.  This allows you to have, for example, a JSON field that can
// accept a name or number.
// TODO: Rename to Int32OrString
//
// +protobuf=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
// +k8s:openapi-gen=true
//
// swagger:model intstrIntOrString
type IntstrIntOrString struct {

	// int val
	IntVal int32 `json:"intVal,omitempty"`

	// str val
	StrVal string `json:"strVal,omitempty"`

	// type
	Type int64 `json:"type,omitempty"`
}

// Validate validates this intstr int or string
func (m *IntstrIntOrString) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this intstr int or string based on context it is used
func (m *IntstrIntOrString) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *IntstrIntOrString) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *IntstrIntOrString) UnmarshalBinary(b []byte) error {
	var res IntstrIntOrString
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}