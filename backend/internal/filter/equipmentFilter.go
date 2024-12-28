package filter

import (
	"fmt"
)

const (
	DataTypeStr   = "string"
	DataTypeInt   = "int"
	DataTypeFloat = "float"

	OperatorEq               = "="
	OperatorLowerThan        = "<"
	OperatorLowerThanEqual   = "<="
	OperatorGreaterThan      = ">"
	OperatorGreaterThanEqual = ">="
	OperatorBetween          = "between" //TODO: think about operators
)

type optionsMap struct {
	options map[string][]Field
}

func NewOptionsMap() *optionsMap {
	return &optionsMap{
		options: make(map[string][]Field),
	}
}

type Field struct {
	Name     string
	Value    string
	Operator string
	Type     string
}

type OptionsMap interface {
	AddField(name, operator, value, dtype string) error
	Fields() map[string][]Field
}

type Options interface {
	Limit() int
	IsToApply() bool
	AddField(name, operator, value, dtype string) error
	Fields() []Field
}

func (o *optionsMap) AddField(name, operator, value, dtype string) error {
	err := validateOperator(operator)
	if err != nil {
		return err
	}
	o.options[name] = append(o.options[name], Field{
		Name:     name,
		Operator: operator,
		Value:    value,
		Type:     dtype,
	})
	return nil
}

func (o *optionsMap) Fields() map[string][]Field {
	return o.options
}

func validateOperator(operator string) error {
	switch operator {
	case OperatorEq:
	case OperatorLowerThan:
	case OperatorLowerThanEqual:
	case OperatorGreaterThan:
	case OperatorGreaterThanEqual:
	default:
		return fmt.Errorf("invalid operator: %s", operator)
	}
	return nil
}
