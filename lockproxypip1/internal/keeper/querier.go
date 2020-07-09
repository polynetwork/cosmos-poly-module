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

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryProxyByOperator:
			return queryProxyByOperator(ctx, req, k)
		case types.QueryRegistry:
			return queryRegistry(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryProxyByOperator(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryProxyByOperatorParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	proxyHash := k.GetLockProxyByOperator(ctx, params.Operator)
	//if proxyHash == nil {
	//	return nil, sdk.ErrInternal(fmt.Sprintf("queryProxyByOperator, operator:%s havenot created lockproxy contract before", params.Operator.String()))
	//}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, proxyHash)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal proxyHash: %x to JSON", proxyHash)
	}

	return bz, nil
}

func queryRegistry(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryRegistryParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}

	isRegistered := k.AssetIsRegistered(ctx, params.LockProxyHash, params.AssetHash, params.NativeChainId, params.NativeLockProxyHash, params.NativeAssetHash)
	result := "unregistered"
	if isRegistered {
		result = "registered"
	}
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, result)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal result", result)
	}

	return bz, nil
}
