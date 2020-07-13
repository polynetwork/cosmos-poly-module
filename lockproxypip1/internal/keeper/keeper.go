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
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
	selfexported "github.com/polynetwork/cosmos-poly-module/lockproxypip1/exported"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CrossChainManager
	selfexported.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, bk types.BankKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CrossChainManager) Keeper {

	// ensure mint module account is set
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
	return ctx.KVStore(k.storeKey).Get((GetBindChainIdKey(toContractAddr, fromChainId))) != nil
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
			sdk.NewAttribute(types.AttributeKeyProxyHash, hex.EncodeToString(creator.Bytes())),
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

func (k Keeper) UpdateRegistry(ctx sdk.Context, lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) error {
	if k.AssetIsRegistered(ctx, lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash) {
		return types.ErrRegistryAlreadyExists(fmt.Sprintf("asset already registered %x, %d, %x, %x", assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash))
	}

	store := ctx.KVStore(k.storeKey)
	registryKey := GetRegistryKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	store.Set(registryKey, []byte("1"))

	bindChainIdKey := GetBindChainIdKey(lockProxyHash, nativeChainId)
	if store.Get(bindChainIdKey) == nil {
		store.Set(bindChainIdKey, []byte("1"))
	}

	return nil
}

func (k Keeper) GetBalance(ctx sdk.Context, balanceKey []byte) sdk.Int {
	store := ctx.KVStore(k.storeKey)
	currentAmount := sdk.ZeroInt()
	currentAmountBz := store.Get(balanceKey)
	if currentAmountBz != nil {
		err := k.cdc.UnmarshalBinaryLengthPrefixed(currentAmountBz, &currentAmount)
		if err != nil {
			panic(err)
		}
	}

	return currentAmount
}

func (k Keeper) StoreBalance(ctx sdk.Context, balanceKey []byte, newAmount sdk.Int) {
	store := ctx.KVStore(k.storeKey)
	newAmountBz, err := k.cdc.MarshalBinaryLengthPrefixed(newAmount)
	if err != nil {
		panic(err)
	}
	store.Set(balanceKey, newAmountBz)
}

func (k Keeper) IncreaseBalance(ctx sdk.Context, lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte, amount sdk.Int) {
	balanceKey := GetBalanceKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	currentAmount := k.GetBalance(ctx, balanceKey)
	newAmount := currentAmount.Add(amount)
	k.StoreBalance(ctx, balanceKey, newAmount)
}

func (k Keeper) DecreaseBalance(ctx sdk.Context, lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte, amount sdk.Int) error {
	balanceKey := GetBalanceKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	currentAmount := k.GetBalance(ctx, balanceKey)
	newAmount := currentAmount.Sub(amount)
	if newAmount.LT(sdk.ZeroInt()) {
		return types.ErrBalance(fmt.Sprintf("insufficient balance, current balance: %s, decrement balance: %s", currentAmount.String(), amount.String()))
	}
	k.StoreBalance(ctx, balanceKey, newAmount)
	return nil
}

func (k Keeper) AssetIsRegistered(ctx sdk.Context, lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetRegistryKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	registryBytes := store.Get(key)
	return len(registryBytes) != 0
}

func (k Keeper) RegisterAsset(ctx sdk.Context, fromChainId uint64, fromContractAddr []byte, toContractAddr []byte, argsBs []byte) error {
	if exist := k.EnsureLockProxyExist(ctx, toContractAddr); !exist {
		return types.ErrRegisterAsset(fmt.Sprintf("lockproxy with hash: %s not created", toContractAddr))
	}

	args := new(types.RegisterAssetTxArgs)
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs)); err != nil {
		return types.ErrUnLock(fmt.Sprintf("unlock, Deserialization args error:%s", err))
	}
	assetHash := args.AssetHash
	// check if denom exists
	nativeAssetHash := args.NativeAssetHash

	if err := k.UpdateRegistry(ctx, toContractAddr, assetHash, fromChainId, fromContractAddr, nativeAssetHash); err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRegisterAsset,
			sdk.NewAttribute(types.AttributeKeyFromChainId, fmt.Sprintf("%d", fromChainId)),
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString(fromContractAddr)),
			sdk.NewAttribute(types.AttributeKeyToContractHash, string(toContractAddr)),
			sdk.NewAttribute(types.AttributeKeyAssetHash, hex.EncodeToString(assetHash)),
			sdk.NewAttribute(types.AttributeKeyNativeAssetHash, string(nativeAssetHash)),
		),
	})
	return nil
}

