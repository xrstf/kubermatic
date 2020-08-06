// Code generated by go-swagger; DO NOT EDIT.

package openstack

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// ListOpenstackTenantsNoCredentialsReader is a Reader for the ListOpenstackTenantsNoCredentials structure.
type ListOpenstackTenantsNoCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListOpenstackTenantsNoCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListOpenstackTenantsNoCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListOpenstackTenantsNoCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListOpenstackTenantsNoCredentialsOK creates a ListOpenstackTenantsNoCredentialsOK with default headers values
func NewListOpenstackTenantsNoCredentialsOK() *ListOpenstackTenantsNoCredentialsOK {
	return &ListOpenstackTenantsNoCredentialsOK{}
}

/*ListOpenstackTenantsNoCredentialsOK handles this case with default header values.

OpenstackTenant
*/
type ListOpenstackTenantsNoCredentialsOK struct {
	Payload []*models.OpenstackTenant
}

func (o *ListOpenstackTenantsNoCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/tenants][%d] listOpenstackTenantsNoCredentialsOK  %+v", 200, o.Payload)
}

func (o *ListOpenstackTenantsNoCredentialsOK) GetPayload() []*models.OpenstackTenant {
	return o.Payload
}

func (o *ListOpenstackTenantsNoCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListOpenstackTenantsNoCredentialsDefault creates a ListOpenstackTenantsNoCredentialsDefault with default headers values
func NewListOpenstackTenantsNoCredentialsDefault(code int) *ListOpenstackTenantsNoCredentialsDefault {
	return &ListOpenstackTenantsNoCredentialsDefault{
		_statusCode: code,
	}
}

/*ListOpenstackTenantsNoCredentialsDefault handles this case with default header values.

errorResponse
*/
type ListOpenstackTenantsNoCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list openstack tenants no credentials default response
func (o *ListOpenstackTenantsNoCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *ListOpenstackTenantsNoCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/openstack/tenants][%d] listOpenstackTenantsNoCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ListOpenstackTenantsNoCredentialsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListOpenstackTenantsNoCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
