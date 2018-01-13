package endeca

import (
	"bytes"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"../utils"
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

//Rules is a cartridgeID and rules definition
type Rules struct {
	cartridgeID string
	rules       []string
}

// SharedContent is a generic struct used for XML walking
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

// MapCartridges allows one to map all cartridges and usages for a given base path
func MapCartridges(basePath string, DisableColor bool, Debug bool) []Cartridge {
	var cartridges []Cartridge
	var endecaRules []Rules

	cartridgeList := getCartridgePaths(basePath+"/templates", DisableColor, Debug)

	endecaRules = getTemplateRules(basePath, DisableColor, Debug)

	for _, cartridge := range cartridgeList {
		// Should be getting descriptions from XML first than accordingly from property files

		templateID, templateDescription, err := getTemplateData(cartridge, basePath+"/templates", DisableColor, Debug)
		if err != nil {
			utils.DisplayError("Couldn't read cartridge "+cartridge, err, DisableColor)
		} else {
			var cartridgeRules []string
			for _, endecaRule := range endecaRules {
				if endecaRule.cartridgeID == cartridge {
					cartridgeRules = endecaRule.rules
				}
			}
			var newCartridge = Cartridge{
				id:          templateID,
				description: templateDescription,
				rules:       cartridgeRules,
			}
			newCartridge = getCartridgeSitePageUsage(basePath, newCartridge, DisableColor, Debug)
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
	utils.DisplayDebug("Starting Endeca template site and page usage scan for "+cartridge.id, Debug, DisableColor)
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
					var pathInfo = strings.Split(path, "/")
					var siteName = pathInfo[2]
					var pageName = strings.Join(pathInfo[3:len(pathInfo)-1], "/")

					if n.XMLName.Local == "TemplateId" {
						cartridgeName := string(n.ContentItem)
						if cartridgeName == cartridge.id {
							utils.DisplayDebug("Found template in "+path+" which means it was in site "+siteName, Debug, DisableColor)
							cartridge = cartridge.addSite(siteName)
							cartridge = cartridge.addPage(pageName)
						}
					}
					if n.XMLName.Local == "String" {
						var stringValue = string(n.ContentItem)
						var oldCartridgeRules = cartridge.rules
						for _, oldCartridgeEndecaRule := range oldCartridgeRules {
							if "/content/"+oldCartridgeEndecaRule == stringValue {
								cartridge = cartridge.addSite(siteName)
								cartridge = cartridge.addPage(pageName)
							}
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

func getTemplateRules(basePath string, DisableColor bool, Debug bool) []Rules {
	var endecaRules []Rules
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
								endecaRules[index] = Rules{
									cartridgeID: endecaRule.cartridgeID,
									rules:       rulesInEndecaRule,
								}
							}
						}
						if !found {
							var paths []string
							endecaRules = append(endecaRules, Rules{
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
		utils.DisplayWarning("Cartridge ID not defined in template. Cartridge name: "+templateName, DisableColor)
	} else {
		templateID = contentTemplate.ID
	}

	if contentTemplate.Description == "${template.description}" {
		templateDescription = getDescriptionFromProperty(basePath+"/"+templateName, DisableColor, Debug)
	} else if strings.TrimSpace(contentTemplate.Description) == "" {
		templateDescription = "No description provided."
		utils.DisplayWarning("Cartridge definition not defined in template. Cartridge name: "+templateName, DisableColor)
	} else {
		templateDescription = contentTemplate.Description
	}

	return templateID, templateDescription, templateError
}

// GetID returns the ID of the cartridge
func (f Cartridge) GetID() string {
	return f.id
}

// GetDescription returns the description of the cartridge
func (f Cartridge) GetDescription() string {
	return f.description
}

// GetPages returns the pages for a given cartridge
func (f Cartridge) GetPages() []string {
	return f.pages
}

// GetSites returns the sites for a given cartridge
func (f Cartridge) GetSites() []string {
	return f.sites
}

// GetPath returns the path for a given cartridge
func (f Cartridge) GetPath() string {
	return f.path
}

// GetRules returns a list of rules for a given cartridge
func (f Cartridge) GetRules() []string {
	return f.rules
}

func (f Cartridge) addPage(page string) Cartridge {
	var foundPage bool
	for _, oldCartridgePage := range f.pages {
		if page == oldCartridgePage {
			foundPage = true
		}
	}
	if !foundPage {
		f.pages = append(f.pages, page)
	}
	return f
}

func (f Cartridge) addSite(site string) Cartridge {
	var foundSite bool
	for _, oldCartridgeSite := range f.sites {
		if site == oldCartridgeSite {
			foundSite = true
		}
	}
	if !foundSite {
		f.sites = append(f.sites, site)
	}
	return f
}
