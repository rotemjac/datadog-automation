package dddeployer

import "fmt"

func GetDatadogDeployer(deployerType string, tenantId string, cluster string, environment string, team string) (IDatadogDeployer, error) {
	switch deployerType {
	case "tenant":
		return NewTenantDeployer(tenantId, environment, team), nil
	case "slo":
		return NewSloDeployer(cluster, environment, team), nil
	default:
		return nil, fmt.Errorf("Wrong datadog deployer type passed")
	}
}
