package bootstrap

import (
	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"os"
)

var datadogApiClient *datadog.APIClient

func main() {

}

func init() {
	// Fetch from secret manager
	//apiKey := ""
	//appKey := ""

	// #############  Get ENV vars ############# //
	apiKey := os.Getenv("API_KEY")
	appKey := os.Getenv("APP_KEY")

	// Initialize the Datadog client
	config := datadog.NewConfiguration()
	config.AddDefaultHeader("DD-API-KEY", apiKey)
	config.AddDefaultHeader("DD-APPLICATION-KEY", appKey)
	datadogApiClient = datadog.NewAPIClient(config)
}

func Client() *datadog.APIClient {
	return datadogApiClient
}
