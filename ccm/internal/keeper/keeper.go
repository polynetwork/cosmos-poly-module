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
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
	hs "github.com/polynetwork/cosmos-poly-module/headersync"
	polycommon "github.com/polynetwork/poly/common"
	polytype "github.com/polynetwork/poly/core/types"
	"github.com/polynetwork/poly/merkle"
	ccmc "github.com/polynetwork/poly/native/service/cross_chain_manager/common"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/libs/log"
	ttype "github.com/tendermint/tendermint/types"
	"strconv"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	paramSpace   params.Subspace
	hsKeeper     types.HeaderSyncKeeper
	supplyKeeper types.SupplyKeeper
	ulKeeperMap  map[string]types.UnlockKeeper
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramSpace params.Subspace, hsk types.HeaderSyncKeeper, supplyKeeper types.SupplyKeeper) Keeper {

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
		hsKeeper:     hsk,
		supplyKeeper: supplyKeeper,
		ulKeeperMap:  nil,
	}
}
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k *Keeper) MountUnlockKeeperMap(ulKeeperMap map[string]types.UnlockKeeper) {
	k.ulKeeperMap = make(map[string]types.UnlockKeeper)
	for key, value := range ulKeeperMap {
		k.ulKeeperMap[key] = value
	}
}

// GetParams returns the total set of ccm parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return params
}

// SetParams sets the total set of ccm parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

func (k Keeper) IfContainToContract(ctx sdk.Context, keystore string, toContractAddr []byte, fromChainId uint64) *types.QueryContainToContractRes {
	unlockKeeper, ok := k.ulKeeperMap[keystore]
	if !ok {
		return &types.QueryContainToContractRes{KeyStore: keystore, Info: "map doesnot contain current keystore"}
	}
	var res types.QueryContainToContractRes
	res.KeyStore = keystore
	res.Exist = unlockKeeper.ContainToContractAddr(ctx, toContractAddr, fromChainId)
	k.Logger(ctx).Info(fmt.Sprintf("key is %+v ", keystore))
	k.Logger(ctx).Info(fmt.Sprintf("IfContains %+v ", unlockKeeper.ContainToContractAddr(ctx, toContractAddr, fromChainId)))

	return &res
}

func (k Keeper) SetDenomCreator(ctx sdk.Context, denom string, creator sdk.AccAddress) {
	ctx.KVStore(k.storeKey).Set(GetDenomToCreatorKey(denom), creator.Bytes())
}

func (k Keeper) GetDenomCreator(ctx sdk.Context, denom string) sdk.AccAddress {
	return ctx.KVStore(k.storeKey).Get(GetDenomToCreatorKey(denom))
}

func (k Keeper) ExistDenom(ctx sdk.Context, denom string) (string, bool) {
	storedSupplyCoins := k.supplyKeeper.GetSupply(ctx).GetTotal()
	//return storedSupplyCoins.AmountOf(denom) != sdk.ZeroInt() || len(k.GetOperator(ctx, denom)) != 0
	if len(k.GetDenomCreator(ctx, denom)) != 0 {
		return fmt.Sprintf("ccmKeeper.GetDenomCreator(ctx,%s) is %s", denom, sdk.AccAddress(k.GetDenomCreator(ctx, denom)).String()), true
	}
	if !storedSupplyCoins.AmountOf(denom).Equal(sdk.ZeroInt()) {
		return fmt.Sprintf("supply.AmountOf(%s) is %s", denom, storedSupplyCoins.AmountOf(denom).String()), true
	}
	return "", false
}

func (k Keeper) GetModuleBalance(ctx sdk.Context, moduleName string) (sdk.Coins, error) {
	//	ensure the passed-in initialAmt is equal to the balance of lockproxy module account
	moduleAcct := k.supplyKeeper.GetModuleAccount(ctx, moduleName)
	if moduleAcct == nil {
		return nil, types.ErrGetModuleBalance(fmt.Sprintf("Module: %s account doesnot exist", moduleName))
	}
	return moduleAcct.GetCoins(), nil
}

