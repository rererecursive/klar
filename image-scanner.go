package main

import (
    "encoding/json"
    "fmt"
    "os"
    "strings"
    b64 "encoding/base64"

    "github.com/aws/aws-sdk-go/service/ecr"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/awserr"
    "github.com/eawsy/aws-cloudformation-go-customres/service/cloudformation/customres"
    "github.com/eawsy/aws-lambda-go-core/service/lambda/runtime"
    "github.com/eawsy/aws-lambda-go-event/service/lambda/runtime/event/cloudformationevt"
)

/*
    "Region": "ap-southeast-2",
    "Action": "Fail",   // Fail | Notify | Ignore
    "SNSFailure": "",
    "Repository": "",
    "Image": "",
    "Tags": "",
    "Threshold": "",
    "EventID": "<random>
*/

var (
    // Handle is the Lambda's entrypoint.
    Handle customres.LambdaHandler
)

func init() {
    customres.Register("ScanRequest", new(ScanRequest))
    Handle = customres.HandleLambda
}

// Happy IDE means happy developer.
func main() {

}

// ScanRequest represents a simple, custom resource.
type ScanRequest struct {
    Action *string `json:"Action"`
    Region *string `json:"Region"`
    Repository *string `json:"Repository"`
    Image *string `json:"Image"`
    Tags *string `json:"Tags"`
    Threshold *string `json:"Threshold"`
}

// Create is invoked when the resource is created.
func (r *ScanRequest) Create(evt *cloudformationevt.Event, ctx *runtime.Context) (string, interface{}, error) {
    // TODO: export each property as an environment variable
    evt.PhysicalResourceID = customres.NewPhysicalResourceID(evt)
    return r.Update(evt, ctx)
}

var (
    defaultExampleThing = "Fail!"
)

// Update is invoked when the resource is updated.
func (r *ScanRequest) Update(evt *cloudformationevt.Event, ctx *runtime.Context) (string, interface{}, error) {
    if err := json.Unmarshal(evt.ResourceProperties, r); err != nil {
        return "", r, err
    }

    if r.Action == nil {
        r.Action = &defaultExampleThing
    }

    username, password := GetDockerLogin(r.Region)

    os.Setenv("DOCKER_USERNAME", username)
    os.Setenv("DOCKER_PASSWORD", password)
    os.Setenv("CLAIR_ADDR", "coreo-Clair-QDQKVCAO6IRL-2075767864.ap-southeast-2.elb.amazonaws.com")
    os.Setenv("CLAIR_OUTPUT", *r.Threshold)

    klar(r.Image)

    return evt.PhysicalResourceID, r, nil
}

// Delete is invoked when the resource is deleted.
func (r *ScanRequest) Delete(*cloudformationevt.Event, *runtime.Context) error {
    return nil
}

func GetDockerLogin(region *string) (username string, password string) {
    svc := ecr.New(session.New(&aws.Config {
        Region: aws.String(*region)},
    ))

    input := &ecr.GetAuthorizationTokenInput{}

    result, err := svc.GetAuthorizationToken(input)
    if err != nil {
        if aerr, ok := err.(awserr.Error); ok {
            switch aerr.Code() {
            case ecr.ErrCodeServerException:
                fmt.Println(ecr.ErrCodeServerException, aerr.Error())
            case ecr.ErrCodeInvalidParameterException:
                fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
            default:
                fmt.Println(aerr.Error())
            }
        } else {
            // Print the error, cast err to awserr.Error to get the Code and
            // Message from an error.
            fmt.Println(err.Error())
        }
        return
    }

    token := *result.AuthorizationData[0].AuthorizationToken
    dec, _ := b64.StdEncoding.DecodeString(token)
    parts := strings.Split(string(dec), ":")
    return parts[0], parts[1]
}
