package main

import (
	"os"

	"gopkg.in/yaml.v2"
)

type UserConfig struct {
	Host              string   `yaml:"host"`
	HostAuthType      string   `yaml:"host_auth_type"`
	Username          string   `yaml:"username"`
	Password          string   `yaml:"password"`
	PasswordEncrypted string   `yaml:"password_encrypted"`
	ProjectSlugs      []string `yaml:"project_slugs"`
}

func readConfigFile(path string) (*UserConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	s := UserConfig{}
	err = yaml.Unmarshal([]byte(b), &s)
	return &s, err
}

func writeConfigFile(cfg *UserConfig, path string) error {
	cfgCopy := *cfg
	cfgCopy.Password = "" // clear password to avoid writing it back to config
	b, err := yaml.Marshal(&cfgCopy)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}
