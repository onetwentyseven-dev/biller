package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kelseyhightower/envconfig"
	ssmconfig "github.com/onetwentyseven-dev/go-ssm-config"
)

var appConfig struct {
	*envConfig
	*ssmConfig
}

type envConfig struct {
	DBHost       string `envconfig:"DB_HOST" required:"true"`
	DBSchema     string `envconfig:"DB_SCHEMA" required:"true"`
	DBUsername   string `envconfig:"DB_USER" required:"true"`
	SSMPrefix    string `envconfig:"SSM_PREFIX" required:"true"`
	AuthClientID string `envconfig:"AUTH_CLIENT_ID" required:"true"`
	AuthAudience string `envconfig:"AUTH_AUDIENCE" required:"true"`
	AuthTenant   string `envconfig:"AUTH_TENANT" required:"true"`
}

type ssmConfig struct {
	DatabasePassword string `ssm:"db_pass" required:"true"`
}

func loadConfig(awsConfig aws.Config) {
	var env = new(envConfig)
	err := envconfig.Process("", env)
	if err != nil {
		panic(fmt.Sprintf("envconfig: %s", err))
	}

	var ssm = new(ssmConfig)
	err = ssmconfig.Process(context.TODO(), awsConfig, env.SSMPrefix, ssm)
	if err != nil {
		panic(fmt.Sprintf("ssmconfig: %s", err))
	}

	appConfig.envConfig = env
	appConfig.ssmConfig = ssm
}
