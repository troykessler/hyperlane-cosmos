package keeper_test

import (
	"encoding/hex"
	"fmt"

	pdTypes "github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"

	ismtypes "github.com/troykessler/hyperlane-cosmos/x/core/01_interchain_security/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	i "github.com/troykessler/hyperlane-cosmos/tests/integration"
	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/keeper"
	"github.com/troykessler/hyperlane-cosmos/x/core/types"
)

/*

TEST CASES - msg_mailbox.go

* CreateMailbox (invalid) with invalid default ISM and without hooks
* CreateMailbox (invalid) with non-existing default ISM and without hooks
* CreateMailbox (invalid) with valid default ISM (Noop) and invalid default hook
* CreateMailbox (invalid) with valid default ISM (Multisig) and invalid default hook
* CreateMailbox (invalid) with valid default ISM (Noop) and non-existent default hook
* CreateMailbox (invalid) with valid default ISM (Multisig) and non-existent default hook
* DispatchMessage (valid) with NoopISM
* DispatchMessage (valid) with MultisigISM
* DispatchMessage (valid) with custom hook
* DispatchMessage (valid)
* ProcessMessage (invalid) with empty message
* ProcessMessage (invalid) with invalid non-hex message
* ProcessMessage (invalid) with invalid metadata (Noop ISM)
* SetMailbox (valid) without hooks
* SetMailbox (valid) without hooks
* SetMailbox (valid)

*/

