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

package btcx

import (
	"github.com/polynetwork/cosmos-poly-module/btcx/client/rest"
	"github.com/polynetwork/cosmos-poly-module/btcx/exported"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	RouterKey         = types.RouterKey

	AttributeValueCategory = types.AttributeValueCategory

	EventTypeCreateDenom         = types.EventTypeCreateDenom
	AttributeKeySourceAssetDenom = types.AttributeKeySourceAssetDenom
	AttributeKeyFromAssetHash    = types.AttributeKeyFromAssetHash
	AttributeKeyRedeemScript     = types.AttributeKeyRedeemScript

	EventTypeBindAsset      = types.EventTypeBindAsset
	AttributeKeyCreator     = types.AttributeKeyCreator
	AttributeKeyToChainId   = types.AttributeKeyToChainId
	AttributeKeyToAssetHash = types.AttributeKeyToAssetHash

	EventTypeLock           = types.EventTypeLock
	AttributeKeyFromAddress = types.AttributeKeyFromAddress
	AttributeKeyToAddress   = types.AttributeKeyToAddress
	AttributeKeyAmount      = types.AttributeKeyAmount

	EventTypeUnlock = types.EventTypeUnlock
)

var (
	// functions aliases
	RegisterCodec = types.RegisterCodec
	NewKeeper     = keeper.NewKeeper
	NewQuerier    = keeper.NewQuerier

	NewMsgCreateDenom   = types.NewMsgCreateDenom
	NewMsgBindAssetHash = types.NewMsgBindAssetHash
	NewMsgLock          = types.NewMsgLock

	// variable aliases
	ModuleCdc = types.ModuleCdc
)

type (
	Keeper              = keeper.Keeper
	DenomCrossChainInfo = types.DenomCrossChainInfo
	DenomInfo           = types.DenomInfo

	MsgBindAssetHash = types.MsgBindAssetHash
	MsgCreateDenom   = types.MsgCreateDenom
	MsgLock          = types.MsgLock
	CreateCoinReq    = rest.CreateCoinReq
	BindAssetHashReq = rest.BindAssetHashReq
	LockReq          = rest.LockReq

	ToBTCArgs = types.ToBTCArgs
	BTCArgs   = types.BTCArgs

	UnlockKeeper = exported.UnlockKeeper
)