func (k Keeper) CreateCoinAndDelegateToProxy(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, lockproxyHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) error {
	if reason, exist := k.ccmKeeper.ExistDenom(ctx, coin.Denom); exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("denom:%s already exist, reason: %s", coin.Denom, reason))
	}
	if exist := k.EnsureLockProxyExist(ctx, lockproxyHash); !exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("lockproxy with hash: %s not created", lockproxyHash))
	}

	k.ccmKeeper.SetDenomCreator(ctx, coin.Denom, creator)

	if err := k.UpdateRegistry(ctx, lockproxyHash, []byte(coin.Denom), nativeChainId, nativeLockProxyHash, nativeAssetHash); err != nil {
		return err
	}

	if err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(coin)); err != nil {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("supplyKeeper.MintCoins Error: %s", err.Error()))
	}

	k.IncreaseBalance(ctx, lockproxyHash, []byte(coin.Denom), nativeChainId, nativeLockProxyHash, nativeAssetHash, coin.Amount)

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

func (k Keeper) Lock(ctx sdk.Context, lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toChainProxyHash []byte, toChainAssetHash []byte, toAddressBs []byte, value sdk.Int, deductFeeInLock bool, feeAmount sdk.Int, feeAddress []byte) error {
	if exist := k.EnsureLockProxyExist(ctx, lockProxyHash); !exist {
		return types.ErrLock(fmt.Sprintf("lockproxy with hash: %s not created", lockProxyHash))
	}

	args := types.TxArgs{
		FromAssetHash: []byte(sourceAssetDenom),
		ToAssetHash:   toChainAssetHash,
		ToAddress:     toAddressBs,
		Amount:        value.BigInt(),
		FeeAmount:     feeAmount.BigInt(),
		FeeAddress:    feeAddress,
	}

	afterFeeAmount := value
	if deductFeeInLock && feeAmount.GT(sdk.ZeroInt()) {
		feeAddressAcc := sdk.AccAddress(args.FeeAddress)
		if feeAddressAcc.Empty() {
			return types.ErrLock("FeeAmount is present but FeeAddress is empty")
		}

		if feeAmount.GT(value) {
			return types.ErrLock(fmt.Sprintf("feeAmount %s is greater than value %s", feeAmount.String(), value.String()))
		}

		afterFeeAmount = value.Sub(feeAmount)
		feeCoins := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, feeAmount))
		k.bankKeeper.SendCoins(ctx, fromAddress, feeAddress, feeCoins)

		args.Amount = afterFeeAmount.BigInt()
		args.FeeAmount = big.NewInt(0)
	}

	// send coin of sourceAssetDenom from fromAddress to module account address
	amountCoins := sdk.NewCoins(sdk.NewCoin(sourceAssetDenom, afterFeeAmount))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, fromAddress, types.ModuleName, amountCoins); err != nil {
		return types.ErrLock(fmt.Sprintf("supplyKeeper.SendCoinsFromAccountToModule Error: from: %s, moduleAccount: %s of moduleName: %s, amount: %s", fromAddress.String(), k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName).GetAddress(), types.ModuleName, amountCoins.String()))
	}

	// get target asset hash from storage
	sink := polycommon.NewZeroCopySink(nil)
	if err := args.Serialization(sink, 32); err != nil {
		return types.ErrLock(fmt.Sprintf("TxArgs Serialization Error:%v", err))
	}
	fromContractHash := lockProxyHash
	if err := k.ccmKeeper.CreateCrossChainTx(ctx, fromAddress, toChainId, fromContractHash, toChainProxyHash, "unlock", sink.Bytes()); err != nil {
		return types.ErrLock(fmt.Sprintf("ccmKeeper.CreateCrossChainTx Error: toChainId: %d, fromContractHash: %x, toChainProxyHash: %x, args: %x", toChainId, fromContractHash, toChainProxyHash, args))
	}
	if amountCoins.AmountOf(sourceAssetDenom).IsNegative() {
		return types.ErrLock(fmt.Sprintf("the coin being crossed has negative amount value, coin:%s", amountCoins.String()))
	}

	if !k.AssetIsRegistered(ctx, lockProxyHash, []byte(sourceAssetDenom), toChainId, toChainProxyHash, toChainAssetHash) {
		return types.ErrLock(fmt.Sprintf("missing asset registry: lockProxyHash: %s, denom: %s, toChainId: %d, toChainProxyHash: %s, toChainAssetHash: %s", string(lockProxyHash), sourceAssetDenom, toChainId, hex.EncodeToString(toChainProxyHash), hex.EncodeToString(toChainAssetHash)))
	}

	k.IncreaseBalance(ctx, lockProxyHash, []byte(sourceAssetDenom), toChainId, toChainProxyHash, toChainAssetHash, afterFeeAmount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeLock,
			sdk.NewAttribute(types.AttributeKeyFromContractHash, hex.EncodeToString([]byte(sourceAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyToChainProxyHash, hex.EncodeToString(toChainProxyHash)),
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString(toChainAssetHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyToAddress, hex.EncodeToString(toAddressBs)),
			sdk.NewAttribute(types.AttributeKeyAmount, value.String()),
			sdk.NewAttribute(types.AttributeKeyLockProxy, hex.EncodeToString(fromContractHash)),
		),
	})

	return nil
}

