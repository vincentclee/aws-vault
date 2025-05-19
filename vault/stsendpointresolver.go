package vault

import (
	"context"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

// GetSTSEndpointResolver returns an endpoint resolver for STS that handles regional endpoints settings
// in accordance with https://docs.aws.amazon.com/credref/latest/refdocs/setting-global-sts_regional_endpoints.html
func GetSTSEndpointResolver(stsRegionalEndpoints string) sts.EndpointResolverV2 {
	return &customSTSEndpointResolver{stsRegionalEndpoints: stsRegionalEndpoints}
}

// Legacy non-exported function maintained for backward compatibility
func getSTSEndpointResolver(stsRegionalEndpoints string) sts.EndpointResolverV2 {
	return GetSTSEndpointResolver(stsRegionalEndpoints)
}

type customSTSEndpointResolver struct {
	stsRegionalEndpoints string
}

func (r *customSTSEndpointResolver) ResolveEndpoint(ctx context.Context, params sts.EndpointParameters) (smithyendpoints.Endpoint, error) {
	region := ""
	if params.Region != nil {
		region = *params.Region
	}

	if r.stsRegionalEndpoints == "legacy" && region != "" {
		if region == "ap-northeast-1" ||
			region == "ap-south-1" ||
			region == "ap-southeast-1" ||
			region == "ap-southeast-2" ||
			region == "aws-global" ||
			region == "ca-central-1" ||
			region == "eu-central-1" ||
			region == "eu-north-1" ||
			region == "eu-west-1" ||
			region == "eu-west-2" ||
			region == "eu-west-3" ||
			region == "sa-east-1" ||
			region == "us-east-1" ||
			region == "us-east-2" ||
			region == "us-west-1" ||
			region == "us-west-2" {
			log.Println("Using legacy STS endpoint sts.amazonaws.com")

			uri, _ := url.Parse("https://sts.amazonaws.com")
			return smithyendpoints.Endpoint{
				URI: *uri,
			}, nil
		}
	}
	return smithyendpoints.Endpoint{}, &aws.EndpointNotFoundError{}
}
