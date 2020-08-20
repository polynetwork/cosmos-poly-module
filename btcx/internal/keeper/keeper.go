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
	"github.com/btcsuite/btcutil"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/btcx/exported"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/types"
	polycommon "github.com/polynetwork/poly/common"
	"github.com/tendermint/tendermint/libs/log"
	"math"
	"math/big"
	"strconv"
)

// Keeper of the btcx store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CCMKeeper
	exported.UnlockKeeper
}

// NewKeeper creates a new btcx Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CCMKeeper) Keeper {

	// ensure btcx module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		bankKeeper:   bk,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) CreateDenom(ctx sdk.Context, creator sdk.AccAddress, denom string, redeemScript string) error {
	if reason, exist := k.ccmKeeper.ExistDenom(ctx, denom); exist {
		return types.ErrCreateDenom(fmt.Sprintf("denom:%s already exist, due to reason:%s", denom, reason))
	}
	k.ccmKeeper.SetDenomCreator(ctx, denom, creator)

	redeemScriptBs, err := hex.DecodeString(redeemScript)
	if err != nil {
		return types.ErrInvalidRedeemScript(fmt.Sprintf("Invalid redeemScript :%s, Error: %s", redeemScript, err))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDenomToCreatorKey(denom), creator)

	scriptHashBs := btcutil.Hash160(redeemScriptBs)
	store.Set(GetCreatorDenomToScriptHashKey(creator, denom), scriptHashBs)
	store.Set(GetScriptHashToRedeemScript(scriptHashBs), redeemScriptBs)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateDenom,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, denom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString([]byte(denom))),
			sdk.NewAttribute(types.AttributeKeyRedeemScript, hex.EncodeToString(redeemScriptBs)),
		),
	})
	k.Logger(ctx).Info(fmt.Sprintf("creator:%s initialized denom: %s ", creator.String(), denom))
	return nil
}

func (k Keeper) BindAssetHash(ctx sdk.Context, creator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) error {
	if !k.ValidCreator(ctx, sourceAssetDenom, creator) {
		return types.ErrBindAssetHash(fmt.Sprintf("BindAssetHash, creator is not valid, expect:%s, got:%s", k.ccmKeeper.GetDenomCreator(ctx, sourceAssetDenom).String(), creator.String()))
	}
	store := ctx.KVStore(k.storeKey)
	if !bytes.Equal(creator.Bytes(), store.Get(GetDenomToCreatorKey(sourceAssetDenom))) {
		return types.ErrBindAssetHash(fmt.Sprintf("BindAssetHash, creator: %s created Denom: %s, yet not in %s module", creator.String(), sourceAssetDenom, types.ModuleName))
	}
	store.Set(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId), toAssetHash)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString(sdk.AccAddress(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString(toAssetHash)),
		),
	})
	return nil
}

