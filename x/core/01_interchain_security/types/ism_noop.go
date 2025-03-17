package types

import (
	"context"

	"github.com/troykessler/hyperlane-cosmos/util"
)

var _ HyperlaneInterchainSecurityModule = &NoopISM{}

func (m *NoopISM) GetId() (util.HexAddress, error) {
	return m.Id, nil
}

func (m *NoopISM) ModuleType() uint8 {
	return INTERCHAIN_SECURITY_MODULE_TYPE_UNUSED
}

func (m *NoopISM) Verify(_ context.Context, _ []byte, _ util.HyperlaneMessage) (bool, error) {
	return true, nil
}
