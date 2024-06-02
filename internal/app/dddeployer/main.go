package dddeployer

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"

	"github.com/rotemjac/datadog-automation/internal/pkg/model"
)

// use a single instance of Validate, it caches struct info
var validate *validator.Validate
var logger zerolog.Logger
var logLevel string
var configFilePath string
var configFileName string

func init() {

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
	logger.Info().Msg("datadog deployer service has started..")

}

func readSyntheticFromFile(configFilePath string, configFileName string) (*model.SyntheticTestsConfig, error) {
	// Construct the full path to the YAML file
	logger.Debug().Msgf("In readSyntheticFromFile..")
	yamlPath := filepath.Join(configFilePath, configFileName)
	logger.Debug().Msgf("Reading Yaml file: (%s)\n", yamlPath)

	// Read the YAML file
	dataInYaml, err := os.ReadFile(yamlPath)
	if err != nil {
		logger.Error().Msgf("Failed to read YAML: %v\n", err)
		return nil, err
	} else {
		logger.Info().Msgf("Succeed to read YAML")
	}

	// Create a new Config object to store the parsed data
	syntheticTestsConfig := model.SyntheticTestsConfig{}
	err2 := yaml.Unmarshal(dataInYaml, &syntheticTestsConfig)

	if err2 != nil {
		logger.Error().Msgf("Failed to unmarshal YAML: %v\n", err2)
		return nil, err2
	} else {
		logger.Info().Msgf("Succeed to unmarshal YAML: %s, content: %v\n", yamlPath, syntheticTestsConfig)
	}
	return &syntheticTestsConfig, nil
}

func prepareSyntheticTests(syntheticTestsConfig *model.SyntheticTestsConfig, jwtToken string, beBaseUrl string, region string, slackChannel string) *model.SyntheticTestsConfig {

	for index, syntheticTest := range syntheticTestsConfig.HttpTests {
		syntheticTestsConfig.HttpTests[index].JwtToken = jwtToken
		syntheticTestsConfig.HttpTests[index].Url = "https://" + beBaseUrl + syntheticTest.Url
		syntheticTestsConfig.HttpTests[index].Region = region
		syntheticTestsConfig.HttpTests[index].SlackChannel = slackChannel
		logger.Info().Msgf("Synthetic http test name: %s\n", syntheticTest.Name)
		logger.Debug().Msgf("Synthetic http test JwtToken: %s\n", syntheticTest.JwtToken)
		logger.Debug().Msgf("Synthetic http test Url: %s\n", syntheticTest.Url)
		logger.Debug().Msgf("Synthetic http test Region: %s\n", syntheticTest.Region)
	}

	for index, syntheticTest := range syntheticTestsConfig.GrpcTests {
		syntheticTestsConfig.GrpcTests[index].JwtToken = jwtToken
		syntheticTestsConfig.GrpcTests[index].Url = beBaseUrl
		syntheticTestsConfig.GrpcTests[index].Host = beBaseUrl
		syntheticTestsConfig.GrpcTests[index].Region = region
		syntheticTestsConfig.GrpcTests[index].SlackChannel = slackChannel
		logger.Info().Msgf("Synthetic grpc test name: %s\n", syntheticTest.Name)
		logger.Debug().Msgf("Synthetic grpc test JwtToken: %s\n", syntheticTest.JwtToken)
		logger.Debug().Msgf("Synthetic grpc test Url: %s\n", syntheticTest.Url)
		logger.Debug().Msgf("Synthetic grpc test Host: %s\n", syntheticTest.Host)
		logger.Debug().Msgf("Synthetic grpc test Region: %s\n", syntheticTest.Region)

	}

	return syntheticTestsConfig

}
