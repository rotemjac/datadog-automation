package dddeployer

type IDatadogDeployer interface {
	ValidateTags() (bool, error)
	Execute(jwtToken string, beBaseUrl string, region string, slackChannel string) (bool, error)
	GetTags() []string
	GetHeaders(jwtToken string, requestBody string) map[string]string
}
