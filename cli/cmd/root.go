package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/cyrildiagne/kuda/pkg/config"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"

	yaml "gopkg.in/yaml.v2"
)

var version = "dev"

var cfgFile string
var cfg config.UserConfig

// RootCmd is the main command.
var RootCmd = &cobra.Command{
	Use:     "kuda",
	Short:   "Kuda - Serverless APIs on remote GPUs",
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	RootCmd.Version = version
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", home+"/.kuda.yaml",
		"Configuration file.")
}

// initConfig reads in the config file.
func loadConfig() {
	// Check if config file exists.
	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		return
	}
	// Load config
	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Loaded config from", cfgFile)

		// Check access token
		i, err := strconv.ParseInt(strconv.Itoa(cfg.Provider.User.Token.ExpirationTime/1000), 10, 64)
		if err != nil {
			panic(err)
		}
		tm := time.Unix(i, 0)

		if tm.Before(time.Now()) {
			fmt.Println("Refreshing token...")
			refreshURL := cfg.Provider.AuthURL + "/refresh"
			refreshToken := cfg.Provider.User.Token.RefreshToken
			res, err := refreshAuthToken(refreshURL, refreshToken)
			if err != nil {
				panic(err)
			}

			cfg.Provider.User.Token.RefreshToken = res.RefreshToken
			cfg.Provider.User.Token.AccessToken = res.AccessToken
			cfg.Provider.User.Token.ExpirationTime = (int(time.Now().Unix()) + res.ExpiresIn) * 1000
			writeConfig(cfg)
		}
	}
}
