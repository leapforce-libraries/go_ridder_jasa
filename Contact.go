package ridder_jasa

import (
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Contact struct {
	RidderID             int32  `json:"RidderId"`
	InsightlyID          int32  `json:"InsightlyId"`
	Person               Person `json:"Person"`
	Email                string `json:"Email"`
	Cellphone            string `json:"Cellphone"`
	Phone                string `json:"Phone"`
	Manual               bool   `json:"Manual"`
	MainContact          bool   `json:"MainContact"`
	MainContactCreditor  bool   `json:"MainContactCreditor"`
	MainContactDebtor    bool   `json:"MainContactDebtor"`
	FunctionName         string `json:"FunctionName"`
	EmploymentTerminated bool   `json:"EmploymentTerminated"`
	OrganizationID       int32  `json:"OrganizationId"`
}

func (service *Service) GetContact(ridderID int32) (*Contact, *errortools.Error) {
	contact := Contact{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		URL:           service.url(fmt.Sprintf("contacts?ridderid=%v", ridderID)),
		ResponseModel: &contact,
	}
	_, _, e := service.httpRequest(&requestConfig)

	return &contact, e
}

func (service *Service) UpdateContact(contact *Contact) (*int32, *errortools.Error) {
	if contact == nil {
		return nil, nil
	}

	ev := service.validateContact(contact)

	contactID := new(int32)

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		URL:           service.url(fmt.Sprintf("contacts/%v", contact.RidderID)),
		BodyModel:     contact,
		ResponseModel: contactID,
	}
	req, res, e := service.httpRequest(&requestConfig)

	if ev != nil {
		ev.SetRequest(req)
		ev.SetResponse(res)
		errortools.CaptureWarning(ev)
	}

	return contactID, e
}

func (service *Service) CreateContact(newContact *Contact) (*int32, *errortools.Error) {
	if newContact == nil {
		return nil, nil
	}

	ev := service.validateContact(newContact)

	contactID := new(int32)

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		URL:           service.url("contacts"),
		BodyModel:     newContact,
		ResponseModel: contactID,
	}
	req, res, e := service.httpRequest(&requestConfig)

	if ev != nil {
		ev.SetRequest(req)
		ev.SetResponse(res)
		errortools.CaptureWarning(ev)
	}

	return contactID, e
}

func (service *Service) validateContact(contact *Contact) *errortools.Error {
	if contact == nil {
		return nil
	}

	errors := []string{}

	service.truncateString("Phone", &(*contact).Phone, maxLengthContactPhone, &errors)
	service.truncateString("Email", &(*contact).Email, maxLengthContactEmail, &errors)
	service.truncateString("FunctionName", &(*contact).FunctionName, maxLengthContactFunctionName, &errors)
	service.truncateString("Cellphone", &(*contact).Cellphone, maxLengthContactCellphone, &errors)
	service.truncateString("LastName", &(*contact).Person.LastName, maxLengthContactLastName, &errors)
	service.truncateString("Initials", &(*contact).Person.Initials, maxLengthContactInitials, &errors)
	service.truncateString("FirstName", &(*contact).Person.FirstName, maxLengthContactFirstName, &errors)

	if len(errors) > 0 {
		return errortools.ErrorMessage(strings.Join(errors, "\n"))
	}

	return nil
}
