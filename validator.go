package validator

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sCuz12/goinputvalidator/types"
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

}

type ConfirmPasswordInfo struct {
	Exist bool
	Value string
}

// Returns pointer to Validator Struct
func New() *Validator {
	return &Validator{
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


func (v *Validator) Validate(input interface{}) []error {
	var errors []error

	done := make(chan ConfirmPasswordInfo)

	val := reflect.ValueOf(input)

	go checkConfirmationPasswordExists(&val,done)

	confirmationPasswordInfo := <-done

	
	//loop through fiels in struct getter from reflect
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i) // the whole field ex: {Title  string validate:"required|max:1" 0 [0] false}
		fieldName := field.Name
		fieldValue := val.Field(i).Interface() //whole value of field ex: This is title
		
		var fieldValueStr interface{} 
		
		switch fieldValue.(type) {
			case int:
				fieldValueStr = fieldValue
			case bool :
				fieldValueStr = fieldValue
			case []string : 
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
						minAssigned := ruleVal
						valueLen 	:= len(fieldValueStr.(string))
				
						if valueLen < minAssigned.(int) {
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
				//Dateformat
				case string(types.DateFormat) : 
					{
						dateFormat := ruleVal.(string)
						
						isValidFormat := isValidDateFormat(fieldValueStr.(string),dateFormat)

						if !isValidFormat {
							errors = append(errors, NewValidationError( fieldName,fmt.Sprintf( "Invalid date format The required format is %s",dateFormat)))
						}
						
					}
				//Between
				case string(types.Between) : 
				{
					isValidBetween := validateBetween(fieldValueStr.(int), ruleVal.(string))
					if !isValidBetween {
						errors = append(errors, NewValidationError( fieldName,fmt.Sprintf( "Number out of range , the number should be between %s ",ruleValue)))
					}

				}
				//In
				case string(types.In) : {
					allowedValues := strings.Split(ruleVal.(string),",")

					isInputIn := validateIN(fieldValueStr.(string),allowedValues)

					if !isInputIn {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Value '%s' is not allowed. Allowed values are: %s", fieldValueStr, allowedValues)))
					}
				}
				//Accepted
				case  string(types.Accepted) : {
					supportedValues := []interface{} {"yes", "on",1,true}

					isAcceptedValue := validateAccepted(fieldValueStr,supportedValues)

					if !isAcceptedValue {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Value '%v' is not allowed. Allowed values are: %v", fieldValueStr, supportedValues)))
					}
					
				}
				//Not in 
				case string(types.NotIn) : {
					notAllowedValues := strings.Split(ruleVal.(string),",")

					isNotIn := validateIN(fieldValueStr.(string),notAllowedValues)
					
					if(isNotIn) {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Value '%s' is not allowed. Must not be one of: %s", fieldValueStr, notAllowedValues)))
					}
				}
				//Size 
				case string(types.Size) : {
					isCorrectSize := validateSize(fieldValueStr,ruleVal.(int))
					
					if !isCorrectSize {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Field '%s' should have exactly %d elements", fieldName, ruleVal.(int))))
					}
				}
				//Confirmed
				case string(types.Confirmed) : {
					//First way is to loop through all fields and check if the confirmation_password exist
					if !confirmationPasswordInfo.Exist {
						errors = append(errors, NewValidationError(fieldName, "Confirmation password is required when using confirmed validation. Please provide a confirmation password for this field."))
						break
					}

					if(fieldValue != confirmationPasswordInfo.Value) {
						errors = append(errors, NewValidationError(fieldName, "Confirmation_password and Password fields should match"))
					}
				}
					//doesnt END
				case string(types.Doesnt_end_with) :{
					suffixchecker := SuffixChecker{}

					canPass := suffixchecker.IsValid(fieldValueStr.(string),ruleVal.(string))
					
					if !canPass {
						errors = append(errors, NewValidationError(fieldName,fmt.Sprintf("Input cannot end with value: %s", ruleVal.(string))))
					}
				}
				//doesnt START
				case string(types.Doesnt_start_with) : {
					prefixchecker := PrefixChecker{}
					
					canPass := prefixchecker.IsValid(fieldValueStr.(string),ruleVal.(string))

					if !canPass {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Input cannot start with value : %s", ruleVal.(string))))
					}

				}
				//mac address format
				case string(types.MacAddressFormat) : {
					isMacFormat := isMacAddressFormat(fieldValueStr.(string))

					if(!isMacFormat) {
						errors = append(errors, NewValidationError(fieldName, "Invalid MAC address format. The MAC address should have the format 'XX:XX:XX:XX:XX:XX' where X is a hexadecimal digit (0-9, A-F, a-f)."))
					}
				}
				//greater than
				case string(types.GreaterThan) : {
					comparator := GreaterThanComparator{}

					canPass:= comparator.IsValid(fieldValueStr.(int), ruleVal.(int))

					if(!canPass) {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Invalid value. The value must be greater than the specified threshold which is : %v",ruleVal)))
					}
				}
				//less than 
				case string(types.LessThan) : {

					comparator := LessThanComparator{}

					canPass := comparator.IsValid(fieldValueStr.(int),ruleVal.(int))

					if !canPass {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Invalid value. The value must be less than the specified floor  which is : %v",ruleVal)))
					}
				}
				//Greater equal than
				case string(types.GreaterEqualThan) : {
					comparator := GreaterEqualThanComparator{}

					canPass := comparator.IsValid(fieldValueStr.(int), ruleVal.(int))
					if !canPass {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Invalid value. The value must be greater OR equal than the specified threshold floor  which is : %v",ruleVal)))
					}
				}
				case string(types.LessThanEqual) : {
					comparator := LessThanEqualComparator{}

					canPass := comparator.IsValid(fieldValueStr.(int),ruleVal.(int))

					if !canPass {
						errors = append(errors, NewValidationError(fieldName, fmt.Sprintf("Invalid value. The value must be less OR equal than the specified floor  which is : %v",ruleVal)))
					}

				}
			}
		}

	}
	return errors
}

