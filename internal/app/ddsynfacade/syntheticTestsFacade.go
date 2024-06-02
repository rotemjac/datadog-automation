package ddsynfacade

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/rs/zerolog"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV1"
	"github.com/rotemjac/datadog-automation/internal/pkg/bootstrap"
	"github.com/rotemjac/datadog-automation/internal/pkg/model"
)

type DatadogFacade interface {
	GetTests() any
	UpdateTests() any
	DeleteTests() any

	GetApiTest() any

	CheckIfSyntheticTestExistsByName() any
	CheckIfSyntheticTestExistsByNameAndOneTag() any
	CheckIfSyntheticTestExistsByNameAndTags() any

	CreateHttpSyntheticTest() any
	UpdateHttpSyntheticTest() any

	CreateGrpcSyntheticTest() any
	UpdateGrpcSyntheticTest() any
}

// ApiClient implements the datadogFacade interface
type MySyntheticsApi struct {
	client *datadogV1.SyntheticsApi
}

var mySyntheticsApi MySyntheticsApi
var logger zerolog.Logger
var logLevel string
var configFilePath string
var configFileName string

func init() {
	// Fetch Datadog API client
	apiClient := bootstrap.Client()

	// Create a new instance of the SyntheticApi
	mySyntheticsApi = MySyntheticsApi{
		client: datadogV1.NewSyntheticsApi(apiClient),
	}

	// #############  Get ENV vars ############# //
	logLevel = os.Getenv("LOG_LEVEL")
	configFilePath = os.Getenv("CONFIG_FILE_PATH")
	configFileName = os.Getenv("CONFIG_FILE_NAME")

	// ############# Logging Setup ############# //
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger = zerolog.New(os.Stderr).With().Logger()
	logger.Info().Msg("Initialize synthetic tests facade..")
}

func (apiClient MySyntheticsApi) GetTests(syntheticTestsConfig *model.SyntheticTestsConfig, tags []string) ([]datadogV1.SyntheticsTestDetails, error) {
	logger.Debug().Msg("In syntheticTests.GetTests")
	// We need to fetch ALL tests because there is no efficient filtering function AFAIK
	allSyntheticTests, err := getAllSyntheticTests()
	if err != nil {
		return nil, err
	} else {
		testsByTag, err := getSyntheticTestsByTags(allSyntheticTests, tags)
		return testsByTag, err
	}
}

func (apiClient MySyntheticsApi) UpdateTests(syntheticTestsConfig *model.SyntheticTestsConfig) error {
	//Do something
	return nil
}

func (apiClient MySyntheticsApi) DeleteTests(syntheticTestsConfig *model.SyntheticTestsConfig, tags []string) error {
	var testsCount = 0
	var failureCount = 0
	var totalErrors = errors.New("")

	logger.Debug().Msg("In syntheticTests.GetTests")
	// We need to fetch ALL tests because there is no efficient filtering function AFAIK
	allSyntheticTests, err := getAllSyntheticTests()
	if err != nil {
		return err
	} else {
		testsToDelete, err := getSyntheticTestsByTags(allSyntheticTests, tags)
		if err != nil {
			return err
		} else {
			for _, testToDelete := range testsToDelete {
				testsCount = testsCount + 1
				err := deleteSyntheticTest(*testToDelete.PublicId)
				if err != nil {
					logger.Error().Msgf("Error when deleting test: %v\n", *testToDelete.PublicId)
					totalErrors = err //errors.Join(totalErrors, err) Failing with Go 1.19
					failureCount = failureCount + 1
				}
			}
			logger.Debug().Msgf("testsCount: %v", testsCount)
			logger.Debug().Msgf("failureCount: %v", failureCount)
			if (testsCount - failureCount) > 1 {
				logger.Info().Msgf("The total number of tests that were failed to delete are: %v", failureCount)
				return nil
			} else {
				logger.Error().Msgf("The total number of tests tat were failed to delete are: %v", failureCount)
				return totalErrors
			}
		}
	}
}

