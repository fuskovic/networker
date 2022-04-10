package encoder

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncoder(t *testing.T) {
	type testObject struct {
		Field1 string `json:"field_1" yaml:"field_1" table:"FIELD_1"`
		Field2 int    `json:"field_2" yaml:"field_2" table:"FIELD_2"`
	}

	for _, test := range []struct {
		name           string
		output         string
		object         testObject
		assertExpected func(b []byte)
	}{
		{
			name:   "encode as json",
			output: "json",
			object: testObject{
				Field1: "a",
				Field2: 1,
			},
			assertExpected: func(b []byte) {
				o := new(testObject)
				require.NoError(t, json.Unmarshal(b, o))
				require.Equal(t, "a", o.Field1)
				require.Equal(t, 1, o.Field2)
			},
		},
		{
			name:   "encode as yaml",
			output: "yaml",
			object: testObject{
				Field1: "b",
				Field2: 2,
			},
			assertExpected: func(b []byte) {
				expected := "- field_1: b\n  field_2: 2\n"
				require.Equal(t, expected, string(b))
			},
		},
		{
			name: "encode as table",
			object: testObject{
				Field1: "c",
				Field2: 3,
			},
			assertExpected: func(b []byte) {
				expected := "FIELD_1    FIELD_2    \nc          3          \n"
				require.Equal(t, expected, string(b))
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			enc := New[testObject](buf, test.output)
			require.NoError(t, enc.Encode(test.object))
			test.assertExpected(buf.Bytes())
		})
	}
}
