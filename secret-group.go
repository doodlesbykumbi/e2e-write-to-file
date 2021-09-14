package main

import (
	"os"
	"strings"
)

const secretGroupPrefix = "conjur.org/conjur-secrets."
const secretGroupFileTemplatePrefix = "conjur.org/secret-file-template."
const secretGroupFilePathPrefix = "conjur.org/secret-file-path."
const modeDefaultWrite os.FileMode = 0660

type SecretGroup struct {
	Label     string
	FilePath     string
	FileTemplate string
	SecretSpecs  []SecretSpec
	SecretsMap  map[string][]byte
}

func (s SecretGroup) ProcessTemplate() error {
	// modeDefaultWrite is the default but can be configured at a later date
	f, _ := os.OpenFile(s.FilePath, os.O_WRONLY|os.O_CREATE, modeDefaultWrite)
	defer func() {
		_ = f.Close()
	}()

	return processTemplate(
		s.Label,
		s.FileTemplate,
		s.SecretsMap,
		f,
	)
}

func NewSecretGroups(annotations map[string]string) ([]*SecretGroup, []error) {
	var secretGroups []*SecretGroup

	var errors []error
	for k, v := range annotations {
		if strings.HasPrefix(k, secretGroupPrefix) {
			groupLabel := strings.TrimPrefix(k, secretGroupPrefix)
			secretSpecs, err := NewSecretSpecs([]byte(v))
			if err != nil {
				// Accumulate errors
				errors = append(errors, err)
				continue
			}

			fileTemplate := annotations[secretGroupFileTemplatePrefix+groupLabel]
			filePath := annotations[secretGroupFilePathPrefix+groupLabel]

			secretGroups = append(secretGroups, &SecretGroup{
				Label:        groupLabel,
				SecretSpecs:  secretSpecs,
				FilePath:     filePath,
				FileTemplate: fileTemplate,
			})
		}
	}

	if len(errors) >0 {
		return nil, errors
	}

	return secretGroups, nil
}
