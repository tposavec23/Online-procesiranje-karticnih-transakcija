package h2h

import (
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/moov-io/iso8583"
)

// Response codes (DE39) listed in the specification. Not every service will use all of them, but
// the specification defines them all.
const (
	CodeApproved              = "00"
	CodeDoNotHonor            = "05"
	CodeFormatError           = "30"
	CodeInsufficientFunds     = "51"
	CodeExpiredCard           = "54"
	CodeExceedsAmountLimit    = "61"
	CodeExceedsFrequencyLimit = "65"
)

// CurrencyEUR is the ISO 4217 numeric code for the euro (DE49).
const CurrencyEUR = 978

const (
	mtiSale            = "1100"
	mtiSaleResponse    = "1110"
	processingCodeSale = "000000"

	maxAmount       = 999_999_999_999 // DE4 and Product.UnitPrice have 12 digits
	maxQuantity     = 999_999         // Product.Quantity has 6 digits
	maxProducts     = 10
	maxCodeLength   = 32
	unitPriceDigits = 12 // DE63 unit price, formatted like DE4
	quantityDigits  = 6  // DE63 number of items
)

// Product is one entry of DE63 "Product Data": a code, a unit price and a quantity. The
// specification calls it a product; elsewhere in the project it is called an item.
type Product struct {
	Code      string
	UnitPrice int64 // in cents
	Quantity  int64
}

// Sale is the Authorization Request (1100) of the Sale transaction.
type Sale struct {
	PAN      string // DE2
	Amount   int64  // DE4, in cents
	STAN     string // DE11
	Time     string // DE12, hhmmss; the specification has no time zone, so none is applied
	Date     string // DE13, MMDD; the specification has no year
	Expiry   string // DE14, YYMM
	RRN      string // DE37
	TID      string // DE41
	MID      string // DE42
	Currency int64  // DE49
	Products []Product
}

// SaleResponse is the Authorization Request Response (1110). It repeats the request fields, which
// the specification requires to come back with the same values.
type SaleResponse struct {
	Sale
	AuthCode     string // DE38, only present when the sale is approved
	ResponseCode string // DE39
}

// fields holds every field of a 1100 and 1110 message with its moov tag. moov leaves out fields
// with an empty value, so a request, which has no DE38 and DE39, does not send them.
type fields struct {
	MTI            string `iso8583:"0"`
	PAN            string `iso8583:"2"`
	ProcessingCode string `iso8583:"3"`
	Amount         int64  `iso8583:"4"`
	STAN           string `iso8583:"11"`
	Time           string `iso8583:"12"`
	Date           string `iso8583:"13"`
	Expiry         string `iso8583:"14"`
	RRN            string `iso8583:"37"`
	AuthCode       string `iso8583:"38"`
	ResponseCode   string `iso8583:"39"`
	TID            string `iso8583:"41"`
	MID            string `iso8583:"42"`
	Currency       int64  `iso8583:"49"`
	Products       string `iso8583:"63"`
}

// Message returns the 1100 message for the sale. The message holds the card number and expiry date
// in clear text, so do not log it: log the Sale instead. Errors never contain either.
func (s Sale) Message() (*iso8583.Message, error) {
	if err := s.validate(); err != nil {
		return nil, err
	}

	return toMessage(s.toFields(mtiSale))
}

func (s Sale) toFields(mti string) fields {
	return fields{
		MTI:            mti,
		PAN:            s.PAN,
		ProcessingCode: processingCodeSale,
		Amount:         s.Amount,
		STAN:           s.STAN,
		Time:           s.Time,
		Date:           s.Date,
		Expiry:         s.Expiry,
		RRN:            s.RRN,
		TID:            s.TID,
		MID:            s.MID,
		Currency:       s.Currency,
		Products:       formatProducts(s.Products),
	}
}

// ParseSale reads a 1100 message. It returns an error if the message is not a valid sale request.
func ParseSale(msg *iso8583.Message) (Sale, error) {
	f, err := fromMessage(msg, mtiSale)
	if err != nil {
		return Sale{}, err
	}
	if f.AuthCode != "" || f.ResponseCode != "" {
		return Sale{}, errors.New("h2h: unexpected authorisation code or response code (DE38, DE39) in a sale request")
	}

	return toSale(f)
}

// Message returns the 1110 message for the response. Like the 1100 message it holds the card number
// and expiry date in clear text, so do not log it.
func (r SaleResponse) Message() (*iso8583.Message, error) {
	if err := r.validate(); err != nil {
		return nil, err
	}

	f := r.Sale.toFields(mtiSaleResponse)
	f.AuthCode = r.AuthCode
	f.ResponseCode = r.ResponseCode

	return toMessage(f)
}

