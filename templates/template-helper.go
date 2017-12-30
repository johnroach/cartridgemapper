package templates

import (
	"html/template"
	"os"

	"github.com/JohnRoach/cartridgemapper/endeca"
	"github.com/JohnRoach/cartridgemapper/utils"
)

//CartridgeOutputHTML receives the cartridges and by using the IndexPage template
//in templates it produces a cool looking HTML page that can be used.
func CartridgeOutputHTML(cartridges []endeca.Cartridge, outputPath string, DisableColor bool, Debug bool) {
	t, parseFileError := template.New("IndexPage").Parse(IndexPage)
	if parseFileError != nil {
		utils.DisplayError("Had a parsefile error", parseFileError, DisableColor)
		return
	}

	fo, createOutputError := os.Create(outputPath + "/index.html")
	if createOutputError != nil {
		utils.DisplayError("Failed to create output", createOutputError, DisableColor)
		return
	}

	templateExecuteError := t.ExecuteTemplate(fo, "IndexPage", cartridges)
	fo.Close()
	if templateExecuteError != nil {
		utils.DisplayError("Had a template execute error", templateExecuteError, DisableColor)
		return
	}
	utils.DisplayInfo("Created index.html file at "+outputPath+"/index.html", DisableColor)
}
