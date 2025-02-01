package pgxtxpool

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	expErr := ErrTxPool
	arrErr := []error{
		ErrTxPoolIDNotFound,
		ErrTxPoolNotFound,
		ErrTxPoolTrxStillExistsInPool,
	}
	for _, err := range arrErr {
		testName := fmt.Sprintf("TEST: %s", err.Error())
		t.Run(testName, func(t *testing.T) {

			// make sure that child errors can be match with parent errors
			if !errors.Is(err, expErr) {
				t.FailNow()
			}

			// make sure that child errors can be unwrapped become parent errors
			if errx := errors.Unwrap(err); !errors.Is(errx, expErr) {
				t.FailNow()
			}
		})
	}

}
