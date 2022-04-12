package lib

import (
	"fmt"
	"reflect"
)

type MustCheck interface {
	IsValid() bool
}

func Check(items []MustCheck) error {
	for _, item := range items {
		if !item.IsValid() {
			return fmt.Errorf("[type: %s]: invalid value: %s", reflect.TypeOf(item), item)
		}
	}
	return nil
}
