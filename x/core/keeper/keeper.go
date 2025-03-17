package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/troykessler/hyperlane-cosmos/util"
	ismkeeper "github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	postdispatchkeeper "github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/keeper"
	"github.com/troykessler/hyperlane-cosmos/x/core/types"
)

type Keeper struct {
	cdc          codec.BinaryCodec
	addressCodec address.Codec

	// authority is the address capable of executing a MsgUpdateParams and other authority-gated messages.
	// Typically this should be the x/gov module account.
	authority string

	// Mailboxes is a map of mailbox IDs to mailboxes
	Mailboxes collections.Map[uint64, types.Mailbox]
	// Messages is a set of tuples. The first key is the mailbox ID, second key is the message ID.
	Messages collections.KeySet[collections.Pair[uint64, []byte]]
	// MailboxesSequence is a monotonically increasing number of mailboxes. The
	// internal ID for a mailbox is the sequence number when it was created.
	MailboxesSequence collections.Sequence

	Schema collections.Schema

	bankKeeper types.BankKeeper

	IsmKeeper          ismkeeper.Keeper
	PostDispatchKeeper postdispatchkeeper.Keeper

	ismRouter          *util.Router[util.InterchainSecurityModule]
	postDispatchRouter *util.Router[util.PostDispatchModule]
	appRouter          *util.Router[util.HyperlaneApp]
}

// NewKeeper creates a new Keeper instance
func NewKeeper(cdc codec.BinaryCodec, addressCodec address.Codec, storeService storetypes.KVStoreService, authority string, bankKeeper types.BankKeeper) Keeper {
	if _, err := addressCodec.StringToBytes(authority); err != nil {
		panic(fmt.Errorf("invalid authority address: %w", err))
	}

	sb := collections.NewSchemaBuilder(storeService)
	k := Keeper{
		cdc:               cdc,
		addressCodec:      addressCodec,
		authority:         authority,
		Mailboxes:         collections.NewMap(sb, types.MailboxesKey, "mailboxes", collections.Uint64Key, codec.CollValue[types.Mailbox](cdc)),
		Messages:          collections.NewKeySet(sb, types.MessagesKey, "messages", collections.PairKeyCodec(collections.Uint64Key, collections.BytesKey)),
		MailboxesSequence: collections.NewSequence(sb, types.MailboxesSequenceKey, "mailboxes_sequence"),
		bankKeeper:        bankKeeper,

		// REFACTORED
		IsmKeeper:          ismkeeper.NewKeeper(cdc, storeService),
		PostDispatchKeeper: postdispatchkeeper.NewKeeper(cdc, storeService, bankKeeper),

		ismRouter:          util.NewRouter[util.InterchainSecurityModule](types.IsmRouterKey, "router_ism", sb),
		postDispatchRouter: util.NewRouter[util.PostDispatchModule](types.PostDispatchRouterKey, "router_post_dispatch", sb),
		appRouter:          util.NewRouter[util.HyperlaneApp](types.AppRouterKey, "router_app", sb),
	}

	k.IsmKeeper.SetCoreKeeper(k)
	k.PostDispatchKeeper.SetCoreKeeper(k)

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}

	k.Schema = schema

	return k
}

func (k Keeper) AppRouter() *util.Router[util.HyperlaneApp] {
	return k.appRouter
}

func (k *Keeper) ReceiverIsmId(ctx context.Context, recipient util.HexAddress) (util.HexAddress, error) {
	handler, err := k.appRouter.GetModule(recipient)
	if err != nil {
		return util.HexAddress{}, err
	}
	ism, err := (*handler).ReceiverIsmId(ctx, recipient)
	if err != nil {
		return util.HexAddress{}, err
	}
	return *ism, nil
}

func (k *Keeper) Handle(ctx context.Context, mailboxId util.HexAddress, message util.HyperlaneMessage) error {
	handler, err := k.appRouter.GetModule(message.Recipient)
	if err != nil {
		return err
	}
	return (*handler).Handle(ctx, mailboxId, message)
}

