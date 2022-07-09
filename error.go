package env

import (
	"fmt"
	"reflect"
)

type ErrInvalidType struct {
	Type reflect.Value
}

func (e *ErrInvalidType) Error() string {
	t := e.Type.Type()
	k := e.Type.Kind()

	if t == nil {
		return "env: invalid type: nil"
	}
	if k != reflect.Pointer {
		return fmt.Sprintf(
			"env: invalid type: non-pointer %s of kind %s",
			t.String(),
			k.String(),
		)
	}

	return fmt.Sprintf("env: invalid type: %s of kind %s", t.String(), k.String())
}
