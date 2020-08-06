// Code generated by go-swagger; DO NOT EDIT.

package settings

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// GetCurrentUserSettingsReader is a Reader for the GetCurrentUserSettings structure.
type GetCurrentUserSettingsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetCurrentUserSettingsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetCurrentUserSettingsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetCurrentUserSettingsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewGetCurrentUserSettingsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetCurrentUserSettingsOK creates a GetCurrentUserSettingsOK with default headers values
func NewGetCurrentUserSettingsOK() *GetCurrentUserSettingsOK {
	return &GetCurrentUserSettingsOK{}
}

/*GetCurrentUserSettingsOK handles this case with default header values.

UserSettings
*/
type GetCurrentUserSettingsOK struct {
	Payload *models.UserSettings
}

func (o *GetCurrentUserSettingsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/me/settings][%d] getCurrentUserSettingsOK  %+v", 200, o.Payload)
}

func (o *GetCurrentUserSettingsOK) GetPayload() *models.UserSettings {
	return o.Payload
}

func (o *GetCurrentUserSettingsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.UserSettings)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetCurrentUserSettingsUnauthorized creates a GetCurrentUserSettingsUnauthorized with default headers values
func NewGetCurrentUserSettingsUnauthorized() *GetCurrentUserSettingsUnauthorized {
	return &GetCurrentUserSettingsUnauthorized{}
}

/*GetCurrentUserSettingsUnauthorized handles this case with default header values.

EmptyResponse is a empty response
*/
type GetCurrentUserSettingsUnauthorized struct {
}

func (o *GetCurrentUserSettingsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /api/v1/me/settings][%d] getCurrentUserSettingsUnauthorized ", 401)
}

func (o *GetCurrentUserSettingsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetCurrentUserSettingsDefault creates a GetCurrentUserSettingsDefault with default headers values
func NewGetCurrentUserSettingsDefault(code int) *GetCurrentUserSettingsDefault {
	return &GetCurrentUserSettingsDefault{
		_statusCode: code,
	}
}

/*GetCurrentUserSettingsDefault handles this case with default header values.

errorResponse
*/
type GetCurrentUserSettingsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the get current user settings default response
func (o *GetCurrentUserSettingsDefault) Code() int {
	return o._statusCode
}

func (o *GetCurrentUserSettingsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/me/settings][%d] getCurrentUserSettings default  %+v", o._statusCode, o.Payload)
}

func (o *GetCurrentUserSettingsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *GetCurrentUserSettingsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
