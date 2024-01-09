package validator

import (
	"fmt"
	"govalidator/pkg/rules"
	"govalidator/types"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)


type ValidationError struct {
	Field string 
	Message string
}

type Validator struct {
	//ex: title can several rules
	//ex : [body : [{max 100} {required}]]
	Rules map[string][]rules.Rule
}

//Returns pointer to Validator Struct
func New() *Validator {
	return &Validator{
		Rules: make(map[string][]rules.Rule),
	}
}

func NewValidationError(field string , message string) *ValidationError {
	return &ValidationError {
		Field: field,
		Message: message,
	}
}

// Error implements the error interface for ValidationError.
func (ve *ValidationError) Error() string {
    return fmt.Sprintf("Validation error in field '%s': %s", ve.Field, ve.Message)
}

func (v *Validator) AddRule(field string , rule rules.Rule) {
	v.Rules[field] = append(v.Rules[field], rule)
}

//Debug -Testing purposes
func (v *Validator) DebugRules() {
	fmt.Println(v.Rules)
	for field,rules := range v.Rules {
		fmt.Printf("The field is %v \n",field)

		for _,rule := range rules {
			fmt.Printf("Rule type %v and rule param %v \n" ,rule.Type , rule.Param)
		}
	}
} 


func (v *Validator) Validate(input interface{})[]error {
	var errors []error

	val := reflect.ValueOf(input)

	//loop through fiels in struct getter from reflect
	for i:=0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // the whole field ex: {Title  string validate:"required|max:1" 0 [0] false}
		fieldName := field.Name

		fieldValue := val.Field(i).Interface()  //whole value of field ex: This is title
		fieldValueStr := fieldValue.(string) 


		//get rules required,email etc
		rulesTags := field.Tag.Get("validate")

		if rulesTags == "" {
			//no validation tags
			continue;
		}

		rulesArr := strings.Split(rulesTags,"|")

		for _,ruleValue := range rulesArr {
			rulesPart := strings.Split(ruleValue,":")

			var ruleType string = "" ;
			var ruleNum  int 	= 0

			if len(rulesPart) > 1 {
				ruleType = rulesPart[0]
				ruleNum,_  = strconv.Atoi(rulesPart[1])
			} else {
				ruleType = rulesPart[0]
			}
	
			switch ruleType {
				//max
				case string(types.Max): {
					maxAssigned := ruleNum

					valueLen := len(fieldValueStr)
					
					if valueLen > maxAssigned {
						//return fmt.Errorf("Field %s exceeds the maximum length of %d", field, maxAssigned) 
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field %s exceeds the maximum length of %d", fieldName, maxAssigned) ))
					}

				}
				//min
				case string(types.Min) : {
					minAssigned := ruleNum

					valueLen := len(fieldValueStr)

					if(valueLen < minAssigned) {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("Field %s is less than the minimum length of %d", fieldName, minAssigned)))
					}
				}
				//required 
				case string(types.Required) : {
		
					if fieldValueStr == "" {
						//return fmt.Errorf("%s is required", field)
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field not found: %s", fieldName)))
					}
				}
				//email 
				case string(types.Email) : {
					
					isEmailValid := isValidEmail(fieldValueStr)

					if(!isEmailValid) {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("Input field %s must be a valid email format",fieldName)))
					}
				}
				//url
				case string(types.Url): {

					isUrl := isURL(fieldValueStr)

					if !isUrl {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("The %s field must be a valid URL.", fieldName)))
					}
				}
				//active URL
				case string(types.ActiveUrl) : {
					
					isActiveUrl := isActiveUrl(fieldValueStr)

					if !isActiveUrl {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("The %s field must be active URL.", fieldName)))
					}
				}
				//ip
				case string(types.IpFormat): {
					isValidIP := isValidIP(fieldValueStr)

					if !isValidIP {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("The %s field must be valid IP format", fieldName)))
					}
				}

			}
			
		}

	}
	return  errors
}


func isValidEmail(email string) bool {
	emailPattern :=  `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

  	// Compile the regular expression pattern
	re := regexp.MustCompile(emailPattern)

    // Use the regular expression to match the email
	return re.MatchString(email)
}


func isURL(givenURL string) bool {
	_,err := url.ParseRequestURI(givenURL)


	if(err != nil) {
		return false
	}

	return true
}

func isActiveUrl(givenURL string) bool {
	response, err := http.Get(givenURL)

	if err != nil {
		//no response -> website is not active
		return false
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return true
	}

	return false
}

func isValidIP(givenIP string) bool {
	ip := net.ParseIP(givenIP)

	return ip != nil
}