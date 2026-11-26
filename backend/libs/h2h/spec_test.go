package h2h_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sevenpay/backend/libs/h2h"
)

// The four examples from the H2H specification, each with the field values it carries.
var specExamples = []struct {
	name   string
	hex    string
	fields map[int]string
}{
	{
		name: "sale request",
		hex:  "00A51100703C000008C080021633363530373838343533333136313733303030303030303030303030303031323939303030303237303833303232313132373237303531324A5357443233324C31465454303030303038544D50303030303030303030303031097800613132333435362D37383930233030303030303030303430303030303030333B4146462D33322D343233233030303030303030303039393030303030313B",
		fields: map[int]string{
			0: "1100", 2: "3650788453316173", 3: "000000", 4: "1299", 11: "000027", 12: "083022", 13: "1127", 14: "2705",
			37: "12JSWD232L1F", 41: "TT000008", 42: "TMP000000000001", 49: "978",
			63: "123456-7890#000000000400000003;AFF-32-423#000000000099000001;",
		},
	},
	{
		name: "sale response",
		hex:  "00AD1110703C00000EC080021633363530373838343533333136313733303030303030303030303030303031323939303030303237303833303232313132373237303531324A5357443233324C314641413333313030305454303030303038544D50303030303030303030303031097800613132333435362D37383930233030303030303030303430303030303030333B4146462D33322D343233233030303030303030303039393030303030313B",
		fields: map[int]string{
			0: "1110", 2: "3650788453316173", 3: "000000", 4: "1299", 11: "000027", 12: "083022", 13: "1127", 14: "2705",
			37: "12JSWD232L1F", 38: "AA3310", 39: "00", 41: "TT000008", 42: "TMP000000000001", 49: "978",
			63: "123456-7890#000000000400000003;AFF-32-423#000000000099000001;",
		},
	},
	{
		name: "void sale request",
		hex:  "00661420703C000008C080001633363530373838343533333136313733303230303030303030303030303031323939303030303238303833303232313132373237303531324A5357443233324C31465454303030303038544D503030303030303030303030310978",
		fields: map[int]string{
			0: "1420", 2: "3650788453316173", 3: "020000", 4: "1299", 11: "000028", 12: "083022", 13: "1127", 14: "2705",
			37: "12JSWD232L1F", 41: "TT000008", 42: "TMP000000000001", 49: "978",
		},
	},
	{
		name: "void sale response",
		hex:  "00681430703C00000AC080001633363530373838343533333136313733303230303030303030303030303031323939303030303238303833303232313132373237303531324A5357443233324C314630305454303030303038544D503030303030303030303030310978",
		fields: map[int]string{
			0: "1430", 2: "3650788453316173", 3: "020000", 4: "1299", 11: "000028", 12: "083022", 13: "1127", 14: "2705",
			37: "12JSWD232L1F", 39: "00", 41: "TT000008", 42: "TMP000000000001", 49: "978",
		},
	},
}

// exampleFrame returns the whole example: the length header and the message after it.
func exampleFrame(t *testing.T, exampleHex string) (header, body []byte) {
	t.Helper()

	frame, err := hex.DecodeString(exampleHex)
	require.NoError(t, err)

	return frame[:2], frame[2:]
}

func TestSpec(t *testing.T) {
	t.Parallel()

	for _, tc := range specExamples {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			header, body := exampleFrame(t, tc.hex)

			length, err := h2h.ReadLength(bytes.NewReader(header))
			require.NoError(t, err)
			assert.Equal(t, len(body), length, "length in the header")

			msg := iso8583.NewMessage(h2h.Spec)
			require.NoError(t, msg.Unpack(body))
			for id, want := range tc.fields {
				got, err := msg.GetString(id)
				require.NoError(t, err, "field %d", id)
				assert.Equal(t, want, got, "field %d", id)
			}
			assert.Len(t, msg.GetFields(), len(tc.fields)+1, "fields present, bitmap included")

			packed, err := msg.Pack()
			require.NoError(t, err)
			assert.Equal(t, body, packed)
		})
	}
}

func TestReadLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []byte
		want    int
		wantErr bool
	}{
		{name: "length of the sale request", input: []byte{0x00, 0xA5}, want: 165},
		{name: "largest allowed length", input: []byte{0x04, 0x00}, want: 1024},
		{name: "above the largest length", input: []byte{0x04, 0x01}, wantErr: true},
		{name: "header cut short", input: []byte{0x00}, wantErr: true},
		{name: "no header", input: nil, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := h2h.ReadLength(bytes.NewReader(tc.input))
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWriteLength(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		length  int
		want    []byte
		wantErr bool
	}{
		{name: "length of the sale request", length: 165, want: []byte{0x00, 0xA5}},
		{name: "largest allowed length", length: 1024, want: []byte{0x04, 0x00}},
		{name: "above the largest length", length: 1025, wantErr: true},
		{name: "negative length", length: -1, wantErr: true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			n, err := h2h.WriteLength(&buf, tc.length)
			if tc.wantErr {
				require.Error(t, err)
				assert.Empty(t, buf.Bytes())
				return
			}
			require.NoError(t, err)
			assert.Equal(t, 2, n)
			assert.Equal(t, tc.want, buf.Bytes())
		})
	}
}
