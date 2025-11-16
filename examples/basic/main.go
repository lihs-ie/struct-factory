package main

import (
	"fmt"

	"github.com/lihs-ie/forge/factory"
)

type User struct {
	name string
	age  int
}

type UserProperties struct {
	name string
	age  int
}

type UserFactory struct{}

func (f *UserFactory) Instantiate(props UserProperties) User {
	return User(props)
}

func (f *UserFactory) Prepare(overrides factory.Partial[UserProperties], seed int64) UserProperties {
	props := UserProperties{
		name: fmt.Sprintf("user-%d", seed),
		age:  int(seed % 100),
	}
	if overrides != nil {
		overrides(&props)
	}
	return props
}

func (f *UserFactory) Retrieve(instance User) UserProperties {
	return UserProperties(instance)
}

func main() {
	builder := factory.Builder(&UserFactory{})

	user := builder.Build(nil)
	fmt.Printf("random: %+v\n", user)

	overridden := builder.Build(factory.Override[UserProperties](map[string]any{
		"name": "alice",
		"age":  30,
	}))
	fmt.Printf("override: %+v\n", overridden)
}
