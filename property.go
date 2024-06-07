package parseme

import "strings"

type PropertyType int

const (
	Value PropertyType = iota
	Boolean
)

type Property struct {
	propertyType PropertyType
	name         string
	value        string
}

func (p *Property) IsBoolean() bool {
	return p.propertyType == Boolean
}

func (p *Property) Name() string {
	return p.name
}

func (p *Property) Value() string {
	return p.value
}

func (p *Property) BooleanValue() (bool, error) {
	if !p.IsBoolean() {
		return false, &PropertyBooleanValueError{}
	}

	if p.value == "true" {
		return true, nil
	} else {
		return false, nil
	}
}

func (p *Property) SetType(propertyType PropertyType) {
	if p.propertyType == Value && propertyType == Boolean {
		if p.value != "true" && p.value != "false" {
			p.value = "true"
		}
	}
	p.propertyType = propertyType
}

func (p *Property) SetName(name string) error {
	trimmed := strings.Trim(name, " ")
	if len(trimmed) == 0 {
		return &PropertyEmptyNameError{}
	}

	if !isValidProperty(trimmed) {
		return &PropertyInvalidNameError{name: trimmed}
	}

	p.name = trimmed
	return nil
}

func (p *Property) SetValue(value string) error {
	if p.propertyType == Boolean {
		if value != "true" && value != "false" {
			return &PropertyInvalidBooleanError{value: value}
		}
	}

	p.value = removeQuotes(value)
	return nil
}

func NewProperty(propertyType PropertyType, name string, value string) *Property {
	property := &Property{}
	property.SetType(propertyType)
	property.SetName(name)
	property.SetValue(value)
	return property
}
