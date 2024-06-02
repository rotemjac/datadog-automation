package dddeployer

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/rotemjac/datadog-automation/internal/app/ddsynfacade"
)

type TenantDeployer struct {
	TenantId    string `validate:"required"`
	Environment string `validate:"required,oneof=dev staging prod"`
	Team        string `validate:"required,oneof=runtime ops"`
}

func NewTenantDeployer(tenantId string, environment string, team string) IDatadogDeployer {
	tenantDeployer := &TenantDeployer{
		TenantId:    tenantId,
		Environment: environment,
		Team:        team,
	}

	return tenantDeployer
}

func (o *TenantDeployer) ValidateTags() (bool, error) {
	validate = validator.New()

	// Validate the values of tagsFilter
	errs := validate.Struct(o)
	if errs != nil {
		return false, errs
	}
	return true, nil
}

func (o *TenantDeployer) Execute(jwtToken string, beBaseUrl string, region string, slackChannel string) (bool, error) {
	var err error

	syntheticTestsConfig, err := readSyntheticFromFile(configFilePath, configFileName)
	if err != nil {
		logger.Error().Msgf("Error in readSyntheticFromFile: %s\n", err)
		return false, fmt.Errorf("Error in readSyntheticFromFile: %s\n", err)
	} else {
		logger.Info().Msgf("syntheticTestsConfig: %v\n", syntheticTestsConfig)
	}

	if o.TenantId == "" {
		return false, fmt.Errorf("Tenant id is empty")
	}

	for index, syntheticTest := range syntheticTestsConfig.HttpTests {
		syntheticTestsConfig.HttpTests[index].Name = syntheticTest.Name + "-" + o.TenantId
	}
	for index, syntheticTest := range syntheticTestsConfig.GrpcTests {
		syntheticTestsConfig.GrpcTests[index].Name = syntheticTest.Name + "-" + o.TenantId
	}
	syntheticTestsConfig = prepareSyntheticTests(syntheticTestsConfig, jwtToken, beBaseUrl, region, slackChannel)

	stClient := ddsynfacade.MySyntheticsApi{}
	for _, syntheticTest := range syntheticTestsConfig.HttpTests {
		publicId, err := stClient.CheckIfSyntheticTestExistsByName(syntheticTest.Name)
		if err != nil {
			return false, err
		}
		if publicId == "" {
			logger.Info().Msgf("Synthetic test: (%s) not exists, creating it...\n", syntheticTest.Name)
			stClient.CreateHttpSyntheticTest(syntheticTest, o.GetTags(), o.GetHeaders(jwtToken, ""))
		} else {
			logger.Info().Msgf("Synthetic test: (%s) already exists, updating it...\n", syntheticTest.Name)
			stClient.UpdateHttpSyntheticTest(syntheticTest, publicId, o.GetTags(), o.GetHeaders(jwtToken, ""))
		}
	}

	for _, syntheticTest := range syntheticTestsConfig.GrpcTests {

		publicId, err := stClient.CheckIfSyntheticTestExistsByName(syntheticTest.Name)
		if err != nil {
			return false, err
		}
		if publicId == "" {
			stClient.CreateGrpcSyntheticTest(syntheticTest, o.GetTags(), o.GetHeaders(jwtToken, ""))
		} else {
			logger.Info().Msgf("Synthetic test: (%s) already exists, updating it...\n", syntheticTest.Name)
			stClient.UpdateGrpcSyntheticTest(syntheticTest, publicId, o.GetTags(), o.GetHeaders(jwtToken, ""))
		}

	}
	return true, nil
}

func (o *TenantDeployer) GetTags() []string {
	tags := []string{
		fmt.Sprintf("env:%s", o.Environment),
		fmt.Sprintf("team:%s", o.Team),
		fmt.Sprintf("tenant-id:%s", o.TenantId),
	}

	return tags
}

func (o *TenantDeployer) GetHeaders(jwtToken string, requestBody string) map[string]string {
	var metadata = make(map[string]string)
	metadata["x-tenant-id"] = o.TenantId

	if jwtToken != "" {
		metadata["Authorization"] = fmt.Sprintf("Bearer {{%v}}", jwtToken)
	}
	if requestBody != "" {
		metadata["content-type"] = "application-json"
	}

	return metadata
}
