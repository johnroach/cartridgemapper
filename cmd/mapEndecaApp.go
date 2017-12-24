package cmd

import (
	//"fmt"
	"github.com/JohnRoach/cartridgemapper/endeca"
	"github.com/JohnRoach/cartridgemapper/utils"
	"github.com/spf13/cobra"
	"os"
)

// mapEndecaAppCmd represents the mapEndecaApp command
var mapEndecaAppCmd = &cobra.Command{
	Use:   "mapEndecaApp [path to application export zip]",
	Short: "mapEndecaApp maps the Endeca cartridges used in an Endeca Application",
	Long: `mapEndecaApp maps the Endeca cartridges used in an Endeca Application.
The cartridges will be go through validation.
For example:
    cartridgemapp mapEndecaApp /full/path/to/endeca/exported/Application.zip
    cartridgemapp mapEndecaApp /full/path/to/endeca/exported/Application.zip --output json
`,
	Example: "cartridgemapp mapEndecaApp /full/path/to/endeca/exported/Application.zip --debug",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var endecaAppPath string = args[0]
		mapEndecaApp(endecaAppPath)
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

func mapEndecaApp(endecaAppPath string) {
	dirError := os.MkdirAll(".remove_me", os.ModePerm)
	if dirError == nil {
		_, error := utils.Unzip(endecaAppPath, ".remove_me")
		if error == nil {
			utils.DisplayInfo("Unzipped exported endeca application file...", DisableColor)

			endeca.MapCartridges(".remove_me", DisableColor, Debug)
			//fmt.Println(cartridges)
			removeDirectory(".remove_me")
			utils.DisplayInfo("Removed temporary directory...", DisableColor)
		} else {
			utils.DisplayError("Couldn't unzip file.", error, DisableColor)
			os.RemoveAll(".remove_me")
		}
	} else {
		utils.DisplayError("Couldn't create test directory.", dirError, DisableColor)
	}
}

func removeDirectory(path string) error {
	removeError := os.RemoveAll(".remove_me")
	if removeError != nil {
		utils.DisplayError("Couldn't remove .remove_me temp folder.", removeError, DisableColor)
	}
	return removeError
}
