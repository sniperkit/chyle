package chyle

import (
	"fmt"

	"github.com/antham/envh"
)

// decorater extends data from commit hashmap with data picked from third part service
type decorater interface {
	decorate(*map[string]interface{}) (*map[string]interface{}, error)
}

// decorate process all defined decorator and apply them
func decorate(decorators *map[string][]decorater, changelog *Changelog) (*Changelog, error) {
	var err error

	datas := []map[string]interface{}{}

	for _, d := range changelog.Datas {
		result := &d

		for _, decorator := range (*decorators)["datas"] {
			result, err = decorator.decorate(&d)

			if err != nil {
				return nil, err
			}
		}

		datas = append(datas, *result)
	}

	changelog.Datas = datas

	metadatas := changelog.Metadatas

	for _, decorator := range (*decorators)["metadatas"] {
		m, err := decorator.decorate(&metadatas)

		if err != nil {
			return nil, err
		}

		metadatas = *m
	}

	changelog.Metadatas = metadatas

	return changelog, nil
}

// createDecorators build decorators from a config
func createDecorators(config *envh.EnvTree) (*map[string][]decorater, error) {
	results := map[string][]decorater{"metadatas": {}, "datas": {}}

	var decType string
	var dec decorater
	var decs []decorater
	var err error
	var subConfig envh.EnvTree

	for _, k := range config.GetChildrenKeys() {
		dec = nil
		decs = []decorater{}
		decType = ""

		switch k {
		case "JIRA":
			decType = "datas"
			subConfig, err = config.FindSubTree("JIRA")

			if err != nil {
				break
			}

			dec, err = buildJiraDecorator(&subConfig)
		case "ENV":
			decType = "metadatas"
			subConfig, err = config.FindSubTree("ENV")

			if err != nil {
				break
			}

			decs, err = buildEnvDecorators(&subConfig)
		default:
			err = fmt.Errorf(`a wrong decorator key containing "%s" was defined`, k)
		}

		if err != nil {
			return nil, err
		}

		if len(decs) == 0 {
			results[decType] = append(results[decType], dec)
		} else {
			results[decType] = append(results[decType], decs...)
		}
	}

	return &results, nil
}