var _ = Describe("msg_mailbox.go", Ordered, func() {
	var s *i.KeeperTestSuite
	var creator i.TestValidatorAddress
	var sender i.TestValidatorAddress

	BeforeEach(func() {
		s = i.NewCleanChain()
		creator = i.GenerateTestValidatorAddress("Creator")
		sender = i.GenerateTestValidatorAddress("Sender")
		err := s.MintBaseCoins(creator.Address, 1_000_000)
		Expect(err).To(BeNil())
	})

	// CreateMailbox
	It("CreateMailbox (invalid) with invalid default ISM and without hooks", func() {
		// Arrange
		defaultIsm, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:      creator.Address,
			DefaultIsm: defaultIsm,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with non-existing default ISM and without hooks", func() {
		// Arrange
		defaultIsm, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:      creator.Address,
			DefaultIsm: defaultIsm,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("ism with id %s does not exist", defaultIsm)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Noop) and invalid default hook", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and invalid default hook", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Noop) and non-existent default hook", func() {
		// Arrange
		ismId := createNoopIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0x004b867052ca9c65e33362112f35fb548f8732c2fe45f07b9c591b38e865def0")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (invalid) with valid default ISM (Multisig) and non-existent default hook", func() {
		// Arrange
		ismId := createMultisigIsm(s, creator.Address)
		igpId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")

		// Act
		_, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:       creator.Address,
			DefaultIsm:  ismId,
			DefaultHook: &igpId,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("hook with id %s does not exist", igpId)))

		verifyInvalidMailboxCreation(s)
	})

	It("CreateMailbox (valid) with NoopISM", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		noopHookId := createNoopHook(s, creator.Address)
		ismId := createNoopIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:        creator.Address,
			DefaultIsm:   ismId,
			RequiredHook: &igpId,
			DefaultHook:  &noopHookId,
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewSingleMailbox(s, res, creator.Address, ismId.String(), igpId.String(), noopHookId.String())
	})

	It("CreateMailbox (valid) with MultisigISM", func() {
		// Arrange
		igpId := createIgp(s, creator.Address)
		noopHookId := createNoopHook(s, creator.Address)
		ismId := createMultisigIsm(s, creator.Address)

		// Act
		res, err := s.RunTx(&types.MsgCreateMailbox{
			Owner:        creator.Address,
			DefaultIsm:   ismId,
			RequiredHook: &igpId,
			DefaultHook:  &noopHookId,
		})

		// Assert
		Expect(err).To(BeNil())

		verifyNewSingleMailbox(s, res, creator.Address, ismId.String(), igpId.String(), noopHookId.String())
	})

	It("DispatchMessage (valid) with custom hook", func() {
		// Arrange
		mailboxId, igpId, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		hexSender, _ := util.DecodeHexAddress(sender.Address)
		recipient, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		body, _ := hex.DecodeString("0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

		_, err = s.App().HyperlaneKeeper.DispatchMessage(
			s.Ctx(),
			mailboxId,
			hexSender,
			sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(1250000))),
			1,
			recipient,
			body,
			util.StandardHookMetadata{
				GasLimit: math.NewInt(50000),
				Address:  sender.AccAddress,
			},
			&igpId,
		)

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	It("DispatchMessage (valid)", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		hexSender, _ := util.DecodeHexAddress(sender.Address)
		recipient, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		body, _ := hex.DecodeString("0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

		_, err = s.App().HyperlaneKeeper.DispatchMessage(
			s.Ctx(),
			mailboxId,
			hexSender,
			sdk.NewCoins(sdk.NewCoin("acoin", math.NewInt(1000000))),
			1,
			recipient,
			body,
			util.StandardHookMetadata{
				GasLimit: math.NewInt(50000),
				Address:  sender.AccAddress,
			},
			nil,
		)

		// Assert
		Expect(err).To(BeNil())

		verifyDispatch(s, mailboxId, 1)
	})

	It("ProcessMessage (invalid) with empty message", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   "",
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid message"))
	})

	It("ProcessMessage (invalid) with invalid non-hex message", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   "test123",
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to decode message"))
	})

	It("ProcessMessage (invalid) with invalid metadata (Noop ISM)", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "xxx",
			Message:   "0xe81bf6f262305f49f318d68f33b04866f092ffdb2ecf9c98469b4a8b829f65e4",
		})

		// Assert
		Expect(err.Error()).To(Equal("failed to decode metadata"))
	})

	PIt("ProcessMessage (valid) (Noop ISM)", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		// TODO: Create token to use as recipient
		recipient := util.CreateMockHexAddress("recipient", 0)

		message := util.HyperlaneMessage{
			Version:     1,
			Nonce:       1,
			Origin:      0,
			Sender:      util.CreateMockHexAddress("sender", 0),
			Destination: 1,
			Recipient:   recipient,
			Body:        nil,
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   message.String(),
		})

		// Assert
		Expect(err).To(BeNil())
	})

	It("SetMailbox (valid) without hooks", func() {
		// Arrange
		mailboxId, requiredHook, defaultHook, ism := createValidMailbox(s, creator.Address, "noop", 1)

		noopIsmId := createNoopIsm(s, sender.Address)
		defaultHookId := createIgp(s, creator.Address)
		requiredHookId := createIgp(s, creator.Address)
		newOwner := "new_owner"

		// Act
		_, err := s.RunTx(&types.MsgSetMailbox{
			Owner:        sender.Address,
			MailboxId:    mailboxId,
			DefaultIsm:   &noopIsmId,
			DefaultHook:  &defaultHookId,
			RequiredHook: &requiredHookId,
			NewOwner:     newOwner,
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("%s does not own mailbox with id %s", sender.Address, mailboxId.String())))

		mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(mailbox.DefaultIsm).To(Equal(ism))
		Expect(mailbox.DefaultHook).To(Equal(&defaultHook))
		Expect(mailbox.RequiredHook).To(Equal(&requiredHook))
		Expect(mailbox.Owner).To(Equal(creator.Address))
	})

	It("SetMailbox (valid) without hooks", func() {
		// Arrange
		mailboxId, requiredHook, defaultHook, _ := createValidMailbox(s, creator.Address, "noop", 1)

		noopIsmId := createNoopIsm(s, sender.Address)
		newOwner := "new_owner"

		// Act
		_, err := s.RunTx(&types.MsgSetMailbox{
			Owner:        creator.Address,
			MailboxId:    mailboxId,
			DefaultIsm:   &noopIsmId,
			DefaultHook:  nil,
			RequiredHook: nil,
			NewOwner:     newOwner,
		})

		// Assert
		Expect(err).NotTo(HaveOccurred())

		mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(mailbox.DefaultIsm).To(Equal(noopIsmId))
		Expect(mailbox.DefaultHook).To(Equal(&defaultHook))
		Expect(mailbox.RequiredHook).To(Equal(&requiredHook))
		Expect(mailbox.Owner).To(Equal(newOwner))
	})

	It("SetMailbox (valid)", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		noopIsmId := createNoopIsm(s, sender.Address)
		defaultHookId := createIgp(s, creator.Address)
		requiredHookId := createIgp(s, creator.Address)
		newOwner := "new_owner"

		// Act
		_, err := s.RunTx(&types.MsgSetMailbox{
			Owner:        creator.Address,
			MailboxId:    mailboxId,
			DefaultIsm:   &noopIsmId,
			DefaultHook:  &defaultHookId,
			RequiredHook: &requiredHookId,
			NewOwner:     newOwner,
		})

		// Assert
		Expect(err).NotTo(HaveOccurred())

		mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
		Expect(err).To(BeNil())
		Expect(mailbox.DefaultIsm).To(Equal(noopIsmId))
		Expect(mailbox.DefaultHook).To(Equal(&defaultHookId))
		Expect(mailbox.RequiredHook).To(Equal(&requiredHookId))
		Expect(mailbox.Owner).To(Equal(newOwner))
	})

	// TODO: ProcessMessage (valid) (Multisig ISM)
})

