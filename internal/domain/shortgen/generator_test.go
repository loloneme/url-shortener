package shortgen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncodeID(t *testing.T) {
	t.Run("encodes to 10 characters", func(t *testing.T) {
		result, err := EncodeID(1)
		require.NoError(t, err)
		assert.Len(t, result, shortLength)
	})

	t.Run("zero id returns error", func(t *testing.T) {
		_, err := EncodeID(0)
		assert.Error(t, err)
	})

	t.Run("different ids produce different shorts", func(t *testing.T) {
		s1, _ := EncodeID(1)
		s2, _ := EncodeID(2)
		assert.NotEqual(t, s1, s2)
	})

	t.Run("too large id returns error", func(t *testing.T) {
		_, err := EncodeID(^uint64(0))
		assert.Error(t, err)
	})

	t.Run("encoded result passes validation", func(t *testing.T) {
		for i := uint64(1); i <= 100; i++ {
			short, err := EncodeID(i)
			require.NoError(t, err)
			assert.NoError(t, Validate(short))
		}
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
}