/* ################################################################# */
/* ####################### Local functions ######################### */
/* ################################################################# */
func (apiClient MySyntheticsApi) CreateHttpSyntheticTest(httpSyntheticTest model.SyntheticTestHttpOptions, tags []string, headers map[string]string) {
	// Create an API GRPC test returns "OK - Returns the created test details." response

	body := getHttpSyntheticTestBody(httpSyntheticTest, tags, headers)
	ctx := datadog.NewDefaultContext(context.Background())
	resp, r, err := mySyntheticsApi.client.CreateSyntheticsAPITest(ctx, body)

	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.CreateSyntheticsAPITest`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
	}
	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	logger.Info().Msgf("Response from `SyntheticsApi.CreateSyntheticsAPITest`:\n%s\n", responseContent)

}
func (apiClient MySyntheticsApi) UpdateHttpSyntheticTest(httpSyntheticTest model.SyntheticTestHttpOptions, pubicId string, tags []string, headers map[string]string) {
	body := getHttpSyntheticTestBody(httpSyntheticTest, tags, headers)
	ctx := datadog.NewDefaultContext(context.Background())
	resp, r, err := mySyntheticsApi.client.UpdateAPITest(ctx, pubicId, body)

	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.UpdateAPITest`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
	}
	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	logger.Info().Msgf("Response from `SyntheticsApi.UpdateAPITest`:\n%s\n", responseContent)
}

// In case body is empty , need to return nil, when we return empty string it cause synthetic tests to fail on http GET
func getHttpBodyType(requestBody string) *datadogV1.SyntheticsTestRequestBodyType {

	if requestBody != "" {
		return datadogV1.SYNTHETICSTESTREQUESTBODYTYPE_APPLICATION_JSON.Ptr()
	}

	return nil
}

// In case body is empty , need to return nil, when we return empty string it cause synthetic tests to fail on http GET
func getHttpBody(requestBody string) *string {

	if requestBody != "" {
		return datadog.PtrString(requestBody)
	}

	return nil
}

func getHttpSyntheticTestBody(httpSyntheticTest model.SyntheticTestHttpOptions, tags []string, headers map[string]string) datadogV1.SyntheticsAPITest {
	body := datadogV1.SyntheticsAPITest{
		Config: datadogV1.SyntheticsAPITestConfig{
			Assertions: getHttpAssertions(httpSyntheticTest),
			Request: &datadogV1.SyntheticsTestRequest{
				Headers: headers,
				Method:  datadog.PtrString(httpSyntheticTest.Method),   //method
				Timeout: datadog.PtrFloat64(httpSyntheticTest.Timeout), //timeout
				Url:     datadog.PtrString(httpSyntheticTest.Url),
				//BodyType: datadogV1.SYNTHETICSTESTREQUESTBODYTYPE_APPLICATION_JSON.Ptr(),
				BodyType: getHttpBodyType(httpSyntheticTest.RequestBody),
				//Body:     datadog.PtrString(httpSyntheticTest.RequestBody),
				Body: getHttpBody(httpSyntheticTest.RequestBody),
			},
		},
		Locations: []string{
			fmt.Sprintf("aws:%v", httpSyntheticTest.Region),
		},
		Message: fmt.Sprintf("Monitor triggered. Notify:  @%v ", httpSyntheticTest.SlackChannel),
		Name:    httpSyntheticTest.Name,
		Options: datadogV1.SyntheticsTestOptions{
			AcceptSelfSigned:   datadog.PtrBool(false),
			AllowInsecure:      datadog.PtrBool(true),
			FollowRedirects:    datadog.PtrBool(false),
			MinFailureDuration: datadog.PtrInt64(10),
			MinLocationFailed:  datadog.PtrInt64(1),
			MonitorName:        datadog.PtrString(""), //Monitor Name
			//MonitorPriority:    datadog.PtrInt32(5),
			MonitorOptions: &datadogV1.SyntheticsTestOptionsMonitorOptions{
				RenotifyInterval: datadog.PtrInt64(httpSyntheticTest.RenotifyInterval),
			},
			Retry: &datadogV1.SyntheticsTestOptionsRetry{
				Count:    datadog.PtrInt64(2),
				Interval: datadog.PtrFloat64(300),
			},
			TickEvery:   datadog.PtrInt64(httpSyntheticTest.TickEverySec),
			HttpVersion: datadogV1.SYNTHETICSTESTOPTIONSHTTPVERSION_ANY.Ptr(),
		},
		Subtype: datadogV1.SYNTHETICSTESTDETAILSSUBTYPE_HTTP.Ptr(),
		// Tags for the test
		Tags: tags,
		Type: datadogV1.SYNTHETICSAPITESTTYPE_API,
	}
	return body
}

