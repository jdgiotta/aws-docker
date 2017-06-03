package login

import (
    "os/user"
    "fmt"
    "github.com/go-ini/ini"
    "os"
    "os/exec"
    "strings"
)

var debug = toggleDebug()

// Do perform the action of getting profile registry IDs and executes the login to ecr services
func Do() {
    usr, err := user.Current()
    must(err)
    if debug {
        fmt.Printf("From user directory %s\n", usr.HomeDir)
    }

    cfg, err := ini.LooseLoad(usr.HomeDir+"/.aws/config")

    must(err)

    if len(os.Args) < 3 {
        Help()
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

    getCmd := exec.Command("aws", "ecr", "get-login", "--region", "us-east-1", fmt.Sprintf("--registry-ids=%s", strings.Join(awsConfigProfiles, " ")))

    getOutResult, err := getCmd.Output()
    must(err)

    if debug {
        fmt.Println("Executing \"" + string(getOutResult)+"\"")
    }
    awsDockerLoginCmd := exec.Command(string(getOutResult))
    awsDockerLoginCmd.Run()
}


// Help prints the help dialog for this module
func Help() {
    fmt.Println("Usage: aws-docker login [profile]... \n\nLog into given AWS ECR registry profile(s)\n\nIn order to use login command edit ~/.aws/config by adding ```ecrregistryid``` to the given profile(s)")
}

func toggleDebug() bool {
    for _, element := range os.Args {
        if element == "-v" {
            return true
        }
    }
    return false
}


func must(err error) {
    if err != nil {
        Help()
        os.Exit(2)
    }
}