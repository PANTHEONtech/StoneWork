package main

import (
	"fmt"
	"os"

	"github.com/buildkite/interpolate"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// EntityVar is a variable of an entity defined with a template to render its
// value.
type EntityVar struct {
	Index int `json:"-"`

	Name        string `json:"name"`
	Description string `json:"description"`
	Value       string `json:"default"`
	Type        string `json:"type"`
}

// Entity is a blueprint for an object defined with a config template of
// related parts.
type Entity struct {
	Origin string `json:"-"`

	Name        string      `json:"name"`
	Plural      string      `json:"plural"`
	Description string      `json:"description"`
	Vars        []EntityVar `json:"vars"`
	Config      string      `json:"config"`
}

func (e Entity) GetName() string {
	return e.Name
}

func (e Entity) GetPlural() string {
	if e.Plural == "" {
		return e.Name + "s"
	}
	return e.Plural
}

func (e Entity) GetVariables() []EntityVar {
	return e.Vars
}

const defaultEntityFile = "entities.yaml"

// EntityFile is a file containing entities loaded during initialization.
type EntityFile struct {
	Entities []Entity `json:"entities"`
}

func loadEntityFile(file string) (*EntityFile, error) {
	logrus.Tracef("loading entity file: %s", file)

	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var entityFile EntityFile
	if err := yaml.Unmarshal(b, &entityFile); err != nil {
		return nil, err
	}
	return &entityFile, nil
}

func loadEntityFiles(files []string) ([]Entity, error) {
	files = append(files, defaultEntityFile)

	var (
		entities []Entity
		uniq     = make(map[string]Entity)
	)

	logrus.Tracef("loading %d entity files", len(files))

	for _, f := range files {
		entityFile, err := loadEntityFile(f)
		if err != nil {
			if os.IsNotExist(err) && f == defaultEntityFile {
				logrus.Debugf("default entity file %s not found", f)
				continue
			}
			return nil, fmt.Errorf("loading entity file %s failed: %w", f, err)
		}

		logrus.Tracef("%d entities loaded from file %s", len(entityFile.Entities), f)

		for _, entity := range entityFile.Entities {
			if err := validateEntity(&entity); err != nil {
				return nil, fmt.Errorf("invalid entity %v: %w", entity.Name, err)
			}
			if entity.Origin == "" {
				entity.Origin = f
			}
			if _, ok := uniq[entity.Name]; ok {
				return nil, fmt.Errorf("duplicate entity %v in file %v", entity.Name, f)
			}

			logrus.Tracef("loaded entity: %v", entity.Name)
			uniq[entity.Name] = entity
			entities = append(entities, entity)
		}
	}

	logrus.Debugf("loaded %d entities from %d file(s)", len(entities), len(files))

	return entities, nil
}

func validateEntity(entity *Entity) error {
	if entity == nil {
		return nil
	}

	// TODO: validate name

	// validate variables
	vars := make(map[string]EntityVar)
	for i, v := range entity.Vars {
		v.Index = i
		if dup, ok := vars[v.Name]; ok {
			return fmt.Errorf("duplicate var %v on index %d, previous on index %d", v.Name, i, dup.Index)
		}
		vars[v.Name] = v

		idents, err := interpolate.Identifiers(v.Value)
		if err != nil {
			return fmt.Errorf("invalid var reference in value of var %v: %w", v.Name, err)
		}
		for _, ident := range idents {
			if isBuiltinVar(ident) {
				continue
			}
			if _, ok := vars[ident]; !ok {
				return fmt.Errorf("undefined var reference %v found in var %v value", ident, v.Name)
			}
		}
	}

	// validate config
	idents, err := interpolate.Identifiers(entity.Config)
	if err != nil {
		return fmt.Errorf("invalid var reference in config: %w", err)
	}
	for _, ident := range idents {
		if isBuiltinVar(ident) {
			continue
		}
		if _, ok := vars[ident]; !ok {
			return fmt.Errorf("undefined var reference %v found in config", ident)
		}
	}
	return nil
}

func isBuiltinVar(name string) bool {
	if name == "ID" || name == "IDX" {
		return true
	}
	return false
}
