// Package h2h describes the H2H protocol (ISO 8583 based) used between the bank-communication
// service and the Bank 1 simulator: the message format, the Sale (1100) and SaleResponse (1110)
// messages, and the response codes. Both sides import it so they cannot disagree on the format.
package h2h

import (
	"fmt"
	"io"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583/encoding"
	"github.com/moov-io/iso8583/field"
	"github.com/moov-io/iso8583/network"
	"github.com/moov-io/iso8583/padding"
	"github.com/moov-io/iso8583/prefix"
)

// MaxMessageSize is the largest message length, in bytes, that the length header may state.
// The specification limits a message to 1024 bytes.
const MaxMessageSize = 1024

// maxProductDataLength is the limit of DE63 in the specification (llvar ans..520), which is exactly
// ten products with the longest code: 10 * (32 + "#" + 18 + ";").
const maxProductDataLength = 520

// Spec is the H2H message format, shared by every message. It must not be modified.
//
// Fields 0 and 1 are the message type and the bitmap; the specification does not number them, moov does.
// The specification's "lvar" (one byte of length) is prefix.BCD.LL and "llvar" (two bytes) is
// prefix.BCD.LLLL, because moov counts the digits of the length.
var Spec = &iso8583.MessageSpec{
	Name: "H2H",
	Fields: map[int]field.Field{
		0:  field.NewString(&field.Spec{Length: 4, Description: "Message Type Identifier", Enc: encoding.BCD, Pref: prefix.BCD.Fixed}),
		1:  field.NewBitmap(&field.Spec{Length: 8, Description: "Bitmap", Enc: encoding.Binary, Pref: prefix.Binary.Fixed}),
		2:  field.NewString(&field.Spec{Length: 19, Description: "Primary Account Number", Enc: encoding.ASCII, Pref: prefix.BCD.LL}),
		3:  asciiString(6, "Processing Code"),
		4:  field.NewNumeric(&field.Spec{Length: 12, Description: "Transaction Amount", Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed, Pad: padding.Left('0')}),
		11: asciiString(6, "Systems Trace Number"),
		12: asciiString(6, "Transaction Time"),
		13: asciiString(4, "Transaction Date"),
		14: asciiString(4, "Card Expiration Date"),
		37: asciiString(12, "Retrieval Reference Number"),
		38: asciiString(6, "Authorisation Code"),
		39: asciiString(2, "Response Code"),
		41: asciiString(8, "Terminal ID"),
		42: asciiString(15, "Merchant ID"),
		49: field.NewNumeric(&field.Spec{Length: 3, Description: "Currency Code", Enc: encoding.BCD, Pref: prefix.BCD.Fixed, Pad: padding.Left('0')}),
		63: field.NewString(&field.Spec{Length: maxProductDataLength, Description: "Product Data", Enc: encoding.ASCII, Pref: prefix.BCD.LLLL}),
	},
}

func asciiString(length int, description string) field.Field {
	return field.NewString(&field.Spec{Length: length, Description: description, Enc: encoding.ASCII, Pref: prefix.ASCII.Fixed})
}

// ReadLength reads the two byte header that precedes every message and returns the length of the
// message that follows. It returns an error for a length above MaxMessageSize.
func ReadLength(r io.Reader) (int, error) {
	header := network.NewBinary2BytesHeader()
	if _, err := header.ReadFrom(r); err != nil {
		return 0, fmt.Errorf("h2h: read length header: %w", err)
	}
	if header.Length() > MaxMessageSize {
		return 0, fmt.Errorf("h2h: message length %d exceeds %d", header.Length(), MaxMessageSize)
	}

	return header.Length(), nil
}

// WriteLength writes the two byte header for a message of the given length and returns the number
// of header bytes written. It returns an error for a length outside 0..MaxMessageSize.
func WriteLength(w io.Writer, length int) (int, error) {
	if length < 0 || length > MaxMessageSize {
		return 0, fmt.Errorf("h2h: message length %d is outside 0..%d", length, MaxMessageSize)
	}

	header := network.NewBinary2BytesHeader()
	if err := header.SetLength(length); err != nil {
		return 0, fmt.Errorf("h2h: set length header: %w", err)
	}
	n, err := header.WriteTo(w)
	if err != nil {
		return n, fmt.Errorf("h2h: write length header: %w", err)
	}

	return n, nil
}
