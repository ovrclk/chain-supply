module github.com/ovrclk/chain-supply

go 1.16

require (
	github.com/alecthomas/kong v0.2.16
	github.com/cosmos/cosmos-sdk v0.41.3
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/tendermint/tendermint v0.34.8
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.3-alpha.regen.1

replace google.golang.org/grpc => google.golang.org/grpc v1.33.2

replace github.com/cosmos/cosmos-sdk => github.com/ovrclk/cosmos-sdk v0.41.4-akash-4

replace github.com/tendermint/tendermint => github.com/ovrclk/tendermint v0.34.9-akash-1
