package validator

import (
	"fmt"
	"govalidator/pkg/rules"
	"govalidator/types"
)

type Validator struct {
	//ex: title can several rules
	Rules map[string][]rules.Rule
}

//Returns pointer to Validator Struct
func New() *Validator {
	return &Validator{
		Rules: make(map[string][]rules.Rule),
	}
}


func (v *Validator) AddRule(field string , rule rules.Rule) {
	v.Rules[field] = append(v.Rules[field], rule)
}

//Debug purposes
func (v *Validator) DebugRules() {
	for field,rules := range v.Rules {
		fmt.Printf("The field is %v \n",field)

		for _,rule := range rules {
			fmt.Printf("Rule type %v and rule param %v \n" ,rule.Type , rule.Param)
		}
	}
} 

func (v *Validator) Validate(input map[string]interface{}) error {
	//loop through existing Rules 
	for field , rules := range v.Rules {
		for _, rule := range rules {
			switch rule.Type {
				case string(types.Required) : {
					value ,ok := input[field] 

					if !ok || value == nil || value == "" {
						return fmt.Errorf("%s is required", field)
					}
				}

				case string(types.Max) : {
					_ , ok := input[field]
					maxAssigned := rule.Param 

					fmt.Println(input)
					if  !ok {
						return fmt.Errorf("Not found")
					}

					if maxAssigned != nil {
						return fmt.Errorf("Param with the max value is required")
					}
					
				}
			}
		}
	}
	return nil
}