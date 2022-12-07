package main

import (
	"fmt"
	"log"
	"os"
	"path"
    
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)


func main() {
    loadConfig()
    app := &cli.App{
        Commands: []*cli.Command{
            {
                Category: "Database",
                Name: "Create a new Database",
                Aliases: []string{"c"},
                Action: func(cCtx *cli.Context) error {
                    CreateLocalDB(cCtx.Args().First())
                    return nil
                },
            },
            
        },
        Action: func(cCtx *cli.Context) error {

            fmt.Printf("Hello %q", cCtx.Args().Get(0))

            return nil
        },
    }


    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

func loadConfig() {
    const DefaultConfigLocation = DefaultConfigLocation
    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    fullConfigPath := path.Join(homeDir, DefaultConfigLocation, AppName)
    fmt.Println(fullConfigPath)
    // init config file
    err = os.MkdirAll(fullConfigPath, os.ModePerm)
    if err != nil {
        log.Fatal(err)
    }
    viper.SetConfigType("toml")
    viper.AddConfigPath(fullConfigPath)
    viper.SafeWriteConfig()
    // if err := viper.SafeWriteConfigAs(fullConfigPath); err != nil {
    //     if os.IsNotExist(err) {
    //         err = viper.WriteConfigAs(fullConfigPath)
    //         if err != nil {
    //             log.Fatal(err)
    //         }
    //     }
    // }
    viper.WatchConfig()
    viper.Set("heftse", "menemtse")
    viper.WriteConfig()

    // config, err := viper.SafeWriteConfig(fullConfigPath)
    // if (err != nil) {
    //     log.Fatal(err)
    // }
    // if not exists create 
}