package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

const defaultEntityFile = "entities.yaml"

type EntityOption struct {
	Name        string `json:"name"`
	Value       string `json:"default"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type Entity struct {
	Name    string         `json:"name"`
	Plural  string         `json:"plural"`
	Options []EntityOption `json:"options"`
	Config  string         `json:"config"`
}

func (e Entity) GetName() string {
	return e.Name
}

func (e Entity) GetPlural() string {
	return e.Plural
}

func (e Entity) GetOptions() []EntityOption {
	return e.Options
}

type EntityFile struct {
	Entities []Entity `json:"entities"`
}

func loadEntityFile(entityFile string) (*EntityFile, error) {
	var file EntityFile

	logrus.Tracef("loading entity file: %s", entityFile)

	b, err := os.ReadFile(entityFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

func loadEntityFiles(files []string) ([]Entity, error) {
	if files == nil {
		files = []string{defaultEntityFile}
	}

	var entities []Entity
	for _, f := range files {
		file, err := loadEntityFile(f)
		if err != nil {
			if f == defaultEntityFile {
				logrus.Debugf("default entity file %s not found", f)
				continue
			}
			return nil, fmt.Errorf("failed to load entity file %s: %w", f, err)
		}

		logrus.Tracef("%d entities loaded from file %s", len(file.Entities), f)
		entities = append(entities, file.Entities...)
	}

	logrus.Debugf("loaded %d entities from %d file(s)", len(entities), len(files))

	return entities, nil
}
