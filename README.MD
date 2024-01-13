# GoValidator

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/yourpackage)](https://goreportcard.com/report/github.com/yourusername/yourpackage)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

GoValidator is a robust and versatile Go package that streamlines the process of validating data through struct information. Designed to simplify and enhance input data validation for your Go applications, GoValidator empowers developers to effortlessly validate data

**Key Features:** 

- **Simple and Intuitive Validation:** Easily validate user input by adding validation tags to your struct fields.

- **Support for Common Validation Rules:** The package supports a variety of common validation rules, including:
  - Required fields
  - Maximum and minimum lengths for strings
  - Email format validation
  - URL format validation
  - Active URL validation
  - IP address format validation
  - Date format validation
  - Between validation 
  - In validation 

## Installation

To install this package, you can use the `go get ` command:

```bash
go get github.com/sCuz12/goinputvalidator
```

## Usage 

```bash
type User struct {
	Title        string `validate:"required|max:10|min:2"`
	Body         string `validate:"required"`
	Email        string `validate:"email|required"`
	Website      string `validate:"url"`
	ActiveURL    string `validate:"active_url"`
	IPAddress    string `validate:"ipformat"`
	Birthdate	 string `validate:"date"`
	AllowAge	 int    `validate:"between:1,40"`
	UserType 	 string `validate:"in:driver,client,admin"`,
	Subscription string`validate:"notIn:free,freebie"`
	SendEmails 	 interface{} `validate:"accepted"`
   Fields		[]string `validate:"size:2"`
   Confirmation_password string 
	Password 	string `validate:"confirmed"`
}
```


``` bash 
func main() {
	v := validator.New()

	user := User{
		Title:  	"This is a test",
		Body:       "Some body content",
		Email:      "johndoe@gmail.com",
		Website:    "https://www.goooogle.com",
		ActiveURL:  "https://georgehadjisavva.dev",
		IPAddress:  "127.0.0.1",
		Birthdate:  "07-10-1995",
		AllowAge:	12,
		UserType: 	"client",
		Subscription : "free",
		SendEmails: 1,
      Fields: []string{"field1", "field2"},
      Confirmation_password : "password",
	   Password: 	"password1", 
	}

	err := v.Validate(user)

	if err != nil {
		for _,singleErr := range err {
			fmt.Println(singleErr)
		}
	}
}
```

---

## Validation Documentation

1. **Required**
   - **Rule Format:** `required`
   - **Example:** `validate:"required|max:20|min:2"`
   - **Description:** Validates that the field is required and has a maximum length of 20 characters, with a minimum length of 2 characters.

2. **Email**
   - **Rule Format:** `email`
   - **Example:** `validate:"email"`
   - **Description:** Validates that the email is in a valid format 

3. **Website**
    - **Rule Format:** `url`
    - **Example:** `validate:"url"`
    - **Description:** Validates that the website is a valid URL.

4. **Active Website**
    - **Rule Format:** `active_url`
    - **Example:** `validate:"active_url"`
    - **Description:** Validates that the website URL is active(accessible).

5. **Ip Format**
    - **Rule Format:** `ipformat`
    - **Example:** `validate:"ipformat"`
    - **Description:** Validates that the IP address is in a valid format.

6. **Date**
    - **Rule Format:** `date`
    - **Example:** `validate:"date"`
    - **Description:** Validates that the birthdate is in the allowed date formats
    - **Notes:** Acceptable formats 
      - "YYYY/MM/DD" format, mapped to "2006/01/02".
      - "YYYY/DD/MM" format, mapped to "2006/02/01".
      - "YYYY-MM-DD" format, mapped to "2006-01-02".
      - "YYYY-DD-MM" format, mapped to "2006-02-01".

7. **Date Format - Specify Format**
    - **Rule Format:** `dateFormat:YYYY-MM-DD`
    - **Example:** `validate:"dateFormat:YYYY-MM-DD"`
    - **Description:** Validates that the birthdate is in the format specified "YYYY-MM-DD".

8. **Between**
   - **Rule Format:** `between:floor,ceiling`
   - **Example:** `validate:"between:1,40"`
   - **Description:** Validates that the input value is between the specified floor and ceiling values (inclusive).

9. **In**
   - **Rule Format:** `in:value1,value2,...,valueN`
   - **Example:** `validate:"in:driver,client,admin"`
   - **Description:** Validates that the input value is one of the specified allowed values.

10. **Accepted**
   - **Rule Format:** `accepted`
   - **Example:** `validate:"accepted"`
   - **Description:** Validates if the input value is one of the accepted values: "yes", "on", 1, or true.
11. **NotIn**
   - **Rule Format:** `notin`
   - **Example:** `validate:"notIn:free,freebie"`
   - **Description:** The NotIn validator ensures that the input value is not one of the specified disallowed values.
12. **Size**
   - **Rule Format:** `size:N`
   - **Example:** `validate:"size:8"`
   - **Description:** Validates that the input data has a specific size. The behavior of this rule varies depending on the data type: 
 
                        - For strings, "N" corresponds to the number of characters in the string.
                        - For arrays or slices, "N" corresponds to the number of elements in the array or slice.
13. **Confirmed**
   - **Rule Format:**  `confirmed`
   - **Example:**  `validate:confirmed`
   - **Description:**  Validates that a field matches its **Confirmation_password** counterpart. For example, it ensures that a **Password** field matches the **Confirmation_password** field to confirm that the user entered the correct password. This rule is  used for password confirmation fields in registration or password reset forms.
