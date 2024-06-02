package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"

	"github.com/rotemjac/datadog-automation/internal/app/dddeployer"
)

var synTestType, tenantId, region, environment, team, jwtToken, beBaseUrl, logLevel, slackChannel, cluster, debugSleep string

// use a single instance of Validate, it caches struct info
// var validate *validator.Validate
var logger zerolog.Logger

func main() {

	// #############  Get ENV vars ############# //
	synTestType = os.Getenv("SYN_TEST_TYPE") //options: "tenant" , "slo"
	tenantId = os.Getenv("TENANT_ID")
	environment = os.Getenv("ENVIRONMENT")
	cluster = os.Getenv("CLUSTER")
	team = os.Getenv("TEAM")
	jwtToken = os.Getenv("JWT_TOKEN")
	beBaseUrl = os.Getenv("BE_BASE_URL")
	region = os.Getenv("REGION")
	logLevel = os.Getenv("LOG_LEVEL")
	slackChannel = os.Getenv("SLACK_CHANNEL")
	debugSleep = os.Getenv("DEBUG_SLEEP")

	// ############# Logging Setup ############# //
	// Set log level
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if logLevel == "DEBUG" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	logger = zerolog.New(os.Stderr).With().Logger()
	logger.Info().Msgf("Datadog Automation deployer service has started, mode (%s)..", synTestType)

	// ############# Validation Setup ############# //
	//validate = validator.New()

	// ####### debug and print synthetic tests for comparison
	//stClient := ddsynfacade.MySyntheticsApi{}
	//stClient.GetApiTest("r9q-zvp-6pi")

	//responseContent, _ := json.MarshalIndent(resp, "", "  ")
	//logger.Info().Msgf("Response from `SyntheticsApi.GetAPITest`:\n%s\n", responseContent)

	//stClient.GetApiTest("rz8-5gj-pue")
	//responseContent2, _ := json.MarshalIndent(resp2, "", "  ")
	//logger.Info().Msgf("Response2 from `SyntheticsApi.GetAPITest`:\n%s\n", responseContent2)
	// ####### debug and print synthetic tests for comparison

	// Getting the relevant datadog deployer from factory pattern based on SYN_TEST_TYPE env var
	synTestDeployer, err := dddeployer.GetDatadogDeployer(synTestType, tenantId, cluster, environment, team)
	if err != nil {
		logger.Fatal().Msgf("FATAL Error: Not a valid synthetic test type (%s)", synTestType)
	} else {
		// Validate inputs before executing
		tagsOk, tagsErr := synTestDeployer.ValidateTags()

		if !tagsOk {
			logger.Fatal().Msgf("FATAL Error: Exit with Failed validation on Tags %s", tagsErr)
		}

		// creating synthetic tests
		executionOk, executionErr := synTestDeployer.Execute(jwtToken, beBaseUrl, region, slackChannel)
		if executionOk {
			logger.Debug().Msg("Execution passed with no errors..")
		} else {
			logger.Fatal().Msgf("FATAL Error: Execution failed: %s", executionErr)
		}
	}

	if logLevel == "DEBUG" && debugSleep == "true" {
		time.Sleep(50 * time.Second)
	}
}
