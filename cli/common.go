package cli

import (
    "fmt"
    "os"
    "os/user"

    "github.com/go-ini/ini"
)

var (
    debug = ToggleDebug()
    config = getAwsConfig()
)

type CommonOpts struct {
    verbose *bool
}

func getAwsConfig () *ini.File {
    usr, _ := user.Current()
    //must(err)
    if debug {
        fmt.Printf("From user directory %s\n", usr.HomeDir)
    }

    c, _ := ini.LooseLoad(usr.HomeDir+"/.aws/config")
    return c

}


// GetProfileSection returns the config chuck of the give profile
func GetProfileSection (profile string) (*ini.Section, error) {
    section, err := config.GetSection(fmt.Sprintf("profile %s", profile))
    return section, err
}

// ToggleDebug turns on stdout if argument flag -v is used
func ToggleDebug() bool {
    for _, element := range os.Args {
        if element == "-v" {
            return true
        }
    }
    return false
}
