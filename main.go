package main

import (
    "os/exec"
    "os"
    "fmt"
    "os/user"
    "github.com/go-ini/ini"
    "strings"
)

var debug = false

func main () {
    for _, element := range os.Args {
        if element == "-v" {
            debug = true
            break
        }
    }
    if len(os.Args) < 2 {
        printHelp()
        os.Exit(2)
    }
    switch os.Args[1] {
    case "login":
        dologin()
    default:
        printHelp()
    }
}

func dologin () {
    usr, err := user.Current()
    must(err)
    if debug {
        fmt.Printf("From user directory %s\n", usr.HomeDir)
    }

    cfg, err := ini.LooseLoad(usr.HomeDir+"/.aws/config")

    must(err)

    if len(os.Args) < 3 {
        printLoginHelp()
        os.Exit(2)
    }

    awsConfigProfiles := []string{}
    for _, profile := range os.Args[2:] {
        section, err := cfg.GetSection(fmt.Sprintf("profile %s", profile))
        if err == nil {
            key, err := section.GetKey("ecrregistryid")
            if err == nil {
                if debug {
                    fmt.Printf("Obtained Key %s\n", key)
                }
                awsConfigProfiles = append(awsConfigProfiles, key.String())
            } else {
                fmt.Printf("No registry ids were found for profiles %s \n", key.String())
            }
        } else {
            fmt.Printf("No profile was found for %s \n", profile)
        }
    }

    loginCmd := exec.Command("aws", "ecr", "get-login", "--region", "us-east-1", fmt.Sprintf("--registry-ids=%s", strings.Join(awsConfigProfiles, " ")))

    loginOutResult, err := loginCmd.Output()
    must(err)

    if debug {
        fmt.Println("Executing \"" + string(loginOutResult)+"\"")
    }
    awsDockerLoginCmd := exec.Command(string(loginOutResult))
    awsDockerLoginCmd.Run()
}

func printHelp() {
    fmt.Println("Usage: aws-docker COMMAND [options]\nCommands:\n\tlogin [profile] Log into given AWS ECR registry profile\n\nIn order to use login command edit ~/.aws/config by adding ```ecrregistryid``` to the given profile")
}

func printLoginHelp () {
    fmt.Println("Usage: aws-docker login [profile]\n\nLog into given AWS ECR registry profile\n\nIn order to use login command edit ~/.aws/config by adding ```ecrregistryid``` to the given profile")
}
func must(err error) {
    if err != nil {
        printHelp()
        os.Exit(2)
    }
}


