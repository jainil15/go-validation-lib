package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

// TODO: checkDate
type Validator[T int | float64 | string] struct {
	value       T
	checkString bool
	checkEmail  bool
	checkMin    bool
	checkMax    bool
	min         *int
	max         *int
}

func (v *Validator[T]) Validate() []string {
	errs := make([]string, 0)
	if v.checkString && !v.CheckString() {
		errs = append(errs, "Not a string")
	}
	if v.checkMin && !v.CheckMin() {
		errs = append(errs, "Greater than min")
	}
	if v.checkMax && !v.CheckMax() {
		errs = append(errs, "Greater than max")
	}
	if v.checkEmail && !v.CheckEmail() {
		errs = append(errs, "Not an Email")
	}
	if len(errs) < 1 {
		return nil
	}
	return errs
}

type Validatable[T int | float64 | string] map[string]Validator[T]

func NewValidator[T int | float64 | string](
	value T,
	checkString, checkEmail, checkMin, checkMax bool,
) Validator[T] {
	return Validator[T]{
		value:       value,
		checkEmail:  checkEmail,
		checkString: checkString,
		checkMin:    checkMin,
		checkMax:    checkMax,
	}
}

func (v *Validator[T]) CheckString() bool {
	if _, ok := any(v.value).(string); ok {
		return true
	}
	return false
}

func (v *Validator[T]) CheckEmail() bool {
	regex := `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
	if str, ok := any(v.value).(string); ok {
		return regexp.MustCompile(regex).MatchString(str)
	}
	return false
}

func (v *Validator[T]) CheckMin() bool {
	if v.min != nil {
		if v.value < any(*v.min).(T) {
			return false
		}
	}
	return true
}

func (v *Validator[T]) CheckMax() bool {
	if v.max != nil {
		if v.value > any(*v.max).(T) {
			return false
		}
	}
	return true
}

// func (validatable *Validatable[T]) Get(k string) Validator[T] {
// 	validator, ok := (*validatable)[k]
// 	if !ok {
// 		return nil
// 	}
// 	return validator
// }

func (validatable *Validatable[T]) Parse(v any) (map[string][]string, error) {
	marshalled, err := json.Marshal(validatable)
	if err != nil {
		return nil, err
	}
	unMarshalled := make(map[string]any)
	json.Unmarshal(marshalled, &unMarshalled)
	errMap := make(map[string][]string)
	for k := range unMarshalled {
		validator, ok := (*validatable)[k]
		if !ok {
			return nil, errors.New(fmt.Sprintf("%v not found", k))
		}
		errMap[k] = validator.Validate()
	}
	return errMap, nil
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func main() {
	user := User{
		Email: "jainil",
		Name:  "jainil",
		Age:   12,
	}
	v := Validatable[string]{
		"email": NewValidator(user.Name, false, true, false, false),
		"name":  NewValidator(user.Email, true, true, false, false),
	}
	b, err := v.Parse(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
}
