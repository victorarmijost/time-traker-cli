package repl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"varmijo/time-tracker/utils"
)

const TemplateDir = "templates"

func init() {
	err := os.MkdirAll(utils.GeAppPath(TemplateDir), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

type Template map[string]string

type TemplateHandler struct {
	name      string
	templates map[string]Template
}

func NewTemplateHandler(name string) *TemplateHandler {
	err := os.MkdirAll(utils.GeAppPath(TemplateDir)+"/"+name, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	return &TemplateHandler{
		name:      name,
		templates: map[string]Template{},
	}
}

func (t *TemplateHandler) add(name string, template Template) {
	t.templates[name] = template
}

func (t *TemplateHandler) Get(name string) (Template, error) {
	if template, ok := t.templates[name]; ok {
		return template, nil
	}

	return nil, fmt.Errorf("template not found")
}

func (t *TemplateHandler) List() []string {
	list := []string{}
	for k, v := range t.templates {
		if des, ok := v["x-description"]; ok {
			list = append(list, fmt.Sprintf("%s: %s [%s] - %s", k, v["x-name"], v["x-category"], des))
			continue
		}

		list = append(list, fmt.Sprintf("%s: %s [%s]", k, v["x-name"], v["x-category"]))

	}
	return list
}

func (t *TemplateHandler) Load() error {
	path := utils.GeAppPath(TemplateDir) + "/" + t.name

	items, err := ioutil.ReadDir(path)

	if err != nil {
		return err
	}

	for _, item := range items {
		if !item.IsDir() {
			basename := item.Name()
			data, err := ioutil.ReadFile(path + "/" + basename)
			if err != nil {
				return err
			}

			temp := make(Template)
			err = json.Unmarshal(data, &temp)

			if err != nil {
				return err
			}

			name := strings.TrimSuffix(basename, filepath.Ext(basename))

			t.add(name, temp)

		}
	}

	return nil
}

func (t *TemplateHandler) Save(name string, template Template) error {
	path := utils.GeAppPath(TemplateDir) + "/" + t.name

	data, err := json.MarshalIndent(template, "", "\t")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path+"/"+name+".json", data, 0644)
	if err != nil {
		return err
	}

	t.add(name, template)

	return nil
}
