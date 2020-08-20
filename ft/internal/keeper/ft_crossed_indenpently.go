/*
 * Copyright (C) 2020 The poly network Authors
 * This file is part of The poly network library.
 *
 * The  poly network  is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * The  poly network  is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 * You should have received a copy of the GNU Lesser General Public License
 * along with The poly network .  If not, see <http://www.gnu.org/licenses/>.
 */

package keeper

import (
	"bytes"
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
	polycommon "github.com/polynetwork/poly/common"
	"strconv"
)

func (k Keeper) CreateDenom(ctx sdk.Context, creator sdk.AccAddress, denom string) error {
	if reason, exist := k.ccmKeeper.ExistDenom(ctx, denom); exist {
		return types.ErrCreateDenom(fmt.Sprintf("denom: %s already exist, due to reason: %s", denom, reason))
	}
	//k.SetOperator(ctx, denom, creator)
	k.ccmKeeper.SetDenomCreator(ctx, denom, creator)
	ctx.KVStore(k.storeKey).Set(GetIndependentCrossDenomKey(denom), []byte(denom))
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCoins,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, denom),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
		),
	})
	k.Logger(ctx).Info(fmt.Sprintf("creator: %s initialized denom: %s ", creator.String(), denom))
	return nil
}

func (k Keeper) BindAssetHash(ctx sdk.Context, creator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) error {
	if !k.ValidCreator(ctx, sourceAssetDenom, creator) {
		return types.ErrBindAssetHash(fmt.Sprintf("creator is not valid, expect: %s, got: %s", k.ccmKeeper.GetDenomCreator(ctx, sourceAssetDenom).String(), creator.String()))
	}

	store := ctx.KVStore(k.storeKey)
	if !bytes.Equal([]byte(sourceAssetDenom), store.Get(GetIndependentCrossDenomKey(sourceAssetDenom))) {
		return types.ErrBindAssetHash(fmt.Sprintf("denom: %s is not designed to be able to be bondAssetHash through this interface", sourceAssetDenom))

	}
	store.Set(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId), toAssetHash)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString(sdk.AccAddress(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAssetHash)),
		),
	})
	return nil
}

func (k Keeper) Lock(ctx sdk.Context, fromAddr sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddr []byte, amount sdk.Int) error {
	sink := polycommon.NewZeroCopySink(nil)
	args := types.TxArgs{
		ToAddress: toAddr,
		Amount:    amount.BigInt(),
	}
	if err := args.Serialization(sink, 32); err != nil {
		return types.ErrLock(fmt.Sprintf("TxArgs Serialization error: %s", err.Error()))
	}

	// get toAssetHash from storage
	toAssetHash := ctx.KVStore(k.storeKey).Get(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId))
	// ensure toAssetHash is not empty
	if toAssetHash == nil {
		return types.ErrLock(fmt.Sprintf("toAssetHash is empty for Denom: %s, toChainId: %d", sourceAssetDenom, toChainId))
	}
	// invoke cross_chain_manager module to construct cosmos proof
	if err := k.ccmKeeper.CreateCrossChainTx(ctx, fromAddr, toChainId, []byte(sourceAssetDenom), toAssetHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrLock(fmt.Sprintf("ccmKeeper.CreateCrossChainTx, toChainId: %d, denom: %s, toAssetHash: %x, args: %x, Error: %s", toChainId, sourceAssetDenom, toAssetHash, args, err.Error()))
	}

	// burn coins from fromAddr
	if err := k.BurnCoins(ctx, fromAddr, sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, amount))); err != nil {
		return types.ErrLock(fmt.Sprintf("ft_crossed_independently.Lock.BurnCoins Error: %s", err.Error()))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddr)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error {

	var args types.TxArgs
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs), 32); err != nil {
		return types.ErrUnLock(fmt.Sprintf("Deserialize args: %x,  Error: %s", argsBs, err.Error()))
	}

	store := ctx.KVStore(k.storeKey)
	denom := string(toContractAddr)
	storedFromAssetHash := store.Get(GetBindAssetHashKey([]byte(denom), fromChainId))
	if !bytes.Equal(fromContractAddr, storedFromAssetHash) {
		return types.ErrUnLock(fmt.Sprintf("fromContractaddr: %x is not the stored assetHash: %x for denom: %s and fromChainId: %d", fromContractAddr, storedFromAssetHash, denom, fromChainId))
	}

	toAccAddr := sdk.AccAddress(args.ToAddress)
	amount := sdk.NewIntFromBigInt(args.Amount)
	if err := k.MintCoins(ctx, toAccAddr, sdk.NewCoins(sdk.NewCoin(denom, amount))); err != nil {
		return types.ErrUnLock(fmt.Sprintf("ft_crossed_independently.Unlock.MintCoins, toAddress: %s, denom: %s, amount: %s, Error: %s", toAccAddr.String(), denom, amount.String(), err.Error()))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString([]byte(denom))),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAccAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denom string) *types.DenomInfo {
	operator := k.ccmKeeper.GetDenomCreator(ctx, denom)
	if len(operator) == 0 {
		return nil
	}
	return &types.DenomInfo{
		Creator:     operator.String(),
		Denom:       denom,
		AssetHash:   hex.EncodeToString([]byte(denom)),
		TotalSupply: k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom),
	}
}

func (k Keeper) GetDenomCrossChainInfo(ctx sdk.Context, denom string, toChainId uint64) *types.DenomCrossChainInfo {
	return &types.DenomCrossChainInfo{
		DenomInfo:   *k.GetDenomInfo(ctx, denom),
		ToChainId:   toChainId,
		ToAssetHash: hex.EncodeToString(ctx.KVStore(k.storeKey).Get(GetBindAssetHashKey([]byte(denom), toChainId))),
	}
}

func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetBindAssetHashKey(toContractAddr, fromChainId))) != nil
}

func (k Keeper) ValidCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) bool {
	//store := ctx.KVStore(k.storeKey)
	//return bytes.Equal(store.Get(GetDenomToOperatorKey(denom)), creator.Bytes())
	return bytes.Equal(k.ccmKeeper.GetDenomCreator(ctx, denom), creator.Bytes())
}
