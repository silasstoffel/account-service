package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type AwsConfig struct {
	Region                 string
	Endpoint               string // keep empty for production
	AccountServiceTopicArn string
	AccountServiceQueueUrl string
	WebhookSenderQueueUrl  string
}

type AppConfig struct {
	ApiPort string
	ApiHost string
}

type AppDatabase struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type Config struct {
	Env string
	Aws AwsConfig
	App AppConfig
	Db  AppDatabase
}

func NewConfigFromEnvVars() *Config {
	awsEndpoint := ""
	apiPort := "80"
	if os.Getenv("ENV") == "" || os.Getenv("ENV") == "development" {
		os.Setenv("ENV", "development")
		if err := godotenv.Load(); err != nil {
			panic(err)
		}
		awsEndpoint = os.Getenv("AWS_ENDPOINT")
	}

	if os.Getenv("API_PORT") != "" {
		apiPort = os.Getenv("API_PORT")
	}

	return &Config{
		Env: os.Getenv("ENV"),
		Aws: AwsConfig{
			Region:                 os.Getenv("AWS_REGION"),
			Endpoint:               awsEndpoint,
			AccountServiceTopicArn: os.Getenv("ACCOUNT_SERVICE_TOPIC_ARN"),
			AccountServiceQueueUrl: os.Getenv("ACCOUNT_SERVICE_QUEUE_URL"),
			WebhookSenderQueueUrl:  os.Getenv("WEBHOOK_SENDER_QUEUE_URL"),
		},
		App: AppConfig{
			ApiPort: apiPort,
			ApiHost: os.Getenv("API_HOST"),
		},
		Db: AppDatabase{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}
