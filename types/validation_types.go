package types


type ValidationTypeRule string

const (
	Required ValidationTypeRule = "required"
	Max ValidationTypeRule = "max"
	Min ValidationTypeRule = "min"
	Email ValidationTypeRule = "email"
)
