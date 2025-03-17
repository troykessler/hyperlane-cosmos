package interchain_security

import (
	"github.com/cosmos/gogoproto/grpc"
	"github.com/spf13/cobra"
	"github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/types"
)

import "github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/client/cli"

// GetTxCmd returns the root command for the core ISMs
func GetTxCmd() *cobra.Command {
	return cli.GetTxCmd()
}

// RegisterMsgServer registers the core ism handler for transactions
func RegisterMsgServer(server grpc.Server, msgServer types.MsgServer) {
	types.RegisterMsgServer(server, msgServer)
}

// RegisterQueryService registers the gRPC query service for api queries
func RegisterQueryService(server grpc.Server, queryServer types.QueryServer) {
	types.RegisterQueryServer(server, queryServer)
}