func checkConfirmationPasswordExists(fields *reflect.Value,done chan<-ConfirmPasswordInfo) {
	var confirmationPasswordInfo ConfirmPasswordInfo

	for i := 0; i < fields.NumField(); i++ {

		field := fields.Type().Field(i) // the whole field ex: {Title  string validate:"required|max:1" 0 [0] false}
		confirmationPasswordVal := fields.Field(i).Interface()

		fieldName := field.Name

		if fieldName == "Confirmation_password" {
	
			confirmationPasswordInfo = ConfirmPasswordInfo{
				Exist: true,
				Value: confirmationPasswordVal.(string),	
			}


			break
		}
	}
	done<-confirmationPasswordInfo
}


func parseRules(ruleParts []string) (string, interface{}) {
	var ruleValue interface{}

	ruleName := ruleParts[0]

	if len(ruleParts) <= 1 {
		return ruleName, nil
	}

	if num, err := strconv.Atoi(ruleParts[1]); err == nil {
		ruleValue = num
	} else {
		// If it's not an integer, treat it as a string
		ruleValue = ruleParts[1]
	}
	return ruleName,ruleValue
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
func validateIN(givenString string , allowedValuesSlice []string) bool {

	for _,value := range allowedValuesSlice {
		if value == givenString {
			 return true
		} 
	}

	return false
}

func validateAccepted(givenInput interface{}, supportedValues []interface{}) bool {

	for _ , val := range supportedValues {
		if val == givenInput {
			return true
		} 
	}
	return false

}

func validateSize(givenInput interface{},allowedSize int) bool {
	switch v := givenInput.(type) {
	case string :
		return len(v) == allowedSize
	case []int : 
		return len(v) == allowedSize

	case []string : 
		return len(v) == allowedSize
	default:
		return false

	}
}

func isMacAddressFormat(input string) bool {
	regexPattern := "^([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})$"
	//compile the regex
	regex := regexp.MustCompile(regexPattern)

	//match the input string 
	return regex.MatchString(input)
}
