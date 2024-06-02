package dddeployer

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestSloDeployerTags(t *testing.T) {
	validate = validator.New()
	// Defining the columns of the table

	var tests = []struct {
		name             string
		Cluster          string
		inputEnvironment string
		inputTeam        string
		desiredRes       bool
	}{
		{"Valid tags - dev", "edge", "dev", "runtime", true},
		{"Valid tags - staging", "gateway", "staging", "runtime", true},
		{"Valid tags - prod", "edge", "prod", "ops", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syn_test_deployer, err := GetDatadogDeployer("slo", "223423", tt.Cluster, tt.inputEnvironment, tt.inputTeam)
			if err != nil {
				//t.Errorf("got %b, want %b", ans, tt.desiredRes)
				t.Errorf("Error! %s", err)
			}

			ans, _ := syn_test_deployer.ValidateTags()
			if ans != tt.desiredRes {
				//t.Errorf("got %b, want %b", ans, tt.desiredRes)
				t.Errorf("Error!")
			}
		})
	}
}
