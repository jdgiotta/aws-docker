package main

import (
    "os"
    "github.com/jdgiotta/aws-docker/login"
    "fmt"
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
    fmt.Println("Usage: aws-docker COMMAND [options]\nCommands:\n\tlogin [profile]... Log into given AWS ECR registry profile(s)\n\nIn order to use login command edit ~/.aws/config by adding ```ecrregistryid``` to the given profile(s)")
}


