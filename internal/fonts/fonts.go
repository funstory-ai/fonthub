package fonts

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// FontMetadata represents the structure of each font's metadata
type FontMetadata struct {
	Name       string   `json:"name"`
	License    string   `json:"license"`
	Version    string   `json:"version"`
	Categories []string `json:"categories"`
	Widths     []string `json:"widths"`
	Weights    []string `json:"weights"`
	Styles     []string `json:"styles"`
	Languages  []string `json:"languages"`
	Source     string   `json:"source"`
	FontsFiles struct {
		Link string `json:"link"`
		Path string `json:"path"`
	} `json:"fontsFiles"`
}

// FontDatabase holds all the font metadata
type FontDatabase struct {
	Fonts         map[string]FontMetadata
	WidthIndex    map[string]map[string]struct{} // width -> font names
	WeightIndex   map[string]map[string]struct{} // weight -> font names
	StyleIndex    map[string]map[string]struct{} // style -> font names
	LanguageIndex map[string]map[string]struct{} // language -> font names
}

var GlobalFontDB FontDatabase

func init() {
	logrus.Info("Initializing fonts database")
	GlobalFontDB = FontDatabase{
		Fonts:         make(map[string]FontMetadata),
		WidthIndex:    make(map[string]map[string]struct{}),
		WeightIndex:   make(map[string]map[string]struct{}),
		StyleIndex:    make(map[string]map[string]struct{}),
		LanguageIndex: make(map[string]map[string]struct{}),
	}
	BuildFontsDatabase()
}

func BuildFontsDatabase() error {
	// Define the metadata directory path
	metadataDir := "metadataset"

	// Walk through all JSON files in the metadata directory
	err := filepath.Walk(metadataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Error("Error walking through metadata directory:", err)
			return err
		}
		logrus.Info("loading metadata file:", path)

		// Skip if not a JSON file
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// Read the JSON file
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// Parse the JSON data
			var metadata FontMetadata
			if err := json.Unmarshal(data, &metadata); err != nil {
				return err
			}

			// Add to the database using the name as the key
			GlobalFontDB.Fonts[metadata.Name] = metadata

			// Build indices
			for _, width := range metadata.Widths {
				if GlobalFontDB.WidthIndex[width] == nil {
					GlobalFontDB.WidthIndex[width] = make(map[string]struct{})
				}
				GlobalFontDB.WidthIndex[width][metadata.Name] = struct{}{}
			}

			for _, weight := range metadata.Weights {
				if GlobalFontDB.WeightIndex[weight] == nil {
					GlobalFontDB.WeightIndex[weight] = make(map[string]struct{})
				}
				GlobalFontDB.WeightIndex[weight][metadata.Name] = struct{}{}
			}

			for _, style := range metadata.Styles {
				if GlobalFontDB.StyleIndex[style] == nil {
					GlobalFontDB.StyleIndex[style] = make(map[string]struct{})
				}
				GlobalFontDB.StyleIndex[style][metadata.Name] = struct{}{}
			}

			for _, lang := range metadata.Languages {
				if GlobalFontDB.LanguageIndex[lang] == nil {
					GlobalFontDB.LanguageIndex[lang] = make(map[string]struct{})
				}
				GlobalFontDB.LanguageIndex[lang][metadata.Name] = struct{}{}
			}
		}

		return nil
	})

	return err
}

// GetAllFonts returns all fonts in the database
func GetAllFonts() []FontMetadata {
	fonts := make([]FontMetadata, 0, len(GlobalFontDB.Fonts))
	for _, font := range GlobalFontDB.Fonts {
		fonts = append(fonts, font)
	}
	return fonts
}

// GetFontByFamily returns a specific font by its family name
func GetFontByFamily(family string) (FontMetadata, bool) {
	font, exists := GlobalFontDB.Fonts[family]
	return font, exists
}

// GetFontsByCategory returns all fonts in a specific category
func GetFontsByCategory(category string) []FontMetadata {
	var fonts []FontMetadata
	for _, font := range GlobalFontDB.Fonts {
		for _, cat := range font.Categories {
			if cat == category {
				fonts = append(fonts, font)
				break
			}
		}
	}
	return fonts
}

// FontSelector represents search criteria for fonts
type FontSelector struct {
	Widths    []string
	Weights   []string
	Styles    []string
	Languages []string
}

// GetFontsBySelector returns fonts matching all specified criteria
func GetFontsBySelector(selector FontSelector) []FontMetadata {
	matches := make(map[string]struct{})
	var firstSet bool

	// Check widths
	if len(selector.Widths) > 0 {
		widthMatches := make(map[string]struct{})
		for _, width := range selector.Widths {
			if fonts, ok := GlobalFontDB.WidthIndex[width]; ok {
				// Union of all fonts matching any of the width criteria
				for name := range fonts {
					widthMatches[name] = struct{}{}
				}
			}
		}
		if !firstSet {
			matches = widthMatches
			firstSet = true
		} else {
			// Intersect with existing matches
			for name := range matches {
				if _, ok := widthMatches[name]; !ok {
					delete(matches, name)
				}
			}
		}
	}

	// Check weights
	if len(selector.Weights) > 0 {
		weightMatches := make(map[string]struct{})
		for _, weight := range selector.Weights {
			if fonts, ok := GlobalFontDB.WeightIndex[weight]; ok {
				for name := range fonts {
					weightMatches[name] = struct{}{}
				}
			}
		}
		if !firstSet {
			matches = weightMatches
			firstSet = true
		} else {
			for name := range matches {
				if _, ok := weightMatches[name]; !ok {
					delete(matches, name)
				}
			}
		}
	}

	// Check styles
	if len(selector.Styles) > 0 {
		styleMatches := make(map[string]struct{})
		for _, style := range selector.Styles {
			if fonts, ok := GlobalFontDB.StyleIndex[style]; ok {
				for name := range fonts {
					styleMatches[name] = struct{}{}
				}
			}
		}
		if !firstSet {
			matches = styleMatches
			firstSet = true
		} else {
			for name := range matches {
				if _, ok := styleMatches[name]; !ok {
					delete(matches, name)
				}
			}
		}
	}

	// Check languages
	if len(selector.Languages) > 0 {
		langMatches := make(map[string]struct{})
		for _, lang := range selector.Languages {
			if fonts, ok := GlobalFontDB.LanguageIndex[lang]; ok {
				for name := range fonts {
					langMatches[name] = struct{}{}
				}
			}
		}
		if !firstSet {
			matches = langMatches
			firstSet = true
		} else {
			for name := range matches {
				if _, ok := langMatches[name]; !ok {
					delete(matches, name)
				}
			}
		}
	}

	// If no criteria were specified, return all fonts
	if !firstSet {
		for name := range GlobalFontDB.Fonts {
			matches[name] = struct{}{}
		}
	}

	// Convert matches to slice of FontMetadata
	result := make([]FontMetadata, 0, len(matches))
	for name := range matches {
		result = append(result, GlobalFontDB.Fonts[name])
	}
	return result
}
