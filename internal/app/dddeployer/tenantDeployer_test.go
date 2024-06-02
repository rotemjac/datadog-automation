package dddeployer

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

func TestTenantDeployerTags(t *testing.T) {
	validate = validator.New()
	// Defining the columns of the table
	var tests = []struct {
		name             string
		inputTenantId    string
		inputEnvironment string
		inputTeam        string
		desiredRes       bool
	}{
		{"Valid tags - dev", "1122334455", "dev", "runtime", true},
		{"Valid tags - staging", "2233445566", "staging", "runtime", true},
		{"Valid tags - prod", "3344556677", "prod", "ops", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			syn_test_deployer, err := GetDatadogDeployer("tenant", tt.inputTenantId, "", tt.inputEnvironment, tt.inputTeam)
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
