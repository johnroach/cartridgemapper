package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// mapEndecaAppCmd represents the mapEndecaApp command
var mapEndecaAppCmd = &cobra.Command{
	Use:   "mapEndecaApp",
	Short: "mapEndecaApp maps the Endeca cartridges used in an Endeca Application",
	Long: `mapEndecaApp maps the Endeca cartridges used in an Endeca Application.
The cartridges will be go through validation.
For example:
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication
    cartridgemapp mapEndecaApp /full/path/to/endeca/App/lication --output json
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mapEndecaApp called")
	},
}

func init() {
	rootCmd.AddCommand(mapEndecaAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mapEndecaAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mapEndecaAppCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
