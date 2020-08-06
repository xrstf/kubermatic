// Code generated by go-swagger; DO NOT EDIT.

package tokens

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// UpdateServiceAccountTokenReader is a Reader for the UpdateServiceAccountToken structure.
type UpdateServiceAccountTokenReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateServiceAccountTokenReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateServiceAccountTokenOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUpdateServiceAccountTokenUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUpdateServiceAccountTokenForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateServiceAccountTokenDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateServiceAccountTokenOK creates a UpdateServiceAccountTokenOK with default headers values
func NewUpdateServiceAccountTokenOK() *UpdateServiceAccountTokenOK {
	return &UpdateServiceAccountTokenOK{}
}

/*UpdateServiceAccountTokenOK handles this case with default header values.

ServiceAccountToken
*/
type UpdateServiceAccountTokenOK struct {
	Payload *models.ServiceAccountToken
}

func (o *UpdateServiceAccountTokenOK) Error() string {
	return fmt.Sprintf("[PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}][%d] updateServiceAccountTokenOK  %+v", 200, o.Payload)
}

func (o *UpdateServiceAccountTokenOK) GetPayload() *models.ServiceAccountToken {
	return o.Payload
}

func (o *UpdateServiceAccountTokenOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ServiceAccountToken)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateServiceAccountTokenUnauthorized creates a UpdateServiceAccountTokenUnauthorized with default headers values
func NewUpdateServiceAccountTokenUnauthorized() *UpdateServiceAccountTokenUnauthorized {
	return &UpdateServiceAccountTokenUnauthorized{}
}

/*UpdateServiceAccountTokenUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type UpdateServiceAccountTokenUnauthorized struct {
}

func (o *UpdateServiceAccountTokenUnauthorized) Error() string {
	return fmt.Sprintf("[PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}][%d] updateServiceAccountTokenUnauthorized ", 401)
}

func (o *UpdateServiceAccountTokenUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateServiceAccountTokenForbidden creates a UpdateServiceAccountTokenForbidden with default headers values
func NewUpdateServiceAccountTokenForbidden() *UpdateServiceAccountTokenForbidden {
	return &UpdateServiceAccountTokenForbidden{}
}

/*UpdateServiceAccountTokenForbidden handles this case with default header values.

EmptyResponse is a empty response
*/
type UpdateServiceAccountTokenForbidden struct {
}

func (o *UpdateServiceAccountTokenForbidden) Error() string {
	return fmt.Sprintf("[PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}][%d] updateServiceAccountTokenForbidden ", 403)
}

func (o *UpdateServiceAccountTokenForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateServiceAccountTokenDefault creates a UpdateServiceAccountTokenDefault with default headers values
func NewUpdateServiceAccountTokenDefault(code int) *UpdateServiceAccountTokenDefault {
	return &UpdateServiceAccountTokenDefault{
		_statusCode: code,
	}
}

/*UpdateServiceAccountTokenDefault handles this case with default header values.

errorResponse
*/
type UpdateServiceAccountTokenDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the update service account token default response
func (o *UpdateServiceAccountTokenDefault) Code() int {
	return o._statusCode
}

func (o *UpdateServiceAccountTokenDefault) Error() string {
	return fmt.Sprintf("[PUT /api/v1/projects/{project_id}/serviceaccounts/{serviceaccount_id}/tokens/{token_id}][%d] updateServiceAccountToken default  %+v", o._statusCode, o.Payload)
}

func (o *UpdateServiceAccountTokenDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *UpdateServiceAccountTokenDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
