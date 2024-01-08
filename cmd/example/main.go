package main

import (
	"fmt"
	"govalidator/pkg/rules"
	"govalidator/pkg/validator"
)

func main () {
	v := validator.New();

	v.AddRule("title",rules.Rule{Type: "required"})
	v.AddRule("title" , rules.Rule{Type: "max",Param: 100})
	v.AddRule("body" , rules.Rule{Type: "max",Param: 100})
	//v.DebugRules();

	err := v.Validate(map[string]interface{}{
		"title": "asda",
    	"body":  "Some body content",
	})

	if err != nil {
		for _,singleErr := range err {
			fmt.Println(singleErr)
		}
	}

	fmt.Println("And the life goes on")


}