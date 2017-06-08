package main

import (
    "fmt"
    "os"

    "github.com/jdgiotta/aws-docker/cmd/login"
)

func main () {
    if len(os.Args) < 2 {
        help()
        os.Exit(2)
    }
    switch os.Args[1] {
    case "login":
        login.Do()
    default:
        help()
    }
}

func help() {
    fmt.Println("\nUsage: aws-docker COMMAND [options]\nCommands:\n\tlogin --profile [profile] Log into given AWS ECR registry profile\n\nIn order to use login command edit ~/.aws/config by adding ```ecrregistryid``` to the given profile then creates a docker login session")
}


