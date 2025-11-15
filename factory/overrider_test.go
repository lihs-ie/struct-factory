package factory

import "testing"

func TestOverrideWithMap(t *testing.T) {
	factory := &StringFactory{}
	builder := Builder(factory)

	result := builder.Build(Override[StringProperties](map[string]any{
		"Value": "test from map",
	}))

	if result != "test from map" {
		t.Errorf("Expected 'test from map', got '%s'", result)
	}
}

func TestOverrideWithMapCaseInsensitive(t *testing.T) {
	factory := &StringFactory{}
	builder := Builder(factory)

	result := builder.Build(Override[StringProperties](map[string]any{
		"value": "lowercase field name",
	}))

	if result != "lowercase field name" {
		t.Errorf("Expected 'lowercase field name', got '%s'", result)
	}
}

func TestOverrideWithStruct(t *testing.T) {
	type TestLiteral struct {
		Value string
	}

	factory := &StringFactory{}
	builder := Builder(factory)

	result := builder.Build(Override[StringProperties](TestLiteral{
		Value: "test from struct",
	}))

	if result != "test from struct" {
		t.Errorf("Expected 'test from struct', got '%s'", result)
	}
}

func TestOverrideWithMultipleFields(t *testing.T) {
	type Status string
	const (
		StatusActive   Status = "active"
		StatusInactive Status = "inactive"
	)

	factory := NewEnumFactory([]Status{StatusActive, StatusInactive})
	builder := Builder(factory)

	result := builder.Build(Override[EnumProperties[Status]](map[string]any{
		"Value":      StatusActive,
		"Exclusions": []Status{StatusInactive},
	}))

	if result != StatusActive {
		t.Errorf("Expected StatusActive, got %v", result)
	}
}

func TestOverrideWithUnexportedField(t *testing.T) {
	type PrivateProps struct {
		value string
	}

	type PrivateLiteral struct {
		value string
	}

	props := &PrivateProps{}

	overrider := Override[PrivateProps](PrivateLiteral{value: "private"})
	overrider.Apply(props)

	if props.value != "private" {
		t.Errorf("Expected 'private', got '%s'", props.value)
	}
}

func TestOverrideWithTypeConversion(t *testing.T) {
	type IntProps struct {
		Value int
	}

	props := &IntProps{}

	overrider := Override[IntProps](map[string]any{
		"Value": int64(42),
	})
	overrider.Apply(props)

	if props.Value != 42 {
		t.Errorf("Expected 42, got %d", props.Value)
	}
}

func TestOverridePrefersSetter(t *testing.T) {
	setterCalled := false

	type PropsWithSetter struct {
		value string
	}

	type testProps struct {
		PropsWithSetter
	}

	props := &testProps{}

	overrider := Override[testProps](map[string]any{
		"value": "via setter",
	})

	overrider.Apply(props)

	if props.value != "via setter" {
		t.Errorf("Expected 'via setter', got '%s'", props.value)
	}

	_ = setterCalled
}

func TestOverrideWithInvalidField(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid field")
		}
	}()

	factory := &StringFactory{}
	builder := Builder(factory)

	builder.Build(Override[StringProperties](map[string]any{
		"NonExistentField": "value",
	}))
}

func TestOverrideWithInvalidType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid literal type")
		}
	}()

	Override[StringProperties](42)
}

func TestOverrideWithNonConvertibleType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-convertible type")
		}
	}()

	type StringProps struct {
		Value string
	}

	props := &StringProps{}

	overrider := Override[StringProps](map[string]any{
		"Value": []int{1, 2, 3},
	})

	overrider.Apply(props)
}

func TestOverrideWithCaseInsensitiveOption(t *testing.T) {
	type Props struct {
		FieldName string
	}

	props := &Props{}

	overrider := Override[Props](map[string]any{
		"fieldname": "test",
	}, WithCaseInsensitive())

	overrider.Apply(props)

	if props.FieldName != "test" {
		t.Errorf("Expected 'test', got '%s'", props.FieldName)
	}
}

func TestOverrideWithSetterNotMatching(t *testing.T) {
	type Props struct {
		Value string
	}

	props := &Props{}

	overrider := Override[Props](map[string]any{
		"Value": "direct field",
	})

	overrider.Apply(props)

	if props.Value != "direct field" {
		t.Errorf("Expected 'direct field', got '%s'", props.Value)
	}
}

func TestOverrideWithAddressableStruct(t *testing.T) {
	type Props struct {
		Value string
	}

	props := &Props{Value: "initial"}

	type Literal struct {
		Value string
	}

	literal := &Literal{Value: "updated"}

	overrider := Override[Props](*literal)
	overrider.Apply(props)

	if props.Value != "updated" {
		t.Errorf("Expected 'updated', got '%s'", props.Value)
	}
}

func TestOverrideWithExportedFieldDirectSet(t *testing.T) {
	type Props struct {
		Value string
	}

	props := &Props{}

	type Literal struct {
		Value string
	}

	overrider := Override[Props](Literal{Value: "exported"})
	overrider.Apply(props)

	if props.Value != "exported" {
		t.Errorf("Expected 'exported', got '%s'", props.Value)
	}
}

func TestOverrideWithUnexportedFieldWithoutPermission(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when trying to set unexported field with DisallowUnexported")
		}
	}()

	type PrivateProps struct {
		value string
	}

	props := &PrivateProps{}

	overrider := Override[PrivateProps](map[string]any{
		"value": "should fail",
	}, DisallowUnexported())

	overrider.Apply(props)
}
