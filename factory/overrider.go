package factory

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
	"unsafe"
)

// Overrider stores a prepared Partial that mutates properties of type P.
type Overrider[P any] struct {
	fn Partial[P]
}

// Apply runs the stored override against the provided properties pointer.
func (o Overrider[P]) Apply(properties *P) {
	if o.fn != nil {
		o.fn(properties)
	}
}

// Func returns the partial function backing this Overrider.
func (o Overrider[P]) Func() Partial[P] {
	return o.fn
}

type overrideOptions struct {
	caseInsensitive bool
	allowUnexported bool
}

// OverrideOption configures how Override applies entries to targets.
type OverrideOption func(*overrideOptions)

// WithCaseInsensitive forces case-insensitive field matching (enabled by default).
func WithCaseInsensitive() OverrideOption {
	return func(opts *overrideOptions) {
		opts.caseInsensitive = true
	}
}

// DisallowUnexported prevents Override from mutating unexported struct fields.
func DisallowUnexported() OverrideOption {
	return func(opts *overrideOptions) {
		opts.allowUnexported = false
	}
}

// Override normalizes a literal (map or struct) into an Overrider for properties P.
func Override[P any](literal any, opts ...OverrideOption) Overrider[P] {
	config := overrideOptions{caseInsensitive: true, allowUnexported: true}
	for _, opt := range opts {
		opt(&config)
	}

	entries, err := parseOverrideLiteral(literal, config.caseInsensitive)
	if err != nil {
		panic(err)
	}

	return Overrider[P]{
		fn: func(properties *P) {
			if err := applyOverrideEntries(properties, entries, config); err != nil {
				panic(err)
			}
		},
	}
}

type literalEntry struct {
	originalName string
	key          string
	value        reflect.Value
}

func parseOverrideLiteral(literal any, caseInsensitive bool) ([]literalEntry, error) {
	if literal == nil {
		return nil, errors.New("override: literal cannot be nil")
	}

	value := reflect.ValueOf(literal)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return nil, errors.New("override: literal pointer cannot be nil")
		}
		value = value.Elem()
	}

	switch value.Kind() {
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return nil, fmt.Errorf("override: map key must be string, got %s", value.Type().Key())
		}

		entries := make([]literalEntry, 0, value.Len())
		iter := value.MapRange()
		for iter.Next() {
			key := iter.Key().String()
			entries = append(entries, literalEntry{
				originalName: key,
				key:          canonicalName(key, caseInsensitive),
				value:        iter.Value(),
			})
		}
		return entries, nil

	case reflect.Struct:
		entries := make([]literalEntry, 0, value.NumField())
		if !value.CanAddr() {
			addr := reflect.New(value.Type())
			addr.Elem().Set(value)
			value = addr.Elem()
		}
		typ := value.Type()
		for i := 0; i < value.NumField(); i++ {
			field := typ.Field(i)
			fieldValue := value.Field(i)
			if !fieldValue.CanInterface() {
				fieldValue = reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr())).Elem()
			}
			entries = append(entries, literalEntry{
				originalName: field.Name,
				key:          canonicalName(field.Name, caseInsensitive),
				value:        fieldValue,
			})
		}
		return entries, nil
	default:
		return nil, fmt.Errorf("override: unsupported literal type %s", value.Kind())
	}
}

func canonicalName(name string, caseInsensitive bool) string {
	if !caseInsensitive {
		return name
	}
	return strings.ToLower(name)
}

func applyOverrideEntries[P any](properties *P, entries []literalEntry, config overrideOptions) error {
	target := reflect.ValueOf(properties)
	if target.Kind() != reflect.Pointer || target.IsNil() {
		return fmt.Errorf("override: target must be a non-nil pointer, got %T", properties)
	}

	elem := target.Elem()
	if elem.Kind() != reflect.Struct {
		return fmt.Errorf("override: target must point to a struct, got %s", elem.Kind())
	}

	for _, entry := range entries {
		if err := applyOverrideEntry(target, elem, entry, config); err != nil {
			return err
		}
	}

	return nil
}

