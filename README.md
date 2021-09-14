# e2e-secrets-to-file

Run end to end happy-path that uses a representative annotations file generated by `./testdata/generate.go`:

```shell
# Generate annotation.txt in ./test_data
# This file contains a representative file of annotations as
# you what you might find in a Kubernetes pod generated via the
# Downwards API
$ go run ./testdata
Generating testdata/annotations.txt

# Run end to end flow.
# Parses annotations, generates secret groups, fetches secrets (mock), write secrets to file
$ go run ./
Parse annotations from downward API
Generate Secret Groups
Fetch secrets for db
Fetch dev/openshift/api-url
Fetch dev/openshift/password
Fetch dev/openshift/username
Fetch secrets for cache
Fetch dev/openshift/api-url
Fetch dev/openshift/password
Fetch dev/openshift/username
Process template for db to ./testdata/db.json
Process template for cache to ./testdata/cache.cfg
```


Run unit tests:
```shell
# Unit tests are written in a table driven fashion where appropriate.
# We make use of the standard Go tooling for unit testing, the only addition
# is an assertion library. 
$ go test -v -count 1 ./...
=== RUN   TestNewAnnotations
=== RUN   TestNewAnnotations/valid_example
=== RUN   TestNewAnnotations/malformed_line_without_equals
=== RUN   TestNewAnnotations/malformed_line_without_quoted_value
--- PASS: TestNewAnnotations (0.00s)
    --- PASS: TestNewAnnotations/valid_example (0.00s)
    --- PASS: TestNewAnnotations/malformed_line_without_equals (0.00s)
    --- PASS: TestNewAnnotations/malformed_line_without_quoted_value (0.00s)
=== RUN   TestNewSecretSpecs
=== RUN   TestNewSecretSpecs/normal_test
=== RUN   TestNewSecretSpecs/valid_example
=== RUN   TestNewSecretSpecs/malformed_not_a_list
=== RUN   TestNewSecretSpecs/malformed_multiple_key-values_in_one_entry
--- PASS: TestNewSecretSpecs (0.00s)
    --- PASS: TestNewSecretSpecs/valid_example (0.00s)
    --- PASS: TestNewSecretSpecs/malformed_not_a_list (0.00s)
    --- PASS: TestNewSecretSpecs/malformed_multiple_key-values_in_one_entry (0.00s)
PASS
ok      e2e-write-to-file       0.015s
```