func (k Keeper) Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error {
	args := new(types.TxArgs)
	if err := args.Deserialization(polycommon.NewZeroCopySource(argsBs), 32); err != nil {
		return types.ErrUnLock(fmt.Sprintf("unlock, Deserialization args error:%s", err))
	}
	fromAssetHash := args.FromAssetHash
	toAssetHash := args.ToAssetHash
	toAddress := args.ToAddress
	amount := sdk.NewIntFromBigInt(args.Amount)
	feeAmount := sdk.NewIntFromBigInt(args.FeeAmount)

	if !k.AssetIsRegistered(ctx, toContractAddr, toAssetHash, fromChainId, fromContractAddr, fromAssetHash) {
		return types.ErrUnLock(fmt.Sprintf("missing asset registry: toContractAddr: %s, toAssetHash: %s, fromChainId: %d, fromContractAddr: %s, fromAssetHash: %s", string(toContractAddr), toAssetHash, fromChainId, hex.EncodeToString(fromContractAddr), hex.EncodeToString(fromAssetHash)))
	}

	// to asset hash should be the hex format string of source asset denom name, NOT Module account address
	toAssetDenom := string(toAssetHash)

	toAcctAddress := make(sdk.AccAddress, len(toAddress))
	copy(toAcctAddress, toAddress)

	if err := k.EnsureAccountExist(ctx, toAddress); err != nil {
		return err
	}

	afterFeeAmount := amount
	if feeAmount.GT(sdk.ZeroInt()) {
		if feeAmount.GT(amount) {
			return types.ErrUnLock(fmt.Sprintf("feeAmount %s is greater than amount %s", feeAmount.String(), amount.String()))
		}

		feeAddressAcc := sdk.AccAddress(args.FeeAddress)
		if feeAddressAcc.Empty() {
			return types.ErrUnLock("FeeAmount is present but FeeAddress is empty")
		}

		afterFeeAmount = afterFeeAmount.Sub(feeAmount)
		feeCoins := sdk.NewCoins(sdk.NewCoin(toAssetDenom, feeAmount))
		if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, feeAddressAcc, feeCoins); err != nil {
			return types.ErrUnLock(fmt.Sprintf("supplyKeeper.SendCoinsFromModuleToAccount, Error: send coins:%s from Module account:%s to receiver account:%s error", feeCoins.String(), k.GetModuleAccount(ctx).GetAddress().String(), feeAddressAcc.String()))
		}

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeUnlock,
				sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString([]byte(toAssetDenom))),
				sdk.NewAttribute(types.AttributeKeyToAddress, feeAddressAcc.String()),
				sdk.NewAttribute(types.AttributeKeyAmount, feeAmount.String()),
			),
		})
	}

	// mint coin of sourceAssetDenom
	amountCoins := sdk.NewCoins(sdk.NewCoin(toAssetDenom, afterFeeAmount))

	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, toAcctAddress, amountCoins); err != nil {
		return types.ErrUnLock(fmt.Sprintf("supplyKeeper.SendCoinsFromModuleToAccount, Error: send coins:%s from Module account:%s to receiver account:%s error", amountCoins.String(), k.GetModuleAccount(ctx).GetAddress().String(), toAcctAddress.String()))
	}

	err := k.DecreaseBalance(ctx, toContractAddr, toAssetHash, fromChainId, fromContractAddr, fromAssetHash, amount)
	if err != nil {
		return err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeUnlock,
			sdk.NewAttribute(types.AttributeKeyToChainAssetHash, hex.EncodeToString([]byte(toAssetDenom))),
			sdk.NewAttribute(types.AttributeKeyToAddress, toAcctAddress.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
		),
	})
	return nil
}
