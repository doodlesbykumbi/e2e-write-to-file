package main

import (
	"io"
	"text/template"

	"gopkg.in/yaml.v3"
)

func processTemplate(name string, templateString string, secretsMap map[string][]byte, writer io.Writer) error {
	secretsMapWithStrings := map[string]string{}
	for k, v := range secretsMap {
		secretsMapWithStrings[k] = string(v)
	}

	t, err := template.New(name).Funcs(template.FuncMap{
		// secret is a custom utility function with more bells and whistles than raw map
		// e.g. it can panic for a secret label that doesn't exist.
		"secret": func(label string) string {
			v, ok := secretsMapWithStrings[label]
			if ok {
				return v
			}

			// Panic here doesn't mean panic in the process. Panic here is captured as an error
			// when the template is executed.
			panic("secret label used in template is not present in the spec")
		},
		"toYAML": func(input interface{}) string {
			// Marshal it to YAML format.
			d, err := yaml.Marshal(&input)
			if err != nil {
				panic(err)
			}

			return string(d)
		},
	}).Parse(templateString)
	if err != nil {
		return err
	}

	// Execute the Go template.
	// secretsMap is passed as second argument but probably better to only provide 'secret'
	// custom utility function (shown above).
	return t.Execute(writer, secretsMapWithStrings)
}
