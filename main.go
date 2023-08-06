package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jszwec/csvutil"
)

type Validator[T any] func(T) error

func Validate[T any](data []T, validators ...Validator[T]) (err error) {
	for i, row := range data {
		var colErrors error
		for _, validator := range validators {
			if err := validator(row); err != nil {
				colErrors = errors.Join(colErrors, err)
			}
		}
		if colErrors != nil {
			err = errors.Join(err, fmt.Errorf("row %d: [%v]", i, colErrors.Error()))
		}
	}
	return
}

type User struct {
	Name string
	Age  int
}

func main() {
	file, err := os.Open("./demo.csv")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    csvData, err := ioutil.ReadFile("./demo.csv")
    if err != nil {
        panic(err)
    }
    var users []User
    if err := csvutil.Unmarshal(csvData, &users); err != nil {
        panic(err)
    }

	if err := Validate(users, checkName, checkAge); err != nil {
		fmt.Println(err)
	}
}

func checkName(user User) error {
	if user.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	return nil
}

func checkAge(user User) error {
	if user.Age > 18 {
		return fmt.Errorf("age must be less than 18")
	}
	return nil
}