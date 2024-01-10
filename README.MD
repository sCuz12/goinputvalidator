# GoValidator

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/yourpackage)](https://goreportcard.com/report/github.com/yourusername/yourpackage)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://opensource.org/licenses/MIT)

GoValidator is a robust and versatile Go package designed to simplify and enhance input data validation for your Go applications. It empowers developers to effortlessly validate and sanitize incoming data, ensuring data integrity and security.

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

- **Customizable Rules:** You can define custom validation rules by adding them to your struct tags.



## Installation

To install this package, you can use the `go get ` command:

```bash
go get https://github.com/sCuz12/goinputvalidator
```

## Usage 

```bash
type User struct {
	Title       string `validate:"required|max:10|min:2"`
	Body        string `validate:"required"`
	Email       string `validate:"email|required"`
	Website     string `validate:"url"`
	ActiveURL   string `validate:"active_url"`
	IPAddress   string `validate:"ipformat"`
	Birthdate	string `validate:"date"`
}
```


``` bash 
func main() {
	v := validator.New()

	user := User{
		Title:  "aasda",
		Body:       "Some body content",
		Email:      "johndoe@gmail.com",
		Website:    "https://www.goooogle.com",
		ActiveURL:  "https://georgehadjisavva.dev",
		IPAddress:  "127.0.0.1",
		Birthdate:  "07-10-1995",
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

#### _Supported Date Formats_

The validation package supports multiple date formats for date fields. These formats determine how dates should be formatted when specifying date-related rules in your validation tags. Below are the supported date formats:

1. **ISO 8601 Format (YYYY-MM-DD)**
   - Format: "2006-01-02"
   - Example: "2022-12-31"

2. **Day-Month-Year (DD-MM-YYYY)**
   - Format: "02-01-2006"
   - Example: "31-12-2022"

3. **Month-Day-Year (MM-DD-YYYY)**
   - Format: "01-02-2006"
   - Example: "12-31-2022"

4. **Year-Month-Day (YYYY-MM-DD)**
   - Format: "2006-01-02"
   - Example: "2022-12-31"

5. **Alternate Format (YYYY/MM/DD)**
   - Format: "2006/01/02"
   - Example: "2022/12/31"

These date formats provide flexibility in how you specify date values in your validation rules. You can use any of these formats when defining date-related rules for your struct fields. The validation package will attempt to parse dates using the specified format(s) and validate them accordingly.

It's important to use the format that matches the date representation you expect in your application. Additionally, when working with date values, ensure consistency in the format you use to avoid unexpected parsing issues.

Please note that the order of day, month, and year components may vary depending on the chosen format. Always use the appropriate format for your specific use case.
