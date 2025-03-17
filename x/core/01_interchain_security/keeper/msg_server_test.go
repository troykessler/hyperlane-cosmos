package keeper_test

import (
	"crypto/ecdsa"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	types2 "github.com/troykessler/hyperlane-cosmos/x/core/types"

	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	i "github.com/troykessler/hyperlane-cosmos/tests/integration"
	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/keeper"
	"github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/types"
	keeper2 "github.com/troykessler/hyperlane-cosmos/x/core/keeper"
)

/*

TEST CASES - msg_server.go

* Create (valid) Noop ISM
* Create (invalid) MessageIdMultisig ISM with less addresses
* Create (invalid) MessageIdMultisig ISM with invalid threshold
* Create (invalid) MessageIdMultisig ISM with duplicate validator addresses
* Create (invalid) MessageIdMultisig ISM with invalid validator addresses
* Create (invalid) MessageIdMultisig ISM with unsorted validator addresses
* Create (valid) MessageIdMultisig ISM
* Create (invalid) MerkleRootMultisig ISM with less addresses
* Create (invalid) MerkleRootMultisig ISM with invalid threshold
* Create (invalid) MerkleRootMultisig ISM with duplicate validator addresses
* Create (invalid) MerkleRootMultisig ISM with invalid validator addresses
* Create (invalid) MerkleRootMultisig ISM with unsorted validator addresses
* Create (valid) MerkleRootMultisig ISM
* AnnounceValidator (invalid) with empty validator
* AnnounceValidator (invalid) with invalid validator
* AnnounceValidator (invalid) with empty storage location
* AnnounceValidator (invalid) with empty signature
* AnnounceValidator (invalid) with invalid signature
* AnnounceValidator (invalid) with invalid signature recovery id
* AnnounceValidator (invalid) same storage location for validator (replay protection)
* AnnounceValidator (invalid) for non-existing Mailbox ID
* AnnounceValidator (invalid) for invalid Mailbox ID
* AnnounceValidator (invalid) for non-matching signature validator pair
* AnnounceValidator (valid)
* AnnounceValidator (valid) add another storage location

*/

