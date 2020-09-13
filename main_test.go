package main

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func Test(t *testing.T) {
	err := prepareDatabase()

	if err != nil {
		t.Fail()
	}
}
