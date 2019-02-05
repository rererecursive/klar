package main

import (
    "encoding/json"
    "os"

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

    os.Setenv("CLAIR_ADDR", "coreo-Clair-QDQKVCAO6IRL-2075767864.ap-southeast-2.elb.amazonaws.com")
    os.Setenv("CLAIR_OUTPUT", *r.Threshold)
    os.Setenv("DOCKER_TOKEN", )
    // TODO: include SDK, https://docs.aws.amazon.com/sdk-for-go/api/service/ecr/#ECR.GetAuthorizationToken
    GetAuthorizationToken
    klar(r.Image)

    return evt.PhysicalResourceID, r, nil
}

// Delete is invoked when the resource is deleted.
func (r *ScanRequest) Delete(*cloudformationevt.Event, *runtime.Context) error {
    return nil
}
