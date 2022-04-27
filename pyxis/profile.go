package pyxis

import (
	"os"

	"gopkg.in/yaml.v3"

	"app/base/utils"
)

const (
	profilesFilePath = "/vuln4shift/pyxis/profiles.yml"
)

var (
	profileMap map[string]map[string]map[string]struct{}
	profile    = utils.Cfg.PyxisProfile
)

type Profile struct {
	Name       string     `yaml:"name"`
	Registries []Registry `yaml:"registries"`
}

type Registry struct {
	Name         string   `yaml:"name"`
	Repositories []string `yaml:"repositories"`
}

func parseProfiles() {
	profilesFile, err := os.ReadFile(profilesFilePath)
	if err != nil {
		logger.Fatalf("Unable to read profiles file: %s", err)
	}

	var profiles []Profile
	err = yaml.Unmarshal(profilesFile, &profiles)
	if err != nil {
		logger.Fatalf("Unable to parse profiles: %s", err)
	}

	// Transform parsed structure to map of maps
	profileMap = make(map[string]map[string]map[string]struct{}, len(profiles))
	for _, profile := range profiles {
		profileMap[profile.Name] = make(map[string]map[string]struct{}, len(profile.Registries))
		for _, registry := range profile.Registries {
			profileMap[profile.Name][registry.Name] = make(map[string]struct{}, len(registry.Repositories))
			for _, repository := range registry.Repositories {
				profileMap[profile.Name][registry.Name][repository] = struct{}{}
			}
		}
	}

	// If the profile is not empty, it must exist in the yaml
	if profile != "" {
		if _, found := profileMap[profile]; !found {
			logger.Fatalf("Unknown profile: %s", profile)
		}
	}
}

func repositoryInProfile(registry, repository string) bool {
	// Empty profile = sync all repositories
	if profile != "" {
		_, found := profileMap[profile][registry][repository]
		return found
	}
	return true
}
