# h2h

Module `sevenpay/backend/libs/h2h` builds and reads the H2H messages that Bank communication service 1 and the bank (Bank simulator 1) exchange over TCP. H2H is the protocol from the specification we received from 7pay, and it is based on ISO 8583.

Both services import this module, so the message format is written in one place and the services cannot disagree on it. The module knows nothing about the network, timeouts or approval rules; the services handle those.

## Communication

Every message starts with a four-digit message type (MTI). A payment uses two of them:

| MTI | Message | Type in the module | Sent by |
|---|---|---|---|
| 1100 | authorization request | `Sale` | Bank communication |
| 1110 | bank response | `SaleResponse` | the bank (simulator) |

The response repeats the request fields and adds the response code (data element 39, written DE39) and, only when the payment is approved, the authorization code (DE38).

```
   ┌──────────────────────┐                       ┌──────────────────────┐
   │  Bank communication  │                       │    Bank simulator    │
   └───────────┬──────────┘                       └───────────┬──────────┘
               │                                              │
 ┌─────────────┴────────────────────┐                         │
 │ 1. sale.Message()                │                         │
 │    Sale → 1100 request           │                         │
 └─────────────┬────────────────────┘                         │
               │                                              │
               │        2. 1100 · authorization request       │
               ├─────────────────────────────────────────────►│
               │                conn.Send(msg)                │
               │                             ┌────────────────┴────────────────┐
               │                             │ 3. h2h.ParseSale(msg)           │
               │                             │    1100 request → Sale          │
               │                             └────────────────┬────────────────┘
               │                             ┌────────────────┴────────────────┐
               │                             │ 4. decide                       │
               │                             │    approve (00) or decline      │
               │                             └────────────────┬────────────────┘
               │                             ┌────────────────┴────────────────┐
               │                             │ 5. resp.Message()               │
               │                             │    SaleResponse → 1110 response │
               │                             └────────────────┬────────────────┘
               │                                              │
               │            6. 1110 · bank response           │
               │◄─────────────────────────────────────────────┤
               │                c.Reply(reply)                │
 ┌─────────────┴────────────────────┐                         │
 │ 7. h2h.ParseSaleResponse(reply)  │                         │
 │    1110 response → SaleResponse  │                         │
 └──────────────────────────────────┘                         │
```

| Step | What happens |
|---|---|
| 1 | Bank communication builds the 1100 request from a `Sale` |
| 2 | moov sends it over TCP (`conn.Send`) |
| 3 | the simulator reads the received message into a `Sale` |
| 4 | the simulator decides whether to approve (`00`) or decline the payment |
| 5 | it builds the 1110 response from a `SaleResponse` |
| 6 | moov sends it back on the same connection (`c.Reply`) |
| 7 | Bank communication reads the response into a `SaleResponse` |

`conn.Send` pairs the 1110 with its 1100 by STAN (DE11), so STAN must not repeat among requests waiting on the same connection. The specification also requires STAN to grow by one for every transaction (999999 rolls over to 000000) and RRN to be unique; Bank communication creates both, the module does not.

### What happens to the message

