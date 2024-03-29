package ridder_jasa

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	errortools "github.com/leapforce-libraries/go_errortools"
	go_http "github.com/leapforce-libraries/go_http"
	utilities "github.com/leapforce-libraries/go_utilities"
)

const (
	apiName                            string = "Ridder"
	maxLengthOrganizationEmail         int    = 255
	maxLengthOrganizationName          int    = 60
	maxLengthOrganizationPhone         int    = 50
	maxLengthOrganizationWebsite       int    = 255
	maxLengthAddressHouseNumber        int    = 50
	maxLengthAddressCity               int    = 50
	maxLengthAddressZipCode            int    = 50
	maxLengthAddressStreet             int    = 50
	maxLengthContactPhone              int    = 50
	maxLengthContactEmail              int    = 255
	maxLengthContactFunctionName       int    = 150
	maxLengthContactCellphone          int    = 50
	maxLengthContactLastName           int    = 127
	maxLengthContactInitials           int    = 50
	maxLengthContactFirstName          int    = 127
	maxLengthOpportunityName           int    = 80
	maxLengthOpportunityInsightlyState int    = 4000
	DateTimeLayout                     string = "2006-01-02T15:04:05"
)

type Service struct {
	apiUrl      string
	apiKey      string
	httpService *go_http.Service
}

type ServiceConfig struct {
	ApiUrl                string
	ApiKey                string
	MaxRetries            *uint
	SecondsBetweenRetries *uint32
}

func NewService(config ServiceConfig) (*Service, *errortools.Error) {
	if config.ApiUrl == "" {
		return nil, errortools.ErrorMessage("Service Api Url not provided")
	}

	if config.ApiKey == "" {
		return nil, errortools.ErrorMessage("Service Api Key not provided")
	}

	httpService, e := go_http.NewService(&go_http.ServiceConfig{})
	if e != nil {
		return nil, e
	}

	return &Service{
		apiUrl:      strings.TrimRight(config.ApiUrl, "/"),
		apiKey:      config.ApiKey,
		httpService: httpService,
	}, nil
}

func (service *Service) httpRequest(requestConfig *go_http.RequestConfig) (*http.Request, *http.Response, *errortools.Error) {
	// add api key header
	header := http.Header{}
	header.Set("X-ApiKey", service.apiKey)

	if !utilities.IsNil(requestConfig.BodyModel) {
		header.Set("Content-Type", "application/json-patch+json")
	}
	(*requestConfig).NonDefaultHeaders = &header

	// add error model
	errorResponse := ErrorResponse{}
	(*requestConfig).ErrorModel = &errorResponse

	request, response, e := service.httpService.HttpRequest(requestConfig)
	if errorResponse.Error != "" {
		e.SetMessage(errorResponse.Error)
	}

	return request, response, e
}

func (service *Service) url(path string) string {
	return fmt.Sprintf("%s/%s", service.apiUrl, path)
}

func (service *Service) truncateString(fieldName string, value *string, maxLength int, errors *[]string) {
	if len(*value) > maxLength {
		*value = (*value)[:maxLength]

		*errors = append(*errors, fmt.Sprintf("%s truncated to %v characters.", fieldName, maxLength))
	}
}

func (service *Service) removeSpecialCharacters(test *string) *errortools.Error {
	if test == nil {
		return nil
	}

	re := regexp.MustCompile(`[\\/:*?"<>|]`)

	removedCount := len(*test) - len(string(re.ReplaceAll([]byte(*test), []byte(""))))

	if removedCount == 0 {
		return nil
	}

	message := fmt.Sprintf("%v special characters in '%s' replaced by a dot", removedCount, *test)
	(*test) = string(re.ReplaceAll([]byte(*test), []byte(".")))

	return errortools.ErrorMessage(message)
}

func (service *Service) ApiName() string {
	return apiName
}

func (service *Service) ApiKey() string {
	return service.apiKey
}

func (service *Service) ApiCallCount() int64 {
	return service.httpService.RequestCount()
}

func (service *Service) ApiReset() {
	service.httpService.ResetRequestCount()
}
