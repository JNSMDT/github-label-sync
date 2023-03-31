package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

type GLSLabels struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

type TomlLabel struct {
	GLSLabels
	Scope string
}

type TomlConfig struct {
	ScopeSeperator string
	Labels         []TomlLabel
}

type GLSConfig []GLSLabels

// Convert the TomlConfigStruct to github-label-sync Format
func convertFormat(config TomlConfig) GLSConfig {
	seperator := config.ScopeSeperator
	tomlLabels := config.Labels

	var glsLabels []GLSLabels
	for _, tomlLabel := range tomlLabels {
		name := tomlLabel.Name
		color := tomlLabel.Color
		description := tomlLabel.Description
		scope := tomlLabel.Scope

		var newName string
		if scope != "" {
			newName = fmt.Sprintf("%s%s%s", scope, seperator, name)
		} else {
			newName = name
		}

		glsLabels = append(glsLabels, GLSLabels{
			Name:        newName,
			Color:       color,
			Description: description,
		})
	}

	return glsLabels
}

func main() {
	pathPointer := flag.String("config", "./labels.toml", "The Path to the config, by default \"./labels.toml\"")

	path := *pathPointer

	var tomlLabelConfig TomlConfig
	toml.DecodeFile(path, &tomlLabelConfig)

	glsLabelConfig := convertFormat(tomlLabelConfig)

	jsonData, err := json.MarshalIndent(glsLabelConfig, "", "  ")

	if err != nil {
		log.Fatalln(err)
	}

	os.WriteFile("labels.json", jsonData, 0644)
}