func (k Keeper) CreateCrossChainTx(ctx sdk.Context, fromAddr sdk.AccAddress, toChainId uint64, fromContractHash, toContractHash []byte, method string, args []byte) error {
	crossChainId, err := k.getCrossChainId(ctx)
	if err != nil {
		return err
	}
	if err := k.setCrossChainId(ctx, crossChainId.Add(sdk.NewInt(1))); err != nil {
		return err
	}

	ttx := make([]byte, len(ctx.TxBytes()))
	copy(ttx, ctx.TxBytes())
	txHash := ttype.Tx(ttx).Hash()
	crossChainIdBs := crossChainId.BigInt().Bytes()
	txParam := ccmc.MakeTxParam{
		TxHash:              txHash,
		CrossChainID:        crossChainIdBs,
		FromContractAddress: fromContractHash,
		ToChainID:           toChainId,
		ToContractAddress:   toContractHash,
		Method:              method,
		Args:                args,
	}
	sink := polycommon.NewZeroCopySink(nil)
	txParam.Serialization(sink)

	store := ctx.KVStore(k.storeKey)

	txParamHash := tmhash.Sum(sink.Bytes())
	store.Set(GetCrossChainTxKey(txParamHash), sink.Bytes())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCrossChainTx,
			sdk.NewAttribute(types.AttributeKeyStatus, "1"),
			sdk.NewAttribute(types.AttributeCrossChainId, crossChainId.String()),
			sdk.NewAttribute(types.AttributeKeyTxParamHash, hex.EncodeToString(txParamHash)),
			sdk.NewAttribute(types.AttributeKeyFromAddress, fromAddr.String()),
			sdk.NewAttribute(types.AttributeKeyFromContract, hex.EncodeToString(fromContractHash)),
			sdk.NewAttribute(types.AttributeKeyToChainId, strconv.FormatUint(toChainId, 10)),
			sdk.NewAttribute(types.AttributeKeyMakeTxParam, hex.EncodeToString(sink.Bytes())),
		),
	})
	return nil
}

func (k Keeper) ProcessCrossChainTx(ctx sdk.Context, fromChainId uint64, proofStr string, headerStr, headerProofStr, curHeaderStr string) error {
	headerToBeVerified := new(polytype.Header)
	headerBs, err := hex.DecodeString(headerStr)
	if err != nil {
		return types.ErrProcessCrossChainTx(fmt.Sprintf("Decode proof hex string: %s to bytes, Error: %s ", headerStr, err.Error()))
	}
	if err := headerToBeVerified.Deserialization(polycommon.NewZeroCopySource(headerBs)); err != nil {
		return types.ErrProcessCrossChainTx(hs.ErrDeserializeHeader(err).Error())
	}

	headerInCurEpoch := new(polytype.Header)
	curHeaderBs, err := hex.DecodeString(curHeaderStr)
	if err != nil {
		headerInCurEpoch = nil
	} else {
		if err := headerInCurEpoch.Deserialization(polycommon.NewZeroCopySource(curHeaderBs)); err != nil {
			headerInCurEpoch = nil
		}
	}

	headerProof, err := hex.DecodeString(headerProofStr)
	if err != nil {
		headerProof = nil
	}

	if err := k.hsKeeper.ProcessHeader(ctx, headerToBeVerified, headerProof, headerInCurEpoch); err != nil {
		return types.ErrProcessCrossChainTx(fmt.Sprintf("ProcessHeader Error, %s", err.Error()))
	}

	proof, err := hex.DecodeString(proofStr)
	if err != nil {
		return types.ErrProcessCrossChainTx(fmt.Sprintf("Decode proof hex string: %s to bytes, Error: %s", proofStr, err.Error()))
	}

	merkleValue, err := k.VerifyToCosmosTx(ctx, proof, headerToBeVerified)
	if err != nil {
		return types.ErrProcessCrossChainTx(fmt.Sprintf("VerifyToCosmostx failed, %s", err.Error()))
	}
	currentChainCrossChainId := k.GetParams(ctx).ChainIdInPolyNet
	if merkleValue.MakeTxParam.ToChainID != currentChainCrossChainId {
		return types.ErrProcessCrossChainTx(fmt.Sprintf("toChainId is not for this chain, expect: %d, got: %d", currentChainCrossChainId, merkleValue.MakeTxParam.ToChainID))
	}
	// check if tocontractAddress is lockproxy module account, if yes, invoke lockproxy.unlock(), otherwise, invoke btcx.unlock
	for key, unlockKeeper := range k.ulKeeperMap {
		k.Logger(ctx).Info(fmt.Sprintf("key is %+v ", key))
		k.Logger(ctx).Info(fmt.Sprintf("IfContains %+v ", unlockKeeper.ContainToContractAddr(ctx, merkleValue.MakeTxParam.ToContractAddress, fromChainId)))

		if unlockKeeper.ContainToContractAddr(ctx, merkleValue.MakeTxParam.ToContractAddress, merkleValue.FromChainID) {
			if err := unlockKeeper.Unlock(ctx, merkleValue.FromChainID, merkleValue.MakeTxParam.FromContractAddress, merkleValue.MakeTxParam.ToContractAddress, merkleValue.MakeTxParam.Args); err != nil {
				return types.ErrProcessCrossChainTx(fmt.Sprintf("Unlock failed, for module: %s, Error: %s", key, err.Error()))
			}
			return nil
		}
	}

	return types.ErrProcessCrossChainTx(fmt.Sprintf("Cannot find any unlock keeper to perform 'unlock' method for toContractAddr:%x, fromChainId:%d", merkleValue.MakeTxParam.ToContractAddress, fromChainId))
}

