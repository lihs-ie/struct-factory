package factory

import (
	"testing"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusClosed   Status = "closed"
)

func TestEnumFactory(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
		StatusInactive,
		StatusClosed,
	})

	builder := Builder(factory)
	status := builder.BuildWith(0, nil)
	found := false

	for _, candidate := range []Status{StatusPending, StatusActive, StatusInactive, StatusClosed} {
		if status == candidate {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Status %v is not a valid candidate", status)
	}

	status2 := builder.BuildWith(1, nil)
	if status2 == "" {
		t.Error("Expected non-empty status")
	}

	status3 := builder.BuildWith(4, nil)
	if status3 == "" {
		t.Error("Expected non-empty status with seed 4")
	}
}

func TestEnumFactoryWithExclusions(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
		StatusInactive,
		StatusClosed,
	})

	builder := Builder(factory)
	status := builder.BuildWith(0, Override[EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusPending, StatusInactive},
	}))

	if status == StatusPending || status == StatusInactive {
		t.Errorf("Status %v should be excluded", status)
	}

	if status != StatusActive && status != StatusClosed {
		t.Errorf("Status %v is not in remaining candidates", status)
	}

	status2 := builder.BuildWith(1, Override[EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusPending, StatusInactive},
	}))

	if status2 == StatusPending || status2 == StatusInactive {
		t.Errorf("Status %v should be excluded", status2)
	}
}

func TestEnumFactoryWithAllExcluded(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
	})

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic when all candidates are excluded")
		}
	}()

	builder := Builder(factory)
	builder.BuildWith(0, Override[EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusPending, StatusActive},
	}))
}

func TestEnumFactoryInstantiate(t *testing.T) {
	factory := NewEnumFactory([]Status{StatusPending, StatusActive})

	properties := EnumProperties[Status]{value: StatusActive}

	result := factory.Instantiate(properties)
	if result != StatusActive {
		t.Errorf("Expected StatusActive, got %v", result)
	}
}

func TestEnumFactoryRetrieve(t *testing.T) {
	factory := NewEnumFactory([]Status{StatusPending, StatusActive})

	properties := factory.Retrieve(StatusPending)
	if properties.value != StatusPending {
		t.Errorf("Expected Value StatusPending, got %v", properties.value)
	}
	if len(properties.exclusions) != 0 {
		t.Errorf("Expected empty Exclusions, got %v", properties.exclusions)
	}
}

func TestEnumFactoryWithBuilder(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
		StatusInactive,
		StatusClosed,
	})

	builder := Builder(factory)

	status1 := builder.Build(nil)
	if status1 == "" {
		t.Error("Expected non-empty status")
	}

	statuses := builder.BuildList(10, nil)
	if len(statuses) != 10 {
		t.Errorf("Expected 10 statuses, got %d", len(statuses))
	}

	for _, status := range statuses {
		found := false
		for _, candidate := range []Status{StatusPending, StatusActive, StatusInactive, StatusClosed} {
			if status == candidate {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Status %v is not a valid candidate", status)
		}
	}

	status2 := builder.BuildWith(1, nil)
	if status2 == "" {
		t.Error("Expected non-empty status with seed 1")
	}
}

func TestEnumFactoryWithExclusionBuilder(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
		StatusInactive,
		StatusClosed,
	})

	builder := Builder(factory)

	statuses := builder.BuildList(5, Override[EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusInactive},
	}))

	for _, status := range statuses {
		if status == StatusInactive {
			t.Error("StatusInactive should be excluded")
		}
		if status != StatusPending && status != StatusActive && status != StatusClosed {
			t.Errorf("Status %v is not in remaining candidates", status)
		}
	}
}

func TestEnumFactoryWithLiteralValueOverride(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusClosed,
	})

	builder := Builder(factory)

	status := builder.Build(Override[EnumProperties[Status]](map[string]any{
		"Value": StatusClosed,
	}))

	if status != StatusClosed {
		t.Errorf("Expected StatusClosed, got %v", status)
	}
}

func TestEnumFactoryWithLiteralStructOverride(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusActive,
	})

	builder := Builder(factory)

	status := builder.Build(Override[EnumProperties[Status]](map[string]any{
		"Value": StatusActive,
	}))

	if status != StatusActive {
		t.Errorf("Expected StatusActive, got %v", status)
	}
}

func TestEnumFactoryWithLiteralExclusions(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
		StatusActive,
		StatusInactive,
		StatusClosed,
	})

	builder := Builder(factory)

	status := builder.Build(Override[EnumProperties[Status]](map[string]any{
		"Exclusions": []Status{StatusPending, StatusInactive, StatusClosed},
	}))

	if status == StatusPending || status == StatusInactive || status == StatusClosed {
		t.Errorf("Status %v should be excluded", status)
	}

	if status != StatusActive {
		t.Errorf("Expected StatusActive, got %v", status)
	}
}

func TestEnumFactoryWithInlineOverride(t *testing.T) {
	factory := NewEnumFactory([]Status{
		StatusPending,
	})

	builder := Builder(factory)

	status := builder.Build(Override[EnumProperties[Status]](map[string]any{
		"Value":      StatusPending,
		"Exclusions": []Status{},
	}))

	if status != StatusPending {
		t.Errorf("Expected StatusPending, got %v", status)
	}
}
