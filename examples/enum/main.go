package main

import (
	"fmt"

	"github.com/lihs-ie/forge/factory"
)

type Status string

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
)

func main() {
	enumFactory := factory.NewEnumFactory([]Status{StatusActive, StatusInactive})
	builder := factory.Builder(enumFactory)

	value := builder.Build(factory.Override[factory.EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusInactive},
	}))

	fmt.Println("enum value:", value)
}
