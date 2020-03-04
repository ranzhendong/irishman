module kvnuts

go 1.13

require (
	datastruck v1.1.2
	errorhandle v1.1.2
	github.com/xujiajun/nutsdb v0.5.0
)

replace (
	datastruck v1.1.2 => ../datastruck
	errorhandle v1.1.2 => ../errorhandle
)
