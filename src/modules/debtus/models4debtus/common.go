package models4debtus

import (
	"errors"
	"fmt"
)

func validateString(errMess, s string, validValues []string) error {
	var ok bool
	for _, validValue := range validValues {
		if s == validValue {
			ok = true
		}
	}
	if !ok {
		return fmt.Errorf("%v: '%v'", errMess, s)
	}
	return nil
}

var ErrNoProperties = errors.New("no properties")

//var checkHasProperties = func(kind string, properties []datastore.Property) {
//	if len(properties) == 0 {
//		panic(errors.WithMessage(ErrNoProperties, fmt.Sprintf("kind="+kind)).Error())
//	}
//}