func applyOverrideEntry(targetPtr, targetValue reflect.Value, entry literalEntry, config overrideOptions) error {
	setterName := buildSetterName(entry.originalName)
	if setterName != "" {
		method := targetPtr.MethodByName(setterName)
		if method.IsValid() && method.Type().NumIn() == 1 {
			arg, err := prepareOverrideValue(entry.value, method.Type().In(0))
			if err != nil {
				return fmt.Errorf("override: cannot assign %q via setter: %w", entry.originalName, err)
			}
			method.Call([]reflect.Value{arg})
			notifyOverride(targetPtr, entry.originalName)
			return nil
		}
	}

	fieldValue, fieldInfo, ok := lookupField(targetValue, entry.key, config.caseInsensitive)
	if !ok {
		return fmt.Errorf("override: unknown field %q on %s", entry.originalName, targetValue.Type())
	}

	prepared, err := prepareOverrideValue(entry.value, fieldValue.Type())
	if err != nil {
		return fmt.Errorf("override: cannot assign %q: %w", fieldInfo.Name, err)
	}

	isExported := isExportedStructField(&fieldInfo)
	if isExported && fieldValue.CanSet() {
		fieldValue.Set(prepared)
		notifyOverride(targetPtr, fieldInfo.Name)
		return nil
	}

	if !isExported {
		if !config.allowUnexported {
			return fmt.Errorf("override: field %q is unexported and DisallowUnexported was provided", fieldInfo.Name)
		}
		if fieldValue.CanAddr() {
			ptr := reflect.NewAt(fieldValue.Type(), unsafe.Pointer(fieldValue.UnsafeAddr()))
			ptr.Elem().Set(prepared)
			notifyOverride(targetPtr, fieldInfo.Name)
			return nil
		}
	}

	return fmt.Errorf("override: field %q cannot be set", fieldInfo.Name)
}

func lookupField(targetValue reflect.Value, key string, caseInsensitive bool) (reflect.Value, reflect.StructField, bool) {
	canonical := canonicalName(key, caseInsensitive)
	return lookupFieldRecursive(targetValue, canonical, caseInsensitive)
}

func lookupFieldRecursive(value reflect.Value, canonical string, caseInsensitive bool) (reflect.Value, reflect.StructField, bool) {
	structType := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)
		if canonicalName(field.Name, caseInsensitive) == canonical {
			return fieldValue, field, true
		}
		//nolint:nestif // Complexity 6 is acceptable for embedded struct field lookup
		if field.Anonymous {
			embedded := fieldValue
			if embedded.Kind() == reflect.Pointer {
				if embedded.IsNil() {
					continue
				}
				embedded = embedded.Elem()
			}
			if embedded.Kind() == reflect.Struct {
				if nestedValue, nestedField, ok := lookupFieldRecursive(embedded, canonical, caseInsensitive); ok {
					return nestedValue, nestedField, true
				}
			}
		}
	}
	return reflect.Value{}, reflect.StructField{}, false
}

func prepareOverrideValue(value reflect.Value, targetType reflect.Type) (reflect.Value, error) {
	if !value.IsValid() {
		if canBeNil(targetType) {
			return reflect.Zero(targetType), nil
		}
		return reflect.Value{}, fmt.Errorf("cannot assign nil to %s", targetType)
	}

	for value.Kind() == reflect.Interface {
		if value.IsNil() {
			if canBeNil(targetType) {
				return reflect.Zero(targetType), nil
			}
			return reflect.Value{}, fmt.Errorf("cannot assign nil to %s", targetType)
		}
		value = value.Elem()
	}

	if value.Type().AssignableTo(targetType) {
		return value, nil
	}

	if value.Type().ConvertibleTo(targetType) {
		return value.Convert(targetType), nil
	}

	return reflect.Value{}, fmt.Errorf("cannot convert %s to %s", value.Type(), targetType)
}

func canBeNil(targetType reflect.Type) bool {
	switch targetType.Kind() {
	case reflect.Interface, reflect.Pointer, reflect.Map, reflect.Slice, reflect.Func, reflect.Chan:
		return true
	default:
		return false
	}
}

type overrideTracker interface {
	noteOverride(field string)
}

func notifyOverride(targetPtr reflect.Value, field string) {
	if tracker, ok := targetPtr.Interface().(overrideTracker); ok {
		tracker.noteOverride(field)
	}
}

func buildSetterName(name string) string {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ""
	}
	return "Set" + exportableName(trimmed)
}

func exportableName(name string) string {
	var builder strings.Builder
	upperNext := true
	for _, r := range name {
		if r == '_' || r == '-' || r == ' ' {
			upperNext = true
			continue
		}
		if upperNext {
			builder.WriteRune(unicode.ToUpper(r))
			upperNext = false
			continue
		}
		builder.WriteRune(r)
	}
	return builder.String()
}

func isExportedStructField(field *reflect.StructField) bool {
	return field.PkgPath == "" && isExportedIdentifier(field.Name)
}

func isExportedIdentifier(name string) bool {
	if name == "" {
		return false
	}
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}
