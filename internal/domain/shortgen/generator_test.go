package shortgen

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator(t *testing.T) {
	testRand := rand.New(rand.NewSource(42))
	gen := NewGenerator(testRand)

	t.Run("length is 10", func(t *testing.T) {
		result := gen.Generate()

		assert.Len(t, result, 10)
	})

	t.Run("result contains only charset chars", func(t *testing.T) {
		result := gen.Generate()

		for _, ch := range result {
			require.True(t, strings.ContainsRune(charset, ch))
		}
	})

	t.Run("same rand seed generates same result", func(t *testing.T) {
		rand1 := rand.New(rand.NewSource(42))
		rand2 := rand.New(rand.NewSource(42))
		gen1 := NewGenerator(rand1)
		gen2 := NewGenerator(rand2)

		result1 := gen1.Generate()
		result2 := gen2.Generate()

		assert.Equal(t, result1, result2)
	})
}

func TestValidate(t *testing.T) {
	t.Run("valid short codes", func(t *testing.T) {
		validCodes := []string{
			"abcdefghij",
			"ABCDEFGHIJ",
			"0123456789",
			"abc_DEF_12",
			"__________",
		}
		for _, code := range validCodes {
			assert.NoError(t, Validate(code), "expected %q to be valid", code)
		}
	})

	t.Run("too short", func(t *testing.T) {
		assert.Error(t, Validate("abc"))
	})

	t.Run("too long", func(t *testing.T) {
		assert.Error(t, Validate("abcdefghijk"))
	})

	t.Run("empty string", func(t *testing.T) {
		assert.Error(t, Validate(""))
	})

	t.Run("invalid characters", func(t *testing.T) {
		invalidCodes := []string{
			"abcde-ghij",
			"abcde ghij",
			"abcde!ghij",
			"abcde/ghij",
		}
		for _, code := range invalidCodes {
			assert.Error(t, Validate(code), "expected %q to be invalid", code)
		}
	})

	t.Run("generated code passes validation", func(t *testing.T) {
		gen := NewGenerator(rand.New(rand.NewSource(99)))
		for i := 0; i < 100; i++ {
			assert.NoError(t, Validate(gen.Generate()))
		}
	})
}
