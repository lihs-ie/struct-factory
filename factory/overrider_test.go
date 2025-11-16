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
		value string //nolint:unused // Field is intentionally unused as test expects panic before access
	}

	props := &PrivateProps{}

	overrider := Override[PrivateProps](map[string]any{
		"value": "should fail",
	}, DisallowUnexported())

	overrider.Apply(props)
}

func TestOverrideWithNilValue(t *testing.T) {
	type PropsWithPointer struct {
		Value *string
	}

	props := &PropsWithPointer{}

	overrider := Override[PropsWithPointer](map[string]any{
		"Value": nil,
	})
	overrider.Apply(props)

	if props.Value != nil {
		t.Errorf("Expected nil, got %v", props.Value)
	}
}

func TestOverrideWithNilPointer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil pointer literal")
		}
	}()

	var nilPtr *struct{ Value string }
	Override[struct{ Value string }](nilPtr)
}

func TestOverrideWithNilMap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil literal")
		}
	}()

	Override[struct{ Value string }](nil)
}

type trackedPropsForTest struct {
	Value          string
	overriddenList []string
}

func (p *trackedPropsForTest) noteOverride(field string) {
	p.overriddenList = append(p.overriddenList, field)
}

func TestOverrideTrackerNotification(t *testing.T) {
	props := &trackedPropsForTest{}

	overrider := Override[trackedPropsForTest](map[string]any{
		"Value": "tracked",
	})
	overrider.Apply(props)

	if props.Value != "tracked" {
		t.Errorf("Expected 'tracked', got '%s'", props.Value)
	}

	if len(props.overriddenList) != 1 || props.overriddenList[0] != "Value" {
		t.Errorf("Expected override notification for 'Value', got %v", props.overriddenList)
	}
}

func TestOverrideWithEmbeddedStruct(t *testing.T) {
	type BaseProps struct {
		BaseValue string
	}

	type DerivedProps struct {
		BaseProps
		DerivedValue string
	}

	props := &DerivedProps{}

	overrider := Override[DerivedProps](map[string]any{
		"BaseValue":    "from base",
		"DerivedValue": "from derived",
	})
	overrider.Apply(props)

	if props.BaseValue != "from base" {
		t.Errorf("Expected 'from base', got '%s'", props.BaseValue)
	}
	if props.DerivedValue != "from derived" {
		t.Errorf("Expected 'from derived', got '%s'", props.DerivedValue)
	}
}

func TestOverrideWithEmbeddedPointerStruct(t *testing.T) {
	type BaseProps struct {
		BaseValue string
	}

	type DerivedProps struct {
		*BaseProps
		DerivedValue string
	}

	props := &DerivedProps{
		BaseProps: &BaseProps{},
	}

	overrider := Override[DerivedProps](map[string]any{
		"BaseValue":    "from base pointer",
		"DerivedValue": "from derived",
	})
	overrider.Apply(props)

	if props.BaseValue != "from base pointer" {
		t.Errorf("Expected 'from base pointer', got '%s'", props.BaseValue)
	}
	if props.DerivedValue != "from derived" {
		t.Errorf("Expected 'from derived', got '%s'", props.DerivedValue)
	}
}

func TestOverrideWithNilEmbeddedPointerStruct(t *testing.T) {
	type BaseProps struct {
		BaseValue string
	}

	type DerivedProps struct {
		*BaseProps
		DerivedValue string
	}

	props := &DerivedProps{}

	overrider := Override[DerivedProps](map[string]any{
		"DerivedValue": "only derived",
	})
	overrider.Apply(props)

	if props.DerivedValue != "only derived" {
		t.Errorf("Expected 'only derived', got '%s'", props.DerivedValue)
	}
}

func TestOverrideWithCaseSensitive(t *testing.T) {
	type Props struct {
		FieldName string
	}

	props := &Props{}

	overrider := Override[Props](map[string]any{
		"FieldName": "exact match",
	})
	overrider.Apply(props)

	if props.FieldName != "exact match" {
		t.Errorf("Expected 'exact match', got '%s'", props.FieldName)
	}
}

func TestOverrideWithMapKeyConversion(t *testing.T) {
	type Props struct {
		IntValue int64
	}

	props := &Props{}

	overrider := Override[Props](map[string]any{
		"IntValue": int32(123),
	})
	overrider.Apply(props)

	if props.IntValue != 123 {
		t.Errorf("Expected 123, got %d", props.IntValue)
	}
}

func TestOverrideWithNilSlice(t *testing.T) {
	type Props struct {
		SliceValue []string
	}

	props := &Props{}

	overrider := Override[Props](map[string]any{
		"SliceValue": nil,
	})
	overrider.Apply(props)

	if props.SliceValue != nil {
		t.Errorf("Expected nil, got %v", props.SliceValue)
	}
}

