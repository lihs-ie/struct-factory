package factory

import (
	"fmt"
	"maps"
	"testing"
)

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
	return User(properties)
}

func (f *UserFactory) Prepare(overrides Partial[UserProperties], seed int64) UserProperties {
	properties := UserProperties{
		ID:   seed,
		Name: fmt.Sprintf("User%d", seed),
		Age:  int(seed % 100),
	}

	if overrides != nil {
		overrides(&properties)
	}

	return properties
}

func (f *UserFactory) Retrieve(instance User) UserProperties {
	return UserProperties(instance)
}

type profile struct {
	first string
	age   int
}

type complexShape struct {
	Profile profile
	Tags    []string
	Meta    map[string]int
}

type complexProperties struct {
	profile profile
	tags    []string
	meta    map[string]int
}

type complexFactory struct{}

func (f *complexFactory) Instantiate(props complexProperties) complexShape {
	return complexShape{
		Profile: props.profile,
		Tags:    append([]string(nil), props.tags...),
		Meta:    maps.Clone(props.meta),
	}
}

func (f *complexFactory) Prepare(overrides Partial[complexProperties], seed int64) complexProperties {
	props := complexProperties{
		profile: profile{first: fmt.Sprintf("profile-%d", seed), age: int(seed % 50)},
		tags:    []string{"alpha", "beta"},
		meta:    map[string]int{"seed": int(seed)},
	}
	if overrides != nil {
		overrides(&props)
	}
	return props
}

func (f *complexFactory) Retrieve(shape complexShape) complexProperties {
	return complexProperties{
		profile: shape.Profile,
		tags:    append([]string(nil), shape.Tags...),
		meta:    maps.Clone(shape.Meta),
	}
}

func TestFactoryMethods(t *testing.T) {
	factory := &UserFactory{}

	properties := factory.Prepare(nil, 54321)
	if properties.ID != 54321 {
		t.Errorf("Expected ID 54321, got %d", properties.ID)
	}

	if properties.Name != "User54321" {
		t.Errorf("Expected Name 'User54321', got '%s'", properties.Name)
	}

	properties2 := factory.Prepare(func(p *UserProperties) {
		p.Name = "OverriddenName"
	}, 99999)
	if properties2.Name != "OverriddenName" {
		t.Errorf("Expected Name 'OverriddenName', got '%s'", properties2.Name)
	}

	user := factory.Instantiate(properties)
	if user.ID != properties.ID {
		t.Errorf("Expected ID %d, got %d", properties.ID, user.ID)
	}
	if user.Name != properties.Name {
		t.Errorf("Expected Name '%s', got '%s'", properties.Name, user.Name)
	}

	retrievedProps := factory.Retrieve(user)
	if retrievedProps.ID != user.ID {
		t.Errorf("Expected ID %d, got %d", user.ID, retrievedProps.ID)
	}
	if retrievedProps.Name != user.Name {
		t.Errorf("Expected Name '%s', got '%s'", user.Name, retrievedProps.Name)
	}
	if retrievedProps.Age != user.Age {
		t.Errorf("Expected Age %d, got %d", user.Age, retrievedProps.Age)
	}
}

func TestBuilder(t *testing.T) {
	factory := &UserFactory{}
	builder := Builder(factory)

	user1 := builder.Build(nil)
	if user1.ID == 0 {
		t.Error("Expected non-zero ID")
	}

	user2 := builder.Build(Override[UserProperties](map[string]any{
		"Age": 30,
	}))
	if user2.Age != 30 {
		t.Errorf("Expected Age 30, got %d", user2.Age)
	}

	users := builder.BuildList(5, nil)
	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}

	seenIDs := make(map[int64]bool)
	for _, user := range users {
		if seenIDs[user.ID] {
			t.Errorf("Duplicate ID found: %d", user.ID)
		}
		seenIDs[user.ID] = true
	}

	user3 := builder.BuildWith(99999, nil)
	if user3.ID != 99999 {
		t.Errorf("Expected ID 99999, got %d", user3.ID)
	}

	usersWithSeeds := builder.BuildListWith(3, 10000, nil)
	if len(usersWithSeeds) != 3 {
		t.Errorf("Expected 3 users, got %d", len(usersWithSeeds))
	}
	for i, user := range usersWithSeeds {
		expectedID := int64(10000 + i)
		if user.ID != expectedID {
			t.Errorf("Expected ID %d, got %d", expectedID, user.ID)
		}
	}

	duplicated := builder.Duplicate(user1, Override[UserProperties](map[string]any{
		"Name": "Duplicated",
	}))
	if duplicated.ID != user1.ID {
		t.Errorf("Expected same ID %d, got %d", user1.ID, duplicated.ID)
	}
	if duplicated.Name != "Duplicated" {
		t.Errorf("Expected Name 'Duplicated', got '%s'", duplicated.Name)
	}
	if duplicated.Age != user1.Age {
		t.Errorf("Expected Age %d, got %d", user1.Age, duplicated.Age)
	}
}

