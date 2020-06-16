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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryHeader:
			return queryHeader(ctx, req, k)
		case types.QueryCurrentHeight:
			return queryCurrentHeight(ctx, req, k)
		case types.QueryKeyHeights:
			return queryKeyHeights(ctx, req, k)
		case types.QueryKeyHeight:
			return queryKeyHeight(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])

		}
	}
}

func queryHeader(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryHeaderParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	header, err := k.GetHeaderByHeight(ctx, params.ChainId, params.Height)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to get header by height: %d for chainId: %d, Error: %s", params.Height, params.ChainId, err)
	}
	sink := polycommon.NewZeroCopySink(nil)
	if e := header.Serialization(sink); e != nil {
		return nil, types.ErrSerializeHeader(err)
	}
	return sink.Bytes(), nil
}
func queryCurrentHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryCurrentHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	height, err := k.GetCurrentHeight(ctx, params.ChainId)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to get current height for chainId: %d, Error: %s", params.ChainId, err)
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, height)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal height: %d to JSON", height)
	}

	return bz, nil
}

func queryKeyHeights(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryKeyHeightsParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	heights := k.GetKeyHeights(ctx, params.ChainId)
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, heights)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal keyHeights: %+v to JSON", heights)
	}

	return bz, nil
}

func queryKeyHeight(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryKeyHeightParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	height, err := k.FindKeyHeight(ctx, params.Height, params.ChainId)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to get key height for chainId: %d, height: %d, Error: %s", params.ChainId, params.Height, err)
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, height)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal height: %d to JSON", height)
	}

	return bz, nil
}
