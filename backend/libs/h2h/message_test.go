package h2h_test

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"testing"

	"github.com/moov-io/iso8583"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sevenpay/backend/libs/h2h"
)

const (
	saleRequestExample  = 0
	saleResponseExample = 1
)

// exampleSale is the sale from the specification's example.
func exampleSale() h2h.Sale {
	return h2h.Sale{
		PAN:      "3650788453316173",
		Amount:   1299,
		STAN:     "000027",
		Time:     "083022",
		Date:     "1127",
		Expiry:   "2705",
		RRN:      "12JSWD232L1F",
		TID:      "TT000008",
		MID:      "TMP000000000001",
		Currency: h2h.CurrencyEUR,
		Products: []h2h.Product{
			{Code: "123456-7890", UnitPrice: 400, Quantity: 3},
			{Code: "AFF-32-423", UnitPrice: 99, Quantity: 1},
		},
	}
}

func exampleSaleResponse() h2h.SaleResponse {
	return h2h.SaleResponse{Sale: exampleSale(), AuthCode: "AA3310", ResponseCode: h2h.CodeApproved}
}

// exampleMessage returns the unpacked message of one of the specification's examples.
func exampleMessage(t *testing.T, example int) *iso8583.Message {
	t.Helper()

	_, body := exampleFrame(t, specExamples[example].hex)
	msg := iso8583.NewMessage(h2h.Spec)
	require.NoError(t, msg.Unpack(body))

	return msg
}

func TestSale_Message(t *testing.T) {
	t.Parallel()

	msg, err := exampleSale().Message()
	require.NoError(t, err)
	packed, err := msg.Pack()
	require.NoError(t, err)

	_, want := exampleFrame(t, specExamples[saleRequestExample].hex)
	assert.Equal(t, want, packed)
}

func TestSale_Message_invalid(t *testing.T) {
	t.Parallel()

	manyProducts := func(n int) []h2h.Product {
		products := make([]h2h.Product, n)
		for i := range products {
			products[i] = h2h.Product{Code: "A", UnitPrice: 1, Quantity: 1}
		}
		return products
	}
	withProduct := func(p h2h.Product) func(*h2h.Sale) {
		return func(s *h2h.Sale) { s.Products = []h2h.Product{p} }
	}

	tests := []struct {
		name   string
		mutate func(*h2h.Sale)
	}{
		{name: "empty pan", mutate: func(s *h2h.Sale) { s.PAN = "" }},
		{name: "pan of 20 digits", mutate: func(s *h2h.Sale) { s.PAN = strings.Repeat("1", 20) }},
		{name: "pan with a letter", mutate: func(s *h2h.Sale) { s.PAN = "36507884533161A3" }},
		{name: "expiry month 13", mutate: func(s *h2h.Sale) { s.Expiry = "2713" }},
		{name: "expiry month 00", mutate: func(s *h2h.Sale) { s.Expiry = "2700" }},
		{name: "expiry of 3 digits", mutate: func(s *h2h.Sale) { s.Expiry = "275" }},
		{name: "zero amount", mutate: func(s *h2h.Sale) { s.Amount = 0 }},
		{name: "amount of 13 digits", mutate: func(s *h2h.Sale) { s.Amount = 1_000_000_000_000 }},
		{name: "stan of 5 digits", mutate: func(s *h2h.Sale) { s.STAN = "00027" }},
		{name: "time of hour 24", mutate: func(s *h2h.Sale) { s.Time = "240000" }},
		{name: "time of 5 digits", mutate: func(s *h2h.Sale) { s.Time = "08302" }},
		{name: "date of 30 february", mutate: func(s *h2h.Sale) { s.Date = "0230" }},
		{name: "rrn of 11 characters", mutate: func(s *h2h.Sale) { s.RRN = "12JSWD232L1" }},
		{name: "rrn with a hash", mutate: func(s *h2h.Sale) { s.RRN = "12JSWD232L1#" }},
		{name: "tid of 7 characters", mutate: func(s *h2h.Sale) { s.TID = "TT00000" }},
		{name: "mid of 14 characters", mutate: func(s *h2h.Sale) { s.MID = "TMP00000000000" }},
		{name: "currency 0", mutate: func(s *h2h.Sale) { s.Currency = 0 }},
		{name: "currency 1000", mutate: func(s *h2h.Sale) { s.Currency = 1000 }},
		{name: "no products", mutate: func(s *h2h.Sale) { s.Products = nil }},
		{name: "eleven products", mutate: func(s *h2h.Sale) { s.Products = manyProducts(11) }},
		{name: "empty product code", mutate: withProduct(h2h.Product{Code: "", UnitPrice: 1, Quantity: 1})},
		{name: "product code of 33 bytes", mutate: withProduct(h2h.Product{Code: strings.Repeat("A", 33), UnitPrice: 1, Quantity: 1})},
		{name: "product code with a hash", mutate: withProduct(h2h.Product{Code: "A#B", UnitPrice: 1, Quantity: 1})},
		{name: "product code with a semicolon", mutate: withProduct(h2h.Product{Code: "A;B", UnitPrice: 1, Quantity: 1})},
		{name: "product code with a non ASCII letter", mutate: withProduct(h2h.Product{Code: "čaj", UnitPrice: 1, Quantity: 1})},
		{name: "negative unit price", mutate: withProduct(h2h.Product{Code: "A", UnitPrice: -1, Quantity: 1})},
		{name: "unit price of 13 digits", mutate: withProduct(h2h.Product{Code: "A", UnitPrice: 1_000_000_000_000, Quantity: 1})},
		{name: "negative quantity", mutate: withProduct(h2h.Product{Code: "A", UnitPrice: 1, Quantity: -1})},
		{name: "quantity of 7 digits", mutate: withProduct(h2h.Product{Code: "A", UnitPrice: 1, Quantity: 1_000_000})},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sale := exampleSale()
			tc.mutate(&sale)

			msg, err := sale.Message()
			require.Error(t, err)
			assert.Nil(t, msg)
			// The error names the field but never repeats the card number or the expiry date.
			for _, secret := range []string{sale.PAN, sale.Expiry} {
				if secret != "" {
					assert.NotContains(t, err.Error(), secret)
				}
			}
		})
	}
}

