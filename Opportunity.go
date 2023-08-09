package ridder_jasa

import (
	"fmt"
	"net/http"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
)

type Workflow string

const (
	WorkflowNone                  Workflow = "None"
	WorkflowReject                Workflow = "Reject"
	WorkflowRejectAndMakeHistoric Workflow = "RejectAndMakeHistoric"
	WorkflowMakeHistoric          Workflow = "MakeHistoric"
	WorkflowReOpen                Workflow = "ReOpen"
)

type Opportunity struct {
	RidderID             int32   `json:"RidderId"`
	InsightlyID          int32   `json:"InsightlyId"`
	InsightlyState       string  `json:"InsightlyState"`
	OfferNumber          int32   `json:"OfferNumber"`
	OpportunityName      string  `json:"OpportunityName"`
	OrganizationID       int32   `json:"OrganizationId"`
	ContactID            int32   `json:"ContactId"`
	Currency             string  `json:"Currency"`
	OpportunityCreated   *string `json:"OpportunityCreated,omitempty"`
	ForecastCloseDate    *string `json:"ForecastCloseDate"`
	ProbabilityOfWinning int32   `json:"ProbabilityOfWinning"`
	SalesPerson          *int32  `json:"SalesPerson,omitempty"`
	ExternalKey          string  `json:"ExternalKey"`
	Revision             int32   `json:"Revision"`
}

type OpportunityResponse struct {
	RidderID    int32 `json:"RidderId"`
	OfferNumber int32 `json:"OfferNumber"`
	Revision    int32 `json:"Revision"`
}

func (service *Service) GetOpportunity(ridderID int32) (*Opportunity, *errortools.Error) {
	opportunity := Opportunity{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodGet,
		Url:           service.url(fmt.Sprintf("opportunities?ridderid=%v", ridderID)),
		ResponseModel: &opportunity,
	}
	_, _, e := service.httpRequest(&requestConfig)

	return &opportunity, e
}

func (service *Service) UpdateOpportunity(opportunity *Opportunity) (*OpportunityResponse, *errortools.Error) {
	if opportunity == nil {
		return nil, nil
	}

	ev := service.validateOpportunity(opportunity)

	opportunityResponse := OpportunityResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url(fmt.Sprintf("opportunities/%v", opportunity.RidderID)),
		BodyModel:     opportunity,
		ResponseModel: &opportunityResponse,
	}
	req, res, e := service.httpRequest(&requestConfig)

	if ev != nil {
		ev.SetRequest(req)
		ev.SetResponse(res)
		errortools.CaptureWarning(ev)
	}

	return &opportunityResponse, e
}

func (service *Service) CreateOpportunity(newOpportunity *Opportunity) (*OpportunityResponse, *errortools.Error) {
	if newOpportunity == nil {
		return nil, nil
	}

	ev := service.validateOpportunity(newOpportunity)

	opportunityResponse := OpportunityResponse{}

	requestConfig := go_http.RequestConfig{
		Method:        http.MethodPost,
		Url:           service.url("opportunities"),
		BodyModel:     newOpportunity,
		ResponseModel: &opportunityResponse,
	}
	req, res, e := service.httpRequest(&requestConfig)

	if ev != nil {
		ev.SetRequest(req)
		ev.SetResponse(res)
		errortools.CaptureWarning(ev)
	}

	return &opportunityResponse, e
}

func (service *Service) WorkflowOpportunity(opportunity *Opportunity, workflow Workflow) *errortools.Error {
	if opportunity == nil {
		return nil
	}

	if workflow == WorkflowNone {
		return nil
	}

	requestConfig := go_http.RequestConfig{
		Method:    http.MethodPost,
		Url:       service.url(fmt.Sprintf("opportunities/%v/%s", opportunity.RidderID, workflow)),
		BodyModel: opportunity,
	}
	_, _, e := service.httpRequest(&requestConfig)

	return e
}

func (service *Service) validateOpportunity(opportunity *Opportunity) *errortools.Error {
	if opportunity == nil {
		return nil
	}

	errors := []string{}

	service.truncateString("OpportunityName", &(*opportunity).OpportunityName, maxLengthOpportunityName, &errors)
	service.truncateString("InsightlyState", &(*opportunity).InsightlyState, maxLengthOpportunityInsightlyState, &errors)

	e := service.removeSpecialCharacters(&(*opportunity).OpportunityName)
	if e != nil {
		errors = append(errors, e.Message())
	}

	if len(errors) > 0 {
		return errortools.ErrorMessage(strings.Join(errors, "\n"))
	}

	return nil
}