// Utils
func createIgp(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateIgp{
		Owner: creator,
		Denom: "acoin",
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateIgpResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	igpId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return igpId
}

func createNoopHook(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&pdTypes.MsgCreateNoopHook{
		Owner: creator,
	})
	Expect(err).To(BeNil())

	var response pdTypes.MsgCreateNoopHookResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	noopHookId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	return noopHookId
}

func createValidMailbox(s *i.KeeperTestSuite, creator string, ism string, destinationDomain uint32) (util.HexAddress, util.HexAddress, util.HexAddress, util.HexAddress) {
	var ismId util.HexAddress
	switch ism {
	case "noop":
		ismId = createNoopIsm(s, creator)
	case "multisig":
		ismId = createMultisigIsm(s, creator)
	}

	igpId := createIgp(s, creator)
	noopId := createNoopHook(s, creator)

	err := setDestinationGasConfig(s, creator, igpId.String(), destinationDomain)
	Expect(err).To(BeNil())

	res, err := s.RunTx(&types.MsgCreateMailbox{
		Owner:        creator,
		LocalDomain:  1,
		DefaultIsm:   ismId,
		DefaultHook:  &noopId,
		RequiredHook: &igpId,
	})
	Expect(err).To(BeNil())

	return verifyNewSingleMailbox(s, res, creator, ismId.String(), igpId.String(), noopId.String()), igpId, noopId, ismId
}

func createMultisigIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismtypes.MsgCreateMerkleRootMultisigIsm{
		Creator: creator,
		Validators: []string{
			"0xa05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xb05b6a0aa112b61a7aa16c19cac27d970692995e",
			"0xd05b6a0aa112b61a7aa16c19cac27d970692995e",
		},
		Threshold: 2,
	})
	Expect(err).To(BeNil())

	var response ismtypes.MsgCreateMerkleRootMultisigIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func createNoopIsm(s *i.KeeperTestSuite, creator string) util.HexAddress {
	res, err := s.RunTx(&ismtypes.MsgCreateNoopIsm{
		Creator: creator,
	})
	Expect(err).To(BeNil())

	var response ismtypes.MsgCreateNoopIsmResponse
	err = proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())

	return response.Id
}

func setDestinationGasConfig(s *i.KeeperTestSuite, creator string, igpId string, domain uint32) error {
	_, err := s.RunTx(&pdTypes.MsgSetDestinationGasConfig{
		Owner: creator,
		IgpId: igpId,
		DestinationGasConfig: &pdTypes.DestinationGasConfig{
			RemoteDomain: 1,
			GasOracle: &pdTypes.GasOracle{
				TokenExchangeRate: math.NewInt(1e10),
				GasPrice:          math.NewInt(1),
			},
			GasOverhead: math.NewInt(200000),
		},
	})

	return err
}

func verifyNewSingleMailbox(s *i.KeeperTestSuite, res *sdk.Result, creator, ismId, requiredHookId, defaultHookId string) util.HexAddress {
	var response types.MsgCreateMailboxResponse
	err := proto.Unmarshal(res.MsgResponses[0].Value, &response)
	Expect(err).To(BeNil())
	mailboxId, err := util.DecodeHexAddress(response.Id)
	Expect(err).To(BeNil())

	mailbox, err := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(err).To(BeNil())
	Expect(mailbox.Owner).To(Equal(creator))
	Expect(mailbox.DefaultIsm.String()).To(Equal(ismId))
	if defaultHookId != "" {
		Expect(mailbox.DefaultHook.String()).To(Equal(defaultHookId))
	} else {
		Expect(mailbox.DefaultHook).To(BeNil())
	}
	if requiredHookId != "" {
		Expect(mailbox.RequiredHook.String()).To(Equal(requiredHookId))
	} else {
		Expect(mailbox.RequiredHook).To(BeNil())
	}
	Expect(mailbox.MessageSent).To(Equal(uint32(0)))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))

	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(1))
	Expect(mailboxes.Mailboxes[0].Owner).To(Equal(creator))

	return mailboxId
}

func verifyInvalidMailboxCreation(s *i.KeeperTestSuite) {
	mailboxes, err := keeper.NewQueryServerImpl(s.App().HyperlaneKeeper).Mailboxes(s.Ctx(), &types.QueryMailboxesRequest{})
	Expect(err).To(BeNil())
	Expect(mailboxes.Mailboxes).To(HaveLen(0))
}

func verifyDispatch(s *i.KeeperTestSuite, mailboxId util.HexAddress, messageSent uint32) {
	mailbox, _ := s.App().HyperlaneKeeper.Mailboxes.Get(s.Ctx(), mailboxId.GetInternalId())
	Expect(mailbox.MessageSent).To(Equal(messageSent))
	Expect(mailbox.MessageReceived).To(Equal(uint32(0)))
}
