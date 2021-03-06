package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentity"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"log"
	"time"
)

type credential struct {
	Version int `json:"Version"`
	AccessKeyId string `json:"AccessKeyId"`
	SecretAccessKey string `json:"SecretAccessKey"`
	SessionToken string `json:"SessionToken"`
	Expiration time.Time `json:"Expiration"`
}

const (
	ClientId string = "re23i75fp5ln62p1vtssjrveo" // Cognito app client ID
	Region string = "us-west-2" // region where Cognito is being hosted
	AccountId string = "0123456789012" // aws account ID
	IdentityPoolId = "us-west-2:070eee28-81b9-3767-b400-b8be70facf40" // Federated ID pool ID
	UserPoolId = "us-west-2_UXqM55clL" // Cognito user pool ID
)
func main() {
	username := flag.String("u", "", "username")
	password := flag.String("p", "", "password")
	flag.Parse()
	sess := session.Must(session.NewSession())
	id := cognitoidentity.New(sess, aws.NewConfig().WithRegion(Region))
	login := cognitoidentityprovider.New(sess, aws.NewConfig().WithRegion(Region))
	authOutput, err := login.InitiateAuth(&cognitoidentityprovider.InitiateAuthInput{
		AuthFlow:          aws.String("USER_PASSWORD_AUTH"),
		AuthParameters:    map[string]*string{"USERNAME":username,"PASSWORD":password},
		ClientId:          aws.String(ClientId),
	})
	if err != nil {
		log.Fatal(err)
	}
	idOutput, err := id.GetId(
		&cognitoidentity.GetIdInput{
				AccountId: aws.String(AccountId),
				IdentityPoolId: aws.String(IdentityPoolId),
				Logins: map[string]*string{"cognito-idp.us-west-2.amazonaws.com/" + UserPoolId:authOutput.AuthenticationResult.IdToken},
			},
		)
	if err != nil {
		log.Fatal(err)
	}
	cred, err := id.GetCredentialsForIdentity(
		&cognitoidentity.GetCredentialsForIdentityInput{
			IdentityId:    idOutput.IdentityId,
			Logins: map[string]*string{"cognito-idp.us-west-2.amazonaws.com/" + UserPoolId:authOutput.AuthenticationResult.IdToken},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	c := &credential{
		Version:         1,
		AccessKeyId:     aws.StringValue(cred.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(cred.Credentials.SecretKey),
		SessionToken:    aws.StringValue(cred.Credentials.SessionToken),
		Expiration:      aws.TimeValue(cred.Credentials.Expiration),
	}
	jo, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jo))
}
