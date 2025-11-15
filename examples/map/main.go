package main

import (
	"fmt"

	"github.com/lihs-ie/forge/factory"
)

func main() {
	mapFactory := factory.NewMapFactory(&factory.StringFactory{}, &factory.StringFactory{Min: 3, Max: 6})
	builder := factory.Builder(mapFactory)

	entries := []factory.MapEntry[string, string]{
		{Key: "id", Value: "123"},
		{Key: "name", Value: "neo"},
	}

	custom := builder.Build(factory.Override[factory.MapProperties[string, string]](map[string]any{
		"Entries": entries,
	}))

	fmt.Println("map entries:", custom)
}
