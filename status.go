package main

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type Status struct {
	Total       sdk.Coin `json:"total"`
	Bonded      sdk.Coin `json:"bonded"`
	Circulating sdk.Coin `json:"circulating"`
}

func getTotalSupply(ctx context.Context, bclient banktypes.QueryClient) (sdk.Coins, error) {
	var totalSupply sdk.Coins
	var nextKey []byte

	for {
		req := &banktypes.QueryTotalSupplyRequest{
			Pagination: &query.PageRequest{
				Key:   nextKey,
				Limit: 1000,
			},
		}

		bres, err := bclient.TotalSupply(ctx, req)
		if err != nil {
			return nil, err
		}

		totalSupply = totalSupply.Add(bres.Supply...)

		if bres.Pagination == nil || len(bres.Pagination.NextKey) == 0 {
			break
		}

		nextKey = bres.Pagination.NextKey
	}

	return totalSupply, nil
}

func getStatus(ctx context.Context, cctx client.Context, denom string, locked sdk.Int) (Status, error) {
	// akash query bank total
	bclient := banktypes.NewQueryClient(cctx)
	totalSupply, err := getTotalSupply(ctx, bclient)
	if err != nil {
		return Status{}, err
	}

	// akash query staking pool
	sclient := stakingtypes.NewQueryClient(cctx)
	sres, err := sclient.Pool(ctx, &stakingtypes.QueryPoolRequest{})
	if err != nil {
		return Status{}, err
	}

	ctotal := totalSupply.AmountOf(denom)
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