func TestUserFactoryWithInlineOverride(t *testing.T) {
	factory := &UserFactory{}
	builder := Builder(factory)

	user := builder.Build(Override[UserProperties](map[string]any{
		"Name": "Alice",
		"Age":  30,
	}))

	if user.Name != "Alice" {
		t.Errorf("Expected Name 'Alice', got '%s'", user.Name)
	}

	if user.Age != 30 {
		t.Errorf("Expected Age 30, got %d", user.Age)
	}
}

func TestUserFactoryWithOverrideLiteral(t *testing.T) {
	factory := &UserFactory{}
	builder := Builder(factory)

	user := builder.Build(Override[UserProperties](map[string]any{
		"Name": "Bob",
		"Age":  25,
	}))

	if user.Name != "Bob" {
		t.Errorf("Expected Name 'Bob', got '%s'", user.Name)
	}

	if user.Age != 25 {
		t.Errorf("Expected Age 25, got %d", user.Age)
	}
}

func TestBuilderGeneratesNestedStructFields(t *testing.T) {
	builder := Builder(&complexFactory{})
	shape := builder.Build(nil)

	if shape.Profile.first == "" {
		t.Fatal("expected profile name to be set")
	}
	if len(shape.Tags) == 0 {
		t.Fatal("expected tags to be generated")
	}
	if len(shape.Meta) == 0 {
		t.Fatal("expected metadata to be generated")
	}

	// ensure retrieving and duplicating round-trips nested structures.
	duplicated := builder.Duplicate(shape, nil)
	if duplicated.Profile.first != shape.Profile.first {
		t.Fatalf("duplicate lost profile: %s vs %s", duplicated.Profile.first, shape.Profile.first)
	}
	if len(duplicated.Tags) != len(shape.Tags) {
		t.Fatalf("duplicate lost tags: %v", duplicated.Tags)
	}
	if len(duplicated.Meta) != len(shape.Meta) {
		t.Fatalf("duplicate lost meta: %v", duplicated.Meta)
	}
}

func TestBuilderOverridesSliceAndMapFields(t *testing.T) {
	builder := Builder(&complexFactory{})
	profileOverride := profile{first: "override", age: 42}
	override := Override[complexProperties](map[string]any{
		"profile": profileOverride,
		"tags":    []string{"go", "forge"},
		"meta":    map[string]int{"answer": 42},
	})

	shape := builder.Build(override)

	if shape.Profile != profileOverride {
		t.Fatalf("expected profile %+v, got %+v", profileOverride, shape.Profile)
	}
	if len(shape.Tags) != 2 || shape.Tags[0] != "go" || shape.Tags[1] != "forge" {
		t.Fatalf("unexpected tags: %v", shape.Tags)
	}
	if len(shape.Meta) != 1 || shape.Meta["answer"] != 42 {
		t.Fatalf("unexpected meta: %v", shape.Meta)
	}
}

func ExampleBuilder() {
	builder := Builder(&UserFactory{})

	user := builder.Build(nil)
	fmt.Printf("User: %+v\n", user)

	users := builder.BuildList(5, nil)
	fmt.Printf("Built %d users\n", len(users))

	userWithSeed := builder.BuildWith(12345, nil)
	fmt.Printf("User with seed: %+v\n", userWithSeed)

	usersWithSeeds := builder.BuildListWith(3, 67890, nil)
	fmt.Printf("Built %d users with sequential seeds\n", len(usersWithSeeds))
}
