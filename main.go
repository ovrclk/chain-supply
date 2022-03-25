package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alecthomas/kong"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmrpc "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	denom   = "uakt"
	rpcNode = "http://135.181.60.250:26657"
)

type runctx struct {
	cctx   client.Context
	denom  string
	locked sdk.Int
}

type showCmd struct{}

func (c *showCmd) Run(ctx *runctx) error {
	status, err := getStatus(context.Background(), ctx.cctx, ctx.denom, ctx.locked)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(status, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(buf))
	return nil
}

func main() {
	var cli struct {
		Node   string `help:"RPC URI" default:"https://rpc.prod.ewr1.akash.farm:443/token/PHAH3PAI/"`
		Denom  string `help:"Denomination" default:"uakt"`
		Locked string `help:"Locked token amount" default:"28903382000000"`

		Show   showCmd   `cmd help:"Show summary" default:"1"`
		Server serverCmd `cmd help:"Run server"`
	}

	ctx := kong.Parse(&cli)

	tmclient, err := tmrpc.New(cli.Node, "")
	ctx.FatalIfErrorf(err)

	cctx := createContext()
	cctx = cctx.WithClient(tmclient)

	locked, ok := sdk.NewIntFromString(cli.Locked)
	if !ok || locked.IsNegative() {
		ctx.Fatalf("invalid locked value: %s", cli.Locked)
	}

	rctx := runctx{
		cctx:   cctx,
		denom:  cli.Denom,
		locked: locked,
	}

	ctx.FatalIfErrorf(ctx.Run(&rctx))
}
