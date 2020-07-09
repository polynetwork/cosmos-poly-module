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

package lockproxypip1

import (
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/exported"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	RouterKey         = types.RouterKey

	AttributeValueCategory = types.AttributeValueCategory

	EventTypeBindProxy           = types.EventTypeBindProxy
	AttributeKeyToChainId        = types.AttributeKeyToChainId
	AttributeKeyToChainProxyHash = types.AttributeKeyToChainProxyHash

	EventTypeCreateAndDelegateCoinToProxy = types.EventTypeCreateAndDelegateCoinToProxy

	EventTypeBindAsset           = types.EventTypeBindAsset
	AttributeKeySourceAssetDenom = types.AttributeKeySourceAssetDenom
	AttributeKeySourceAssetHash  = types.AttributeKeySourceAssetHash
	AttributeKeyToChainAssetHash = types.AttributeKeyToChainAssetHash

	EventTypeLock           = types.EventTypeLock
	AttributeKeyFromAddress = types.AttributeKeyFromAddress
	AttributeKeyToAddress   = types.AttributeKeyToAddress
	AttributeKeyAmount      = types.AttributeKeyAmount

	AttributeKeyFromChainId = types.AttributeKeyFromChainId

	EventTypeUnlock              = types.EventTypeUnlock
	AttributeKeyFromContractHash = types.AttributeKeyFromContractHash
	AttributeKeyToAssetDenom     = types.AttributeKeyToAssetDenom

	EventTypeCreateLockProxy = types.EventTypeCreateLockProxy
	AttributeKeyCreator      = types.AttributeKeyCreator
	AttributeKeyProxyHash    = types.AttributeKeyProxyHash
	AttributeKeyLockProxy    = types.AttributeKeyLockProxy
)

var (
	// functions aliases
	RegisterCodec                      = types.RegisterCodec
	NewKeeper                          = keeper.NewKeeper
	NewQuerier                         = keeper.NewQuerier
	NewMsgCreateCoinAndDelegateToProxy = types.NewMsgCreateCoinAndDelegateToProxy

	NewMsgLock            = types.NewMsgLock
	NewMsgCreateLockProxy = types.NewMsgCreateLockProxy
	// variable aliases
	ModuleCdc         = types.ModuleCdc
	OperatorKey       = types.OperatorKey
	ErrInvalidChainId = types.ErrInvalidChainId

	// query balance path
	QueryProxyByOperator = types.QueryProxyByOperator
)

type (
	Keeper                          = keeper.Keeper
	MsgCreateLockProxy              = types.MsgCreateLockProxy
	MsgCreateCoinAndDelegateToProxy = types.MsgCreateCoinAndDelegateToProxy
	MsgLock                         = types.MsgLock

	TxArgs       = types.TxArgs
	UnlockKeeper = exported.UnlockKeeper
)