func getHttpAssertions(httpSyntheticTest model.SyntheticTestHttpOptions) []datadogV1.SyntheticsAssertion {
	assertions := []datadogV1.SyntheticsAssertion{
		datadogV1.SyntheticsAssertion{
			SyntheticsAssertionTarget: &datadogV1.SyntheticsAssertionTarget{
				Operator: datadogV1.SYNTHETICSASSERTIONOPERATOR_IS,
				Type:     datadogV1.SYNTHETICSASSERTIONTYPE_STATUS_CODE,
				Target:   200,
			}},
	}

	if httpSyntheticTest.AssertionJsonPath != "" {
		assert := datadogV1.SyntheticsAssertion{
			SyntheticsAssertionJSONPathTarget: &datadogV1.SyntheticsAssertionJSONPathTarget{
				Operator: datadogV1.SYNTHETICSASSERTIONJSONPATHOPERATOR_VALIDATES_JSON_PATH,
				Target: &datadogV1.SyntheticsAssertionJSONPathTargetTarget{
					JsonPath:    datadog.PtrString(httpSyntheticTest.AssertionJsonPath),
					Operator:    datadog.PtrString(httpSyntheticTest.AssertionJsonOperator),
					TargetValue: httpSyntheticTest.AssertionTargetValue,
				},
				Type: datadogV1.SYNTHETICSASSERTIONTYPE_BODY,
			}}
		assertions = append(assertions, assert)
	}

	return assertions
}

func (apiClient MySyntheticsApi) CreateGrpcSyntheticTest(grpcSyntheticTest model.SyntheticTestGrpcOptions, tags []string, headers map[string]string) {
	// Create an API GRPC test returns "OK - Returns the created test details." response
	body := getGrpcSyntheticTestBody(grpcSyntheticTest, tags, headers)
	ctx := datadog.NewDefaultContext(context.Background())
	resp, r, err := mySyntheticsApi.client.CreateSyntheticsAPITest(ctx, body)

	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.CreateSyntheticsAPITest`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
	}
	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	logger.Info().Msgf("Response from `SyntheticsApi.CreateSyntheticsAPITest`:\n%s\n", responseContent)
}

func (apiClient MySyntheticsApi) UpdateGrpcSyntheticTest(grpcSyntheticTest model.SyntheticTestGrpcOptions, publicId string, tags []string, headers map[string]string) {
	// Create an API GRPC test returns "OK - Returns the created test details." response
	body := getGrpcSyntheticTestBody(grpcSyntheticTest, tags, headers)
	ctx := datadog.NewDefaultContext(context.Background())
	resp, r, err := mySyntheticsApi.client.UpdateAPITest(ctx, publicId, body)

	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.UpdateAPITest`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
	}
	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	logger.Info().Msgf("Response from `SyntheticsApi.UpdateAPITest`:\n%s\n", responseContent)
}

