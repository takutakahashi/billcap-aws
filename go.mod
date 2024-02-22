module github.com/takutakahashi/billcap-aws

go 1.20

replace github.com/takutakahashi/billcap-schema => ../billcap-schema

require (
	github.com/aws/aws-sdk-go-v2 v1.25.0
	github.com/aws/aws-sdk-go-v2/config v1.27.1
	github.com/aws/aws-sdk-go-v2/service/costexplorer v1.34.2
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cobra v1.8.0
	github.com/takutakahashi/billcap-schema v0.0.0-00010101000000-000000000000
)

require (
	github.com/aws/aws-sdk-go-v2/credentials v1.17.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.15.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.11.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.19.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.22.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.27.1 // indirect
	github.com/aws/smithy-go v1.20.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/sys v0.16.0 // indirect
)
