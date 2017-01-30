package draft

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"
)

// BatchUserCreateReader is a Reader for the BatchUserCreate structure.
type BatchUserCreateReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *BatchUserCreateReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewBatchUserCreateOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 409:
		result := NewBatchUserCreateConflict()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewBatchUserCreateOK creates a BatchUserCreateOK with default headers values
func NewBatchUserCreateOK() *BatchUserCreateOK {
	return &BatchUserCreateOK{}
}

/*BatchUserCreateOK handles this case with default header values.

Пользователи созданы успешно.

*/
type BatchUserCreateOK struct {
}

func (o *BatchUserCreateOK) Error() string {
	return fmt.Sprintf("[POST /batch/user/create][%d] batchUserCreateOK ", 200)
}

func (o *BatchUserCreateOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewBatchUserCreateConflict creates a BatchUserCreateConflict with default headers values
func NewBatchUserCreateConflict() *BatchUserCreateConflict {
	return &BatchUserCreateConflict{}
}

/*BatchUserCreateConflict handles this case with default header values.

Хотя бы один пользователь уже присутсвует в базе данных.

*/
type BatchUserCreateConflict struct {
}

func (o *BatchUserCreateConflict) Error() string {
	return fmt.Sprintf("[POST /batch/user/create][%d] batchUserCreateConflict ", 409)
}

func (o *BatchUserCreateConflict) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}