func getGrpcSyntheticTestBody(grpcSyntheticTest model.SyntheticTestGrpcOptions, tags []string, headers map[string]string) datadogV1.SyntheticsAPITest {
	body := datadogV1.SyntheticsAPITest{
		Config: datadogV1.SyntheticsAPITestConfig{
			Assertions: []datadogV1.SyntheticsAssertion{
				datadogV1.SyntheticsAssertion{
					SyntheticsAssertionJSONPathTarget: &datadogV1.SyntheticsAssertionJSONPathTarget{
						Operator: datadogV1.SYNTHETICSASSERTIONJSONPATHOPERATOR_VALIDATES_JSON_PATH,
						Target: &datadogV1.SyntheticsAssertionJSONPathTargetTarget{
							JsonPath:    datadog.PtrString("$.status"),
							Operator:    datadog.PtrString("contains"),
							TargetValue: "SERVING",
						},
						Type: datadogV1.SYNTHETICSASSERTIONTYPE_GRPC_PROTO,
					}},
			},
			Request: &datadogV1.SyntheticsTestRequest{
				Host:                     datadog.PtrString(grpcSyntheticTest.Host),
				Port:                     datadog.PtrInt64(grpcSyntheticTest.Port),
				Service:                  datadog.PtrString(grpcSyntheticTest.ServiceToCheck),
				Method:                   datadog.PtrString(grpcSyntheticTest.Method),
				Message:                  datadog.PtrString(grpcSyntheticTest.Message),
				Timeout:                  datadog.PtrFloat64(grpcSyntheticTest.Timeout), //timeout
				CallType:                 datadogV1.SYNTHETICSTESTCALLTYPE_UNARY.Ptr(),
				Metadata:                 headers,
				CompressedJsonDescriptor: datadog.PtrString(grpcSyntheticTest.CompressedJsonDescriptor),
			},
		},
		Locations: []string{
			fmt.Sprintf("aws:%v", grpcSyntheticTest.Region),
		},
		Message: fmt.Sprintf("Monitor triggered. Notify:  @%v ", grpcSyntheticTest.SlackChannel),
		Name:    grpcSyntheticTest.Name,
		Options: datadogV1.SyntheticsTestOptions{
			MinFailureDuration: datadog.PtrInt64(0),
			MinLocationFailed:  datadog.PtrInt64(1),
			MonitorOptions: &datadogV1.SyntheticsTestOptionsMonitorOptions{
				RenotifyInterval: datadog.PtrInt64(grpcSyntheticTest.RenotifyInterval),
			},
			Retry: &datadogV1.SyntheticsTestOptionsRetry{
				Count:    datadog.PtrInt64(2),
				Interval: datadog.PtrFloat64(300),
			},
			MonitorName: datadog.PtrString(""),
			TickEvery:   datadog.PtrInt64(grpcSyntheticTest.TickEverySec),
		},
		Subtype: datadogV1.SYNTHETICSTESTDETAILSSUBTYPE_GRPC.Ptr(),
		Tags:    tags,
		Type:    datadogV1.SYNTHETICSAPITESTTYPE_API,
	}
	return body
}

func getAllSyntheticTests() (datadogV1.SyntheticsListTestsResponse, error) {
	// Create a new context without specific parameters
	specificContext := context.Background()

	// Get all synthetic tests
	allSyntheticTests, _, err := mySyntheticsApi.client.ListTests(specificContext)
	if err != nil {
		logger.Error().Msgf("Failed with client.ListTests: %v", err)
	}
	return allSyntheticTests, err
}

func (apiClient MySyntheticsApi) CheckIfSyntheticTestExistsByName(name string) (string, error) {
	// Create a new context without specific parameters
	specificContext := context.Background()

	// Get all synthetic tests
	allSyntheticTests, _, err := mySyntheticsApi.client.ListTests(specificContext)

	if err != nil {
		return "", err
	}

	for _, test := range allSyntheticTests.GetTests() {
		testName := test.GetName()
		publicID := test.GetPublicId()
		if testName == name {
			logger.Debug().Msgf("synthetic test:(%s) already exists, public id:(%s)!", testName, publicID)
			return publicID, nil
		}
	}

	return "", nil
}

