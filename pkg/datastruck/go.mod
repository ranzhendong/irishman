module datastruck

go 1.13

require (
	govalidators v1.1.2
	reconstruct v1.1.2
)

replace (
	govalidators v1.1.2 => ../govalidators
	reconstruct v1.1.2 => ../reconstruct
)
