package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"cosmossdk.io/math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	i "github.com/troykessler/hyperlane-cosmos/tests/integration"
	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/types"
)

/*

TEST CASES - logic_message.go

* DispatchMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with non-existing Mailbox ID
* ProcessMessage (invalid) with invalid hex message
* ProcessMessage (invalid) already processed message (replay protection)
* ProcessMessage (invalid) with invalid message: non-registered recipient

*/

var _ = Describe("logic_message.go", Ordered, func() {
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

	It("DispatchMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		hexSender, _ := util.DecodeHexAddress(sender.Address)
		recipient, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		body, _ := hex.DecodeString("0x6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b")

		_, err = s.App().HyperlaneKeeper.DispatchMessage(
			s.Ctx(),
			nonExistingMailboxId,
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
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))

		verifyDispatch(s, mailboxId, 0)
	})

	It("ProcessMessage (invalid) with non-existing Mailbox ID", func() {
		// Arrange
		nonExistingMailboxId, _ := util.DecodeHexAddress("0xd7194459d45619d04a5a0f9e78dc9594a0f37fd6da8382fe12ddda6f2f46d647")
		createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      1337,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: nonExistingMailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to find mailbox with id: %s", nonExistingMailboxId)))
	})

	It("ProcessMessage (invalid) with wrong destination domain", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      1337,
			Sender:      senderHex,
			Destination: 2,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("message destination %v does not match local domain %v", 2, 1)))
	})

	It("ProcessMessage (invalid) with invalid hex message", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String()[:util.BodyOffset-1],
		})

		// Assert
		Expect(err.Error()).To(Equal("invalid hyperlane message"))
	})

	PIt("ProcessMessage (invalid) already processed message (replay protection)", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte(""),
		}

		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})
		Expect(err).To(BeNil())

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("already received messsage with id %s", hypMsg.Id())))
	})

	// TODO rework test, once warp is refactored to use router
	PIt("ProcessMessage (invalid) with invalid message: non-registered recipient", func() {
		// Arrange
		mailboxId, _, _, _ := createValidMailbox(s, creator.Address, "noop", 1)

		err := s.MintBaseCoins(sender.Address, 1_000_000)
		Expect(err).To(BeNil())

		senderHex := util.CreateMockHexAddress("test", 0)
		recipientHex := util.CreateMockHexAddress("test", 0)

		localDomain, err := s.App().HyperlaneKeeper.LocalDomain(s.Ctx(), mailboxId)
		Expect(err).To(BeNil())

		hypMsg := util.HyperlaneMessage{
			Version:     3,
			Nonce:       0,
			Origin:      localDomain,
			Sender:      senderHex,
			Destination: 1,
			Recipient:   recipientHex,
			Body:        []byte("test123"),
		}

		// Act
		_, err = s.RunTx(&types.MsgProcessMessage{
			MailboxId: mailboxId,
			Relayer:   sender.Address,
			Metadata:  "",
			Message:   hypMsg.String(),
		})

		// Assert
		Expect(err.Error()).To(Equal(fmt.Sprintf("failed to get receiver ism address for recipient: %s", recipientHex)))
	})
})
