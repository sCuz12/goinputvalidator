package main

import (
	"fmt"
	"govalidator/pkg/rules"
	"govalidator/pkg/validator"
)

func main () {
	v := validator.New();

	v.AddRule("title",rules.Rule{Type: "required"})
	v.AddRule("title" , rules.Rule{Type: "max",Param: 2})
	v.AddRule("body" , rules.Rule{Type: "max",Param: 2})
	//v.DebugRules();

	err := v.Validate(map[string]interface{}{
		"title": "asda",
    	"body":  "Some body content",
	})

	if err != nil {
		fmt.Println(err)
	}


}