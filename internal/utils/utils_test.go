package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMap(t *testing.T) {
	t.Run("struct to struct", func(t *testing.T) {
		type s1 struct {
			field1 string
			field2 int
		}
		type s2 struct {
			field3 string
			field4 int
		}

		input := []s1{
			{"Alice", 30},
			{"Bob", 25},
		}
		fn := func(s s1) s2 { return s2{field3: s.field1, field4: s.field2} }
		expected := []s2{{"Alice", 30}, {"Bob", 25}}

		result := Map(input, fn)
		assert.Equal(t, expected, result)
	})
}
