package goroonga

import "testing"

func TestInitAndTerminate(t *testing.T) {
	err := Init()
	if err != nil {
		t.Errorf("failed to initialize with error: %s", err)
	}
	defer func() {
		err := Terminate()
		if err != nil {
			t.Errorf("failed to initialize with error: %s", err)
		}
	}()
}
