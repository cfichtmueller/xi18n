package main

type LanguageMap = map[string]string
type KeyMap = map[string]string
type KeysMap = map[string]KeyMap

type Model struct {
	Languages   LanguageMap `yaml:"languages,omitempty"`
	StringsFile string      `yaml:"stringsFile,omitempty"`
	Keys        KeysMap     `yaml:"keys,omitempty"`
}
