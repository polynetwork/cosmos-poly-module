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
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryIfContainContract:
			return queryContractToContractAddr(ctx, req, k)
		case types.QueryParameters:
			return queryParams(ctx, k)
		case types.QueryModuleBalance:
			return queryModuleBalance(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
		}
	}
}

func queryContractToContractAddr(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryContainToContractParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	resInfo := k.IfContainToContract(ctx, params.KeyStore, params.ToContractAddr, params.FromChainId)

	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, resInfo)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal value: %+v to JSON", resInfo)
	}

	return bz, nil
}

func queryParams(ctx sdk.Context, k Keeper) ([]byte, error) {

	params := k.GetParams(ctx)
	bz, e := codec.MarshalJSONIndent(types.ModuleCdc, params)
	if e != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal value: %+v to JSON", params)
	}

	return bz, nil
}

func queryModuleBalance(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryModuleBalanceParam

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	balance, err := k.GetModuleBalance(ctx, params.ModuleName)
	if err != nil {
		return nil, err
	}

	bz, err := codec.MarshalJSONIndent(types.ModuleCdc, balance)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONMarshal, "could not marshal module: %s balance: %s  to JSON", params.ModuleName, balance)
	}

	return bz, nil
}
