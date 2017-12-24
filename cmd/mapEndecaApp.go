package cmd

import (
	"archive/zip"
	"fmt"
	"github.com/fatih/color"
	"github.com/magiconair/properties"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

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
			_, error := Unzip(endecaAppPath, ".remove_me")
			if error == nil {
				color.Green("Unzipped exported endeca application file...")

				cartridges := mapCartridges(".remove_me")

				fmt.Println(cartridges)

				//removeDirectory(".remove_me")
				color.Green("Removed temporary directory...")
			} else {
				color.Red("Couldn't unzip file.")
				log.Fatal(error)
				os.RemoveAll(".remove_me")
			}
		} else {
			color.Red("Couldn't create test directory.")
			log.Fatal(dirError)
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
		var newCartridge = Cartridge{
			id:          cartridge,
			description: getDescription(basePath + "/templates/" + cartridge),
		}
		cartridges = append(cartridges, newCartridge)
	}

	return cartridges
}

func getDescription(path string) string {
	var description string
	p, err := properties.LoadFile(path+"/locales/Resources_en.properties", properties.UTF8)
	if err == nil {
		description = p.GetString("template.description", "")
		if description == "" {
			color.Red("Description doesn't exist for template in " + path)
			return "No description specified."
		}
	} else {
		color.Red("locale file for description doesn't exist for template in " + path)
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
			color.Red("Couldn't read file stats.")
			log.Fatal(err)
		}
	}

	return cartridgePaths
}

func removeDirectory(path string) error {
	removeError := os.RemoveAll(".remove_me")
	if removeError != nil {
		color.Red("Couldn't remove .remove_me temp folder.")
	}
	return removeError
}

// Unzip will un-compress a zip archive,
// moving all files and folders to an output directory
func Unzip(src, dest string) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return filenames, err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return filenames, err
			}

		}
	}
	return filenames, nil
}