The same 1100 request in its three forms, using the example from the specification (12.99 € for two cart items). The steps marked moov are done by [moov iso8583](https://github.com/moov-io/iso8583) and [iso8583-connection](https://github.com/moov-io/iso8583-connection); the services never call them by hand.

Format, from the specification:

- `n` only digits, `an` digits, letters and `_`, `ans` any printable character (the specification writes digits stored as ASCII as `(a)n`);
- the number is the exact length (`n 6`), and `..` means a maximum length (`n ..19`);
- ASCII stores one character per byte, BCD two digits per byte;
- a field of variable length has its length in front (1 or 2 bytes, BCD);
- `Amount` (DE4) and `Currency` (DE49) are Go numbers, and moov pads them with zeros on the left (`1299` becomes `000000001299`); every other fixed-length field is a string and must already have its full length (STAN `"000027"`, not `"27"`).

```
 Sale (Go struct, fields by name)
   PAN:  "3650788453316173"   Amount: 1299            STAN:   "000027"
   Time: "083022"             Date:   "1127"          Expiry: "2705"
   RRN:  "12JSWD232L1F"       TID:    "TT000008"      MID:    "TMP000000000001"
   Currency: 978
   Products: {Code: "123456-7890", UnitPrice: 400, Quantity: 3},
             {Code: "AFF-32-423",  UnitPrice: 99,  Quantity: 1}

        │
        │  ↓ send:     sale.Message()     checks every field, copies the Sale into the private struct fields
        │                                 (DE3 = 000000, Products as the DE63 text), then moov's Marshal()
        │                                 puts each struct field into the message field named in its iso8583 tag
        │
        │  ↑ receive:  h2h.ParseSale()    moov's Unmarshal() copies the message back into fields,
        │                                 every field is checked again, and the result is a Sale
        │
        ▼ ▲

 *iso8583.Message (moov, fields by number)
   field value              format                         meaning
   MTI   1100               n 4, BCD                       message type: authorization request
   DE2   3650788453316173   n ..19, ASCII, length 1 B      card number (PAN)
   DE3   000000             n 6, ASCII                     processing code, always 000000 for a sale
   DE4   1299               n 12, ASCII                    amount in cents (12.99 €)
   DE11  000027             n 6, ASCII                     STAN, sequence number of the transaction
   DE12  083022             n 6, ASCII                     time, hhmmss
   DE13  1127               n 4, ASCII                     date, MMDD
   DE14  2705               n 4, ASCII                     card expiry date, YYMM
   DE37  12JSWD232L1F       an 12, ASCII                   RRN, reference number of the transaction
   DE41  TT000008           an 8, ASCII                    terminal ID
   DE42  TMP000000000001    an 15, ASCII                   merchant ID
   DE49  978                n 3, BCD                       currency (978 = euro)
   DE63  see below          ans ..520, ASCII, length 2 B   cart items
         123456-7890#000000000400000003;AFF-32-423#000000000099000001;
         each item: code # unit price (12 digits) quantity (6 digits) ;

        │
        │  ↓ send:     Pack()             moov: writes every field as h2h.Spec says
        │              h2h.WriteLength    called by moov: puts the message length in front
        │
        │  ↑ receive:  h2h.ReadLength     called by moov: reads how many bytes follow
        │              Unpack()           moov: reads every field as h2h.Spec says
        │
        ▼ ▲

 bytes on the network (167 bytes, hex)
   00 A5                         length (binary): 0x00A5 = 165 bytes follow
   11 00                         MTI 1100 (BCD, 2 digits per byte)
   70 3C 00 00 08 C0 80 02       bitmap: fields 2, 3, 4, 11–14, 37, 41, 42, 49, 63 are present
   16 33 36 35 30 … 31 37 33     DE2: length 16 (BCD), then "3650788453316173" (ASCII)
   30 30 30 30 30 30             DE3: "000000"
   30 30 30 30 30 30 30 30 31 32 39 39
                                 DE4: "000000001299", 12 digits with leading zeros
   …                             DE11, DE12, DE13, DE14, DE37, DE41, DE42 (ASCII)
   09 78                         DE49: 978 (BCD)
   00 61 31 32 33 … 31 3B        DE63: length 61 (BCD), then the items (ASCII)
```

The 1110 response goes the same way, with MTI 1110 and two more fields: DE39 (response code) and, when the payment is approved, DE38 (authorization code).

### Why a private struct and Marshal()

moov can fill a message from a struct whose fields carry the field number in a tag, and read it back the same way:

```go
type fields struct {
	MTI      string `iso8583:"0"`
	PAN      string `iso8583:"2"`
	Amount   int64  `iso8583:"4"`
	// … every field of the 1100 and 1110 messages
	Products string `iso8583:"63"`
}
```

So no field is set or read by hand by its number. `fields` mirrors the message one to one and stays private; `Sale` and `SaleResponse` are what the services use, without the MTI and DE3, and with the cart items as a list instead of the DE63 text.

## Rough examples

Bank communication sends a request and reads the response:

```go
msg, err := sale.Message()                 // Sale → 1100 request
reply, err := conn.Send(msg)               // moov: sends it and waits for the response
resp, err := h2h.ParseSaleResponse(reply)  // 1110 response → SaleResponse
```

Bank simulator reads a request and responds:

```go
sale, err := h2h.ParseSale(msg)            // received 1100 request → Sale
resp := h2h.SaleResponse{Sale: sale, ResponseCode: h2h.CodeApproved, AuthCode: "AB12CD"}
reply, err := resp.Message()               // → 1110 response
err = c.Reply(reply)                       // moov: sends the response on the same connection
```

A service imports the module through its `go.mod`:

```
require sevenpay/backend/libs/h2h v0.0.0
replace sevenpay/backend/libs/h2h => ../../libs/h2h
```

## Validation

`Message()` checks every field before building a message and `Parse…` after reading one, because moov silently leaves out an empty field.

| Field | Rule |
|---|---|
| MTI | 1100 for `ParseSale`, 1110 for `ParseSaleResponse` |
| `PAN` (DE2) | 1–19 digits |
| DE3 | always `000000` |
| `Amount` (DE4) | 1 to 999,999,999,999 cents |
| `STAN` (DE11) | 6 digits |
| `Time` (DE12), `Date` (DE13) | valid time hhmmss and date MMDD |
| `Expiry` (DE14) | YYMM, month 01–12 |
| `RRN` (DE37), `TID` (DE41), `MID` (DE42) | 12, 8 and 15 characters of `0-9 A-Z a-z _` |
| `Currency` (DE49) | 1–999 |
| `Products` (DE63) | 1–10 items; code of 1–32 printable ASCII characters without `#` and `;`; price 0–999,999,999,999; quantity 0–999,999; when reading, every item must be `code#`, 12 + 6 digits and `;` |
| `ResponseCode` (DE39) | not in a request; in a response 2 characters of `0-9 A-Z a-z _`, any value (the service rejects unknown codes) |
| `AuthCode` (DE38) | not in a request; in a response 6 characters with `00`, otherwise empty |

The module checks only the shape of a message, not its meaning: not a valid card number, not items adding up to the amount, and not that a 1110 repeats the fields of its 1100; the services check those themselves. Its errors name the field but never contain the card number or the expiry date.

## Logging

Never log or print `*iso8583.Message`: it holds the card number and the expiry date in clear text. Log `Sale` or `SaleResponse` instead.

- `LogValue()` is used by `slog`: the card number is masked (first 6 and last 4 digits), the expiry date is left out, and only the number of cart items is shown;
- `String()` is used by `fmt` (`%v`, `Println`, error messages) and prints a shorter form, masked the same way:
  `SaleResponse{pan=365078******6173 amount=1299 stan=000027 rrn=12JSWD232L1F code=00}`.

```go
logger.Info("sale approved", "sale", resp)
```

```json
"sale":{"pan":"365078******6173","amount":1299,"stan":"000027","rrn":"12JSWD232L1F","products":2,"response_code":"00","auth_code":"FLRRUB"}
```