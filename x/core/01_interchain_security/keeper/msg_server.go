package keeper

import (
	"bytes"
	"context"

	"cosmossdk.io/errors"

	"cosmossdk.io/collections"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/troykessler/hyperlane-cosmos/util"

	"github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/types"
)

type msgServer struct {
	k *Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the module MsgServer interface.
func NewMsgServerImpl(keeper *Keeper) types.MsgServer {
	return &msgServer{k: keeper}
}

// AnnounceValidator lets a validator store a string in the state, which is queryable.
// The string should contain the storage location for the proofs (e.g. an S3 bucket)
// The Relayer uses this information to fetch the signatures for messages.
func (m msgServer) AnnounceValidator(ctx context.Context, req *types.MsgAnnounceValidator) (*types.MsgAnnounceValidatorResponse, error) {
	if req.Validator == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "validator cannot be empty")
	}

	if req.StorageLocation == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "storage location cannot be empty")
	}

	if req.Signature == "" {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "signature cannot be empty")
	}

	sig, err := util.DecodeEthHex(req.Signature)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "invalid signature")
	}

	mailboxId, err := util.DecodeHexAddress(req.MailboxId)
	if err != nil {
		return nil, errors.Wrap(types.ErrMailboxDoesNotExist, "invalid mailbox id")
	}

	found, err := m.k.coreKeeper.MailboxIdExists(ctx, mailboxId)
	if err != nil || !found {
		return nil, errors.Wrapf(types.ErrMailboxDoesNotExist, "failed to find mailbox with id: %s", mailboxId.String())
	}

	localDomain, err := m.k.coreKeeper.LocalDomain(ctx, mailboxId)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	announcementDigest := types.GetAnnouncementDigest(req.StorageLocation, localDomain, mailboxId.Bytes())
	ethSigningHash := util.GetEthSigningHash(announcementDigest[:])

	recoveredPubKey, err := util.RecoverEthSignature(ethSigningHash[:], sig)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidSignature, err.Error())
	}

	validatorAddress, err := util.DecodeEthHex(req.Validator)
	if err != nil {
		return nil, errors.Wrap(types.ErrInvalidAnnounce, "invalid validator address")
	}

	recoveredAddress := crypto.PubkeyToAddress(*recoveredPubKey)

	if !bytes.Equal(recoveredAddress[:], validatorAddress) {
		return nil, errors.Wrapf(types.ErrInvalidSignature, "validator %s doesn't match signature. recovered address: %s", util.EncodeEthHex(validatorAddress), util.EncodeEthHex(recoveredAddress[:]))
	}

	// Check if validator already exists.
	exists, err := m.k.storageLocations.Has(ctx, collections.Join3(mailboxId.GetInternalId(), validatorAddress, uint64(0)))
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	var storageLocationIndex uint64 = 0
	if exists {
		rng := collections.NewSuperPrefixedTripleRange[uint64, []byte, uint64](mailboxId.GetInternalId(), validatorAddress)

		iter, err := m.k.storageLocations.Iterate(ctx, rng)
		if err != nil {
			return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
		}

		storageLocations, err := iter.Values()
		if err != nil {
			return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
		}

		// It is assumed that a validator announces a reasonable amount of storage locations.
		// Otherwise, one would need to store the hash in a separate lookup table which adds more complexity.
		for _, location := range storageLocations {
			if location == req.StorageLocation {
				return nil, errors.Wrapf(types.ErrInvalidAnnounce, "validator %s already announced storage location %s", req.Validator, req.StorageLocation)
			}
		}
		storageLocationIndex = uint64(len(storageLocations))
	}

	if err = m.k.storageLocations.Set(ctx, collections.Join3(mailboxId.GetInternalId(), validatorAddress, storageLocationIndex), req.StorageLocation); err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	return &types.MsgAnnounceValidatorResponse{}, nil
}

func (m msgServer) CreateMessageIdMultisigIsm(ctx context.Context, req *types.MsgCreateMessageIdMultisigIsm) (*types.MsgCreateMessageIdMultisigIsmResponse, error) {
	ismId, err := m.k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.MessageIdMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidMultisigConfiguration, err.Error())
	}

	if err = m.k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	return &types.MsgCreateMessageIdMultisigIsmResponse{Id: ismId}, nil
}

func (m msgServer) CreateMerkleRootMultisigIsm(ctx context.Context, req *types.MsgCreateMerkleRootMultisigIsm) (*types.MsgCreateMerkleRootMultisigIsmResponse, error) {
	ismId, err := m.k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.MerkleRootMultisigISM{
		Id:         ismId,
		Owner:      req.Creator,
		Validators: req.Validators,
		Threshold:  req.Threshold,
	}

	if err = newIsm.Validate(); err != nil {
		return nil, errors.Wrap(types.ErrInvalidMultisigConfiguration, err.Error())
	}

	if err = m.k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	return &types.MsgCreateMerkleRootMultisigIsmResponse{Id: ismId}, nil
}

func (m msgServer) CreateNoopIsm(ctx context.Context, ism *types.MsgCreateNoopIsm) (*types.MsgCreateNoopIsmResponse, error) {
	ismId, err := m.k.coreKeeper.IsmRouter().GetNextSequence(ctx, types.INTERCHAIN_SECURITY_MODULE_TYPE_UNUSED)
	if err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	newIsm := types.NoopISM{
		Id:    ismId,
		Owner: ism.Creator,
	}

	// no validation needed, as there are no params to this ism

	if err = m.k.isms.Set(ctx, ismId.GetInternalId(), &newIsm); err != nil {
		return nil, errors.Wrap(types.ErrUnexpectedError, err.Error())
	}

	return &types.MsgCreateNoopIsmResponse{Id: ismId}, nil
}
