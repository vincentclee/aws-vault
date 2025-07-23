package vault

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// getSTSClientOptions returns STS client options based on regional endpoint configuration
// in accordance with https://docs.aws.amazon.com/credref/latest/refdocs/setting-global-sts_regional_endpoints.html
func getSTSClientOptions(stsRegionalEndpoints string) []func(*sts.Options) {
	if stsRegionalEndpoints == "legacy" {
		log.Println("Using legacy STS endpoint sts.amazonaws.com")
		return []func(*sts.Options){
			func(o *sts.Options) {
				// Set BaseEndpoint to use the global STS endpoint
				globalEndpoint := "https://sts.amazonaws.com"
				o.BaseEndpoint = &globalEndpoint
			},
		}
	}

	// For regional endpoints (default), return empty options to use default regional behavior
	return []func(*sts.Options){}
}
