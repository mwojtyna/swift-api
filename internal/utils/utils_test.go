package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSwiftCodeHq(t *testing.T) {
	var tests = []struct {
		code   string
		isHq   bool
		hqCode string
	}{
		{"ABCDEFGHXXX", true, ""},
		{"ABCDEFGHIJKL", false, "ABCDEFGHXXX"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint(tt), func(t *testing.T) {
			isHq, hqCode := IsSwiftCodeHq(tt.code)
			assert.Equal(t, tt.isHq, isHq)
			assert.Equal(t, tt.hqCode, hqCode)
		})
	}
}
