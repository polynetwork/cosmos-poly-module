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

package headersync

import (
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
)

const (
	ModuleName                    = types.ModuleName
	DefaultParamspace             = types.DefaultParamspace
	StoreKey                      = types.StoreKey
	QuerierRoute                  = types.QuerierRoute
	QueryConsensusPeers           = types.QueryConsensusPeers
	RouterKey                     = types.RouterKey
	AttributeValueCategory        = types.AttributeValueCategory
	EventTypeSyncHeader           = types.EventTypeSyncHeader
	AttributeKeyChainId           = types.AttributeKeyChainId
	AttributeKeyHeight            = types.AttributeKeyHeight
	AttributeKeyBlockHash         = types.AttributeKeyBlockHash
	AttributeKeyNativeChainHeight = types.AttributeKeyNativeChainHeight
)

var (
	ModuleCdc                    = types.ModuleCdc
	RegisterCodec                = types.RegisterCodec
	NewQuerier                   = keeper.NewQuerier
	NewKeeper                    = keeper.NewKeeper
	NewMsgSyncGenesisParam       = types.NewMsgSyncGenesisParam
	NewMsgSyncHeadersParam       = types.NewMsgSyncHeadersParam
	NewQueryConsensusPeersParams = types.NewQueryConsensusPeersParams
	GetConsensusPeerKey          = keeper.GetConsensusPeerKey
	ErrDeserializeHeader         = types.ErrDeserializeHeader
	ErrMarshalSpecificTypeFail   = types.ErrMarshalSpecificTypeFail
	ErrUnmarshalSpecificTypeFail = types.ErrUnmarshalSpecificTypeFail
	ConsensusPeerPrefix          = keeper.ConsensusPeerPrefix
	KeyHeaderHashPrefix          = keeper.KeyHeaderHashPrefix
)

type (
	Keeper              = keeper.Keeper
	ConsensusPeers      = types.ConsensusPeers
	MsgSyncGenesisParam = types.MsgSyncGenesisParam
	MsgSyncHeadersParam = types.MsgSyncHeadersParam
	QueryHeaderParams   = types.QueryConsensusPeersParams
)
