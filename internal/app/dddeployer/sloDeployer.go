package dddeployer

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/rotemjac/datadog-automation/internal/app/ddsynfacade"
)

type SloDeployer struct {
	Cluster     string `validate:"required"`
	Environment string `validate:"required,oneof=dev staging prod"`
	Team        string `validate:"required,oneof=cwp-runtime ops"`
}

func NewSloDeployer(cluster string, environment string, team string) IDatadogDeployer {
	sloDeployer := &SloDeployer{
		Cluster:     cluster,
		Environment: environment,
		Team:        team,
	}

	return sloDeployer
}

func (o *SloDeployer) ValidateTags() (bool, error) {
	validate = validator.New()

	// Validate the values of tagsFilter
	errs := validate.Struct(o)
	if errs != nil {
		return false, errs
	}
	return true, nil
}

func (o *SloDeployer) Execute(jwtToken string, beBaseUrl string, region string, slackChannel string) (bool, error) {
	var err error
	tags := []string{
		fmt.Sprintf("cluster:%s", o.Cluster),
		fmt.Sprintf("env:%s", o.Environment),
	}

	syntheticTestsConfig, err := readSyntheticFromFile(configFilePath, configFileName)
	if err != nil {
		logger.Error().Msgf("Error in readSyntheticFromFile: %s\n", err)
	} else {
		logger.Info().Msgf("SyntheticTestsConfig: %v\n", syntheticTestsConfig)
	}

	syntheticTestsConfig = prepareSyntheticTests(syntheticTestsConfig, jwtToken, beBaseUrl, region, slackChannel)

	logger.Debug().Msgf("syntheticTestsConfig after preparation: (%v)...\n", syntheticTestsConfig)

	for _, syntheticTest := range syntheticTestsConfig.HttpTests {
		logger.Info().Msgf("syntheticTest url: %v\n", syntheticTest.Url)
	}

	stClient := ddsynfacade.MySyntheticsApi{}
	for _, syntheticTest := range syntheticTestsConfig.HttpTests {
		publicId, err := stClient.CheckIfSyntheticTestExistsByNameAndTags(syntheticTest.Name, tags)
		if err != nil {
			return false, err
		}
		if publicId == "" {
			logger.Info().Msgf("Synthetic test: (%s) not exists, creating it...\n", syntheticTest.Name)
			stClient.CreateHttpSyntheticTest(syntheticTest, o.GetTags(), o.GetHeaders(jwtToken, "true"))
		} else {
			logger.Info().Msgf("Synthetic test: (%s) already exists, updating it...\n", syntheticTest.Name)
			stClient.UpdateHttpSyntheticTest(syntheticTest, publicId, o.GetTags(), o.GetHeaders(jwtToken, "true"))
		}
	}

	for _, syntheticTest := range syntheticTestsConfig.GrpcTests {

		publicId, err := stClient.CheckIfSyntheticTestExistsByNameAndTags(syntheticTest.Name, tags)
		if err != nil {
			return false, err
		}
		if publicId == "" {
			stClient.CreateGrpcSyntheticTest(syntheticTest, o.GetTags(), o.GetHeaders(jwtToken, "true"))
		} else {
			logger.Info().Msgf("Synthetic test: (%s) already exists, updating it...\n", syntheticTest.Name)
			stClient.UpdateGrpcSyntheticTest(syntheticTest, publicId, o.GetTags(), o.GetHeaders(jwtToken, "true"))
		}

	}

	return true, nil
}

func (o *SloDeployer) GetTags() []string {
	tags := []string{
		fmt.Sprintf("env:%s", o.Environment),
		fmt.Sprintf("team:%s", o.Team),
		fmt.Sprintf("cluster:%s", o.Cluster),
	}

	return tags
}

func (o *SloDeployer) GetHeaders(jwtToken string, requestBody string) map[string]string {
	var metadata = make(map[string]string)
	if jwtToken != "" {
		metadata["Authorization"] = fmt.Sprintf("Bearer {{%v}}", jwtToken)
	}
	if requestBody != "" {
		metadata["content-type"] = "application-json"
	}

	return metadata
}
