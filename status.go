package main

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Status struct {
	Total       sdk.Coin `json:"total"`
	Bonded      sdk.Coin `json:"bonded"`
	Circulating sdk.Coin `json:"circulating"`
}

func getStatus(ctx context.Context, cctx client.Context, denom string, locked sdk.Int) (Status, error) {

	// akash query bank total
	bclient := banktypes.NewQueryClient(cctx)
	bres, err := bclient.TotalSupply(ctx, &banktypes.QueryTotalSupplyRequest{})
	if err != nil {
		return Status{}, err
	}

	// akash query staking pool
	sclient := stakingtypes.NewQueryClient(cctx)
	sres, err := sclient.Pool(ctx, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return Status{}, err
	}

	ctotal := bres.Supply.AmountOf(denom)
	cbonded := sres.Pool.BondedTokens

	return Status{
		Total:       sdk.NewCoin(denom, ctotal),
		Bonded:      sdk.NewCoin(denom, cbonded),
		Circulating: sdk.NewCoin(denom, ctotal.Sub(locked)),
	}, nil
}

func createContext() client.Context {
	iregistry := codectypes.NewInterfaceRegistry()
	banktypes.RegisterInterfaces(iregistry)
	stakingtypes.RegisterInterfaces(iregistry)
	cryptocodec.RegisterInterfaces(iregistry)

	cctx := client.Context{}
	cctx = cctx.WithOffline(false)
	cctx = cctx.WithInterfaceRegistry(iregistry)

	return cctx
}
