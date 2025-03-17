package keeper_test

import (
	"fmt"
	"testing"

	"github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPostDispatchKeeper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, fmt.Sprintf("x/%s Keeper Test Suite", types.SubModuleName))
}
