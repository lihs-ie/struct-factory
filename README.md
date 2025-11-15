# forge

A lightweight, type-safe library for generating Go struct instances for testing purposes. Built with generics to provide compile-time safety and excellent IDE support.

## Features

- Type-safe factory pattern with Go generics
- Deterministic instance generation based on seeds
- Flexible override system using maps or structs
- Built-in factories for common types (strings, enums, maps)
- Case-insensitive field matching
- Support for unexported fields (opt-in)
- Builder pattern for convenient usage

## Installation

```bash
go get github.com/lihs-ie/forge
```

## Requirements

- Go 1.25.3 or later

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/lihs-ie/forge/factory"
)

type User struct {
    name     string
    age      uint
    NickName string
}

type UserProperties struct {
    name     string
    age      uint
    NickName string
}

type UserFactory struct{}

func (f *UserFactory) Instantiate(properties UserProperties) User {
    return User{
        name:     properties.name,
        age:      properties.age,
        NickName: properties.NickName,
    }
}

func (f *UserFactory) Prepare(overrides factory.Partial[UserProperties], seed int64) UserProperties {
    properties := UserProperties{
        name:     fmt.Sprintf("User%d", seed),
        age:      uint(seed % 100),
        NickName: fmt.Sprintf("nick%d", seed),
    }

    if overrides != nil {
        overrides(&properties)
    }

    return properties
}

func (f *UserFactory) Retrieve(instance User) UserProperties {
    return UserProperties{
        name:     instance.name,
        age:      instance.age,
        NickName: instance.NickName,
    }
}

func main() {
    builder := factory.Builder(&UserFactory{})

    // Generate a random user
    user := builder.Build(nil)
    fmt.Printf("%+v\n", user)
    // {name:User1234 age:34 NickName:nick1234} random

    // Override with custom values (unexported fields are supported by default)
    customUser := builder.Build(factory.Override[UserProperties](map[string]any{
        "name":     "Alice",
        "age":      uint(30),
        "NickName": "ally",
    }))
    fmt.Printf("%+v\n", customUser)
    // {name:Alice age:30 NickName:ally}
}
```

## Core Concepts

### Factory Interface

Every factory implements the `Factory[T, P]` interface:

```go
type Factory[T any, P any] interface {
    Instantiate(properties P) T
    Prepare(overrides Partial[P], seed int64) P
    Retrieve(instance T) P
}
```

- `T`: The output type (e.g., `string`, `User`)
- `P`: The properties type used internally
- `Prepare`: Creates properties with optional overrides
- `Instantiate`: Converts properties to the final type
- `Retrieve`: Extracts properties from an existing instance

### Builder Pattern

The builder provides a convenient API:

```go
type BuilderHandle[T any, P any] interface {
    Build(overrides any) T
    BuildList(size int, overrides any) []T
    BuildWith(seed int64, overrides any) T
    BuildListWith(size int, seed int64, overrides any) []T
    Duplicate(instance T, overrides any) T
}
```

## Override System

The `Override` function provides a flexible way to customize generated instances:

### Using Maps

```go
builder.Build(factory.Override[StringProperties](map[string]any{
    "Value": "specific value",
}))
```

### Using Structs

```go
type Literal struct {
    Value string
}

builder.Build(factory.Override[StringProperties](Literal{
    Value: "from struct",
}))
```

### Case-Insensitive Matching

Field names are matched case-insensitively by default:

```go
builder.Build(factory.Override[StringProperties](map[string]any{
    "value": "lowercase works too",
}))
```

### Unexported Fields

By default, unexported fields can be set. To disable this:

```go
builder.Build(factory.Override[StringProperties](
    map[string]any{"value": "test"},
    factory.DisallowUnexported(),
))
```

## Built-in Factories

### StringFactory

Generates random strings with customizable length and character sets:

```go
builder := factory.Builder(&factory.StringFactory{})
randomString := builder.Build(factory.Override[factory.StringProperties](map[string]any{
    "min":        5,
    "max":        20,
    "characters": factory.Characters.Alpha,
}))
```

Available character sets:

- `factory.Characters.Alphanumeric` (default)
- `factory.Characters.Alpha`
- `factory.Characters.Numeric`
- `factory.Characters.Symbol`

### EnumFactory

Selects values from a predefined set:

```go
type Status string
const (
    StatusActive   Status = "active"
    StatusInactive Status = "inactive"
    StatusPending  Status = "pending"
)

