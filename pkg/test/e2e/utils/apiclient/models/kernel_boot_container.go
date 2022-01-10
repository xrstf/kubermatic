// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// KernelBootContainer If set, the VM will be booted from the defined kernel / initrd.
//
// swagger:model KernelBootContainer
type KernelBootContainer struct {

	// Image that contains initrd / kernel files.
	Image string `json:"image,omitempty"`

	// ImagePullSecret is the name of the Docker registry secret required to pull the image. The secret must already exist.
	// +optional
	ImagePullSecret string `json:"imagePullSecret,omitempty"`

	// the fully-qualified path to the ramdisk image in the host OS
	// +optional
	InitrdPath string `json:"initrdPath,omitempty"`

	// The fully-qualified path to the kernel image in the host OS
	// +optional
	KernelPath string `json:"kernelPath,omitempty"`

	// image pull policy
	ImagePullPolicy PullPolicy `json:"imagePullPolicy,omitempty"`
}

// Validate validates this kernel boot container
func (m *KernelBootContainer) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateImagePullPolicy(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *KernelBootContainer) validateImagePullPolicy(formats strfmt.Registry) error {
	if swag.IsZero(m.ImagePullPolicy) { // not required
		return nil
	}

	if err := m.ImagePullPolicy.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("imagePullPolicy")
		}
		return err
	}

	return nil
}

// ContextValidate validate this kernel boot container based on the context it is used
func (m *KernelBootContainer) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateImagePullPolicy(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *KernelBootContainer) contextValidateImagePullPolicy(ctx context.Context, formats strfmt.Registry) error {

	if err := m.ImagePullPolicy.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("imagePullPolicy")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *KernelBootContainer) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *KernelBootContainer) UnmarshalBinary(b []byte) error {
	var res KernelBootContainer
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}