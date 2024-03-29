package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type IdentityTokenRetriever struct{}

func (i IdentityTokenRetriever) GetIdentityToken() ([]byte, error) {
	file := os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE")
	return ioutil.ReadFile(file)
}

func main() {
	t_ := os.Getenv("IRSA_TOKENGEN_DURATION_IN_HOURS")
	if t_ == "" {
		t_ = "1"
	}
	t, err := strconv.Atoi(t_)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	stsClient := sts.NewFromConfig(cfg)
	jwtGetter := IdentityTokenRetriever{}
	webIdentityRoleProvider := stscreds.NewWebIdentityRoleProvider(stsClient, os.Getenv("AWS_ROLE_ARN"), jwtGetter, func(o *stscreds.WebIdentityRoleOptions) {
		o.Duration = time.Duration(t) * time.Hour
	})
	credentials, err := webIdentityRoleProvider.Retrieve(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s=%s\n", "AWS_ACCESS_KEY_ID", credentials.AccessKeyID)
	fmt.Printf("%s=%s\n", "AWS_SECRET_ACCESS_KEY", credentials.SecretAccessKey)
	fmt.Printf("%s=%s\n", "AWS_SESSION_TOKEN", credentials.SessionToken)
}
