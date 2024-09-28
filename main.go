package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

// TODO: checkDate
type Validator[T any] struct {
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
		errs = append(errs, fmt.Sprintf("Lower than min %v", *v.min))
	}
	if v.checkMax && !v.CheckMax() {
		errs = append(errs, fmt.Sprintf("Greater than max %v", *v.max))
	}
	if v.checkEmail && !v.CheckEmail() {
		errs = append(errs, "Not an Email")
	}
	if len(errs) < 1 {
		return nil
	}
	return errs
}

type Validatable[T any] map[string]Validator[T]

func NewValidator[T any](
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
	if reflect.TypeOf(v.value) == reflect.TypeOf("") {
		return true
	}
	// if _, ok := any(v.value).(string); ok {
	// 	return true
	// }
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
		if str, ok := any(v.value).(string); ok {
			if len(str) < *v.min {
				return false
			}
		}
	}
	return true
}

func (v *Validator[T]) CheckMax() bool {
	if v.min != nil {
		if str, ok := any(v.value).(string); ok {
			if len(str) > *v.max {
				return false
			}
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
		if errs := validator.Validate(); errs != nil {
			errMap[k] = errs
		}
	}
	return errMap, nil
}

func (v Validator[T]) Email() Validator[T] {
	v.checkEmail = true
	return v
}

func (v Validator[T]) Min(min int) Validator[T] {
	v.checkMin = true
	v.min = &min
	return v
}

func (v Validator[T]) Max(max int) Validator[T] {
	v.checkMax = true
	v.max = &max
	return v
}

func (v Validator[T]) String() Validator[T] {
	v.checkString = true
	return v
}

type User struct {
	Email string `json:"email"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

func New[T int | string](value T) Validator[any] {
	return Validator[any]{
		value: value,
	}
}

func main() {
	user := User{
		Email: "jainil@gmail.com",
		Name:  "jainil",
		Age:   12,
	}

	v := Validatable[any]{
		"email": New(user.Email).String().Min(503),
		"name":  New(user.Name).String(),
		"age":   New(user.Age).String().Min(503),
	}
	b, err := v.Parse(user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v\n", b)
}
