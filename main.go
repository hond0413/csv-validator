package main

import (
	"errors"
	"fmt"
	"io/ioutil"

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

func userTest() {
    csvData, err := ioutil.ReadFile("./demo.csv")
    if err != nil {
        panic(err)
    }
    var users []User
    if err := csvutil.Unmarshal(csvData, &users); err != nil {
        panic(err)
    }

	if err := Validate(users, checkUserName, checkUserAge); err != nil {
		fmt.Println(err)
	}
}

func checkUserName(user User) error {
	if user.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	return nil
}

func checkUserAge(user User) error {
	if user.Age > 18 {
		return fmt.Errorf("age must be less than 18")
	}
	return nil
}

type ItemType string

const (
    Food  ItemType = "food"
    Drink ItemType = "drink"
)

type Item struct {
	Name string
	Price int
	Type ItemType
}

func itemTest() {
    csvData, err := ioutil.ReadFile("./demo1.csv")
    if err != nil {
        panic(err)
    }
    var items []Item
    if err := csvutil.Unmarshal(csvData, &items); err != nil {
        panic(err)
    }

	if err := Validate(items, checkItemName, checkItemPrice, checkItemType); err != nil {
		fmt.Println(err)
	}
}

func checkItemName(item Item) error {
	if item.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	return nil
}

func checkItemPrice(item Item) error {
	if item.Price < 0 {
		return fmt.Errorf("price must be greater than 0")
	}
	return nil
}

func checkItemType(item Item) error {
	if item.Type != Food && item.Type != Drink {
		return fmt.Errorf("type must be food or drink")
	}
	return nil
}

func main() {
	fmt.Println("User Test")
	userTest()
	fmt.Println("Item Test")
	itemTest()
}