package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
)

type ValidationType int64
const (
	Path ValidationType = iota
	KeepPassword 
)

func isValidStr(str string, validationType ValidationType) bool {
	switch validationType {
	case Path:
		//TODO handle path validation
		return true 
	case KeepPassword:
		str = strings.ToLower(str)
		return str == "y" || str == "n"
	}
	return false
}

type Config struct {
	DBPath string `toml:"db-path"`
	KeepDBPassword string `toml:"keep-db-password"`
	CurrentDB string `toml:"current-db"`
	CurrentEnv string `toml:current-env`
}

func CreateConfigFile() (config *Config, err error) {
	homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    fullConfigPath := path.Join(homeDir, DefaultConfigLocation, AppName)
	err = os.MkdirAll(fullConfigPath, os.ModePerm)
    if err != nil {
        log.Fatal(err)
    }
    viper.SetConfigType("toml")
    viper.AddConfigPath(fullConfigPath)
    viper.SafeWriteConfig()

    viper.WatchConfig()
	// TODO set config first timme
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Hey! It seems thats your first time here")
	fmt.Println("Let's do some configs!")
	var input string

	// set DBs location
	for (!isValidStr(input, Path)) {
		fmt.Println("Where whould you like to save your DB(s)? (leave empty for default location)")
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal("Something wrong with your path. try again")
		}
		input = strings.TrimSuffix(input, "\n")
	}
	viper.Set(K_KEEPASS_DB_PATH, input)
	input = ""
	// set whether or not to keep password as environment variable in the current session
	for (isValidStr(input, KeepPassword)) {
		fmt.Println("Whould you like that your DB password will be saved as environment variable in session? (y/n)")
		input, err = reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		input = strings.TrimSuffix(input, "\n")
	}
	viper.Set(K_KEEP_DB_PASSWORD, input)

    viper.WriteConfig()
	err = viper.Unmarshal(&config)
	if err != nil {
		return nil, err
	}
	return config, nil
}


func (Config) LoadConfigFile() (config *Config, err error) {
	homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    fullConfigPath := path.Join(homeDir, DefaultConfigLocation, AppName)
	viper.AddConfigPath(fullConfigPath)
	viper.SetConfigName(AppName)
	viper.SetConfigType(DefaultConfigType)

	err = viper.ReadInConfig()
    if err != nil {
        return
    }

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatal(err)
	}
    return config, nil
}