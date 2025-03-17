package types

import (
	"context"
	"encoding/binary"
	"fmt"
	"slices"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/troykessler/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &MerkleRootMultisigISM{}

func (m *MerkleRootMultisigISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

func (m *MerkleRootMultisigISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG
}

func (m *MerkleRootMultisigISM) Verify(_ context.Context, rawMetadata []byte, message util.HyperlaneMessage) (bool, error) {
	metadata, err := NewMerkleRootMultisigMetadata(rawMetadata)
	if err != nil {
		return false, err
	}

	if metadata.MessageIndex > metadata.SignedIndex {
		return false, fmt.Errorf("invalid signed index")
	}

	digest := metadata.Digest(&message)

	return VerifyMultisig(m.Validators, m.Threshold, metadata.Signatures, digest)
}

func (m *MerkleRootMultisigISM) GetThreshold() uint32 {
	return m.Threshold
}

func (m *MerkleRootMultisigISM) GetValidators() []string {
	return m.Validators
}

func (m *MerkleRootMultisigISM) Validate() error {
	return ValidateNewMultisig(m)
}

type MerkleRootMultisigMetadata struct {
	MerkleTreeHook  [32]byte
	MessageIndex    uint32
	MerkleProof     [32][32]byte
	SignedIndex     uint32
	SignatureCount  uint32
	SignedMessageId [32]byte
	Signatures      [][]byte
}

var (
	// originMerkleTreeOffset := 0
	messageIndexOffset = 32
	messageIdOffset    = 36
	merkleProofOffset  = 68
	merkleProofLength  = 32 * 32
	signedIndexOffset  = 1092
	signaturesOffset   = 1096
	signatureLength    = 65
)

// NewMerkleRootMultisigMetadata validates and creates a new metadata object
func NewMerkleRootMultisigMetadata(metadata []byte) (MerkleRootMultisigMetadata, error) {
	/*
	 * Format of metadata:
	 * [   0:  32] Origin merkle tree address
	 * [  32:  36] Index of message ID in merkle tree
	 * [  36:  68] Signed checkpoint message ID
	 * [  68:1092] Merkle proof
	 * [1092:1096] Signed checkpoint index (computed from proof and index)
	 * [1096:????] Validator signatures (length := threshold * 65)
	 */

	if len(metadata) < signaturesOffset {
		return MerkleRootMultisigMetadata{}, fmt.Errorf("invalid metadata length: got %v, expected at least %v bytes", len(metadata), signaturesOffset)
	}

	signaturesLen := len(metadata) - signaturesOffset
	signatureCount := uint32(signaturesLen / signatureLength)

	if signaturesLen%signatureLength != 0 {
		return MerkleRootMultisigMetadata{}, fmt.Errorf("invalid signatures length in metadata")
	}

	var signatures [][]byte
	for i := 0; i < int(signatureCount); i++ {
		start := signaturesOffset + (i * signatureLength)
		sig := make([]byte, signatureLength)
		copy(sig, metadata[start:start+signatureLength])
		signatures = append(signatures, sig)
	}

	var merkleTreeHook [32]byte
	copy(merkleTreeHook[:], (metadata)[:32])

	proof := (metadata)[merkleProofOffset : merkleProofOffset+merkleProofLength]
	// proof is a 32 element long array of hashes encoded as 32 byte long arrays
	var merkleProof [32][32]byte
	for i := 0; i < 32; i++ {
		copy(merkleProof[i][:], proof[i*32:(i+1)*32])
	}

	var signedMessageId [32]byte
	copy(signedMessageId[:], (metadata)[messageIdOffset:messageIdOffset+32])

	return MerkleRootMultisigMetadata{
		MerkleTreeHook:  merkleTreeHook,
		MessageIndex:    binary.BigEndian.Uint32((metadata)[messageIndexOffset:]),
		MerkleProof:     merkleProof,
		SignedIndex:     binary.BigEndian.Uint32((metadata)[signedIndexOffset:]),
		SignatureCount:  uint32(signaturesLen / signaturesOffset),
		SignedMessageId: signedMessageId,
		Signatures:      signatures,
	}, nil
}

func (m *MerkleRootMultisigMetadata) Digest(message *util.HyperlaneMessage) [32]byte {
	messageId := message.Id()
	signedRoot := util.BranchRoot(messageId, m.MerkleProof, m.MessageIndex)

	return checkpointDigest(
		message.Origin,
		m.MerkleTreeHook,
		signedRoot,
		m.SignedIndex,
		m.SignedMessageId,
	)
}

func (m *MerkleRootMultisigMetadata) Bytes() []byte {
	messageIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(messageIndex, m.MessageIndex)

	var merkleProofBytes []byte
	for _, proof := range m.MerkleProof {
		merkleProofBytes = append(merkleProofBytes, proof[:]...)
	}

	signedIndex := make([]byte, 4)
	binary.BigEndian.PutUint32(signedIndex, m.SignedIndex)

	var signaturesBytes []byte
	for _, sig := range m.Signatures {
		signaturesBytes = append(signaturesBytes, sig...)
	}

	return slices.Concat(
		m.MerkleTreeHook[:],
		messageIndex,
		m.SignedMessageId[:],
		merkleProofBytes,
		signedIndex,
		signaturesBytes,
	)
}

func checkpointDigest(origin uint32, merkleTreeHook, checkpointRoot [32]byte, checkpointIndex uint32, messageId [32]byte) [32]byte {
	hash := domainHash(origin, merkleTreeHook)

	bytes := make([]byte, 0, 32+32+4+32)
	bytes = append(bytes, hash[:]...)
	bytes = append(bytes, checkpointRoot[:]...)
	bytes = binary.BigEndian.AppendUint32(bytes, checkpointIndex)
	bytes = append(bytes, messageId[:]...)

	return util.GetEthSigningHash(crypto.Keccak256(bytes))
}

func domainHash(origin uint32, merkleTreeHook [32]byte) [32]byte {
	bytes := make([]byte, 0, 46)

	bytes = binary.BigEndian.AppendUint32(bytes, origin)
	bytes = append(bytes, merkleTreeHook[:]...)
	bytes = append(bytes, []byte("HYPERLANE")...)

	return crypto.Keccak256Hash(bytes)
}