var _ = Describe("msg_server.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	It("Create (valid) Noop ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateNoopIsm{
			Creator: creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateNoopIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		var ism types.NoopISM
		typeUrl := queryISM(&ism, s, response.Id.String())
		Expect(typeUrl).To(Equal("/hyperlane.core.interchain_security.v1.NoopISM"))
		Expect(ism.Owner).To(Equal(creator.Address))
		Expect(ism.Id.String()).To(Equal(response.Id.String()))
	})

	It("Create (invalid) MessageIdMultisig ISM with less addresses", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses less than threshold: invalid multisig configuration"))
	})

	It("Create (invalid) MessageIdMultisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  0,
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid multisig configuration"))
	})

	It("Create (invalid) MessageIdMultisig ISM with duplicate validator addresses", func() {
		// Arrange
		invalidAddress := []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
		}

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Validators: invalidAddress,
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicate validator address: %v: invalid multisig configuration", invalidAddress[0])))
	})

	It("Create (invalid) MessageIdMultisig ISM with invalid validator address", func() {
		// Arrange
		validValidatorAddress := "0xa04b6a0aa112b61a7aa16c19cac27d970692995e"
		invalidAddress := []string{
			// one character more
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995ef",
			// one character less
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995",
			// invalid character included (`t`)
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995t",
		}

		for _, invalidKey := range invalidAddress {
			// Act
			_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
				Creator: creator.Address,
				Validators: []string{
					validValidatorAddress,
					invalidKey,
				},
				Threshold: 2,
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator address: %s: invalid multisig configuration", invalidKey)))
		}
	})

	It("Create (invalid) MessageIdMultisig ISM with unsorted validator addresses", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses are not sorted correctly in ascending order: invalid multisig configuration"))
	})

	It("Create (valid) MessageIdMultisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMessageIdMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMessageIdMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		var ism types.MessageIdMultisigISM
		typeURL := queryISM(&ism, s, response.Id.String())

		Expect(typeURL).To(Equal("/hyperlane.core.interchain_security.v1.MessageIdMultisigISM"))
		Expect(ism.Owner).To(Equal(creator.Address))
		Expect(ism.Threshold).To(Equal(uint32(2)))
		Expect(ism.Validators).To(HaveLen(3))
		Expect(ism.Validators[0]).To(Equal("0xa05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.Validators[1]).To(Equal("0xb05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.Validators[2]).To(Equal("0xd05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.ModuleType()).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TYPE_MESSAGE_ID_MULTISIG))
	})

	It("Create (invalid) MerkleRootMultisig ISM with less addresses", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses less than threshold: invalid multisig configuration"))
	})

	It("Create (invalid) MerkleRootMultisig ISM with invalid threshold", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator:    creator.Address,
			Validators: []string{},
			Threshold:  0,
		})

		// Assert
		Expect(err.Error()).To(Equal("threshold must be greater than zero: invalid multisig configuration"))
	})

	It("Create (invalid) MerkleRootMultisig ISM with duplicate validator addresses", func() {
		// Arrange
		invalidAddress := []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
		}

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Validators: invalidAddress,
			Threshold:  2,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("duplicate validator address: %v: invalid multisig configuration", invalidAddress[0])))
	})

	It("Create (invalid) MerkleRootMultisig ISM with invalid validator address", func() {
		// Arrange
		validValidatorAddress := "0xa04b6a0aa112b61a7aa16c19cac27d970692995e"
		invalidAddress := []string{
			// one character more
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995ef",
			// one character less
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995",
			// invalid character included (`t`)
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995t",
		}

		for _, invalidKey := range invalidAddress {
			// Act
			_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
				Creator: creator.Address,
				Validators: []string{
					validValidatorAddress,
					invalidKey,
				},
				Threshold: 2,
			})

			// Assert
			Expect(err.Error()).To(Equal(fmt.Sprintf("invalid validator address: %s: invalid multisig configuration", invalidKey)))
		}
	})

	It("Create (invalid) MerkleRootMultisig ISM with unsorted validator addresses", func() {
		// Arrange

		// Act
		_, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator addresses are not sorted correctly in ascending order: invalid multisig configuration"))
	})

	It("Create (valid) MerkleRootMultisig ISM", func() {
		// Arrange

		// Act
		res, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
			Creator: creator.Address,
			Validators: []string{
				"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
				"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
			},
			Threshold: 2,
		})

		// Assert
		Expect(err).To(BeNil())

		var response types.MsgCreateMerkleRootMultisigIsmResponse
		err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
		Expect(err).To(BeNil())

		var ism types.MerkleRootMultisigISM
		typeURL := queryISM(&ism, s, response.Id.String())

		Expect(typeURL).To(Equal("/hyperlane.core.interchain_security.v1.MerkleRootMultisigISM"))
		Expect(ism.Owner).To(Equal(creator.Address))
		Expect(ism.Threshold).To(Equal(uint32(2)))
		Expect(ism.Validators).To(HaveLen(3))
		Expect(ism.Validators[0]).To(Equal("0xa05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.Validators[1]).To(Equal("0xb05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.Validators[2]).To(Equal("0xd05b6a0aa112b61a7aa16c19cac27d970692995e"))
		Expect(ism.ModuleType()).To(Equal(types.INTERCHAIN_SECURITY_MODULE_TYPE_MERKLE_ROOT_MULTISIG))
	})

	It("AnnounceValidator (invalid) with empty validator", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       "",
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("validator cannot be empty: invalid announce"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid validator", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		invalidValidatorAddress := "0x0b1caf89d1edb9ee161093test94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       invalidValidatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid validator address: invalid announce"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with empty storage location", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: "",
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("storage location cannot be empty: invalid announce"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with empty signature", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		storageLocation := "aws://key.pub"

		// Act
		_, err := s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       "",
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("signature cannot be empty: invalid announce"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid signature", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		storageLocation := "aws://key.pub"

		// Act
		_, err := s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       "0x0b1caf89d1edb9ee161093b1ec94ca75611dtest",
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid signature: invalid announce"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) with invalid signature recovery id", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, true)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid signature recovery id: invalid signature"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) same storage location for validator (replay protection)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s already announced storage location %s: invalid announce", validatorAddress, storageLocation)))
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{storageLocation})
	})

	It("AnnounceValidator (invalid) for non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, err := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		Expect(err).To(BeNil())

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		signature := announce(validatorPrivKey, storageLocation, nonExistingMailboxId, 1337, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       nonExistingMailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s: mailbox does not exist", nonExistingMailboxId.String())))
	})

	It("AnnounceValidator (invalid) for invalid Mailbox ID", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String() + "test",
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid mailbox id: mailbox does not exist"))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (invalid) for non-matching signature validator pair", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		correctValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		wrongValidatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db392"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       wrongValidatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("validator %s doesn't match signature. recovered address: %s: invalid signature", wrongValidatorAddress, correctValidatorAddress)))
		validateAnnouncement(s, mailboxId.String(), "", []string{})
	})

	It("AnnounceValidator (valid)", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{
			storageLocation,
		})
	})

	It("AnnounceValidator (valid) add another storage location", func() {
		// Arrange
		mailboxId, _, _ := createValidMailbox(s, creator.Address, "noop")

		validatorAddress := "0x0b1caf89d1edb9ee161093b1ec94ca75611db492"
		validatorPrivKey := "38430941d3ea0e70f9a16192a833dbbf3541b3170781042067173bfe6cba4508"
		storageLocation := "aws://key.pub"

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		signature := announce(validatorPrivKey, storageLocation, mailboxId, localDomain, false)

		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})
		Expect(err).To(BeNil())

		storageLocation2 := "aws://key2.pub"
		signature = announce(validatorPrivKey, storageLocation2, mailboxId, localDomain, false)

		// Act
		_, err = s.RunTx(&types.MsgAnnounceValidator{
			Validator:       validatorAddress,
			StorageLocation: storageLocation2,
			Signature:       signature,
			MailboxId:       mailboxId.String(),
			Creator:         creator.Address,
		})

		// Assert
		Expect(err).To(BeNil())
		validateAnnouncement(s, mailboxId.String(), validatorAddress, []string{
			storageLocation,
			storageLocation2,
		})
	})
})

