package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	_, b, _, _ = runtime.Caller(0)
	workingdir, _ = os.Getwd()
	basepath   = filepath.Dir(b)
)

// Taken from Kubernetes source used to generate annotation file for downward API
// https://github.com/kubernetes/kubernetes/blob/master/pkg/fieldpath/fieldpath.go#L28-L41
func FormatMap(m map[string]string) (fmtStr string) {
	for k, v := range m {
		fmtStr += fmt.Sprintf("%v=%q\n", k, v)
	}
	fmtStr = strings.TrimSuffix(fmtStr, "\n")

	return
}

func main() {
	annotations := FormatMap(map[string]string{
		"conjur.org/conjur-secrets.cache": `- dev/openshift/api-url
- admin-password: dev/openshift/password
- admin-username: dev/openshift/username
`,
		"conjur.org/secret-file-path.cache": "./testdata/cache.cfg",
		"conjur.org/secret-file-template.cache": `"cache": {
	"url": {{ index . "api-url" }},
	"password": {{ printf "%q" (secret "admin-password") }},
	"username": {{ secret "admin-username" }},
	"port": 123456,
}`,
		"conjur.org/conjur-secrets.db": `- dev/openshift/api-url
- admin-password: dev/openshift/password
- admin-username: dev/openshift/username
`,
		"conjur.org/secret-file-path.db": "./testdata/db.json",
		"conjur.org/secret-file-template.db": `{
	"url": {{ printf "%q" (secret "api-url") }},
	"password": {{ printf "%q" (secret "admin-password") }},
	"username": {{ printf "%q" (secret "admin-username") }},
	"port": 123456
}`,
	})

	annotationsFilePath, _ := filepath.Rel(workingdir, filepath.Join(basepath, "./annotations.txt"))
	fmt.Println("Generating " + annotationsFilePath)
	_ = ioutil.WriteFile(annotationsFilePath, []byte(annotations), 0644)
}