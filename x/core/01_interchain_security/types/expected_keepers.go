package types

import (
	"context"

	"github.com/troykessler/hyperlane-cosmos/util"
)

type CoreKeeper interface {
	LocalDomain(ctx context.Context, mailboxId util.HexAddress) (uint32, error)
	MailboxIdExists(ctx context.Context, mailboxId util.HexAddress) (bool, error)
	IsmRouter() *util.Router[util.InterchainSecurityModule]
}