func (apiClient MySyntheticsApi) CheckIfSyntheticTestExistsByNameAndOneTag(name string, tag string) (string, error) {
	// Create a new context without specific parameters
	specificContext := context.Background()

	// Get all synthetic tests
	allSyntheticTests, _, err := mySyntheticsApi.client.ListTests(specificContext)

	if err != nil {
		return "", err
	}

	for _, test := range allSyntheticTests.GetTests() {
		testName := test.GetName()
		publicID := test.GetPublicId()
		testTags := test.GetTags()
		tagExists := checkIfOneTagExist(tag, testTags)
		if testName == name && tagExists {
			logger.Debug().Msgf("synthetic test:(%s) already exists, public id:(%s)!", testName, publicID)
			return publicID, nil
		}
	}

	return "", nil
}

func (apiClient MySyntheticsApi) CheckIfSyntheticTestExistsByNameAndTags(name string, tags []string) (string, error) {
	// Create a new context without specific parameters
	specificContext := context.Background()

	// Get all synthetic tests
	allSyntheticTests, _, err := mySyntheticsApi.client.ListTests(specificContext)

	if err != nil {
		return "", err
	}

	for _, test := range allSyntheticTests.GetTests() {
		testName := test.GetName()
		publicID := test.GetPublicId()
		testTags := test.GetTags()
		allTagsExist := checkIfAllTagsExist(testTags, tags)
		if testName == name && allTagsExist {
			logger.Debug().Msgf("synthetic test:(%s) already exists, public id:(%s)!", testName, publicID)
			return publicID, nil
		}
	}

	return "", nil
}

func getSyntheticTestsByTags(allSyntheticTests datadogV1.SyntheticsListTestsResponse, tags []string) ([]datadogV1.SyntheticsTestDetails, error) {
	allRelevantTests := []datadogV1.SyntheticsTestDetails{}
	filteredTags := tags

	// We need to fetch ALL tests because there is no efficient filtering function AFAIK
	for _, test := range allSyntheticTests.GetTests() {
		testTags := test.GetTags()
		allTagsExist := checkIfAllTagsExist(testTags, filteredTags)
		if allTagsExist {
			allRelevantTests = append(allRelevantTests, test)
			logger.Debug().Msgf("Got synthetic test with ID %s", test.GetTags())
		}
	}
	return allRelevantTests, nil //Non of internal function returns an error
}

func deleteSyntheticTest(publicId string) error {
	ctx := context.Background() // Create a new context without specific parameters
	body := datadogV1.SyntheticsDeleteTestsPayload{
		PublicIds: []string{
			publicId,
		},
	}
	resp, r, err := mySyntheticsApi.client.DeleteTests(ctx, body)
	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.DeleteTests`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
		return err
	} else {
		responseContent, _ := json.MarshalIndent(resp, "", "  ")
		logger.Debug().Msgf("Response from `SyntheticsApi.DeleteTests`:\n%s\n", responseContent)
		return nil
	}
}

func (apiClient MySyntheticsApi) GetApiTest(publicId string) (*datadogV1.SyntheticsAPITest, error) {

	ctx := datadog.NewDefaultContext(context.Background())
	resp, r, err := mySyntheticsApi.client.GetAPITest(ctx, publicId)

	if err != nil {
		logger.Error().Msgf("Error when calling `SyntheticsApi.GetAPITest`: %v\n", err)
		logger.Error().Msgf("Full HTTP response: %v\n", r)
		return &datadogV1.SyntheticsAPITest{}, err
	}

	//responseContent, _ := json.MarshalIndent(resp, "", "  ")
	//fmt.Fprintf(os.Stdout, "Response from `SyntheticsApi.GetAPITest`:\n%s\n", responseContent)

	return &resp, nil

}
