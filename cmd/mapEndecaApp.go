package cmd

import (
	"github.com/JohnRoach/cartridgemapper/endeca"
	"github.com/JohnRoach/cartridgemapper/utils"
	"github.com/spf13/cobra"
	"html/template"
	"os"
)

var outputPath string
var templatePath string

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
	mapEndecaAppCmd.Flags().StringVarP(&outputPath, "outputPath", "o", ".", "Output path for the endeca map")
	mapEndecaAppCmd.Flags().StringVarP(&templatePath, "templatePath", "", "", "Template path for the endeca map")
}

func mapEndecaApp(endecaAppPath string) {
	dirError := os.MkdirAll(".remove_me", os.ModePerm)
	if dirError == nil {
		_, error := utils.Unzip(endecaAppPath, ".remove_me")
		if error == nil {
			utils.DisplayInfo("Unzipped exported endeca application file...", DisableColor)

			var cartridges []endeca.Cartridge = endeca.MapCartridges(".remove_me", DisableColor, Debug)
			cartridgeOutputHTML(cartridges, outputPath, DisableColor, Debug)
			//removeDirectory(".remove_me")
			utils.DisplayInfo("Removed temporary directory...", DisableColor)
		} else {
			utils.DisplayError("Couldn't unzip file.", error, DisableColor)
			removeDirectory(".remove_me")
		}
	} else {
		utils.DisplayError("Couldn't create test directory.", dirError, DisableColor)
	}
}

func cartridgeOutputHTML(cartridges []endeca.Cartridge, outputPath string, DisableColor bool, Debug bool) {
	t, parseFileError := template.ParseFiles("all_cartridges.html")
	if parseFileError != nil {
		utils.DisplayError("Had a parsefile error", parseFileError, DisableColor)
		return
	}

	fo, createOutputError := os.Create(outputPath + "/index.html")
	if createOutputError != nil {
		utils.DisplayError("Failed to create output", createOutputError, DisableColor)
		return
	}

	templateExecuteError := t.ExecuteTemplate(fo, "all_cartridges.html", cartridges)
	fo.Close()
	if templateExecuteError != nil {
		utils.DisplayError("Had a template execute error", templateExecuteError, DisableColor)
		return
	}
	utils.DisplayInfo("Created index.html file at "+outputPath+"/index.html", DisableColor)
}

func removeDirectory(path string) error {
	removeError := os.RemoveAll(".remove_me")
	if removeError != nil {
		utils.DisplayError("Couldn't remove .remove_me temp folder.", removeError, DisableColor)
	}
	return removeError
}
