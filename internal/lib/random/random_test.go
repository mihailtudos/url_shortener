package random

import (
	"strings"
	"testing"
)

func TestGenerateRandomString(t *testing.T) {
	lenght := 10
	got, err := GenerateRandomString(lenght)
	if err != nil {
		t.Error(err)
	}
	
	t.Log("res,", got)
	if len(strings.TrimSpace(got)) != lenght {
		t.Errorf("lenght got = %v, want %v", len(got), lenght)
	}
}