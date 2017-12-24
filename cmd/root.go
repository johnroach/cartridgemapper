package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// AppVersion is the application version
var AppVersion = "1.0.0"

// cfgFile is the config file path
var cfgFile string

// Debug allows printing of bunch of DEBUG values in log
var Debug bool

// DisableColor disables the color for a given run
var DisableColor bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cartridgemapper",
	Short: "Endeca cartridge mapper maps the Endeca application cartridge usage.",
	Long: `Endeca cartridge mapper maps the Endeca application cartridge usage.
This is very useful in understanding how cartridges are used and which cartridges
are available. One would use this tool to point to a given Endeca application and
get an ouput of a certain format.

For example:
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication --output json
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.AddCommand(versionCmd)

	// Persistent flags. Flags that will live for all subcommands.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cartridgemapper.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "adding debug to the logging")
	rootCmd.PersistentFlags().BoolVarP(&DisableColor, "disable-color", "", false, "disable color for logging output")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cartridgemapper" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cartridgemapper")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of cartridgemapper",
	Long:  `All software has versions. This is cartridgemappers's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Cartridge mapper v" + printVersion())
	},
}

// printVersion prints out version which is defined globally
func printVersion() string {
	if Debug {
		fmt.Println("Getting the app version")
	}
	return AppVersion
}
