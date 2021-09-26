package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type serverCmd struct {
	Port string `help:"listen port" default:":8080"`
}

func (c *serverCmd) Run(ctx *runctx) error {

	r := mux.NewRouter()

	showSummary := func(w http.ResponseWriter, r *http.Request) {
		status, err := getStatus(r.Context(), ctx.cctx, ctx.denom, ctx.locked)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		buf, err := json.Marshal(status)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/javascript")
		w.Write(buf)
	}

	r.HandleFunc("/summary", showSummary)

	r.HandleFunc("/", showSummary)

	r.HandleFunc("/circulating", func(w http.ResponseWriter, r *http.Request) {
		status, err := getStatus(r.Context(), ctx.cctx, ctx.denom, ctx.locked)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, formatAmount(status.Circulating.Amount))
	})

	r.HandleFunc("/total", func(w http.ResponseWriter, r *http.Request) {
		status, err := getStatus(r.Context(), ctx.cctx, ctx.denom, ctx.locked)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, formatAmount(status.Total.Amount))
	})

	r.HandleFunc("/bonded", func(w http.ResponseWriter, r *http.Request) {
		status, err := getStatus(r.Context(), ctx.cctx, ctx.denom, ctx.locked)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", "text/plain")
		fmt.Fprint(w, formatAmount(status.Bonded.Amount))
	})

	server := handlers.LoggingHandler(os.Stdout, r)

	fmt.Printf("running server on port %v\n\n", c.Port)

	return http.ListenAndServe(c.Port, server)
}

func formatAmount(amount sdk.Int) string {
	whole := amount.QuoRaw(1000000).Uint64()
	frac := amount.ModRaw(1000000).Uint64()
	return fmt.Sprintf("%d.%06d", whole, frac)
}