func createValidMailbox(s *i.KeeperTestSuite, creator string, ism string) (util.HexAddress, util.HexAddress, util.HexAddress) {
	var ismId util.HexAddress
	switch ism {
	case "noop":
		ismId = createNoopIsm(s, creator)
	case "multisig":
		ismId = createMultisigIsm(s, creator)
	}

	noopPostDispatchMock := i.CreateNoopDispatchHookHandler(s.App().HyperlaneKeeper.PostDispatchRouter())
	hook, err := noopPostDispatchMock.CreateHook(s.Ctx())
	Expect(err).To(BeNil())

	res, err := s.RunTx(&types2.MsgCreateMailbox{
		Owner:        creator,
		DefaultIsm:   ismId,
		DefaultHook:  &hook,
		RequiredHook: &hook,
	})
	Expect(err).To(BeNil())

	return verifyNewMailbox(s, res, creator, ismId.String(), hook.String(), hook.String()), hook, ismId
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&types.MsgCreateMerkleRootMultisigIsm{
		Creator: creator,
		Validators: []string{
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
		},
		Threshold: 2,
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateMerkleRootMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&types.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response types.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func verifyNewMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, defaultIsm, defaultHook, requiredHook string) util.HexAddress {
	var response types2.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(err).To(BeNil())
	Expect(mailbox.Owner).To(Equal(creator))
	Expect(mailbox.DefaultIsm.String()).To(Equal(defaultIsm))
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))
	Expect(mailbox.DefaultHook.String()).To(Equal(defaultHook))
	Expect(mailbox.RequiredHook.String()).To(Equal(requiredHook))

	mailboxes, err := keeper2.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types2.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Owner).To(Equal(creator))

	Expect(mailboxes.Mailboxes[0].DefaultIsm.String()).To(Equal(defaultIsm))
	Expect(mailboxes.Mailboxes[0].MessageSent).To(Equal(uint32(0)))
	Expect(mailboxes.Mailboxes[0].MessageReceived).To(Equal(uint32(0)))
	Expect(mailboxes.Mailboxes[0].DefaultHook.String()).To(Equal(defaultHook))
	Expect(mailboxes.Mailboxes[0].RequiredHook.String()).To(Equal(requiredHook))

	return mailboxId
}

func announce(privKey, storageLocation string, mailboxId util.HexAddress, localDomain uint32, skipRecoveryId bool) string {
	announcementDigest := types.GetAnnouncementDigest(storageLocation, localDomain, mailboxId.Bytes())

	ethDigest := util.GetEthSigningHash(announcementDigest[:])

	privateKey, err := crypto.HexToECDSA(privKey)
	Expect(err).To(BeNil())

	publicKey := privateKey.Public()
	_, ok := publicKey.(*ecdsa.PublicKey)
	Expect(ok).To(BeTrue())

	signedAnnouncement, err := crypto.Sign(ethDigest[:], privateKey)
	Expect(err).To(BeNil())

	if !skipRecoveryId {
		// Required for recovery ID
		// https://eips.ethereum.org/EIPS/eip-155
		signedAnnouncement[64] += 27
	}

	return util.EncodeEthHex(signedAnnouncement)
}

func validateAnnouncement(s *i.KeeperTestSuite, mailboxId, validatorAddress string, storageLocations []string) {
	if validatorAddress == "" {
		announcedStorageLocations, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).AnnouncedStorageLocations(s.Ctx(), &types.QueryAnnouncedStorageLocationsRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())

		Expect(announcedStorageLocations.StorageLocations).To(HaveLen(0))
	} else {
		announcedStorageLocations, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).AnnouncedStorageLocations(s.Ctx(), &types.QueryAnnouncedStorageLocationsRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())

		Expect(announcedStorageLocations.StorageLocations).To(Equal(storageLocations))

		latestAnnouncedStorageLocation, err := keeper.NewQueryServerImpl(&s.App().HyperlaneKeeper.IsmKeeper).LatestAnnouncedStorageLocation(s.Ctx(), &types.QueryLatestAnnouncedStorageLocationRequest{MailboxId: mailboxId, ValidatorAddress: validatorAddress})
		Expect(err).To(BeNil())
		Expect(latestAnnouncedStorageLocation.StorageLocation).To(Equal(storageLocations[len(storageLocations)-1]))
	}
}
