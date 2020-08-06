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

// GetClusterV2Reader is a Reader for the GetClusterV2 structure.
type GetClusterV2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetClusterV2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetClusterV2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetClusterV2Unauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewGetClusterV2Forbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetClusterV2Default(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetClusterV2OK creates a GetClusterV2OK with default headers values
func NewGetClusterV2OK() *GetClusterV2OK {
	return &GetClusterV2OK{}
}

/*GetClusterV2OK handles this case with default header values.

Cluster
*/
type GetClusterV2OK struct {
	Payload *models.Cluster
}

func (o *GetClusterV2OK) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/clusters/{cluster_id}][%d] getClusterV2OK  %+v", 200, o.Payload)
}

func (o *GetClusterV2OK) GetPayload() *models.Cluster {
	return o.Payload
}

func (o *GetClusterV2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Cluster)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetClusterV2Unauthorized creates a GetClusterV2Unauthorized with default headers values
func NewGetClusterV2Unauthorized() *GetClusterV2Unauthorized {
	return &GetClusterV2Unauthorized{}
}

/*GetClusterV2Unauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type GetClusterV2Unauthorized struct {
}

func (o *GetClusterV2Unauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/clusters/{cluster_id}][%d] getClusterV2Unauthorized ", 401)
}

func (o *GetClusterV2Unauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetClusterV2Forbidden creates a GetClusterV2Forbidden with default headers values
func NewGetClusterV2Forbidden() *GetClusterV2Forbidden {
	return &GetClusterV2Forbidden{}
}

/*GetClusterV2Forbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type GetClusterV2Forbidden struct {
}

func (o *GetClusterV2Forbidden) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/clusters/{cluster_id}][%d] getClusterV2Forbidden ", 403)
}

func (o *GetClusterV2Forbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetClusterV2Default creates a GetClusterV2Default with default headers values
func NewGetClusterV2Default(code int) *GetClusterV2Default {
	return &GetClusterV2Default{
		_statusCode: code,
	}
}

/*GetClusterV2Default handles this case with default header values.

errorResponse
*/
type GetClusterV2Default struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get cluster v2 default response
func (o *GetClusterV2Default) Code() int {
	return o._statusCode
}

func (o *GetClusterV2Default) Error() string {
	return fmt.Sprintf("[GET /api/v2/projects/{project_id}/clusters/{cluster_id}][%d] getClusterV2 default  %+v", o._statusCode, o.Payload)
}

func (o *GetClusterV2Default) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetClusterV2Default) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
