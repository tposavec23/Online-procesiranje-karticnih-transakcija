package simulator_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/moov-io/iso8583-connection"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"sevenpay/backend/libs/h2h"
	"sevenpay/backend/services/bank-simulator/internal/simulator"
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

func startServer(t *testing.T) string {
	t.Helper()

	srv, err := simulator.NewH2HServer(slog.New(slog.DiscardHandler))
	require.NoError(t, err)
	require.NoError(t, srv.Start("127.0.0.1:0"))
	t.Cleanup(srv.Close)

	return srv.Addr
}

// connect opens a client connection to the simulator, the way bank-communication does.
func connect(t *testing.T, addr string) *connection.Connection {
	t.Helper()

	conn, err := connection.New(addr, h2h.Spec, h2h.ReadLength, h2h.WriteLength, connection.SendTimeout(time.Second))
	require.NoError(t, err)
	require.NoError(t, conn.Connect())
	t.Cleanup(func() {
		assert.NoError(t, conn.Close())
	})

	return conn
}

func TestNewH2HServer_approvesSale(t *testing.T) {
	t.Parallel()

	conn := connect(t, startServer(t))
	sale := exampleSale()
	msg, err := sale.Message()
	require.NoError(t, err)

	reply, err := conn.Send(msg)
	require.NoError(t, err)
	resp, err := h2h.ParseSaleResponse(reply)
	require.NoError(t, err)

	assert.Equal(t, h2h.CodeApproved, resp.ResponseCode)
	assert.Equal(t, sale, resp.Sale, "the response must repeat the request fields")
}

func TestNewH2HServer_ignoresInvalidSale(t *testing.T) {
	t.Parallel()

	conn := connect(t, startServer(t))
	invalid, err := exampleSale().Message()
	require.NoError(t, err)
	require.NoError(t, invalid.Field(38, "AA3310")) // a request must not hold an authorisation code

	_, err = conn.Send(invalid, connection.SendTimeout(200*time.Millisecond))
	require.ErrorIs(t, err, connection.ErrSendTimeout)

	next := exampleSale()
	next.STAN = "000028"
	msg, err := next.Message()
	require.NoError(t, err)
	_, err = conn.Send(msg)
	require.NoError(t, err, "the connection must stay usable after an invalid message")
}
