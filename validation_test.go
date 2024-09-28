package main

import (
	"log"
	"testing"
)

func Test_Validation(t *testing.T) {
	t.Run("Test Case 1", func(t *testing.T) {
		val := Validatable{
			"password": New().String().Min(10).Max(100),
			"age":      New().Optional(),
			"email":    New().Email().Max(255),
		}
		errs, err := val.Parse(map[string]interface{}{
			"password": "1234",
			"age":      21,
			"email":    "jainl@",
		})
		if errs == nil {
			log.Println("Empty errors")
			t.Fail()
		}
		if err != nil {
			log.Println(err)
			t.Fail()
		}
		log.Println(errs)
	})
}