func (k Keeper) Lock(ctx sdk.Context, fromAddr sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddr []byte, amount sdk.Int) error {
	// transfer back to btc
	store := ctx.KVStore(k.storeKey)
	toAssetHash := store.Get(GetBindAssetHashKey([]byte(sourceAssetDenom), toChainId))
	if toAssetHash == nil {
		return types.ErrLock(fmt.Sprintf("Invoke Lock of `btcx` module for denom: %s is illgeal due to toAssetHash empty for chainId: %d", sourceAssetDenom, toChainId))
	}
	if amount.GTE(sdk.NewIntFromUint64(math.MaxUint64)) {
		return types.ErrLock(fmt.Sprintf("Invoke Lock of `btcx` module, amount: %s too big than MaxUint64", amount.String()))
	}
	sink := polycommon.NewZeroCopySink(nil)
	// construct args bytes
	if toChainId == types.BtcChainId {
		creator := k.ccmKeeper.GetDenomCreator(ctx, sourceAssetDenom)
		if creator.Empty() {
			return types.ErrLock(fmt.Sprintf("Creator of denom: %s is Empty", sourceAssetDenom))
		}
		redeemScript := store.Get(GetScriptHashToRedeemScript(store.Get(GetCreatorDenomToScriptHashKey(creator, sourceAssetDenom))))
		if len(redeemScript) == 0 {
			return types.ErrLock(fmt.Sprintf("Invoke Lock of `btcx` module, denom: %s, toChainId: %d, redeemScript not exist", sourceAssetDenom, toChainId))
		}
		toBtcArgs := types.ToBTCArgs{
			ToBtcAddress: toAddr,
			Amount:       amount.BigInt().Uint64(),
			RedeemScript: redeemScript,
		}
		if err := toBtcArgs.Serialization(sink); err != nil {
			return types.ErrLock(fmt.Sprintf("ToBTCArgs Serialization Error:%v", err))
		}
	} else {
		args := types.BTCArgs{
			ToBtcAddress: toAddr,
			Amount:       amount.BigInt().Uint64(),
		}
		if err := args.Serialization(sink); err != nil {
			return types.ErrLock(fmt.Sprintf("BTCArgs Serialization Error:%v", err))
		}
	}

	// invoke cross_chain_manager module to construct cosmos proof
	if err := k.ccmKeeper.CreateCrossChainTx(ctx, fromAddr, toChainId, []byte(sourceAssetDenom), toAssetHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrLock(fmt.Sprintf("Lock, CreateCrossChainTx Error:%s", err.Error()))
	}
	// burn coins from fromAddr
	if err := k.BurnCoins(ctx, fromAddr, sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, amount))); err != nil {
		return types.ErrLock(fmt.Sprintf("BurnCoins failed, Error: %s", err.Error()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString(toAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddr)),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error {

	var args types.BTCArgs
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs)); err != nil {
		return types.ErrUnLock(fmt.Sprintf("Deserialize args Error: %s", err))
	}
	store := ctx.KVStore(k.storeKey)
	fromAssetHash := store.Get(GetBindAssetHashKey(toContractAddr, fromChainId))
	if len(fromAssetHash) == 0 {
		return types.ErrUnLock(fmt.Sprintf("fromAssetHash not exist for toContractAddr: %x, fromChainId: %d", toContractAddr, fromChainId))
	}
	if !bytes.Equal(fromContractAddr, fromAssetHash) {
		return types.ErrUnLock(fmt.Sprintf("fromContractaddr:%x is not the stored assetHash:%x", fromContractAddr, fromAssetHash))
	}
	toDenom := string(toContractAddr)

	toAccAddr := sdk.AccAddress(args.ToBtcAddress)
	amount := sdk.NewIntFromBigInt(big.NewInt(0).SetUint64(args.Amount))
	if err := k.MintCoins(ctx, toAccAddr, sdk.NewCoins(sdk.NewCoin(toDenom, amount))); err != nil {
		return types.ErrUnLock(fmt.Sprintf("MintCoins from Addr:%s, Error:%v", toAccAddr.String(), err))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString(toContractAddr)),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAccAddr.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}

func (k Keeper) GetDenomInfo(ctx sdk.Context, denom string) *types.DenomInfo {

	store := ctx.KVStore(k.storeKey)
	operator := k.ccmKeeper.GetDenomCreator(ctx, denom)
	if len(operator) == 0 {
		return nil
	}
	denomInfo := new(types.DenomInfo)
	denomInfo.Creator = operator.String()
	denomInfo.Denom = denom
	denomInfo.AssetHash = hex.EncodeToString([]byte(denom))
	denomInfo.TotalSupply = k.supplyKeeper.GetSupply(ctx).GetTotal().AmountOf(denom)
	redeemHash := store.Get(GetCreatorDenomToScriptHashKey(store.Get(GetDenomToCreatorKey(denom)), denom))
	denomInfo.RedeemScriptHash = hex.EncodeToString(redeemHash)
	denomInfo.RedeemScipt = hex.EncodeToString(store.Get(GetScriptHashToRedeemScript(redeemHash)))
	return denomInfo
}

func (k Keeper) GetDenomCrossChainInfo(ctx sdk.Context, denom string, toChainId uint64) *types.DenomCrossChainInfo {
	denomInfo := new(types.DenomCrossChainInfo)
	denomInfo.DenomInfo = *k.GetDenomInfo(ctx, denom)
	denomInfo.ToChainId = toChainId
	store := ctx.KVStore(k.storeKey)
	denomInfo.ToAssetHash = hex.EncodeToString(store.Get(GetBindAssetHashKey([]byte(denom), toChainId)))
	return denomInfo
}
func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetBindAssetHashKey(toContractAddr, fromChainId))) != nil
}

func (k Keeper) ValidCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) bool {
	return bytes.Equal(k.ccmKeeper.GetDenomCreator(ctx, denom), creator.Bytes())
}

func (k Keeper) MintCoins(ctx sdk.Context, toAcct sdk.AccAddress, amt sdk.Coins) error {
	_, err := k.bankKeeper.AddCoins(ctx, toAcct, amt)
	if err != nil {
		return types.ErrMintCoins(fmt.Sprintf("bankkeeper.AddCoins, Error: %s", err.Error()))
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Inflate(amt)

	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted coin:%s to account:%s ", amt.String(), toAcct.String()))

	return nil
}

func (k Keeper) BurnCoins(ctx sdk.Context, fromAcct sdk.AccAddress, amt sdk.Coins) error {

	_, err := k.bankKeeper.SubtractCoins(ctx, fromAcct, amt)
	if err != nil {
		return types.ErrBurnCoins(fmt.Sprintf("bankkeeper.SubtractCoins, Error: %s", err.Error()))
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Deflate(amt)
	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned coin:%s from account:%s ", amt.String(), fromAcct.String()))
	return nil
}
