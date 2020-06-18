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

package ft

import (
	"github.com/polynetwork/cosmos-poly-module/ft/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	RouterKey         = types.RouterKey

	AttributeValueCategory = types.AttributeValueCategory

	AttributeKeyToChainId        = types.AttributeKeyToChainId
	AttributeKeyToChainProxyHash = types.AttributeKeyToChainProxyHash

	EventTypeBindAsset           = types.EventTypeBindAsset
	AttributeKeySourceAssetDenom = types.AttributeKeySourceAssetDenom
	AttributeKeyFromAssetHash    = types.AttributeKeyFromAssetHash
	AttributeKeyToChainAssetHash = types.AttributeKeyToChainAssetHash

	EventTypeLock           = types.EventTypeLock
	AttributeKeyFromAddress = types.AttributeKeyFromAddress
	AttributeKeyToAddress   = types.AttributeKeyToAddress
	AttributeKeyAmount      = types.AttributeKeyAmount

	EventTypeCreateCrossChainTx = types.EventTypeCreateCrossChainTx
	AttributeCrossChainId       = types.AttributeCrossChainId
	AttributeKeyTxParamHash     = types.AttributeKeyTxParamHash
	AttributeKeyMakeTxParam     = types.AttributeKeyMakeTxParam

	EventTypeVerifyToCosmosProof                        = types.EventTypeVerifyToCosmosProof
	AttributeKeyMerkleValueTxHash                       = types.AttributeKeyMerkleValueTxHash
	AttributeKeyMerkleValueMakeTxParamTxHash            = types.AttributeKeyMerkleValueMakeTxParamTxHash
	AttributeKeyMerkleValueMakeTxParamToContractAddress = types.AttributeKeyMerkleValueMakeTxParamToContractAddress
	AttributeKeyFromChainId                             = types.AttributeKeyFromChainId

	EventTypeUnlock              = types.EventTypeUnlock
	AttributeKeyFromContractHash = types.AttributeKeyFromContractHash
	AttributeKeyToAssetDenom     = types.AttributeKeyToAssetDenom

	EventTypeSetRedeemScript = types.EventTypeSetRedeemScript
	AttributeKeyRedeemKey    = types.AttributeKeyRedeemKey
	AttributeKeyRedeemScript = types.AttributeKeyRedeemScript
)

var (
	// functions aliases
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier

	NewMsgLock                         = types.NewMsgLock
	NewMsgCreateDenom                  = types.NewMsgCreateDenom
	NewMsgCreateCoinAndDelegateToProxy = types.NewMsgCreateCoinAndDelegateToProxy
	NewMsgBindAssetHash                = types.NewMsgBindAssetHash
	NewMsgCreateCoins                  = types.NewMsgCreateCoins

	// key function

	ModuleCdc   = types.ModuleCdc
	OperatorKey = types.OperatorKey

	ErrInvalidChainId = types.ErrInvalidChainId

	// query balance path

)

type (
	Keeper = keeper.Keeper

	MsgBindAssetHash                = types.MsgBindAssetHash
	MsgLock                         = types.MsgLock
	MsgCreateDenom                  = types.MsgCreateDenom
	MsgCreateCoinAndDelegateToProxy = types.MsgCreateCoinAndDelegateToProxy
	DenomInfo                       = types.DenomInfo
	DenomCrossChainInfo             = types.DenomCrossChainInfo
	TxArgs                          = types.TxArgs
)
