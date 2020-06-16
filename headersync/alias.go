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
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	RouterKey         = types.RouterKey

	AttributeValueCategory        = types.AttributeValueCategory
	EventTypeSyncHeader           = types.EventTypeSyncHeader
	AttributeKeyChainId           = types.AttributeKeyChainId
	AttributeKeyHeight            = types.AttributeKeyHeight
	AttributeKeyBlockHash         = types.AttributeKeyBlockHash
	AttributeKeyNativeChainHeight = types.AttributeKeyNativeChainHeight
)

var (
	// functions aliases
	RegisterCodec = types.RegisterCodec
	NewQuerier    = keeper.NewQuerier

	NewKeeper              = keeper.NewKeeper
	NewMsgSyncGenesisParam = types.NewMsgSyncGenesisParam
	NewMsgSyncHeadersParam = types.NewMsgSyncHeadersParam

	NewQueryHeaderParams        = types.NewQueryHeaderParams
	NewQueryCurrentHeightParams = types.NewQueryCurrentHeightParams

	// key function
	GetBlockHeaderKey    = keeper.GetBlockHeaderKey
	GetBlockHashKey      = keeper.GetBlockHashKey
	GetBlockCurHeightKey = keeper.GetBlockCurHeightKey
	GetConsensusPeerKey  = keeper.GetConsensusPeerKey
	GetKeyHeightsKey     = keeper.GetKeyHeightsKey

	// variable aliases
	ModuleCdc                    = types.ModuleCdc
	ErrDeserializeHeader         = types.ErrDeserializeHeader
	ErrMarshalSpecificTypeFail   = types.ErrMarshalSpecificTypeFail
	ErrUnmarshalSpecificTypeFail = types.ErrUnmarshalSpecificTypeFail

	BlockHeaderPrefix   = keeper.BlockHeaderPrefix
	BlockHashPrefix     = keeper.BlockHashPrefix
	ConsensusPeerPrefix = keeper.ConsensusPeerPrefix
	KeyHeightPrefix     = keeper.KeyHeightPrefix

	BlockCurrentHeightKey = keeper.BlockCurrentHeightKey
	QueryHeader           = types.QueryHeader
	QueryCurrentHeight    = types.QueryCurrentHeight
	QueryKeyHeights       = types.QueryKeyHeights
	QueryKeyHeight        = types.QueryKeyHeight
)

type (
	Keeper                   = keeper.Keeper
	MsgSyncGenesisParam      = types.MsgSyncGenesisParam
	MsgSyncHeadersParam      = types.MsgSyncHeadersParam
	QueryHeaderParams        = types.QueryHeaderParams
	QueryCurrentHeightParams = types.QueryCurrentHeightParams
	QueryKeyHeightsParams    = types.QueryKeyHeightsParams
)
