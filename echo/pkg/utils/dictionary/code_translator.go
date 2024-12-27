package dictionary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"modular_chassis/echo/pkg"
	"path/filepath"
	"sync"
)

func init() {
	GetCodeTranslator().ParseEmbeddedFiles("services")
}

type codeTranslator struct {
	messages map[string]string
}

var (
	codeTranslatorOnce sync.Once
	codeTranslatorIns  *codeTranslator
)

func GetCodeTranslator() *codeTranslator {
	codeTranslatorOnce.Do(func() {
		if codeTranslatorIns == nil {
			codeTranslatorIns = &codeTranslator{
				messages: make(map[string]string),
			}
		}
	})
	return codeTranslatorIns
}

type Language string

const (
	English Language = "en"
	Farsi   Language = "fa"
)

func (c *codeTranslator) TranslateResponseCode(lang Language, code int) string {
	return c.Get(lang, code)
}

func (c *codeTranslator) Populate(data []byte) {
	var mjs []struct {
		Lang    string `json:"lang"`
		Message string `json:"msg"`
		Key     int    `json:"key"`
	}
	decoder := json.NewDecoder(bytes.NewBuffer(data))
	err := decoder.Decode(&mjs)
	if err != nil {
		return
	}

	for _, m := range mjs {
		c.Put(Language(m.Lang), m.Key, m.Message)
	}
	return
}

func (c *codeTranslator) Put(lang Language, key int, message string) {
	k := fmt.Sprintf("%s:%d", lang, key)
	if _, ok := c.messages[k]; ok {
		return
	}
	c.messages[k] = message
}

func (c *codeTranslator) Get(lang Language, key int) string {
	k := fmt.Sprintf("%s:%d", lang, key)
	if val, ok := c.messages[k]; ok {
		return val
	}
	return ""
}

func (c *codeTranslator) ParseEmbeddedFiles(dirName string) {
	dir, err := pkg.EmbeddedFiles.ReadDir(dirName)
	if err != nil {
		return
	}
	for _, file := range dir {
		if file.IsDir() {
			c.ParseEmbeddedFiles(filepath.Join(dirName, file.Name()))
		} else {
			if file.Name() != "codes.json" {
				continue
			}
			data, err := pkg.EmbeddedFiles.ReadFile(filepath.Join(dirName, file.Name()))
			if err != nil {
				return
			}
			c.Populate(data)
		}
	}
	return
}
