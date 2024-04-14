package helper

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	appConfig "github.com/silasstoffel/account-service/configs"
	"github.com/silasstoffel/account-service/internal/exception"
)

func BuildAwsConfig(app *appConfig.Config) (cfg aws.Config, er error) {
	var awsCfg aws.Config
	var err error

	if app.Env == "development" {
		awsCfg, err = buildDevEnvironment(app)
	} else {
		awsRegion := app.Aws.Region
		awsCfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(awsRegion),
		)
	}

	if err != nil {
		message := "Error when creating AWS client"
		log.Println(message, err.Error())
		return aws.Config{}, exception.NewUnknownError(&err)
	}

	return awsCfg, nil
}

func buildDevEnvironment(app *appConfig.Config) (cfg aws.Config, err error) {
	awsRegion := app.Aws.Region
	awsEndpoint := app.Aws.Endpoint

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     "localstack",
				SecretAccessKey: "localstack",
			}, nil
		})),
	)

	if err != nil {
		return aws.Config{}, err
	}

	return awsCfg, nil
}