var StatusFactory = factory.NewEnumFactory([]Status{
    StatusActive,
    StatusInactive,
    StatusPending,
})

builder := factory.Builder(StatusFactory)
status := builder.Build(nil)
```

Exclude specific values:

```go
status := builder.Build(factory.Override[factory.EnumProperties[Status]](map[string]any{
    "Exclusions": []Status{StatusPending},
}))
```

### MapFactory

Generates maps with random entries:

```go
builder := factory.Builder(factory.NewMapFactory(
    &factory.StringFactory{},
    &factory.StringFactory{},
))

randomMap := builder.Build(nil)
```

Custom entries:

```go
customMap := builder.Build(factory.Override[factory.MapProperties[string, string]](map[string]any{
    "Entries": []factory.MapEntry[string, string]{
        {Key: "key1", Value: "value1"},
        {Key: "key2", Value: "value2"},
    },
}))
```

## Creating Custom Factories

Implement the `Factory[T, P]` interface:

```go
type User struct {
    ID   int64
    Name string
    Age  int
}

type UserProperties struct {
    ID   int64
    Name string
    Age  int
}

type UserFactory struct{}

func (f *UserFactory) Instantiate(properties UserProperties) User {
    return User{
        ID:   properties.ID,
        Name: properties.Name,
        Age:  properties.Age,
    }
}

func (f *UserFactory) Prepare(overrides factory.Partial[UserProperties], seed int64) UserProperties {
    properties := UserProperties{
        ID:   seed,
        Name: factory.Builder(&factory.StringFactory{}).Build(nil),
        Age:  int(seed % 100),
    }

    if overrides != nil {
        overrides(&properties)
    }

    return properties
}

func (f *UserFactory) Retrieve(instance User) UserProperties {
    return UserProperties{
        ID:   instance.ID,
        Name: instance.Name,
        Age:  instance.Age,
    }
}
```

Usage:

```go
builder := factory.Builder(&UserFactory{})

user := builder.Build(factory.Override[UserProperties](map[string]any{
    "Name": "Alice",
    "Age":  30,
}))
```

## Advanced Features

### Deterministic Generation

Use seeds for reproducible results:

```go
user1 := builder.BuildWith(12345, nil)
user2 := builder.BuildWith(12345, nil)
// user1 and user2 have identical values
```

### Batch Generation

Generate multiple instances:

```go
users := builder.BuildList(10, nil) // Generate 10 users
```

With deterministic seeds:

```go
users := builder.BuildListWith(10, 10000, nil)
// Seeds: 10000, 10001, 10002, ..., 10009
```

### Duplication

Clone an existing instance with modifications:

```go
original := builder.Build(nil)
modified := builder.Duplicate(original, factory.Override[UserProperties](map[string]any{
    "Name": "Modified Name",
}))
// modified has the same ID and Age as original, but different Name
```

## Testing

Run all tests:

```bash
go test ./...
```

Check test coverage:

```bash
go test -cover ./factory/...
```

Verify override usage:

```bash
./scripts/check_override_usage.sh
```

## API Reference

### Core Types

- `Factory[T, P]`: Interface for all factories
- `BuilderHandle[T, P]`: Builder interface
- `Overrider[P]`: Type-safe override container
- `Partial[P]`: Function type for property modifications

### Functions

- `Builder[T, P](factory Factory[T, P]) BuilderHandle[T, P]`: Create a builder
- `Override[P](literal any, opts ...OverrideOption) Overrider[P]`: Create an override
- `WithCaseInsensitive() OverrideOption`: Enable case-insensitive matching
- `DisallowUnexported() OverrideOption`: Prevent unexported field access

### Built-in Factories

- `StringFactory`: instantiate via `&factory.StringFactory{}` and configure per-build with `Override()`
- `NewEnumFactory[T](candidates []T) *EnumFactory[T]`
- `NewMapFactory[K, KP, V, VP](keyFactory, valueFactory) *MapFactory[K, KP, V, VP]`

## License

MIT License

Copyright (c) 2025 lihs

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