func TestSale_Message_longestMessage(t *testing.T) {
	t.Parallel()

	sale := exampleSale()
	sale.PAN = strings.Repeat("9", 19)
	sale.Amount = 999_999_999_999
	sale.Products = make([]h2h.Product, 10)
	for i := range sale.Products {
		sale.Products[i] = h2h.Product{Code: strings.Repeat("X", 32), UnitPrice: 999_999_999_999, Quantity: 999_999}
	}

	msg, err := sale.Message()
	require.NoError(t, err)
	packed, err := msg.Pack()
	require.NoError(t, err)

	unpacked := iso8583.NewMessage(h2h.Spec)
	require.NoError(t, unpacked.Unpack(packed))
	got, err := h2h.ParseSale(unpacked)
	require.NoError(t, err)
	assert.Equal(t, sale, got)
}

func TestParseSale(t *testing.T) {
	t.Parallel()

	got, err := h2h.ParseSale(exampleMessage(t, saleRequestExample))
	require.NoError(t, err)
	assert.Equal(t, exampleSale(), got)
}

func TestParseSale_invalid(t *testing.T) {
	t.Parallel()

	validEntry := "A#000000000100000001;"
	tests := []struct {
		name   string
		mutate func(t *testing.T, msg *iso8583.Message)
	}{
		{name: "missing pan", mutate: func(_ *testing.T, msg *iso8583.Message) { msg.UnsetField(2) }},
		{name: "missing product data", mutate: func(_ *testing.T, msg *iso8583.Message) { msg.UnsetField(63) }},
		{name: "processing code of a void", mutate: func(t *testing.T, msg *iso8583.Message) { require.NoError(t, msg.Field(3, "020000")) }},
		{name: "message type of a response", mutate: func(_ *testing.T, msg *iso8583.Message) { msg.MTI("1110") }},
		{name: "products without a hash", mutate: setProducts("ABC;")},
		{name: "products with a short price", mutate: setProducts("A#12345;")},
		{name: "products with text after the last semicolon", mutate: setProducts(validEntry + "X")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := exampleMessage(t, saleRequestExample)
			tc.mutate(t, msg)

			_, err := h2h.ParseSale(msg)
			require.Error(t, err)
		})
	}
}

func TestParseSale_nil(t *testing.T) {
	t.Parallel()

	_, err := h2h.ParseSale(nil)
	require.Error(t, err)
}

func setProducts(data string) func(t *testing.T, msg *iso8583.Message) {
	return func(t *testing.T, msg *iso8583.Message) {
		t.Helper()
		require.NoError(t, msg.Field(63, data))
	}
}