func (k Keeper) IsmRouter() *util.Router[util.InterchainSecurityModule] {
	return k.ismRouter
}

func (k *Keeper) Verify(ctx context.Context, ismId util.HexAddress, metadata []byte, message util.HyperlaneMessage) (bool, error) {
	handler, err := k.ismRouter.GetModule(ismId)
	if err != nil {
		return false, err
	}
	return (*handler).Verify(ctx, ismId, metadata, message)
}

func (k *Keeper) IsmExists(ctx context.Context, ismId util.HexAddress) (bool, error) {
	handler, err := k.ismRouter.GetModule(ismId)
	if err != nil {
		return false, err
	}
	return (*handler).Exists(ctx, ismId)
}

func (k *Keeper) AssertIsmExists(ctx context.Context, id util.HexAddress) error {
	ismExists, err := k.IsmExists(ctx, id)
	if err != nil || !ismExists {
		return fmt.Errorf("ism with id %s does not exist", id.String())
	}

	return nil
}

func (k Keeper) PostDispatchRouter() *util.Router[util.PostDispatchModule] {
	return k.postDispatchRouter
}

func (k *Keeper) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	handler, err := k.postDispatchRouter.GetModule(hookId)
	if err != nil {
		return sdk.NewCoins(), err
	}
	return (*handler).PostDispatch(ctx, mailboxId, hookId, metadata, message, maxFee)
}

func (k *Keeper) PostDispatchHookExists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	handler, err := k.postDispatchRouter.GetModule(hookId)
	if err != nil {
		return false, err
	}
	return (*handler).Exists(ctx, hookId)
}

func (k *Keeper) QuoteDispatch(ctx context.Context, mailboxId, overwriteHookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage) (sdk.Coins, error) {
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), fmt.Errorf("failed to find mailbox with id %s", mailboxId.String())
	}

	// check for valid mailbox state
	if mailbox.RequiredHook == nil {
		return sdk.NewCoins(), types.ErrRequiredHookNotSet
	}
	if mailbox.DefaultHook == nil {
		return sdk.NewCoins(), types.ErrDefaultHookNotSet
	}

	calculateGasPayment := func(hookId util.HexAddress) (sdk.Coins, error) {
		handler, err := k.postDispatchRouter.GetModule(hookId)
		if err != nil {
			return sdk.NewCoins(), err
		}

		return (*handler).QuoteDispatch(ctx, mailboxId, hookId, metadata, message)
	}

	requiredGasPayment, err := calculateGasPayment(*mailbox.RequiredHook)
	if err != nil {
		return sdk.NewCoins(), err
	}

	var defaultHookId util.HexAddress
	if overwriteHookId.IsZeroAddress() {
		defaultHookId = *mailbox.DefaultHook
	} else {
		defaultHookId = overwriteHookId
	}

	defaultGasPayment, err := calculateGasPayment(defaultHookId)
	if err != nil {
		return sdk.NewCoins(), err
	}

	return sdk.Coins.Add(requiredGasPayment, defaultGasPayment...), nil
}

func (k *Keeper) AssertPostDispatchHookExists(ctx context.Context, id util.HexAddress) error {
	hookExists, err := k.PostDispatchHookExists(ctx, id)
	if err != nil || !hookExists {
		return fmt.Errorf("hook with id %s does not exist", id.String())
	}
	return nil
}

func (k Keeper) LocalDomain(ctx context.Context, mailboxId util.HexAddress) (uint32, error) {
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return 0, err
	}
	return mailbox.LocalDomain, err
}

func (k Keeper) MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error) {
	mailbox, err := k.Mailboxes.Has(ctx, mailboxId.GetInternalId())
	if err != nil {
		return false, err
	}
	return mailbox, nil
}

func (k Keeper) GetMailbox(ctx context.Context, mailboxId util.HexAddress) (types.Mailbox, error) {
	mailbox, err := k.Mailboxes.Get(ctx, mailboxId.GetInternalId())
	if err != nil {
		return types.Mailbox{}, err
	}
	return mailbox, nil
}