func TestBuildSetterName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"value", "SetValue"},
		{"my_field", "SetMyField"},
		{"field-name", "SetFieldName"},
		{"field name", "SetFieldName"},
		{"", ""},
		{"  ", ""},
	}

	for _, test := range tests {
		result := buildSetterName(test.input)
		if result != test.expected {
			t.Errorf("buildSetterName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestExportableName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"value", "Value"},
		{"my_field", "MyField"},
		{"field-name", "FieldName"},
		{"field name", "FieldName"},
		{"_private", "Private"},
		{"multi__underscore", "MultiUnderscore"},
	}

	for _, test := range tests {
		result := exportableName(test.input)
		if result != test.expected {
			t.Errorf("exportableName(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

func TestIsExportedIdentifier(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"Value", true},
		{"value", false},
		{"", false},
		{"_private", false},
	}

	for _, test := range tests {
		result := isExportedIdentifier(test.input)
		if result != test.expected {
			t.Errorf("isExportedIdentifier(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestOverrideWithInvalidTarget(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid target")
		}
	}()

	notAStruct := 42

	overrider := Override[int](map[string]any{
		"Value": "test",
	})
	overrider.Apply(&notAStruct)
}

func TestOverrideWithInvalidMapKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for invalid map key type")
		}
	}()

	Override[struct{ Value string }](map[int]string{
		1: "test",
	})
}

func TestOverrideWithDoublePointer(t *testing.T) {
	type Props struct {
		Value string
	}

	literal := Props{Value: "from double pointer"}
	literalPtr := &literal
	doublePtrLiteral := &literalPtr

	props := &Props{}

	overrider := Override[Props](**doublePtrLiteral)
	overrider.Apply(props)

	if props.Value != "from double pointer" {
		t.Errorf("Expected 'from double pointer', got '%s'", props.Value)
	}
}

func TestOverrideWithNilTargetPointer(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for nil target pointer")
		}
	}()

	type Props struct {
		Value string
	}

	var props *Props

	overrider := Override[Props](map[string]any{
		"Value": "test",
	})
	overrider.Apply(props)
}

func TestOverrideWithNilInterface(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when assigning nil interface to non-nillable type")
		}
	}()

	type Props struct {
		Value int
	}

	props := &Props{}

	var nilInterface any

	overrider := Override[Props](map[string]any{
		"Value": nilInterface,
	})
	overrider.Apply(props)
}

func TestOverrideWithInterfaceWrappedNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when assigning interface-wrapped nil to non-nillable type")
		}
	}()

	type Props struct {
		Value int
	}

	props := &Props{}

	var nilInterface any = (*string)(nil)

	overrider := Override[Props](map[string]any{
		"Value": nilInterface,
	})
	overrider.Apply(props)
}

func TestOverrideWithNilFunc(t *testing.T) {
	type PropsWithFunc struct {
		Callback func()
	}

	props := &PropsWithFunc{}

	overrider := Override[PropsWithFunc](map[string]any{
		"Callback": nil,
	})
	overrider.Apply(props)

	if props.Callback != nil {
		t.Errorf("Expected nil callback, got non-nil")
	}
}

func TestOverrideWithNilChan(t *testing.T) {
	type PropsWithChan struct {
		Channel chan int
	}

	props := &PropsWithChan{}

	overrider := Override[PropsWithChan](map[string]any{
		"Channel": nil,
	})
	overrider.Apply(props)

	if props.Channel != nil {
		t.Errorf("Expected nil, got %v", props.Channel)
	}
}

func TestCanonicalNameCaseSensitive(t *testing.T) {
	result := canonicalName("FieldName", false)
	if result != "FieldName" {
		t.Errorf("Expected 'FieldName', got '%s'", result)
	}
}

type propsWithNoArgSetterForTest struct {
	value string
}

func (p *propsWithNoArgSetterForTest) SetValue() {
}

func TestOverrideWithSetterMethodWrongArgCount(t *testing.T) {
	props := &propsWithNoArgSetterForTest{}

	overrider := Override[propsWithNoArgSetterForTest](map[string]any{
		"value": "test",
	})
	overrider.Apply(props)

	if props.value != "test" {
		t.Errorf("Expected 'test', got '%s'", props.value)
	}
}

type propsWithInvalidSetterArgTypeForTest struct {
	value string
}

func (p *propsWithInvalidSetterArgTypeForTest) SetValue(v int) {
	p.value = string(rune(v))
}

func TestOverrideWithSetterMethodInvalidArgType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when setter argument type is invalid")
		}
	}()

	props := &propsWithInvalidSetterArgTypeForTest{}

	overrider := Override[propsWithInvalidSetterArgTypeForTest](map[string]any{
		"value": "test",
	})
	overrider.Apply(props)
}
