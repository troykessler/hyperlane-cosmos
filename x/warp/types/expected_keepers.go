package types

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/troykessler/hyperlane-cosmos/util"
	coreTypes "github.com/troykessler/hyperlane-cosmos/x/core/types"
)

type BankKeeper interface {
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	BurnCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	MintCoins(ctx context.Context, moduleName string, amt sdk.Coins) error
	GetSupply(ctx context.Context, denom string) sdk.Coin
}

type CoreKeeper interface {
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
	GetMailbox(ctx context.Context, mailboxId util.HexAddress) (coreTypes.Mailbox, error)
	AppRouter() *util.Router[util.HyperlaneApp]
	DispatchMessage(
		ctx sdk.Context,
		originMailboxId util.HexAddress,
		// sender address on the origin chain (e.g. token id)
		sender util.HexAddress,
		// the maximum amount of tokens the dispatch is allowed to cost
		maxFee sdk.Coins,
		destinationDomain uint32,
		// Recipient address on the destination chain (e.g. smart contract)
		recipient util.HexAddress,
		body []byte,
		// Metadata for postDispatch Hook
		metadata util.StandardHookMetadata,
		postDispatchHookId *util.HexAddress,
	) (messageId util.HexAddress, error error)
	QuoteDispatch(ctx context.Context, mailboxId util.HexAddress, overwriteHookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage) (sdk.Coins, error)
	AssertPostDispatchHookExists(ctx context.Context, id util.HexAddress) error
	AssertIsmExists(ctx context.Context, id util.HexAddress) error
}
