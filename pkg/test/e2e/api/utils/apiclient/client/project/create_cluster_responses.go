// Code generated by go-swagger; DO NOT EDIT.

package project

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// CreateClusterReader is a Reader for the CreateCluster structure.
type CreateClusterReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CreateClusterReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewCreateClusterCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewCreateClusterUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewCreateClusterForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewCreateClusterDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCreateClusterCreated creates a CreateClusterCreated with default headers values
func NewCreateClusterCreated() *CreateClusterCreated {
	return &CreateClusterCreated{}
}

/*CreateClusterCreated handles this case with default header values.

Cluster
*/
type CreateClusterCreated struct {
	Payload *models.Cluster
}

func (o *CreateClusterCreated) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters][%d] createClusterCreated  %+v", 201, o.Payload)
}

func (o *CreateClusterCreated) GetPayload() *models.Cluster {
	return o.Payload
}

func (o *CreateClusterCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Cluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCreateClusterUnauthorized creates a CreateClusterUnauthorized with default headers values
func NewCreateClusterUnauthorized() *CreateClusterUnauthorized {
	return &CreateClusterUnauthorized{}
}

/*CreateClusterUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type CreateClusterUnauthorized struct {
}

func (o *CreateClusterUnauthorized) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters][%d] createClusterUnauthorized ", 401)
}

func (o *CreateClusterUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateClusterForbidden creates a CreateClusterForbidden with default headers values
func NewCreateClusterForbidden() *CreateClusterForbidden {
	return &CreateClusterForbidden{}
}

/*CreateClusterForbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type CreateClusterForbidden struct {
}

func (o *CreateClusterForbidden) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters][%d] createClusterForbidden ", 403)
}

func (o *CreateClusterForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCreateClusterDefault creates a CreateClusterDefault with default headers values
func NewCreateClusterDefault(code int) *CreateClusterDefault {
	return &CreateClusterDefault{
		_statusCode: code,
	}
}

/*CreateClusterDefault handles this case with default header values.

errorResponse
*/
type CreateClusterDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the create cluster default response
func (o *CreateClusterDefault) Code() int {
	return o._statusCode
}

func (o *CreateClusterDefault) Error() string {
	return fmt.Sprintf("[POST /api/v1/projects/{project_id}/dc/{dc}/clusters][%d] createCluster default  %+v", o._statusCode, o.Payload)
}

func (o *CreateClusterDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *CreateClusterDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
