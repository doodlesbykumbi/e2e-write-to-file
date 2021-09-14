package main

import (
	"fmt"
)

func main() {
	// Parse annotations from downward API
	fmt.Println("Parse annotations from downward API")
	annotations, err := NewAnnotationsFromFile("./testdata/annotations.txt")
	if err != nil {
		panic(err)
	}

	// Generate Secret Groups
	fmt.Println("Generate Secret Groups")
	secretGroups, errs := NewSecretGroups(annotations)
	if len(errs) > 0 {
		panic(errs)
	}

	// Fetch secrets MOCK!
	// TODO: Should we make sure any secret id is only fetched once ?
	// TODO: Should we zeroize the secrets map when we're done ?
	// TODO: Secret fetching should be concurrent and where possible parallel
	for _, group := range secretGroups {
		fmt.Println("Fetch secrets for", group.Label)
		secretsMap := map[string][]byte{}

		for _, spec := range group.SecretSpecs {
			fmt.Println("Fetch", spec.Id)
			secretsMap[spec.Alias] = []byte(fmt.Sprintf("s-%s", spec.Id))
		}

		group.SecretsMap = secretsMap
	}

	// Write secrets to file
	for _, group := range secretGroups {
		fmt.Printf("Process template for %s to %s\n", group.Label, group.FilePath)

		err := group.ProcessTemplate()
		if err != nil {
			// TODO: Accumulate errors instead
			panic(err)
		}
	}
}
