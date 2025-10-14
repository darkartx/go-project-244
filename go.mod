module github.com/darkartx/go-project-244

go 1.24.5

require (
	github.com/darkartx/go-project-244/formatters v0.0.0-00010101000000-000000000000
	github.com/darkartx/go-project-244/shared v0.0.0
	github.com/go-yaml/yaml v2.1.0+incompatible
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v3 v3.4.1
)

replace (
	github.com/darkartx/go-project-244/formatters => ./formatters
	github.com/darkartx/go-project-244/shared => ./shared
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
