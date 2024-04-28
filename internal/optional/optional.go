package optional

import (
	"encoding/json"
	"fmt"
)

type Optional[T any] struct {
	value T
	isSet bool
}

func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value: value,
		isSet: true,
	}
}

func None[T any]() Optional[T] {
	var empty T

	return Optional[T]{
		value: empty,
		isSet: false,
	}
}

func (o *Optional[T]) Get() (T, bool) {
	if o.isSet {
		return o.value, true
	}

	var empty T
	return empty, false
}

func (o *Optional[T]) MustGet() T {
	if o.isSet {
		return o.value
	}

	panic("value is not set")
}

func (o *Optional[T]) GetOrElse(defaultValue T) T {
	if o.isSet {
		return o.value
	}

	return defaultValue
}

func (o *Optional[T]) IsSet() bool {
	return o.isSet
}

func (o *Optional[T]) String() string {
	if o.isSet {
		return fmt.Sprintf("Some[%T](%v)", *new(T), o.value)
	}

	return fmt.Sprintf("None[%T]", *new(T))
}

func (o *Optional[T]) MarshalJSON() ([]byte, error) {
	if o.isSet {
		return json.Marshal(o.value)
	}

	return []byte("null"), nil
}

func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		o.isSet = false
		return nil
	}

	var tmp T
	if err := json.Unmarshal(data, &tmp); err != nil {
		return fmt.Errorf("failed to unmarshal to underlying value: %w", err)
	}

	o.value = tmp
	o.isSet = true

	return nil
}

func (o Optional[T]) PointerValue() *Optional[T] {
	return &o
}
