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
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	selfexported "github.com/polynetwork/cosmos-poly-module/lockproxy/exported"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
	polycommon "github.com/polynetwork/poly/common"
	"strconv"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	authKeeper   types.AccountKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CrossChainManager
	selfexported.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CrossChainManager) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return types.ErrAccountNotExist(fmt.Sprintf("account %s does not exist", addr.String()))
	}
	return nil
}

func (k Keeper) ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool {
	return ctx.KVStore(k.storeKey).Get((GetBindProxyKey(toContractAddr, fromChainId))) != nil
}

func (k Keeper) CreateLockProxy(ctx sdk.Context, creator sdk.AccAddress) error {
	if k.EnsureLockProxyExist(ctx, creator) {
		return types.ErrCreateLockProxy(fmt.Sprintf("creator:%s already created lockproxy contract with hash:%x", creator.String(), creator.Bytes()))
	}
	store := ctx.KVStore(k.storeKey)
	store.Set(GetOperatorToLockProxyKey(creator), creator.Bytes())
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateLockProxy,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(creator.Bytes())),
		),
	})
	ctx.Logger().With("module", fmt.Sprintf("creator:%s initialized a lockproxy contract with hash: %x", creator.String(), creator.Bytes()))
	return nil
}

func (k Keeper) EnsureLockProxyExist(ctx sdk.Context, creator sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return bytes.Equal(store.Get(GetOperatorToLockProxyKey(creator)), creator)
}

func (k Keeper) GetLockProxyByOperator(ctx sdk.Context, operator sdk.AccAddress) []byte {
	store := ctx.KVStore(k.storeKey)
	proxyBytes := store.Get(GetOperatorToLockProxyKey(operator))
	if len(proxyBytes) == 0 || !bytes.Equal(operator.Bytes(), proxyBytes) {
		return nil
	}
	return proxyBytes
}

func (k Keeper) CreateCoinAndDelegateToProxy(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, lockproxyHash []byte) error {

	if reason, exist := k.ccmKeeper.ExistDenom(ctx, coin.Denom); exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("denom:%s already exist, reason: %s", coin.Denom, reason))
	}
	if exist := k.EnsureLockProxyExist(ctx, lockproxyHash); !exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("lockproxy with hash: %s not created", lockproxyHash))

	}
	k.ccmKeeper.SetDenomCreator(ctx, coin.Denom, creator)

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("supplyKeeper.MintCoins Error: %s", err.Error()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAndDelegateCoinToProxy,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, coin.Denom),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, coin.Amount.String()),
		),
	})
	return nil
}

func (k Keeper) BindProxyHash(ctx sdk.Context, operator sdk.AccAddress, toChainId uint64, toProxyHash []byte) error {
	if !k.EnsureLockProxyExist(ctx, operator) {
		return types.ErrBindProxyHash(fmt.Sprintf("operator:%s have NOT created lockproxy contract: %s", operator.String(), operator.Bytes()))
	}
	store := ctx.KVStore(k.storeKey)

	store.Set(GetBindProxyKey(operator, toChainId), toProxyHash)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindProxy,
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(operator.Bytes())),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(toProxyHash)),
		),
	})
	return nil
}

func (k Keeper) GetProxyHash(ctx sdk.Context, operator sdk.AccAddress, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindProxyKey(operator, toChainId))
}

func (k Keeper) BindAssetHash(ctx sdk.Context, operator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) error {
	// ensure the operator has created the lockproxy contract
	if !k.EnsureLockProxyExist(ctx, operator) {
		return types.ErrBindAssetHash(fmt.Sprintf("operator:%s have NOT created lockproxy contract: %s", operator.String(), operator.Bytes()))
	}
	// ensure the sourceAssetDenom has already been created with non-zero supply
	if _, exist := k.ccmKeeper.ExistDenom(ctx, sourceAssetDenom); !exist {
		return types.ErrBindAssetHash(fmt.Sprintf("sourceAssetDenom: %s not exist", sourceAssetDenom))
	}
	store := ctx.KVStore(k.storeKey)
	// store the to asset hash based on the lockproxy contract (operator) and sourceAssetHash + toChainId
	store.Set(GetBindAssetHashKey(operator, []byte(sourceAssetDenom), toChainId), toAssetHash)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeBindAsset,
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(operator.Bytes())),
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, sourceAssetDenom),
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString(toAssetHash)),
		),
	})
	return nil
}

func (k Keeper) GetAssetHash(ctx sdk.Context, lockProxyHash []byte, sourceAssetDenom string, toChainId uint64) []byte {
	store := ctx.KVStore(k.storeKey)
	return store.Get(GetBindAssetHashKey(lockProxyHash, []byte(sourceAssetDenom), toChainId))
}

