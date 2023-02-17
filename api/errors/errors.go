package errors

import (
	"net/http"
)

type ErrorCode string

const (
	BadRequest                      = "ERR_API_BAD_REQUEST"
	InvalidTicket                   = "ERR_API_INVALID_TICKET"
	UnableToTrack                   = "ERR_API_UNABLE_TO_TRACK"
	NoRecordsFound                  = "ERR_API_NO_RECORDS_FOUND"
	NoCompanyRecordsFoundByName     = "ERR_API_NO_COMPANY_RECORDS_FOUNDFOR_GIVEN_NAME"
	NoCompanyRecordsFoundByID       = "ERR_API_NO_COMPANY_RECORDS_FOUND_FOR_GIVEN_ID"
	RecordAlreadyExistsForGivenName = "ERR_API_RECORD_ALREADY_EXISTE_FOR_GIVEN_NAME"
	UnableToCreateCompany           = "ERR_API_UNABLE_TO_CREATE_COMPANY"
	InternalServerError             = "ERR_API_SERVER_ERROR"
	UnableToFetchCompany            = "ERR_API_UNABLE_TO_FETCH_COMPANY"
	UnableToDeleteCompany           = "ERR_API_UNABLE_TO_DELETE_COMPANY"
	UnableToUpdateCompany           = "ERR_API_UNABLE_TO_UPDATE_COMPANY"
)

var ApiErrors = map[ErrorCode]string{
	BadRequest:                      "Invalid request body",
	NoRecordsFound:                  "No records found",
	NoCompanyRecordsFoundByName:     "No records found for given name",
	NoCompanyRecordsFoundByID:       "No records found for given ID",
	RecordAlreadyExistsForGivenName: "Record already exist for given name",
	UnableToCreateCompany:           "Unable to create company",
	InternalServerError:             "Internal server error",
	UnableToFetchCompany:            "Unable to fetch company",
	UnableToDeleteCompany:           "Unable to delete company",
	UnableToUpdateCompany:           "Unable to update company",
}

type ErrorResponse struct {
	HttpStatusCode int       `json:"status"`
	ErrorCode      ErrorCode `json:"error_code,omitempty"`
	ErrorMessage   string    `json:"error_message,omitempty"`
}

// Create new error responses
func NewErrorResponse(httpStatusCode int, errorCode ErrorCode, errorMessage string) *ErrorResponse {
	return &ErrorResponse{
		HttpStatusCode: httpStatusCode,
		ErrorCode:      errorCode,
		ErrorMessage:   errorMessage,
	}
}

func (e ErrorResponse) Error() string {
	return e.ErrorMessage
}

var ErrBadRequest = NewErrorResponse(http.StatusBadRequest, BadRequest, ApiErrors[BadRequest])
var ErrNoRecordsFound = NewErrorResponse(http.StatusBadRequest, NoRecordsFound, ApiErrors[NoRecordsFound])
var ErrNoCompanyRecordsFoundByName = NewErrorResponse(http.StatusBadRequest, NoCompanyRecordsFoundByName, ApiErrors[NoCompanyRecordsFoundByName])
var ErrNoCompanyRecordsFoundByID = NewErrorResponse(http.StatusBadRequest, NoCompanyRecordsFoundByID, ApiErrors[NoCompanyRecordsFoundByID])
var ErrRecordAlreadyExistsForGivenName = NewErrorResponse(http.StatusBadRequest, RecordAlreadyExistsForGivenName, ApiErrors[RecordAlreadyExistsForGivenName])
var ErrUnableToCreateCompany = NewErrorResponse(http.StatusInternalServerError, UnableToCreateCompany, ApiErrors[UnableToCreateCompany])
var ErrInternalServerError = NewErrorResponse(http.StatusInternalServerError, InternalServerError, ApiErrors[InternalServerError])
var ErrUnableToFetchCompany = NewErrorResponse(http.StatusInternalServerError, UnableToFetchCompany, ApiErrors[UnableToFetchCompany])
var ErrUnableToDeleteCompany = NewErrorResponse(http.StatusInternalServerError, UnableToDeleteCompany, ApiErrors[UnableToDeleteCompany])
var ErrUnableToUpdateCompany = NewErrorResponse(http.StatusInternalServerError, UnableToUpdateCompany, ApiErrors[UnableToUpdateCompany])
