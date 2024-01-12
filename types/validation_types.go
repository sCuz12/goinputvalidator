package types


type ValidationTypeRule string

const (
	Required ValidationTypeRule   = "required"
	Max ValidationTypeRule 		  = "max"
	Min ValidationTypeRule 		  = "min"
	Email ValidationTypeRule 	  = "email"
	Url ValidationTypeRule	      = "url"
	ActiveUrl ValidationTypeRule  = "active_url"
	IpFormat ValidationTypeRule   = "ipformat"
	Date	ValidationTypeRule    = "date"
	DateFormat ValidationTypeRule = "dateFormat"
	Between ValidationTypeRule 	  = "between"
	In ValidationTypeRule		  = "in"
	Accepted ValidationTypeRule   = "accepted"
)
