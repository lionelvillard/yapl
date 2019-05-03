module github.com/lionelvillard/yapl

go 1.12

require (
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.2
	go.starlark.net v0.0.0-00010101000000-000000000000

	gopkg.in/yaml.v2 v2.2.2
)

replace go.starlark.net => github.com/lionelvillard/starlark-go v0.1.0
