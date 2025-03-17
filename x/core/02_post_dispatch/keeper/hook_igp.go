package keeper

import (
	"context"
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/troykessler/hyperlane-cosmos/util"
	"github.com/troykessler/hyperlane-cosmos/x/core/02_post_dispatch/types"
)

type InterchainGasPaymasterHookHandler struct {
	k Keeper
}

var _ util.PostDispatchModule = InterchainGasPaymasterHookHandler{}

func (i InterchainGasPaymasterHookHandler) HookType() uint8 {
	return types.POST_DISPATCH_HOOK_TYPE_INTERCHAIN_GAS_PAYMASTER
}

func (i InterchainGasPaymasterHookHandler) PostDispatch(ctx context.Context, _, hookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage, maxFee sdk.Coins) (sdk.Coins, error) {
	return i.PayForGas(ctx, hookId, metadata.Address.String(), message.Id().String(), message.Destination, metadata.GasLimit, maxFee)
}

// QuoteDispatch returns the required Interchain Gas Payment for a certain message.
func (i InterchainGasPaymasterHookHandler) QuoteDispatch(ctx context.Context, _, hookId util.HexAddress, metadata util.StandardHookMetadata, message util.HyperlaneMessage) (sdk.Coins, error) {
	return i.QuoteGasPayment(ctx, hookId, message.Destination, metadata.GasLimit)
}

func (i InterchainGasPaymasterHookHandler) Exists(ctx context.Context, hookId util.HexAddress) (bool, error) {
	has, err := i.k.Igps.Has(ctx, hookId.GetInternalId())
	if err != nil {
		return false, err
	}
	return has, nil
}

// PayForGas executes an InterchainGasPayment using `QuoteGasPayment` beforehand and returns the charged fees.
func (i InterchainGasPaymasterHookHandler) PayForGas(ctx context.Context, hookId util.HexAddress, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, maxFee sdk.Coins) (sdk.Coins, error) {
	if maxFee.Empty() {
		return sdk.NewCoins(), fmt.Errorf("maxFee is required")
	}

	requiredPayment, err := i.QuoteGasPayment(ctx, hookId, destinationDomain, gasLimit)
	if err != nil {
		return sdk.NewCoins(), err
	}

	if requiredPayment.IsAllGT(maxFee) {
		return sdk.NewCoins(), fmt.Errorf("required payment exceeds max hyperlane fee: %v", requiredPayment)
	}

	return requiredPayment, i.PayForGasWithoutQuote(ctx, hookId, sender, messageId, destinationDomain, gasLimit, requiredPayment)
}

// PayForGasWithoutQuote executes an InterchainGasPayment without using `QuoteGasPayment`.
// This is used in the `MsgPayForGas` transaction, as the main purpose is paying an exact
// amount for e.g. re-funding a certain message-id as the first payment wasn't enough.
func (i InterchainGasPaymasterHookHandler) PayForGasWithoutQuote(ctx context.Context, hookId util.HexAddress, sender string, messageId string, destinationDomain uint32, gasLimit math.Int, amount sdk.Coins) error {
	igp, err := i.k.Igps.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return fmt.Errorf("igp does not exist: %s", hookId.String())
	}

	if amount.IsZero() {
		return fmt.Errorf("amount must be greater than zero")
	}

	if messageId == "" {
		return fmt.Errorf("message id cannot be empty")
	}

	senderAcc, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		return err
	}

	// TODO use core-types module name or create sub-account
	err = i.k.bankKeeper.SendCoinsFromAccountToModule(ctx, senderAcc, "hyperlane", amount)
	if err != nil {
		return err
	}

	igp.ClaimableFees = igp.ClaimableFees.Add(amount...)

	err = i.k.Igps.Set(ctx, igp.Id.GetInternalId(), igp)
	if err != nil {
		return err
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	_ = sdkCtx.EventManager().EmitTypedEvent(&types.GasPayment{
		MessageId:   messageId,
		Destination: destinationDomain,
		GasAmount:   gasLimit.String(),
		Payment:     amount.String(),
		IgpId:       hookId.String(),
	})

	return nil
}

// QuoteGasPayment uses the IGP's  DestinationGasConfig to determine the required payment.
func (i InterchainGasPaymasterHookHandler) QuoteGasPayment(ctx context.Context, hookId util.HexAddress, destinationDomain uint32, gasLimit math.Int) (sdk.Coins, error) {
	igp, err := i.k.Igps.Get(ctx, hookId.GetInternalId())
	if err != nil {
		return sdk.NewCoins(), fmt.Errorf("igp does not exist: %s", hookId.String())
	}

	destinationGasConfig, err := i.k.IgpDestinationGasConfigs.Get(ctx, collections.Join(igp.Id.GetInternalId(), destinationDomain))
	if err != nil {
		return sdk.NewCoins(), fmt.Errorf("remote domain %v is not supported", destinationDomain)
	}

	gasLimit = gasLimit.Add(destinationGasConfig.GasOverhead)

	destinationCost := gasLimit.Mul(destinationGasConfig.GasOracle.GasPrice)

	amount := (destinationCost.Mul(destinationGasConfig.GasOracle.TokenExchangeRate)).Quo(types.TokenExchangeRateScale)

	coin := sdk.Coin{
		Denom:  igp.Denom,
		Amount: amount,
	}

	if err = coin.Validate(); err != nil {
		return sdk.NewCoins(), err
	}

	return sdk.NewCoins(coin), nil
}
