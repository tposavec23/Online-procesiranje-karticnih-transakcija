// Package simulator simulates Bank 1. It accepts H2H (ISO 8583) connections from the
// bank-communication service and answers each sale request (1100) with a sale response (1110).
package simulator

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/moov-io/iso8583"
	"github.com/moov-io/iso8583-connection"
	"github.com/moov-io/iso8583-connection/server"

	"sevenpay/backend/libs/h2h"
)

// NewH2HServer returns a server that answers every H2H sale request (1100) with an approved
// sale response (1110); an invalid request gets no response. The caller starts it with Start
// and stops it with Close.
func NewH2HServer(logger *slog.Logger) (*server.Server, error) {
	srv := server.New(h2h.Spec, h2h.ReadLength, h2h.WriteLength,
		connection.InboundMessageHandler(handleSale(logger)),
		connection.ErrorHandler(func(err error) { logConnectionError(logger, err) }),
	)
	err := srv.SetOptions(server.WithErrorHandler(func(err error) {
		logger.Error("h2h server", "error", err)
	}))
	if err != nil {
		return nil, fmt.Errorf("set h2h server options: %w", err)
	}
	srv.AddConnectionHandler(func(conn net.Conn) {
		logger.Info("h2h connection accepted", "remote_addr", conn.RemoteAddr().String())
	})

	return srv, nil
}

// moov calls handleSale concurrently, in a new goroutine for every received message.
func handleSale(logger *slog.Logger) func(*connection.Connection, *iso8583.Message) {
	return func(c *connection.Connection, msg *iso8583.Message) {
		sale, err := h2h.ParseSale(msg)
		if err != nil {
			// A 1110 must repeat the request fields, so an invalid request gets no response
			// and bank-communication times out.
			logger.Warn("h2h message ignored", "error", err)
			return
		}

		resp := h2h.SaleResponse{Sale: sale, AuthCode: newAuthCode(), ResponseCode: h2h.CodeApproved}
		reply, err := resp.Message()
		if err != nil {
			logger.Error("build h2h sale response", "sale", sale, "error", err)
			return
		}
		if err := c.Reply(reply); err != nil {
			logger.Warn("send h2h sale response", "sale", resp, "error", err)
			return
		}
		logger.Info("sale approved", "sale", resp)
	}
}

// newAuthCode returns an authorisation code (DE38): six characters of A-Z and 2-7, all allowed
// by the specification's "an" type.
func newAuthCode() string {
	return rand.Text()[:6]
}

func logConnectionError(logger *slog.Logger, err error) {
	if errors.Is(err, io.EOF) {
		logger.Info("h2h connection closed by client")
		return
	}
	logger.Warn("h2h connection", "error", err)
}