// ParseSaleResponse reads a 1110 message. It returns an error if the message is not a valid sale
// response. The response code is not checked against the known codes.
func ParseSaleResponse(msg *iso8583.Message) (SaleResponse, error) {
	f, err := fromMessage(msg, mtiSaleResponse)
	if err != nil {
		return SaleResponse{}, err
	}
	sale, err := toSale(f)
	if err != nil {
		return SaleResponse{}, err
	}

	r := SaleResponse{Sale: sale, AuthCode: f.AuthCode, ResponseCode: f.ResponseCode}
	if err := r.validateResponse(); err != nil {
		return SaleResponse{}, err
	}

	return r, nil
}

func toMessage(f fields) (*iso8583.Message, error) {
	msg := iso8583.NewMessage(Spec)
	if err := msg.Marshal(&f); err != nil {
		return nil, fmt.Errorf("h2h: marshal message: %w", err)
	}

	return msg, nil
}

func fromMessage(msg *iso8583.Message, wantMTI string) (fields, error) {
	if msg == nil {
		return fields{}, errors.New("h2h: message is nil")
	}

	var f fields
	if err := msg.Unmarshal(&f); err != nil {
		return fields{}, fmt.Errorf("h2h: unmarshal message: %w", err)
	}
	if f.MTI != wantMTI {
		return fields{}, fmt.Errorf("h2h: unexpected message type %q, want %q", f.MTI, wantMTI)
	}

	return f, nil
}

func toSale(f fields) (Sale, error) {
	if f.ProcessingCode != processingCodeSale {
		return Sale{}, errors.New("h2h: unexpected processing code (DE3)")
	}
	products, err := parseProducts(f.Products)
	if err != nil {
		return Sale{}, err
	}

	s := Sale{
		PAN:      f.PAN,
		Amount:   f.Amount,
		STAN:     f.STAN,
		Time:     f.Time,
		Date:     f.Date,
		Expiry:   f.Expiry,
		RRN:      f.RRN,
		TID:      f.TID,
		MID:      f.MID,
		Currency: f.Currency,
		Products: products,
	}
	if err := s.validate(); err != nil {
		return Sale{}, err
	}

	return s, nil
}

// validate checks every field. moov silently leaves out empty values, so without it a missing
// mandatory field would produce an incomplete message instead of an error.
func (s Sale) validate() error {
	switch {
	case !isDigits(s.PAN, 1, 19):
		return errors.New("h2h: invalid pan (DE2): want 1-19 digits")
	case !isValidExpiry(s.Expiry):
		return errors.New("h2h: invalid expiry date (DE14): want YYMM")
	case s.Amount < 1 || s.Amount > maxAmount:
		return errors.New("h2h: invalid amount (DE4): want 1-999999999999 cents")
	case !isDigits(s.STAN, 6, 6):
		return errors.New("h2h: invalid stan (DE11): want 6 digits")
	case !isValidTime(s.Time):
		return errors.New("h2h: invalid transaction time (DE12): want hhmmss")
	case !isValidDate(s.Date):
		return errors.New("h2h: invalid transaction date (DE13): want MMDD")
	case !isAlphanumeric(s.RRN, 12):
		return errors.New("h2h: invalid rrn (DE37): want 12 characters of 0-9 A-Z a-z _")
	case !isAlphanumeric(s.TID, 8):
		return errors.New("h2h: invalid terminal id (DE41): want 8 characters of 0-9 A-Z a-z _")
	case !isAlphanumeric(s.MID, 15):
		return errors.New("h2h: invalid merchant id (DE42): want 15 characters of 0-9 A-Z a-z _")
	case s.Currency < 1 || s.Currency > 999:
		return errors.New("h2h: invalid currency code (DE49): want a numeric currency code")
	}

	return validateProducts(s.Products)
}

func (r SaleResponse) validate() error {
	if err := r.Sale.validate(); err != nil {
		return err
	}

	return r.validateResponse()
}

func (r SaleResponse) validateResponse() error {
	if !isAlphanumeric(r.ResponseCode, 2) {
		return errors.New("h2h: invalid response code (DE39): want 2 characters")
	}

	// The specification sends the authorisation code only when the sale is approved.
	switch {
	case r.ResponseCode != CodeApproved && r.AuthCode != "":
		return errors.New("h2h: unexpected authorisation code (DE38): only sent when approved")
	case r.ResponseCode == CodeApproved && !isAlphanumeric(r.AuthCode, 6):
		return errors.New("h2h: invalid authorisation code (DE38): want 6 characters when approved")
	}

	return nil
}

func validateProducts(products []Product) error {
	if len(products) < 1 || len(products) > maxProducts {
		return fmt.Errorf("h2h: invalid product count (DE63): want 1-%d", maxProducts)
	}
	for _, p := range products {
		switch {
		case !isValidProductCode(p.Code):
			return errors.New("h2h: invalid product code (DE63): want 1-32 printable ascii characters without # and ;")
		case p.UnitPrice < 0 || p.UnitPrice > maxAmount:
			return errors.New("h2h: invalid product unit price (DE63): want 0-999999999999 cents")
		case p.Quantity < 0 || p.Quantity > maxQuantity:
			return errors.New("h2h: invalid product quantity (DE63): want 0-999999")
		}
	}

	return nil
}

