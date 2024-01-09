package validator

import (
	"fmt"
	"govalidator/pkg/rules"
	"govalidator/types"
	"regexp"
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

func (v *Validator) Validate(input map[string]interface{}) []error {
	var errors []error

	//loop through existing Rules 
	//ex: loop rules of body or title
	for field , rules := range v.Rules {
		//ex : title can have multiple validations like required,max etc 
		for _, rule := range rules {
			//ex : switch max,required,email
			switch rule.Type {
				case string(types.Required) : {
					value ,ok := input[field] 

					if !ok || value == nil || value == "" {
						//return fmt.Errorf("%s is required", field)
						errors = append(errors, NewValidationError(field, fmt.Sprintf("Field not found: %s", field)))
					}
				}

				case string(types.Max) : {
					value , ok := input[field]
					maxAssigned := rule.Param.(int)
					
					if  !ok {
						//return fmt.Errorf("Not found")
						errors = append(errors, NewValidationError(field, fmt.Sprintf("Not found")))
					}

					valueStr, okStr := value.(string) // Assuming the field's value is expected to be a string

					if !okStr {
						// return fmt.Errorf("Field %s must be a string for max validation", field)
						errors = append(errors, NewValidationError(field, fmt.Sprintf("Field %s must be a string for max validation", field)))
					}
						
					valueLen := len(valueStr)

					if valueLen > maxAssigned {
						//return fmt.Errorf("Field %s exceeds the maximum length of %d", field, maxAssigned) 
						errors = append(errors, NewValidationError(field, fmt.Sprintf("Field %s exceeds the maximum length of %d", field, maxAssigned) ))
					}
				}
				case string(types.Email) : {
					value , _ := input[field]
					
					isEmailValid := isValidEmail(value.(string))
					

					if(!isEmailValid) {
						errors = append(errors, NewValidationError(field,fmt.Sprintf("Input field %s must be a valid email format",field)))
					}
				}

				case string(types.Min) : {
					value, _ 	:= input[field]
					minAssigned := rule.Param.(int)

					valueStr := value.(string)
					valueLen := len(valueStr)

					if(valueLen < minAssigned) {
						errors = append(errors, NewValidationError(field,fmt.Sprintf("Field %s is less than the minimum length of %d", field, minAssigned)))
					}
				}
			}
		}
	}

	return errors 
}


func isValidEmail(email string) bool {
	emailPattern :=  `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

  	// Compile the regular expression pattern
	re := regexp.MustCompile(emailPattern)

    // Use the regular expression to match the email
	return re.MatchString(email)
}