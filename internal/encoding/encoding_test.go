package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEncoding(t *testing.T) {
	t.Parallel()
	t.Run("RoundTrip", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo\n@enduml"
		encoded, err := Encode(input)
		require.NoError(t, err)
		require.NotEmpty(t, encoded)
		decoded, err := Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, input, decoded)
	})
	t.Run("RoundTripEmpty", func(t *testing.T) {
		t.Parallel()
		encoded, err := Encode("")
		require.NoError(t, err)
		decoded, err := Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, "", decoded)
	})
	t.Run("RoundTripUnicode", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Föö\n@enduml"
		encoded, err := Encode(input)
		require.NoError(t, err)
		decoded, err := Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, input, decoded)
	})
	t.Run("EncodedIsURLSafe", func(t *testing.T) {
		t.Parallel()
		input := "@startuml\nclass Foo {\n+name : String\n}\n@enduml"
		encoded, err := Encode(input)
		require.NoError(t, err)
		for _, c := range encoded {
			assert.Contains(t, alphabet, string(c))
		}
	})
	t.Run("DecodeInvalidChar", func(t *testing.T) {
		t.Parallel()
		_, err := Decode("!!!!")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})
}
