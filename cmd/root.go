package cmd

import (
	"fmt"
	"os"

	"github.com/mikesmitty/grb/grb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	build       string
	cfgFile     string
	downloadDir string
	patchDir    string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "grb",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := checkConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}

		buildVersion, err := grb.GetVersion(build)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}

		err = grb.GetTarball(buildVersion.URL, downloadDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			os.Exit(1)
		}
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.grb.yaml)")
	RootCmd.PersistentFlags().StringP("build", "b", "", "version of go to build")
	RootCmd.PersistentFlags().StringP("download", "d", "", "download directory for tarballs")
	RootCmd.PersistentFlags().StringP("patch", "p", "", "patch directory")
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".grb")  // name of config file (without extension)
	viper.AddConfigPath("$HOME") // adding home directory as first search path
	viper.AutomaticEnv()         // read in environment variables that match
	viper.ReadInConfig()
}

func checkConfig() error {
	build = viper.GetString("build")
	downloadDir = viper.GetString("download")
	patchDir = viper.GetString("patch")

	if downloadDir == "" {
		return fmt.Errorf("error: download directory not set")
	}
	if patchDir == "" {
		return fmt.Errorf("error: patch directory not set")
	}

	return nil
}