func (k Keeper) VerifyToCosmosTx(ctx sdk.Context, proof []byte, header *polytype.Header) (*ccmc.ToMerkleValue, error) {
	value, err := merkle.MerkleProve(proof, header.CrossStateRoot[:])
	if err != nil {
		return nil, types.ErrVerifyToCosmosTx(fmt.Sprintf("merkle.MerkleProve verify failed, Error: %s", err.Error()))
	}

	merkleValue := new(ccmc.ToMerkleValue)
	if err := merkleValue.Deserialization(polycommon.NewZeroCopySource(value)); err != nil {
		return nil, types.ErrVerifyToCosmosTx(fmt.Sprintf("ToMerkeValue Deserialization Error: %s", err.Error()))
	}

	if err := k.checkDoneTx(ctx, merkleValue.FromChainID, merkleValue.MakeTxParam.CrossChainID); err != nil {
		return nil, types.ErrVerifyToCosmosTx(fmt.Sprintf("check if this tx has been done, Error: %s", err.Error()))
	}

	k.putDoneTx(ctx, merkleValue.FromChainID, merkleValue.MakeTxParam.CrossChainID)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeVerifyToCosmosProof,
			sdk.NewAttribute(types.AttributeKeyMerkleValueTxHash, hex.EncodeToString(merkleValue.TxHash)),
			sdk.NewAttribute(types.AttributeKeyMerkleValueMakeTxParamTxHash, hex.EncodeToString(merkleValue.MakeTxParam.TxHash)),
			sdk.NewAttribute(types.AttributeKeyFromChainId, strconv.FormatUint(merkleValue.FromChainID, 10)),
			sdk.NewAttribute(types.AttributeKeyMerkleValueMakeTxParamToContractAddress, hex.EncodeToString(merkleValue.MakeTxParam.ToContractAddress)),
		),
	})
	return merkleValue, nil

}

func (k Keeper) checkDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) error {
	store := ctx.KVStore(k.storeKey)
	value := store.Get(GetDoneTxKey(fromChainId, crossChainId))
	if value != nil {
		return fmt.Errorf("checkDoneTx, tx already done")
	}
	return nil
}
func (k Keeper) putDoneTx(ctx sdk.Context, fromChainId uint64, crossChainId []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetDoneTxKey(fromChainId, crossChainId), crossChainId)
}

func (k Keeper) getCrossChainId(ctx sdk.Context) (sdk.Int, error) {
	store := ctx.KVStore(k.storeKey)
	idBs := store.Get(CrossChainIdKey)
	if idBs == nil {
		return sdk.NewInt(0), nil
	}
	var crossChainId sdk.Int
	if err := k.cdc.UnmarshalBinaryLengthPrefixed(idBs, &crossChainId); err != nil {
		return sdk.NewInt(0), types.ErrUnmarshalSpecificTypeFail(crossChainId, err)
	}

	return crossChainId, nil
}
func (k Keeper) setCrossChainId(ctx sdk.Context, crossChainId sdk.Int) error {
	store := ctx.KVStore(k.storeKey)
	idBs, err := k.cdc.MarshalBinaryLengthPrefixed(crossChainId)
	if err != nil {
		return types.ErrMarshalSpecificTypeFail(crossChainId, err)
	}
	store.Set(CrossChainIdKey, idBs)
	return nil
}
