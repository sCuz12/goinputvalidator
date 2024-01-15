package types


type ValidationTypeRule string

const (
	Required ValidationTypeRule   	     = "required"
	Max ValidationTypeRule 		 	     = "max"
	Min ValidationTypeRule 		  	     = "min"
	Email ValidationTypeRule 	  	     = "email"
	Url ValidationTypeRule	      	     = "url"
	ActiveUrl ValidationTypeRule  	     = "active_url"
	IpFormat ValidationTypeRule   	     = "ipformat"
	Date	ValidationTypeRule    	     = "date"
	DateFormat ValidationTypeRule 	     = "dateFormat"
	Between ValidationTypeRule 	  	     = "between"
	In ValidationTypeRule		         = "in"
	Accepted ValidationTypeRule          = "accepted"
	NotIn	ValidationTypeRule	         = "notIn"
	Size ValidationTypeRule 	         = "size"	
	Confirmed ValidationTypeRule  	     = "confirmed"
	Doesnt_end_with ValidationTypeRule   = "doesnt_end_with"
	Doesnt_start_with ValidationTypeRule = "doesnt_start_with"
	MacAddressFormat ValidationTypeRule	 = "macAddress"
	GreaterThan		ValidationTypeRule   = "gt"
	LessThan ValidationTypeRule 		 = "lt"
	GreaterEqualThan ValidationTypeRule  = "gte"
	LessThanEqual ValidationTypeRule     = "lte"
)
