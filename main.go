package main

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"sort"
)

const (
	projectFile = "xi18n.yml"
)

func main() {
	initialize()
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "xi18n",
	Short: "xi18n generates Localizable.strings files",
	Long:  `Generate Localizable.strings files for your Xcode projects`,
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new project",
	Run: func(c *cobra.Command, args []string) {
		if existsProjectFile() {
			fail("This directory already contains a project.")
		}
		model := &Model{
			Languages: LanguageMap{
				"en": "Project/en.lproj/Localizable.strings",
				"de": "Project/de.lproj/Localizable.strings",
			},
			StringsFile: "Project/Strings.swift",
			Keys: KeysMap{
				"Hello": {
					"en": "Hello",
					"de": "Hallo",
				},
				"World": {
					"en": "World",
					"de": "Welt",
				},
			},
		}

		d, err := yaml.Marshal(model)
		if err != nil {
			log.Fatalf("failed to write model: %v", err)
		}

		if err := os.WriteFile(projectFile, d, 0755); err != nil {
			log.Fatalf("failed to write file: %v", err)
		}
	},
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate files",
	Run: func(c *cobra.Command, args []string) {
		if !existsProjectFile() {
			fail("Didn't find project file")
		}
		d, err := os.ReadFile(projectFile)
		if err != nil {
			fail("error: %v", err)
		}
		var model Model
		if err := yaml.Unmarshal(d, &model); err != nil {
			fail("Couldn't read project file: %v", err)
		}

		langFiles := make(map[string]*os.File)
		langKeys := make([]string, 0, len(model.Languages))
		for k := range model.Languages {
			langKeys = append(langKeys, k)
		}
		sort.Strings(langKeys)

		for key, fileName := range model.Languages {
			file, err := os.OpenFile(fileName, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
			if err != nil {
				fail("Couldn't open language file: %v", err)
			}
			langFiles[key] = file
		}

		keys := make([]string, 0, len(model.Keys))
		for k := range model.Keys {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, key := range keys {
			keyContent := model.Keys[key]
			comment, ok := keyContent["comment"]
			if !ok {
				comment = key
			}

			for _, lang := range langKeys {
				file, ok := langFiles[lang]
				if !ok {
					warn("Found invalid language key: %v", lang)
					continue
				}
				content, ok := keyContent[lang]
				if !ok {
					warn("Key %s is missing value for language %s", key, lang)
					continue
				}
				write(file, createComment(comment))
				write(file, createEntry(key, content))
				write(file, "\n")
			}
		}

		for _, file := range langFiles {
			if err := file.Close(); err != nil {
				fail("Couldn't close file %s: %v", file.Name(), err)
			}
		}

		if len(model.StringsFile) == 0 {
			return
		}

		stringsFile, err := os.OpenFile(model.StringsFile, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
		if err != nil {
			fail("Couldn't open strings file")
		}
		write(stringsFile, "import Foundation\n\nclass Strings {\n\n")
		for _, key := range keys {
			write(stringsFile, fmt.Sprintf("    static let %s = NSLocalizedString(\"%s\", comment: \"%s\")\n", key, key, key))
		}
		write(stringsFile, "\n    private init() {}\n}\n")
	},
}

func initialize() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(genCmd)
}

func existsProjectFile() bool {
	_, err := os.Stat(projectFile)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false
		}
		log.Fatalf("error: %v", err)
	}
	return true
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fail("%v", err)
	}
}
