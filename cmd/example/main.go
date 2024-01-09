package main

import (
	"fmt"
	"govalidator/pkg/validator"
)
type User struct {
	Title       string `json:"name" validate:"required|max:10|min:2"`
	Body        string `validate:"required"`
	Email       string `validate:"email|required"`
	Website     string `validate:"url"`
	WebsiteURL  string `validate:"active_url"`
	IPAddress   string `validate:"ipformat"`
}

func main() {
	// testValidate()

	testValidateSchema();

}


func testValidateSchema() {
	v := validator.New()

	user := User{
		Title:  	"JohnDoe",
		Body:       "Some body content",
		Email:      "johndoe@gmail.com",
		Website:    "https://www.google.com",
		WebsiteURL: "https://georgehadjisavva.dev",
		IPAddress:  "127.0.0.1.0.0.0",
	}

	err := v.Validate(user)

	if err != nil {
		for _,singleErr := range err {
			fmt.Println(singleErr)
		}
	}

}