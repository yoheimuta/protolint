package lint

import (
	"fmt"

	"github.com/yoheimuta/protolint/linter/autodisable"
)

type autoDisableFlag struct {
	raw             string
	autoDisableType autodisable.PlacementType
}

func (f *autoDisableFlag) String() string {
	return fmt.Sprint(f.raw)
}

func (f *autoDisableFlag) Set(value string) error {
	if f.autoDisableType != 0 {
		return fmt.Errorf("auto_disable is already set")
	}

	r, err := GetAutoDisableType(value)
	if err != nil {
		return err
	}
	f.raw = value
	f.autoDisableType = r
	return nil
}

// GetAutoDisableType returns a type from the specified key.
func GetAutoDisableType(value string) (autodisable.PlacementType, error) {
	rs := map[string]autodisable.PlacementType{
		"next": autodisable.Next,
		"this": autodisable.ThisThenNext,
	}
	if r, ok := rs[value]; ok {
		return r, nil
	}
	return autodisable.Noop, fmt.Errorf(`available auto_disable are "next" and "this"`)
}
