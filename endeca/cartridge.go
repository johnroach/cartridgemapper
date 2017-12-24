package endeca

import (
	"encoding/xml"
	"github.com/JohnRoach/cartridgemapper/utils"
	"github.com/magiconair/properties"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Cartridge is a combination of all relevant data for a cartridge
type Cartridge struct {
	// path is for path to cartridge
	path string
	// description for given cartridge
	description string
	// template id for cartridge
	id string
	// sites in which the cartridge is used
	sites []string
	// pages in which the cartridge is used
	pages []string
	// rules in which the cartridge is used
	rules []string
}

// ContentTemplate is a cartridge definition defined in a template
type ContentTemplate struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"Description"`
}

func MapCartridges(basePath string, DisableColor bool, Debug bool) []Cartridge {
	var cartridges []Cartridge

	cartridgeList := getCartridgePaths(basePath+"/templates", DisableColor, Debug)

	for _, cartridge := range cartridgeList {
		// Should be getting descriptions from XML first than accordingly from property files

		templateID, templateDescription, err := getTemplateData(cartridge, basePath+"/templates", DisableColor, Debug)
		if err != nil {
			utils.DisplayError("Couldn't read cartridge "+cartridge, err, DisableColor)
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

func getDescriptionFromProperty(path string, DisableColor bool, Debug bool) string {
	var description string
	p, err := properties.LoadFile(path+"/locales/Resources_en.properties", properties.UTF8)
	if err == nil {
		description = p.GetString("template.description", "")
		if description == "" {
			utils.DisplayError("Description doesn't exist for template in "+path, nil, DisableColor)
			return "No description specified."
		}
	} else {
		utils.DisplayError("locale file for description doesn't exist for template in "+path, err, DisableColor)
		return "No description specified."
	}
	return description
}

// getCartridgePaths gets the cartridge paths
func getCartridgePaths(path string, DisableColor bool, Debug bool) []string {
	var cartridgePaths []string

	files, err := filepath.Glob(path + "/*")
	if err != nil {
		utils.DisplayError("Could list directory.", err, DisableColor)
	}

	for _, file := range files {
		fileInfo, err := os.Stat(file)
		if err == nil {
			if fileInfo.IsDir() {
				_, cartridgeName := filepath.Split(file)
				cartridgePaths = append(cartridgePaths, cartridgeName)
			}
		} else {
			utils.DisplayError("Couldn't read file stats.", err, DisableColor)
		}
	}

	return cartridgePaths
}

func getTemplateData(templateName string, basePath string, DisableColor bool, Debug bool) (string, string, error) {
	var templateDescription string
	var templateID string
	var templateError error
	utils.DisplayDebug("Starting work on "+templateName, Debug, DisableColor)

	xmlFile, err := os.Open(basePath + "/" + templateName + "/template.xml")
	if err != nil {
		utils.DisplayError("Error opening file:", err, DisableColor)
		return templateName, "No description.", err
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	var contentTemplate ContentTemplate
	xml.Unmarshal(b, &contentTemplate)

	if contentTemplate.ID == "" {
		templateID = templateName
		utils.LogWarning("Cartridge ID not defined in template. Cartridge name: "+templateName, DisableColor)
	} else {
		templateID = contentTemplate.ID
	}

	if contentTemplate.Description == "${template.description}" {
		templateDescription = getDescriptionFromProperty(basePath+"/"+templateName, DisableColor, Debug)
	} else if strings.TrimSpace(contentTemplate.Description) == "" {
		templateDescription = "No description provided."
		utils.LogWarning("Cartridge definition not defined in template. Cartridge name: "+templateName, DisableColor)
	} else {
		templateDescription = contentTemplate.Description
	}

	return templateID, templateDescription, templateError
}