// formatProducts returns DE63: for each product its code, a "#" delimiter, the unit price, the
// quantity and a ";" delimiter. The products must already be validated.
func formatProducts(products []Product) string {
	var b strings.Builder
	for _, p := range products {
		fmt.Fprintf(&b, "%s#%012d%06d;", p.Code, p.UnitPrice, p.Quantity)
	}

	return b.String()
}

func parseProducts(data string) ([]Product, error) {
	if data == "" {
		return nil, errors.New("h2h: product data (DE63) is missing")
	}

	var products []Product
	for rest := data; rest != ""; {
		entry, after, found := strings.Cut(rest, ";")
		if !found {
			return nil, errors.New("h2h: malformed product data (DE63): entry does not end with ;")
		}
		code, numbers, found := strings.Cut(entry, "#")
		numbersLength := unitPriceDigits + quantityDigits
		if !found || !isDigits(numbers, numbersLength, numbersLength) {
			return nil, errors.New("h2h: malformed product data (DE63): want code#, 12 price digits, 6 quantity digits")
		}

		// isDigits keeps a sign out of the numbers, so ParseInt below only fails on a programming error.
		unitPrice, err := strconv.ParseInt(numbers[:unitPriceDigits], 10, 64)
		if err != nil {
			return nil, errors.New("h2h: malformed product data (DE63): unit price is not a number")
		}
		quantity, err := strconv.ParseInt(numbers[unitPriceDigits:], 10, 64)
		if err != nil {
			return nil, errors.New("h2h: malformed product data (DE63): quantity is not a number")
		}
		products = append(products, Product{Code: code, UnitPrice: unitPrice, Quantity: quantity})
		rest = after
	}

	return products, nil
}

func isDigits(s string, minLength, maxLength int) bool {
	if len(s) < minLength || len(s) > maxLength {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}

	return true
}

// isAlphanumeric reports whether s has exactly length characters of the specification's "an" type.
func isAlphanumeric(s string, length int) bool {
	if len(s) != length {
		return false
	}
	for i := 0; i < len(s); i++ {
		c := s[i]
		isAllowed := c >= '0' && c <= '9' || c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c == '_'
		if !isAllowed {
			return false
		}
	}

	return true
}

func isValidProductCode(code string) bool {
	if len(code) < 1 || len(code) > maxCodeLength {
		return false
	}
	for i := 0; i < len(code); i++ {
		c := code[i]
		if c < 0x20 || c > 0x7E || c == '#' || c == ';' {
			return false
		}
	}

	return true
}

func isValidExpiry(expiry string) bool {
	if !isDigits(expiry, 4, 4) {
		return false
	}
	month := expiry[2:]

	return month >= "01" && month <= "12"
}

func isValidTime(value string) bool {
	if !isDigits(value, 6, 6) {
		return false
	}
	_, err := time.Parse("150405", value) // time.Parse layout: 15 hour, 04 minute, 05 second

	return err == nil
}

func isValidDate(value string) bool {
	if !isDigits(value, 4, 4) {
		return false
	}
	_, err := time.Parse("0102", value) // time.Parse layout: 01 month, 02 day

	return err == nil
}

// maskPAN keeps the first six and last four digits, the most PCI DSS allows to be shown. A number
// too short to hide anything in is masked completely.
func maskPAN(pan string) string {
	if len(pan) < 11 {
		return strings.Repeat("*", len(pan))
	}

	return pan[:6] + strings.Repeat("*", len(pan)-10) + pan[len(pan)-4:]
}

// TODO: %#v and json.Marshal still print the card number and the expiry date in clear text.
// Consider types for PAN and Expiry with masked String, GoString, MarshalJSON and LogValue.
func (s Sale) attrs() []slog.Attr {
	return []slog.Attr{
		slog.String("pan", maskPAN(s.PAN)),
		slog.Int64("amount", s.Amount),
		slog.String("stan", s.STAN),
		slog.String("rrn", s.RRN),
		slog.Int("products", len(s.Products)),
	}
}

// LogValue implements slog.LogValuer. The card number is masked and the expiry date is left out.
func (s Sale) LogValue() slog.Value {
	return slog.GroupValue(s.attrs()...)
}

// String describes the sale like LogValue does.
func (s Sale) String() string {
	return fmt.Sprintf("Sale{pan=%s amount=%d stan=%s rrn=%s}", maskPAN(s.PAN), s.Amount, s.STAN, s.RRN)
}

// LogValue implements slog.LogValuer like Sale.LogValue, with the response and authorisation codes.
func (r SaleResponse) LogValue() slog.Value {
	attrs := append(r.Sale.attrs(),
		slog.String("response_code", r.ResponseCode),
		slog.String("auth_code", r.AuthCode),
	)

	return slog.GroupValue(attrs...)
}

// String describes the response like LogValue does.
func (r SaleResponse) String() string {
	return fmt.Sprintf("SaleResponse{pan=%s amount=%d stan=%s rrn=%s code=%s}",
		maskPAN(r.PAN), r.Amount, r.STAN, r.RRN, r.ResponseCode)
}
