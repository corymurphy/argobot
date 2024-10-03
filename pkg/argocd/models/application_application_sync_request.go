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

// ApplicationApplicationSyncRequest ApplicationSyncRequest is a request to apply the config state to live state
//
// swagger:model applicationApplicationSyncRequest
type ApplicationApplicationSyncRequest struct {

	// app namespace
	AppNamespace string `json:"appNamespace,omitempty"`

	// dry run
	DryRun bool `json:"dryRun,omitempty"`

	// infos
	Infos []*V1alpha1Info `json:"infos"`

	// manifests
	Manifests []string `json:"manifests"`

	// name
	Name string `json:"name,omitempty"`

	// project
	Project string `json:"project,omitempty"`

	// prune
	Prune bool `json:"prune,omitempty"`

	// resources
	Resources []*V1alpha1SyncOperationResource `json:"resources"`

	// retry strategy
	RetryStrategy *V1alpha1RetryStrategy `json:"retryStrategy,omitempty"`

	// revision
	Revision string `json:"revision,omitempty"`

	// revisions
	Revisions []string `json:"revisions"`

	// source positions
	SourcePositions []string `json:"sourcePositions"`

	// strategy
	Strategy *V1alpha1SyncStrategy `json:"strategy,omitempty"`

	// sync options
	SyncOptions *ApplicationSyncOptions `json:"syncOptions,omitempty"`
}

// Validate validates this application application sync request
func (m *ApplicationApplicationSyncRequest) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateInfos(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateResources(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRetryStrategy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateStrategy(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateSyncOptions(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationApplicationSyncRequest) validateInfos(formats strfmt.Registry) error {
	if swag.IsZero(m.Infos) { // not required
		return nil
	}

	for i := 0; i < len(m.Infos); i++ {
		if swag.IsZero(m.Infos[i]) { // not required
			continue
		}

		if m.Infos[i] != nil {
			if err := m.Infos[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("infos" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("infos" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) validateResources(formats strfmt.Registry) error {
	if swag.IsZero(m.Resources) { // not required
		return nil
	}

	for i := 0; i < len(m.Resources); i++ {
		if swag.IsZero(m.Resources[i]) { // not required
			continue
		}

		if m.Resources[i] != nil {
			if err := m.Resources[i].Validate(formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("resources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("resources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) validateRetryStrategy(formats strfmt.Registry) error {
	if swag.IsZero(m.RetryStrategy) { // not required
		return nil
	}

	if m.RetryStrategy != nil {
		if err := m.RetryStrategy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("retryStrategy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("retryStrategy")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) validateStrategy(formats strfmt.Registry) error {
	if swag.IsZero(m.Strategy) { // not required
		return nil
	}

	if m.Strategy != nil {
		if err := m.Strategy.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("strategy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("strategy")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) validateSyncOptions(formats strfmt.Registry) error {
	if swag.IsZero(m.SyncOptions) { // not required
		return nil
	}

	if m.SyncOptions != nil {
		if err := m.SyncOptions.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("syncOptions")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("syncOptions")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this application application sync request based on the context it is used
func (m *ApplicationApplicationSyncRequest) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateInfos(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateResources(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateRetryStrategy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateStrategy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateSyncOptions(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ApplicationApplicationSyncRequest) contextValidateInfos(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Infos); i++ {

		if m.Infos[i] != nil {

			if swag.IsZero(m.Infos[i]) { // not required
				return nil
			}

			if err := m.Infos[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("infos" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("infos" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) contextValidateResources(ctx context.Context, formats strfmt.Registry) error {

	for i := 0; i < len(m.Resources); i++ {

		if m.Resources[i] != nil {

			if swag.IsZero(m.Resources[i]) { // not required
				return nil
			}

			if err := m.Resources[i].ContextValidate(ctx, formats); err != nil {
				if ve, ok := err.(*errors.Validation); ok {
					return ve.ValidateName("resources" + "." + strconv.Itoa(i))
				} else if ce, ok := err.(*errors.CompositeError); ok {
					return ce.ValidateName("resources" + "." + strconv.Itoa(i))
				}
				return err
			}
		}

	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) contextValidateRetryStrategy(ctx context.Context, formats strfmt.Registry) error {

	if m.RetryStrategy != nil {

		if swag.IsZero(m.RetryStrategy) { // not required
			return nil
		}

		if err := m.RetryStrategy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("retryStrategy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("retryStrategy")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) contextValidateStrategy(ctx context.Context, formats strfmt.Registry) error {

	if m.Strategy != nil {

		if swag.IsZero(m.Strategy) { // not required
			return nil
		}

		if err := m.Strategy.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("strategy")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("strategy")
			}
			return err
		}
	}

	return nil
}

func (m *ApplicationApplicationSyncRequest) contextValidateSyncOptions(ctx context.Context, formats strfmt.Registry) error {

	if m.SyncOptions != nil {

		if swag.IsZero(m.SyncOptions) { // not required
			return nil
		}

		if err := m.SyncOptions.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("syncOptions")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("syncOptions")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ApplicationApplicationSyncRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ApplicationApplicationSyncRequest) UnmarshalBinary(b []byte) error {
	var res ApplicationApplicationSyncRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
