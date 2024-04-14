package secure_test

import (
	"bytes"
	"testing"

	"github.com/benderr/keypass/internal/secure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEncoder_String(t *testing.T) {
	var buf bytes.Buffer
	masterKey := "123"
	content := "my text"
	enc := secure.NewEncoder(&buf)
	err := enc.Encode(content, masterKey)

	require.NoError(t, err, "failed encode")

	var bufOut bytes.Buffer
	bufOut.Write(buf.Bytes())

	dec := secure.NewDecoder(&bufOut)

	var outStr string
	err2 := dec.Decode(&outStr, masterKey)

	require.NoError(t, err2, "failed decode")

	assert.Equal(t, content, outStr)
}

type testModel struct {
	Login    string
	Password string
	Code     int
}

func TestNewEncoder_Struct(t *testing.T) {
	var buf bytes.Buffer
	masterKey := "123"
	content := testModel{Login: "Login", Password: "Password", Code: 123}
	enc := secure.NewEncoder(&buf)
	err := enc.Encode(content, masterKey)

	require.NoError(t, err, "failed encode")

	var bufOut bytes.Buffer
	bufOut.Write(buf.Bytes())

	dec := secure.NewDecoder(&bufOut)

	outModel := new(testModel)
	err2 := dec.Decode(outModel, masterKey)

	require.NoError(t, err2, "failed decode")

	assert.EqualValues(t, content, *outModel)
}
