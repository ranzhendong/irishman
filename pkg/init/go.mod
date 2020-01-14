module init

go 1.13

require (
	datastruck v1.1.2
	errorhandle v1.1.2
	github.com/spf13/viper v1.6.1
)

replace (
	datastruck v1.1.2 => ../datastruck
	errorhandle v1.1.2 => ../errorhandle
)
