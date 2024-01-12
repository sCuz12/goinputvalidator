package main

import (
	"fmt"
	"goinputvalidator/pkg/validator"
)
type User struct {
	Title       string `json:"name" validate:"required|max:20|min:2"`
	Body        string `validate:"required"`
	Email       string `validate:"email|required"`
	Website     string `validate:"url"`
	WebsiteURL  string `validate:"active_url"`
	IPAddress   string `validate:"ipformat"`
	//Birthdate	string `validate:"date"`
	// Birthdate	string `validate:"dateFormat:YYYY/MM/DD"`
	Birthdate	string `validate:"dateFormat:YYYY-MM-DD"`
	AllowAge	int `validate:"between:1,40"`
	UserType 	string `validate:"in:driver,client,admin"`
	SendEmails 	interface{} `validate:"accepted"`
	Subscription string `validate:"notIn:free,freebie"`
	Fields		[]string `validate:"size:2"`
}

func main() {
	// testValidate()

	testValidateSchema();

}


func testValidateSchema() {
	v := validator.New()

	user := User{
		Title:  	"Kostas rmanto",
		Body:       "Some body content",
		Email:      "johndoe@gmail.com",
		Website:    "https://www.google.com",
		WebsiteURL: "https://georgehadjisavva.dev",
		IPAddress:  "127.0.0.1",
		Birthdate:  "2022-02-15",
		AllowAge:	12,
		UserType: "client",
		SendEmails: 1,
		Subscription : "asda",
		Fields: []string{"field1", "field2"},
	}

	err := v.Validate(user)

	if err != nil {
		for _,singleErr := range err {
			fmt.Println(singleErr)
		}
	}

}