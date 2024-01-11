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
	"time"
)

/*Allowed Date Format*/
var formatMap = map[string]string{
	"YYYY/MM/DD" : "2006/01/02",
	"YYYY/DD/MM" : "2006/02/01",
	"YYYY-MM-DD" : "2006-01-02",
	"YYYY-DD-MM" : "2006-02-01",
}


type ValidationError struct {
	Field   string
	Message string
}

type Validator struct {
	//ex: title can several rules
	//ex : [body : [{max 100} {required}]]
	Rules map[string][]rules.Rule
}

// Returns pointer to Validator Struct
func New() *Validator {
	return &Validator{
		Rules: make(map[string][]rules.Rule),
	}
}

func NewValidationError(field string, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// Error implements the error interface for ValidationError.
func (ve *ValidationError) Error() string {
	return fmt.Sprintf("Validation error in field '%s': %s", ve.Field, ve.Message)
}

func (v *Validator) AddRule(field string, rule rules.Rule) {
	v.Rules[field] = append(v.Rules[field], rule)
}

// Debug -Testing purposes
func (v *Validator) DebugRules() {
	fmt.Println(v.Rules)
	for field, rules := range v.Rules {
		fmt.Printf("The field is %v \n", field)

		for _, rule := range rules {
			fmt.Printf("Rule type %v and rule param %v \n", rule.Type, rule.Param)
		}
	}
}

func (v *Validator) Validate(input interface{}) []error {
	var errors []error

	val := reflect.ValueOf(input)
	
	//loop through fiels in struct getter from reflect
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // the whole field ex: {Title  string validate:"required|max:1" 0 [0] false}
		fieldName := field.Name
		fieldValue := val.Field(i).Interface() //whole value of field ex: This is title
		
		var fieldValueStr interface{} 

		switch fieldValue.(type) {
			case int:
				fieldValueStr = fieldValue
			default:
				fieldValueStr = fieldValue.(string)
		}
		
		//get rules required,email etc
		rulesTags := field.Tag.Get("validate")

		if rulesTags == "" {
			//no validation tags
			continue
		}

		rulesArr := strings.Split(rulesTags, "|")

		for _, ruleValue := range rulesArr {
			rulesParts := strings.Split(ruleValue, ":")

			var ruleType string = ""
			var ruleNum int = 0
		
			ruleType, ruleVal := parseRules(rulesParts) // get rule type and value ex : max , 100  or email , nil

			switch ruleType {
			//max
				case string(types.Max):
					{
						maxAssigned := ruleVal

						valueLen := len(fieldValueStr.(string))

						if valueLen > maxAssigned.(int) {
							//return fmt.Errorf("Field %s exceeds the maximum length of %d", field, maxAssigned)
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field %s exceeds the maximum length of %d", fieldName, maxAssigned)))
						}

					}
				//min
				case string(types.Min):
					{
						minAssigned := ruleNum

						valueLen := len(fieldValueStr.(string))

						if valueLen < minAssigned {
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field %s is less than the minimum length of %d", fieldName, minAssigned)))
						}
					}
				//required
				case string(types.Required):
					{

						if fieldValueStr.(string) == "" {
							//return fmt.Errorf("%s is required", field)
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field not found: %s", fieldName)))
						}
					}
				//email
				case string(types.Email):
					{

						isEmailValid := isValidEmail(fieldValueStr.(string))

						if !isEmailValid {
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Input field %s must be a valid email format", fieldName)))
						}
					}
				//url
				case string(types.Url):
					{

						isUrl := isURL(fieldValueStr.(string))

						if !isUrl {
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("The %s field must be a valid URL.", fieldName)))
						}
					}
				//active URL
				case string(types.ActiveUrl):
					{

						isActiveUrl := isActiveUrl(fieldValueStr.(string))

						if !isActiveUrl {
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("The %s field must be active URL.", fieldName)))
						}
					}
				//ip
				case string(types.IpFormat):
					{
						isValidIP := isValidIP(fieldValueStr.(string))

						if !isValidIP {
							errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("The %s field must be valid IP format", fieldName)))
						}
					}
				//date
				case string(types.Date):
					{

						isValidDate := isValidDate(fieldValueStr.(string))

						if !isValidDate {
							errors = append(errors, NewValidationError(fieldName, "Invalid date format"))
						}
					}

				case string(types.DateFormat) : 
					{
						dateFormat := ruleVal.(string)
						
						isValidFormat := isValidDateFormat(fieldValueStr.(string),dateFormat)

						if !isValidFormat {
							errors = append(errors, NewValidationError( fieldName,fmt.Sprintf( "Invalid date format The required format is %s",dateFormat)))
						}
						
					}

				case string(types.Between) : 
				{
					isValidBetween := validateBetween(fieldValueStr.(int), ruleVal.(string))

					if !isValidBetween {
						errors = append(errors, NewValidationError( fieldName,fmt.Sprintf( "Number out of range , the number should be between %s ",ruleValue)))
					}

				}
			}
		}

	}
	return errors
}

func parseRules(ruleParts []string) (string, interface{}) {

	ruleName := ruleParts[0]

	if len(ruleParts) <= 1 {
		return ruleName, nil
	}

	if num, err := strconv.Atoi(ruleParts[1]); err == nil {
		return ruleName, num
	} else {
		// If it's not an integer, treat it as a string
		return ruleName, ruleParts[1]
	}
}

func isValidEmail(email string) bool {
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Compile the regular expression pattern
	re := regexp.MustCompile(emailPattern)

	// Use the regular expression to match the email
	return re.MatchString(email)
}

func isURL(givenURL string) bool {
	_, err := url.ParseRequestURI(givenURL)

	if err != nil {
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

func isValidDate(givenDate string) bool {

	dateFormats := []string{"2006-01-02", "02-01-2006", "02/01/2006", "2006/02/01"}
	for _, date := range dateFormats {
		_, err := time.Parse(date, givenDate)
		if err == nil {
			return true
		}
	}

	return false
}

func isValidDateFormat(givenDate string, dateFormat string ) bool { 
	_,err := time.Parse(formatMap[dateFormat],givenDate)
	return err == nil
}

func validateBetween(givenNumber int, ruleValue string) bool {

	parts := strings.Split(ruleValue,",")

	floorNumber, err := strconv.Atoi(parts[0])
    if err != nil {
        return false
    }

    ceilingNumber, err := strconv.Atoi(parts[1])
    if err != nil {
        return false
    }
	
	if givenNumber > ceilingNumber || givenNumber < floorNumber {
		return false
	}

	return true

}