func TestSaleResponse_Message(t *testing.T) {
	t.Parallel()

	msg, err := exampleSaleResponse().Message()
	require.NoError(t, err)
	packed, err := msg.Pack()
	require.NoError(t, err)

	_, want := exampleFrame(t, specExamples[saleResponseExample].hex)
	assert.Equal(t, want, packed)
}

func TestSaleResponse_Message_declinedWithoutAuthCode(t *testing.T) {
	t.Parallel()

	response := exampleSaleResponse()
	response.ResponseCode = h2h.CodeExpiredCard
	response.AuthCode = ""

	msg, err := response.Message()
	require.NoError(t, err)
	packed, err := msg.Pack()
	require.NoError(t, err)

	// The bitmap is only filled in when the message is packed, so look at what was sent.
	sent := iso8583.NewMessage(h2h.Spec)
	require.NoError(t, sent.Unpack(packed))
	assert.False(t, sent.Bitmap().IsSet(38), "DE38 must not be sent for a declined sale")
	assert.True(t, sent.Bitmap().IsSet(39))

	got, err := h2h.ParseSaleResponse(sent)
	require.NoError(t, err)
	assert.Equal(t, response, got)
}

func TestSaleResponse_Message_invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(*h2h.SaleResponse)
	}{
		{name: "approved without authorisation code", mutate: func(r *h2h.SaleResponse) { r.AuthCode = "" }},
		{name: "declined with a malformed authorisation code", mutate: func(r *h2h.SaleResponse) {
			r.ResponseCode = h2h.CodeDoNotHonor
			r.AuthCode = "AA3"
		}},
		{name: "empty response code", mutate: func(r *h2h.SaleResponse) { r.ResponseCode = "" }},
		{name: "invalid sale part", mutate: func(r *h2h.SaleResponse) { r.PAN = "" }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			response := exampleSaleResponse()
			tc.mutate(&response)

			msg, err := response.Message()
			require.Error(t, err)
			assert.Nil(t, msg)
		})
	}
}

func TestParseSaleResponse(t *testing.T) {
	t.Parallel()

	got, err := h2h.ParseSaleResponse(exampleMessage(t, saleResponseExample))
	require.NoError(t, err)
	assert.Equal(t, exampleSaleResponse(), got)
}

func TestParseSaleResponse_keepsUnknownResponseCode(t *testing.T) {
	t.Parallel()

	msg := exampleMessage(t, saleResponseExample)
	require.NoError(t, msg.Field(39, "91"))
	msg.UnsetField(38) // a bank sends DE38 only when it approves

	got, err := h2h.ParseSaleResponse(msg)
	require.NoError(t, err)
	assert.Equal(t, "91", got.ResponseCode)
}

func TestParseSaleResponse_invalid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(msg *iso8583.Message)
	}{
		{name: "approved without authorisation code", mutate: func(msg *iso8583.Message) { msg.UnsetField(38) }},
		{name: "missing response code", mutate: func(msg *iso8583.Message) { msg.UnsetField(39) }},
		{name: "missing rrn", mutate: func(msg *iso8583.Message) { msg.UnsetField(37) }},
		{name: "message type of a request", mutate: func(msg *iso8583.Message) { msg.MTI("1100") }},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := exampleMessage(t, saleResponseExample)
			tc.mutate(msg)

			_, err := h2h.ParseSaleResponse(msg)
			require.Error(t, err)
		})
	}
}

func TestSale_hidesCardNumber(t *testing.T) {
	t.Parallel()

	sale := exampleSale()
	response := exampleSaleResponse()

	var logged bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logged, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if len(groups) == 0 && a.Key == slog.TimeKey {
				return slog.Attr{}
			}
			return a
		},
	}))
	logger.Info("authorizing", "sale", sale)
	logger.Info("authorized", "response", response)

	outputs := map[string]string{
		"slog":        logged.String(),
		"sale %v":     fmt.Sprintf("%v", sale),
		"response %v": fmt.Sprintf("%v", response),
	}
	for name, output := range outputs {
		assert.NotContains(t, output, sale.PAN, name)
		assert.NotContains(t, output, sale.Expiry, name)
	}
	assert.Contains(t, logged.String(), "365078******6173")
	assert.Contains(t, logged.String(), "response_code=00")
}
