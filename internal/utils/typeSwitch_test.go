package utils

import "testing"

type TestTypeBaseBytes []byte

func (TestTypeBaseBytes) SomeMethod() string {
	return "SomeMethod"
}

func TestTypeSwitch(t *testing.T) {
	switch ttbb := any(TestTypeBaseBytes{}).(type) {
	case []byte:
		t.Log("is []byte")
	default:
		t.Logf("%T\n", ttbb)
	}
}
