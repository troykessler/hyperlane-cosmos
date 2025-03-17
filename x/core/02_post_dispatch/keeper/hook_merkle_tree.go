package keeper

import (
	"context"

	"cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

type MerkleTreeHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = MerkleTreeHookHandler{}

func (i MerkleTreeHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := i.k.merkleTreeHooks.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

func (i MerkleTreeHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE
}

// PostDispatch inserts a message ID into the Merkle tree of a mailbox.
// It updates the tree, emits an event, and stores the changes.
func (i MerkleTreeHookHandler) PostDispatch(ctx context.Context, mailboxId, hookId util.HexAddress, _ util.StandardHookMetadata, message util.HyperlaneMessage, _ sdk.Coins) (sdk.Coins, error) {
	merkleTreeHook, err := i.k.merkleTreeHooks.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return nil, err
	}

	if merkleTreeHook.MailboxId != mailboxId.String() {
		return nil, errors.Wrapf(types.ErrSenderIsNotDesignatedMailbox, "required mailbox id: %s, sender mailbox id: %s", merkleTreeHook.MailboxId, mailboxId.String())
	}

	tree, err := types.TreeFromProto(merkleTreeHook.Tree)
	if err != nil {
		return nil, err
	}

	count := tree.GetCount()

	if err = tree.Insert(message.Id()); err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.InsertedIntoTree{
		MessageId:        message.Id().String(),
		Index:            count,
		MerkleTreeHookId: merkleTreeHook.Id.String(),
	})

	merkleTreeHook.Tree = types.ProtoFromTree(tree)

	if err := i.k.merkleTreeHooks.Set(ctx, hookId.GetInternalId(), merkleTreeHook); err != nil {
		return nil, err
	}

	return sdk.NewCoins(), nil
}

func (i MerkleTreeHookHandler) QuoteDispatch(_ context.Context, _, _ util.HexAddress, _ util.StandardHookMetadata, _ util.HyperlaneMessage) (sdk.Coins, error) {
	return sdk.NewCoins(), nil
}
