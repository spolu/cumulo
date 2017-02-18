package app

import (
	"github.com/spolu/cumulo/api/endpoint"
	goji "goji.io"
	"goji.io/pat"
)

// Controller binds the API
type Controller struct{}

// Bind registers the API routes.
func (c *Controller) Bind(
	mux *goji.Mux,
) {
	// Public
	mux.HandleFunc(pat.Post("/users"), endpoint.HandlerFor(endpoint.EndPtCreateUser))
	// mux.HandleFunc(pat.Post("/users/:user/login"), endpoint.HandlerFor(endpoint.EndPtLoginUser))
	//   send sms
	// mux.HandleFunc(pat.Post("/users/:user/verify"), endpoint.HandlerFor(endpoint.EndPtVerifyUser))
	//  check that
	//    run lnd -> move to a middleware + lndpool

	// Authenticated.
	// mux.HandleFunc(pat.Get("/balances"), endpoint.HandlerFor(endpoint.EndPtListBalances))

	// mux.HandleFunc(pat.Post("/receivers"), endpoint.HandlerFor(endpoint.EndPtCreateReceiver))
	// generate an address to fund the wallet

	// mux.HandleFunc(pat.Post("/transfers"), endpoint.HandlerFor(endpoint.EndPtCreateTransfer))
	// send bitcoin:
	// - to phone number
	// - lightning network address
	// - lightning network payment request

	// mux.HandleFunc(pat.Get("/transfers"), endpoint.HandlerFor(endpoint.EndPtListTransfers))
	// list all transfers:
	// - the ones you created
	// - the ones received
}
