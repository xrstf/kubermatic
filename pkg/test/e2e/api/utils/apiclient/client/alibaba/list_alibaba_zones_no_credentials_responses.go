// Code generated by go-swagger; DO NOT EDIT.

package alibaba

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"k8c.io/kubermatic/v2/pkg/test/e2e/api/utils/apiclient/models"
)

// ListAlibabaZonesNoCredentialsReader is a Reader for the ListAlibabaZonesNoCredentials structure.
type ListAlibabaZonesNoCredentialsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *ListAlibabaZonesNoCredentialsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewListAlibabaZonesNoCredentialsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewListAlibabaZonesNoCredentialsDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewListAlibabaZonesNoCredentialsOK creates a ListAlibabaZonesNoCredentialsOK with default headers values
func NewListAlibabaZonesNoCredentialsOK() *ListAlibabaZonesNoCredentialsOK {
	return &ListAlibabaZonesNoCredentialsOK{}
}

/*ListAlibabaZonesNoCredentialsOK handles this case with default header values.

AlibabaZoneList
*/
type ListAlibabaZonesNoCredentialsOK struct {
	Payload models.AlibabaZoneList
}

func (o *ListAlibabaZonesNoCredentialsOK) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/alibaba/zones][%d] listAlibabaZonesNoCredentialsOK  %+v", 200, o.Payload)
}

func (o *ListAlibabaZonesNoCredentialsOK) GetPayload() models.AlibabaZoneList {
	return o.Payload
}

func (o *ListAlibabaZonesNoCredentialsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewListAlibabaZonesNoCredentialsDefault creates a ListAlibabaZonesNoCredentialsDefault with default headers values
func NewListAlibabaZonesNoCredentialsDefault(code int) *ListAlibabaZonesNoCredentialsDefault {
	return &ListAlibabaZonesNoCredentialsDefault{
		_statusCode: code,
	}
}

/*ListAlibabaZonesNoCredentialsDefault handles this case with default header values.

errorResponse
*/
type ListAlibabaZonesNoCredentialsDefault struct {
	_statusCode int

	Payload *models.ErrorResponse
}

// Code gets the status code for the list alibaba zones no credentials default response
func (o *ListAlibabaZonesNoCredentialsDefault) Code() int {
	return o._statusCode
}

func (o *ListAlibabaZonesNoCredentialsDefault) Error() string {
	return fmt.Sprintf("[GET /api/v1/projects/{project_id}/dc/{dc}/clusters/{cluster_id}/providers/alibaba/zones][%d] listAlibabaZonesNoCredentials default  %+v", o._statusCode, o.Payload)
}

func (o *ListAlibabaZonesNoCredentialsDefault) GetPayload() *models.ErrorResponse {
	return o.Payload
}

func (o *ListAlibabaZonesNoCredentialsDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ErrorResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
