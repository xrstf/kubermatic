// Code generated by go-swagger; DO NOT EDIT.

package gke

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewListGKEDiskTypesParams creates a new ListGKEDiskTypesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListGKEDiskTypesParams() *ListGKEDiskTypesParams {
	return &ListGKEDiskTypesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListGKEDiskTypesParamsWithTimeout creates a new ListGKEDiskTypesParams object
// with the ability to set a timeout on a request.
func NewListGKEDiskTypesParamsWithTimeout(timeout time.Duration) *ListGKEDiskTypesParams {
	return &ListGKEDiskTypesParams{
		timeout: timeout,
	}
}

// NewListGKEDiskTypesParamsWithContext creates a new ListGKEDiskTypesParams object
// with the ability to set a context for a request.
func NewListGKEDiskTypesParamsWithContext(ctx context.Context) *ListGKEDiskTypesParams {
	return &ListGKEDiskTypesParams{
		Context: ctx,
	}
}

// NewListGKEDiskTypesParamsWithHTTPClient creates a new ListGKEDiskTypesParams object
// with the ability to set a custom HTTPClient for a request.
func NewListGKEDiskTypesParamsWithHTTPClient(client *http.Client) *ListGKEDiskTypesParams {
	return &ListGKEDiskTypesParams{
		HTTPClient: client,
	}
}

/* ListGKEDiskTypesParams contains all the parameters to send to the API endpoint
   for the list g k e disk types operation.

   Typically these are written to a http.Request.
*/
type ListGKEDiskTypesParams struct {
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list g k e disk types params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListGKEDiskTypesParams) WithDefaults() *ListGKEDiskTypesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list g k e disk types params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListGKEDiskTypesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list g k e disk types params
func (o *ListGKEDiskTypesParams) WithTimeout(timeout time.Duration) *ListGKEDiskTypesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list g k e disk types params
func (o *ListGKEDiskTypesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list g k e disk types params
func (o *ListGKEDiskTypesParams) WithContext(ctx context.Context) *ListGKEDiskTypesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list g k e disk types params
func (o *ListGKEDiskTypesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list g k e disk types params
func (o *ListGKEDiskTypesParams) WithHTTPClient(client *http.Client) *ListGKEDiskTypesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list g k e disk types params
func (o *ListGKEDiskTypesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WriteToRequest writes these params to a swagger request
func (o *ListGKEDiskTypesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}