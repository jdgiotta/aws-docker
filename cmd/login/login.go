package login

import (
    "encoding/base64"
    "flag"
    "fmt"
    "os"
    "strings"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/ecr"
    "github.com/docker/docker/api/types"
    dockerclient "github.com/docker/docker/client"
    "github.com/jdgiotta/aws-docker/cli"
    "golang.org/x/net/context"
)

//var debug = cli.ToggleDebug()
var opts = options{}

type options struct {
    verbose *bool
    profiles *string
}

type AWSConfigProperties struct {
    Profile *string
    ID *string
    Region *string
}

// Do perform the action of getting profile registry IDs and executes the login to ecr services
func Do() {
    //TODO var opts options

    loginCommand := flag.NewFlagSet("login", flag.ExitOnError)

    opts.profiles = loginCommand.String("profile", "", "Profile the registry belongs to")
    opts.verbose = loginCommand.Bool("v", false, "Verbose output")

    if hasRequiredArgs(loginCommand) == false {
        loginCommand.PrintDefaults()
        os.Exit(1)
    }

    awsConfigProfiles := []*string{}
    region := "us-east-1"
    //for _, profile := range *opts.profiles {
        section, err := cli.GetProfileSection(*opts.profiles)
        if err == nil {
            registryId, err := section.GetKey("ecrregistryid")
            if err == nil {
                if *opts.verbose {
                    fmt.Printf("Obtained Registry ID %s\n", registryId)
                }
                k := registryId.String()
                awsConfigProfiles = append(awsConfigProfiles, &k)
            } else {
                fmt.Printf("No registry ids were found for profiles %s \n", registryId.String())
            }
            regionKey, err := section.GetKey("region")
            if err == nil {
                if *opts.verbose {
                    fmt.Printf("Obtained Region %s\n", regionKey)
                }
                region = regionKey.String()
            }

        } else {
            fmt.Printf("No profile was found for %s \n", *opts.profiles)
        }
    //}

    input := &ecr.GetAuthorizationTokenInput{
        RegistryIds: awsConfigProfiles,
    }
    sess := session.Must(session.NewSession(&aws.Config{Region:&region}))
    svc := ecr.New(sess)

    output, err := svc.GetAuthorizationToken(input)
    if err != nil {
        fmt.Println(err)
    }

    for _, a := range output.AuthorizationData {
        data, err := base64.StdEncoding.DecodeString(*a.AuthorizationToken)
        split := strings.Split(string(data), ":")
        auth := types.AuthConfig{
            Username:      split[0],
            Password:      split[1],
            Email:         "none",
            ServerAddress: *a.ProxyEndpoint,

        }

        dockerCli, err := dockerclient.NewEnvClient()
        must(err)
        okBody, err := dockerCli.RegistryLogin(context.Background(),auth)
        must(err)
        fmt.Println(okBody.Status)
    }
}

func hasRequiredArgs(cmd *flag.FlagSet) bool {
    cmd.Parse(os.Args[2:])
    if cmd.Parsed() {
        if *opts.profiles != "" {
            return true
        }
    }
    return false
}

func must(err error) {
    if err != nil {
        fmt.Println(err)
        os.Exit(2)
    }
}