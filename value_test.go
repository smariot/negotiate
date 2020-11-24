package negotiate

import (
	"errors"
	"testing"
)

func TestMust(t *testing.T) {
	err := errors.New("fake error")

	defer func() {
		if err != recover() {
			t.Errorf("didn't panic with expected error")
		}
	}()

	Must(nil, err)
}
