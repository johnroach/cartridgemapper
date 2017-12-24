package cmd

import (
	//"fmt"
	"github.com/JohnRoach/cartridgemapper/utils"
	"github.com/magiconair/properties"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
)

// Cartridge is a combination of all relevant data for a cartridge
type Cartridge struct {
	path        string
	description string
	id          string
	sites       []string
	pages       []string
	rules       []string
}

// mapEndecaAppCmd represents the mapEndecaApp command
var mapEndecaAppCmd = &cobra.Command{
	Use:   "mapEndecaApp",
	Short: "mapEndecaApp maps the Endeca cartridges used in an Endeca Application",
	Long: `mapEndecaApp maps the Endeca cartridges used in an Endeca Application.
The cartridges will be go through validation.
For example:
    cartridgemapp mapEndecaApp /full/path/to/endeca/exported/Application.zip
    cartridgemapp mapEndecaApp /full/path/to/endeca/exported/Application.zip --output json
`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var endecaAppPath string = args[0]
		dirError := os.MkdirAll(".remove_me", os.ModePerm)
		if dirError == nil {
			_, error := utils.Unzip(endecaAppPath, ".remove_me")
			if error == nil {
				displayInfo("Unzipped exported endeca application file...")

				mapCartridges(".remove_me")
				//fmt.Println(cartridges)
				removeDirectory(".remove_me")
				displayInfo("Removed temporary directory...")
			} else {
				displayError("Couldn't unzip file.", error)
				os.RemoveAll(".remove_me")
			}
		} else {
			displayError("Couldn't create test directory.", dirError)
		}
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

func mapCartridges(basePath string) []Cartridge {
	var cartridges []Cartridge

	cartridgeList := getCartridgePaths(basePath + "/templates")

	for _, cartridge := range cartridgeList {
		// Should be getting descriptions from XML first than accordingly from property files

		templateID, templateDescription, err := getTemplateData(cartridge)
		if err != nil {
			displayError("Couldn't read cartridge "+cartridge, err)
		} else {
			var newCartridge = Cartridge{
				id:          templateID,
				description: templateDescription,
			}
			cartridges = append(cartridges, newCartridge)
		}
	}

	return cartridges
}

func getTemplateData(templateName string) (string, string, error) {
	var templateDescription string
	var templateID string
	var templateError error
	displayDebug("Starting work on " + templateName)
	return templateID, templateDescription, templateError
}

func getDescriptionFromProperty(path string) string {
	var description string
	p, err := properties.LoadFile(path+"/locales/Resources_en.properties", properties.UTF8)
	if err == nil {
		description = p.GetString("template.description", "")
		if description == "" {
			displayError("Description doesn't exist for template in "+path, nil)
			return "No description specified."
		}
	} else {
		displayError("locale file for description doesn't exist for template in "+path, err)
		return "No description specified."
	}
	return description
}

// getCartridgePaths gets the cartridge paths
func getCartridgePaths(path string) []string {
	var cartridgePaths []string

	files, err := filepath.Glob(path + "/*")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err == nil {
			if fileInfo.IsDir() {
				_, cartridgeName := filepath.Split(file)
				cartridgePaths = append(cartridgePaths, cartridgeName)
			}
		} else {
			displayError("Couldn't read file stats.", err)
		}
	}

	return cartridgePaths
}

func removeDirectory(path string) error {
	removeError := os.RemoveAll(".remove_me")
	if removeError != nil {
		displayError("Couldn't remove .remove_me temp folder.", removeError)
	}
	return removeError
}

func displayError(message string, err error) {
	var errorMessage string
	if err != nil {
		errorMessage = err.Error()
	}
	utils.LogError(message+" "+errorMessage, DisableColor)
}

func displayDebug(message string) {
	if Debug {
		utils.LogDebug(message, DisableColor)
	}
}

func displayInfo(message string) {
	utils.LogInfo(message, DisableColor)
}
