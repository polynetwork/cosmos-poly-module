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

package lockproxy

import (
	"github.com/polynetwork/cosmos-poly-module/lockproxy/exported"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
)

const (
	ModuleName                            = types.ModuleName
	DefaultParamspace                     = types.DefaultParamspace
	StoreKey                              = types.StoreKey
	QuerierRoute                          = types.QuerierRoute
	QueryParameters                       = types.QueryParameters
	RouterKey                             = types.RouterKey
	AttributeValueCategory                = types.AttributeValueCategory
	EventTypeCreateLockProxy              = types.EventTypeCreateLockProxy
	EventTypeCreateAndDelegateCoinToProxy = types.EventTypeCreateAndDelegateCoinToProxy
	EventTypeBindProxy                    = types.EventTypeBindProxy
	EventTypeBindAsset                    = types.EventTypeBindAsset
	EventTypeLock                         = types.EventTypeLock
	EventTypeUnlock                       = types.EventTypeUnlock
	AttributeKeyCreator                   = types.AttributeKeyCreator
	AttributeKeyLockProxy                 = types.AttributeKeyLockProxy
	AttributeKeyToChainId                 = types.AttributeKeyToChainId
	AttributeKeyToChainProxyHash          = types.AttributeKeyToChainProxyHash
	AttributeKeySourceAssetDenom          = types.AttributeKeySourceAssetDenom
	AttributeKeyFromAssetHash             = types.AttributeKeyFromAssetHash
	AttributeKeyToAssetHash               = types.AttributeKeyToAssetHash
	AttributeKeyFromAddress               = types.AttributeKeyFromAddress
	AttributeKeyToAddress                 = types.AttributeKeyToAddress
	AttributeKeyAmount                    = types.AttributeKeyAmount
)

var (
	RegisterCodec                      = types.RegisterCodec
	NewKeeper                          = keeper.NewKeeper
	NewQuerier                         = keeper.NewQuerier
	ModuleCdc                          = types.ModuleCdc
	NewMsgCreateLockProxy              = types.NewMsgCreateLockProxy
	NewMsgCreateCoinAndDelegateToProxy = types.NewMsgCreateCoinAndDelegateToProxy
	NewMsgBindAssetHash                = types.NewMsgBindAssetHash
	NewMsgBindProxyHash                = types.NewMsgBindProxyHash
	NewMsgLock                         = types.NewMsgLock
	ErrInvalidChainId                  = types.ErrInvalidChainId
	ErrMsgBindAssetHash                = types.ErrMsgBindAssetHash
	ErrMsgLock                         = types.ErrMsgLock
	ErrAccountNotExist                 = types.ErrAccountNotExist
	ErrCreateLockProxy                 = types.ErrCreateLockProxy
	ErrBindProxyHash                   = types.ErrBindProxyHash
	ErrBindAssetHash                   = types.ErrBindAssetHash
	ErrLock                            = types.ErrLock
	ErrUnLock                          = types.ErrUnLock
	ErrMsgBindProxyHash                = types.ErrMsgBindProxyHash
	ErrCreateCoinAndDelegateToProxy    = types.ErrCreateCoinAndDelegateToProxy
	OperatorToLockProxyKey             = keeper.OperatorToLockProxyKey
	BindProxyPrefix                    = keeper.BindProxyPrefix
	BindAssetPrefix                    = keeper.BindAssetPrefix
	GetOperatorToLockProxyKey          = keeper.GetOperatorToLockProxyKey
	GetBindProxyKey                    = keeper.GetBindProxyKey
	GetBindAssetHashKey                = keeper.GetBindAssetHashKey
	QueryProxyByOperator               = types.QueryProxyByOperator
	QueryProxyHash                     = types.QueryProxyHash
	QueryAssetHash                     = types.QueryAssetHash
	NewQueryProxyByOperatorParam       = types.NewQueryProxyByOperatorParam
	NewQueryProxyHashParam             = types.NewQueryProxyHashParam
	NewQueryAssetHashParam             = types.NewQueryAssetHashParam
)

type (
	Keeper                          = keeper.Keeper
	MsgCreateLockProxy              = types.MsgCreateLockProxy
	MsgCreateCoinAndDelegateToProxy = types.MsgCreateCoinAndDelegateToProxy
	MsgBindProxyHash                = types.MsgBindProxyHash
	MsgBindAssetHash                = types.MsgBindAssetHash
	MsgLock                         = types.MsgLock
	TxArgs                          = types.TxArgs
	UnlockKeeper                    = exported.UnlockKeeper
)
