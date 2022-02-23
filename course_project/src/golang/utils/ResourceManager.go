package utils

import (
	"github.com/magiconair/properties"
	log "github.com/sirupsen/logrus"
)

type ResourceManager struct {
	propertyFiles []string
	props         *properties.Properties
}

func NewResourceManager(propertyFiles []string) *ResourceManager {
	props, err := properties.LoadFiles(propertyFiles, properties.UTF8, false)
	if err != nil {
		log.Fatalf("Error in load resources file: %s", err)
	}
	return &ResourceManager{propertyFiles: propertyFiles, props: props}
}

func (m ResourceManager) GetProp(name string) string {
	prop, ok := m.props.Get(name)
	if !ok {
		return ""
	} else {
		return prop
	}
}
