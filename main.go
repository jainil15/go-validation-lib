package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// TODO: checkDate
type Validator struct {
	checkString bool
	checkEmail  bool
	checkMin    bool
	checkMax    bool
	optional    bool
	min         *int
	max         *int
}

func (v *Validator) Validate(val any) []string {
	if !v.optional && val == nil {
		return []string{"Required"}
	}
	if v.optional && val == nil {
		return nil
	}
	errs := make([]string, 0)
	if v.checkString && !v.CheckString(val) {
		errs = append(errs, "Not a string")
	}
	if v.checkMin && !v.CheckMin(val) {
		errs = append(errs, fmt.Sprintf("Lower than min %v", *v.min))
	}
	if v.checkMax && !v.CheckMax(val) {
		errs = append(errs, fmt.Sprintf("Greater than max %v", *v.max))
	}
	if v.checkEmail && !v.CheckEmail(val) {
		errs = append(errs, "Not an Email")
	}
	if len(errs) < 1 {
		return nil
	}
	return errs
}

type Validatable map[string]Validator

func NewValidator[T any](
	checkString, checkEmail, checkMin, checkMax bool,
) Validator {
	return Validator{
		checkEmail:  checkEmail,
		checkString: checkString,
		checkMin:    checkMin,
		checkMax:    checkMax,
	}
}

func (v *Validator) CheckString(val any) bool {
	if reflect.TypeOf(val) == reflect.TypeOf("") {
		return true
	}
	return false
}

func (v *Validator) CheckEmail(val any) bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	if str, ok := any(val).(string); ok {
		return regexp.MustCompile(regex).MatchString(str)
	}
	return false
}

func (v *Validator) CheckMin(val any) bool {
	if v.min != nil {
		if str, ok := any(val).(string); ok {
			if len(str) < *v.min {
				return false
			}
		}
	}
	return true
}

func (v *Validator) CheckMax(val any) bool {
	if v.min != nil {
		if str, ok := any(val).(string); ok {
			if len(str) > *v.max {
				return false
			}
		}
	}
	return true
}

func (validatable *Validatable) Parse(v any) (map[string][]string, error) {
	marshalled, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	unMarshalled := make(map[string]any)
	json.Unmarshal(marshalled, &unMarshalled)
	errMap := make(map[string][]string)
	for k := range *validatable {
		validator, ok := (*validatable)[k]
		if !ok {
			return nil, errors.New(fmt.Sprintf("%v not found", k))
		}
		if errs := validator.Validate(unMarshalled[k]); errs != nil {
			errMap[k] = errs
		}
	}
	if len(errMap) == 0 {
		return nil, nil
	}
	return errMap, nil
}

func (v Validator) Email() Validator {
	v.checkEmail = true
	return v
}

func (v Validator) Min(min int) Validator {
	v.checkMin = true
	v.min = &min
	return v
}

func (v Validator) Max(max int) Validator {
	v.checkMax = true
	v.max = &max
	return v
}

func (v Validator) String() Validator {
	v.checkString = true
	return v
}

func (v Validator) Optional() Validator {
	v.optional = true
	return v
}

type User struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

func New() Validator {
	return Validator{}
}

func main() {
	user := User{
		Email:    "jainil@gmail.com",
		Name:     "jainil",
		Age:      12,
		Password: "1244",
	}

	v := Validatable{
		"email":     New().String().Min(2),
		"name":      New().String().Min(2),
		"age":       New().Max(2).Min(20),
		"password":  New().String().Min(10).Max(20).Optional(),
		"Interface": New().String(),
	}
	userMap := map[string]interface{}{
		"email": "jainil@gmail.com",
		"name":  "jainil",
		"age":   12,
	}
	b, err := v.Parse(user)
	if err != nil {
		fmt.Println(err)
	}
	b, err = v.Parse(userMap)
	if err != nil {
		fmt.Println(err)
	}
	if b != nil {
		fmt.Printf("%v\n", b)
	}
}
