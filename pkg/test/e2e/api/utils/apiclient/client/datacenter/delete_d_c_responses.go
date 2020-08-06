// Code generated by go-swagger; DO NOT EDIT.

package datacenter

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// DeleteDCReader is a Reader for the DeleteDC structure.
type DeleteDCReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDCReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteDCOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewDeleteDCUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewDeleteDCForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteDCDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteDCOK creates a DeleteDCOK with default headers values
func NewDeleteDCOK() *DeleteDCOK {
	return &DeleteDCOK{}
}

/*DeleteDCOK handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteDCOK struct {
}

func (o *DeleteDCOK) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/seed/{seed_name}/dc/{dc}][%d] deleteDCOK ", 200)
}

func (o *DeleteDCOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDCUnauthorized creates a DeleteDCUnauthorized with default headers values
func NewDeleteDCUnauthorized() *DeleteDCUnauthorized {
	return &DeleteDCUnauthorized{}
}

/*DeleteDCUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteDCUnauthorized struct {
}

func (o *DeleteDCUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/seed/{seed_name}/dc/{dc}][%d] deleteDCUnauthorized ", 401)
}

func (o *DeleteDCUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDCForbidden creates a DeleteDCForbidden with default headers values
func NewDeleteDCForbidden() *DeleteDCForbidden {
	return &DeleteDCForbidden{}
}

/*DeleteDCForbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type DeleteDCForbidden struct {
}

func (o *DeleteDCForbidden) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/seed/{seed_name}/dc/{dc}][%d] deleteDCForbidden ", 403)
}

func (o *DeleteDCForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDCDefault creates a DeleteDCDefault with default headers values
func NewDeleteDCDefault(code int) *DeleteDCDefault {
	return &DeleteDCDefault{
		_statusCode: code,
	}
}

/*DeleteDCDefault handles this case with default header values.

errorResponse
*/
type DeleteDCDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the delete d c default response
func (o *DeleteDCDefault) Code() int {
	return o._statusCode
}

func (o *DeleteDCDefault) Error() string {
	return fmt.Sprintf("[DELETE /api/v1/seed/{seed_name}/dc/{dc}][%d] deleteDC default  %+v", o._statusCode, o.Payload)
}

func (o *DeleteDCDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *DeleteDCDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
