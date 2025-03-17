package cli

import (
	"fmt"

	pdmodule "github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	ism "github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security"
	"github.com/troykessler/hyperlane-cosmos/x/core/types"
)

var (
	// SetMailbox
	defaultIsm   string
	defaultHook  string
	requiredHook string
	newOwner     string
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewMailboxCmd(),
		ism.GetTxCmd(),
		pdmodule.GetTxCmd(),
	)

	return txCmd
}