func (k Keeper) Lock(ctx sdk.Context, lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddressBs []byte, value sdk.Int) error {
	// send coin of sourceAssetDenom from fromAddress to module account address
	amt := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, value))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, amt); err != nil {
		return types.ErrLock(fmt.Sprintf("supplyKeeper.SendCoinsFromAccountToModule Error: from: %s, moduleAccount: %s of moduleName: %s, amount: %s", fromAddress.String(), k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress(), types.ModuleName, amt.String()))
	}

	store := ctx.KVStore(k.storeKey)
	sourceAssetHash := []byte(sourceAssetDenom)
	toChainAssetHash := store.Get(GetBindAssetHashKey(lockProxyHash, sourceAssetHash, toChainId))
	if len(toChainAssetHash) == 0 {
		return types.ErrLock(fmt.Sprintf("toChainAssetHash not exist for lockproxyHash: %x / %s, denom: %s, toChainId: %d", lockProxyHash, sdk.AccAddress(lockProxyHash).String(), sourceAssetDenom, toChainId))
	}
	// get target asset hash from storage
	sink := polycommon.NewZeroCopySink(nil)
	args := types.TxArgs{
		ToAssetHash: toChainAssetHash,
		ToAddress:   toAddressBs,
		Amount:      value.BigInt(),
	}
	if err := args.Serialization(sink, 32); err != nil {
		return types.ErrLock(fmt.Sprintf("TxArgs Serialization Error:%v", err))
	}
	// get target chain proxy hash from storage
	toChainProxyHash := store.Get(GetBindProxyKey(lockProxyHash, toChainId))
	if toChainProxyHash == nil {
		return types.ErrLock(fmt.Sprintf("toChainProxyHash is empty"))
	}
	fromContractHash := lockProxyHash
	if err := k.ccmKeeper.CreateCrossChainTx(ctx, fromAddress, toChainId, fromContractHash, toChainProxyHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrLock(fmt.Sprintf("ccmKeeper.CreateCrossChainTx Error: toChainId: %d, fromContractHash: %x, toChainProxyHash: %x, args: %x", toChainId, fromContractHash, toChainProxyHash, args))
	}
	if amt.AmountOf(sourceAssetDenom).IsNegative() {
		return types.ErrLock(fmt.Sprintf("the coin being crossed has negative amount value, coin:%s", amt.String()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyFromAssetHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString(toChainAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddressBs)),
			sdk.NewAttribute(types.AttributeKeyAmount, value.String()),
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(fromContractHash)),
		),
	})
	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error {

	fromProxyHash := k.GetProxyHash(ctx, toContractAddr, fromChainId)
	if len(fromProxyHash) == 0 {
		return types.ErrUnLock(fmt.Sprintf("the proxyHash is empty with chainId: %d", fromChainId))
	}
	if !bytes.Equal(fromProxyHash, fromContractAddr) {
		return types.ErrUnLock(fmt.Sprintf("stored proxyHash correlated with lockproxyHash: %x and chainId: %d is not equal to fromContractAddress, expect:%x, got:%x", toContractAddr, fromChainId, fromProxyHash, fromContractAddr))
	}
	args := new(types.TxArgs)
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs), 32); err != nil {
		return types.ErrUnLock(fmt.Sprintf("unlock, Deserialization args error:%s", err))
	}
	toAssetHash := args.ToAssetHash
	toAddress := args.ToAddress
	amount := args.Amount

	// to asset hash should be the hex format string of source asset denom name, NOT Module account address
	toAssetDenom := string(toAssetHash)
	if len(k.GetAssetHash(ctx, toContractAddr, toAssetDenom, fromChainId)) == 0 {
		return types.ErrUnLock(fmt.Sprintf("toAssetHash: %x of denom: %s doesnot belong to the current lock proxy hash: %x", toAssetHash, toAssetDenom, toContractAddr))
	}

	// mint coin of sourceAssetDenom
	amt := sdk.NewCoins(sdk.NewCoin(toAssetDenom, sdk.NewIntFromBigInt(amount)))

	toAcctAddress := make(sdk.AccAddress, len(toAddress))
	copy(toAcctAddress, toAddress)

	if err := k.EnsureAccountExist(ctx, toAddress); err != nil {
		return err
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcctAddress, amt); err != nil {
		return types.ErrUnLock(fmt.Sprintf("supplyKeeper.SendCoinsFromModuleToAccount, Error: send coins:%s from Module account:%s to receiver account:%s error", amt.String(), k.GetModuleAccount(ctx).GetAddress().String(), toAcctAddress.String()))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyToAssetHash, hex.EncodeToString([]byte(toAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAcctAddress.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}
