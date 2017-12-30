package endeca

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/JohnRoach/cartridgemapper/utils"
	"github.com/magiconair/properties"
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

type EndecaRules struct {
	cartridgeID string
	rules       []string
}

type SharedContent struct {
	XMLName       xml.Name
	ContentItem   []byte          `xml:",innerxml"`
	SharedContent []SharedContent `xml:",any"`
}

// ContentTemplate is a cartridge definition defined in a template
type ContentTemplate struct {
	ID          string `xml:"id,attr"`
	Description string `xml:"Description"`
}

func MapCartridges(basePath string, DisableColor bool, Debug bool) []Cartridge {
	var cartridges []Cartridge
	var endecaRules []EndecaRules

	cartridgeList := getCartridgePaths(basePath+"/templates", DisableColor, Debug)

	endecaRules = getTemplateEndecaRules(basePath, DisableColor, Debug)

	for _, cartridge := range cartridgeList {
		// Should be getting descriptions from XML first than accordingly from property files

		templateID, templateDescription, err := getTemplateData(cartridge, basePath+"/templates", DisableColor, Debug)
		if err != nil {
			utils.DisplayError("Couldn't read cartridge "+cartridge, err, DisableColor)
		} else {
			var cartridgeEndecaRules []string
			for _, endecaRule := range endecaRules {
				if endecaRule.cartridgeID == cartridge {
					cartridgeEndecaRules = endecaRule.rules
				}
			}
			var newCartridge = Cartridge{
				id:          templateID,
				description: templateDescription,
				rules:       cartridgeEndecaRules,
			}
			getCartridgeSitePageUsage(basePath, newCartridge, DisableColor, Debug)
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

func getCartridgeSitePageUsage(basePath string, cartridge Cartridge, DisableColor bool, Debug bool) Cartridge {
	var endecaSitePath = basePath + "/pages"
	utils.DisplayInfo("Starting Endeca template site and page usage scan...", DisableColor)
	err := filepath.Walk(endecaSitePath, func(path string, f os.FileInfo, walkError error) error {
		if strings.Contains(path, "content.xml") {
			xmlFile, xmlErr := os.Open(path)
			if xmlErr != nil {
				utils.DisplayError("Couldn't read XML for site and page scan at path "+path, xmlErr, DisableColor)
			} else {
				b, _ := ioutil.ReadAll(xmlFile)
				buf := bytes.NewBuffer(b)
				dec := xml.NewDecoder(buf)
				var n SharedContent
				xmlReadErr := dec.Decode(&n)
				if xmlReadErr != nil {
					panic(xmlReadErr)
				}
				walk([]SharedContent{n}, func(n SharedContent) bool {
					if n.XMLName.Local == "TemplateId" {
						cartridgeName := string(n.ContentItem)
						if cartridgeName == cartridge.id {
							utils.DisplayDebug("Found template in "+path, Debug, DisableColor)
						}
					}
					return true
				})
			}
			xmlFile.Close()
		}
		return walkError
	})

	if err != nil {
		utils.DisplayError("Could not walk through site path", err, DisableColor)
	}
	return cartridge
}

func getTemplateEndecaRules(basePath string, DisableColor bool, Debug bool) []EndecaRules {
	var endecaRules []EndecaRules
	var endecaRulesPath = basePath + "/content"
	utils.DisplayInfo("Starting Endeca shared content scan.", DisableColor)
	err := filepath.Walk(endecaRulesPath, func(path string, f os.FileInfo, err error) error {
		if strings.Contains(path, "content.xml") {
			xmlFile, xmlErr := os.Open(path)

			if xmlErr != nil {
				utils.DisplayError("Error opening file: "+path, xmlErr, DisableColor)
			} else {
				b, _ := ioutil.ReadAll(xmlFile)
				buf := bytes.NewBuffer(b)
				dec := xml.NewDecoder(buf)
				var n SharedContent
				xmlReadErr := dec.Decode(&n)
				if xmlReadErr != nil {
					panic(xmlReadErr)
				}

				walk([]SharedContent{n}, func(n SharedContent) bool {
					if n.XMLName.Local == "TemplateId" {
						cartridgeName := string(n.ContentItem)
						var found bool
						var endecaRulePath = strings.Replace(path, ".remove_me/content/", "", -1)
						endecaRulePath = strings.Replace(endecaRulePath, "/content.xml", "", -1)
						for index, endecaRule := range endecaRules {
							if endecaRule.cartridgeID == cartridgeName {
								found = true
								rulesInEndecaRule := endecaRule.rules
								rulesInEndecaRule = append(rulesInEndecaRule, endecaRulePath)
								endecaRules[index] = EndecaRules{
									cartridgeID: endecaRule.cartridgeID,
									rules:       rulesInEndecaRule,
								}
							}
						}
						if !found {
							var paths []string
							endecaRules = append(endecaRules, EndecaRules{
								cartridgeID: cartridgeName,
								rules:       append(paths, endecaRulePath),
							})
						}
					}
					return true
				})
			}
			xmlFile.Close()
		}
		return err
	})
	if err != nil {
		utils.DisplayError("Could not scan content directory.", err, DisableColor)
	}
	utils.DisplayInfo("Finished scanning Endeca rules.", DisableColor)
	return endecaRules
}

func walk(nodes []SharedContent, f func(SharedContent) bool) {
	for _, n := range nodes {
		if f(n) {
			walk(n.SharedContent, f)
		}
	}
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

// ID returns the ID of the cartridge
func (f Cartridge) ID() string {
	return f.id
}

// Description returns the description of the cartridge
func (f Cartridge) Description() string {
	return f.description
}

func (f Cartridge) Pages() []string {
	return f.pages
}

func (f Cartridge) Path() string {
	return f.path
}

func (f Cartridge) Rules() []string {
	return f.rules
}
