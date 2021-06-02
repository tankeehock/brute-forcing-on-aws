package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Guesser verifies the existance of a specific access key and secret key combination
type Guesser struct {
	awsSession *session.Session
}

func (g *Guesser) verifyKey(accessKey, secretKey string) (err error) {
	if g.awsSession == nil {
		g.awsSession, err = session.NewSession(aws.NewConfig().WithRegion(region))
		if err != nil {
			return
		}
	}
	g.awsSession.Config.Credentials = credentials.NewStaticCredentials(accessKey, secretKey, "")

	svc := sts.New(g.awsSession)
	_, err = svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	return
}
