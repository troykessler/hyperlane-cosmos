package keeper

import (
	"context"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

func (ms msgServer) CreateMerkleTreeHook(ctx context.Context, msg *types.MsgCreateMerkleTreeHook) (*types.MsgCreateMerkleTreeHookResponse, error) {
	mailboxId, err := util.DecodeHexAddress(msg.MailboxId)
	if err != nil {
		return nil, errors.Wrapf(types.ErrMailboxDoesNotExist, "invalid mailbox id %s", msg.MailboxId)
	}

	if exists, err := ms.k.coreKeeper.MailboxIdExists(ctx, mailboxId); !exists || err != nil {
		return nil, errors.Wrapf(types.ErrMailboxDoesNotExist, "%s", msg.MailboxId)
	}

	nextId, err := ms.k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_MERKLE_TREE)
	if err != nil {
		return nil, err
	}
	merkleTreeHook := types.MerkleTreeHook{
		Id:        nextId,
		MailboxId: mailboxId.String(),
		Owner:     msg.Owner,
		Tree:      types.ProtoFromTree(util.NewTree(util.ZeroHashes, 0)),
	}

	err = ms.k.merkleTreeHooks.Set(ctx, merkleTreeHook.Id.GetInternalId(), merkleTreeHook)
	if err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateMerkleTreeHook{
		Id:        merkleTreeHook.Id.String(),
		MailboxId: merkleTreeHook.MailboxId,
		Owner:     merkleTreeHook.Owner,
	})

	return &types.MsgCreateMerkleTreeHookResponse{
		Id: nextId.String(),
	}, nil
}

func (ms msgServer) CreateNoopHook(ctx context.Context, msg *types.MsgCreateNoopHook) (*types.MsgCreateNoopHookResponse, error) {
	nextId, err := ms.k.coreKeeper.PostDispatchRouter().GetNextSequence(ctx, types.POST_DISPATCH_HOOK_TYPE_UNUSED)
	if err != nil {
		return nil, err
	}
	noopHook := types.NoopHook{
		Id:    nextId,
		Owner: msg.Owner,
	}

	err = ms.k.noopHooks.Set(ctx, nextId.GetInternalId(), noopHook)
	if err != nil {
		return nil, err
	}

	_ = sdk.UnwrapSDKContext(ctx).EventManager().EmitTypedEvent(&types.EventCreateNoopHook{
		Id:    noopHook.String(),
		Owner: noopHook.Owner,
	})

	return &types.MsgCreateNoopHookResponse{
		Id: nextId.String(),
	}, nil
}
