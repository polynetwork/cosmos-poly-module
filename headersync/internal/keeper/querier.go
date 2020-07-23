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
	polycommon "github.com/polynetwork/poly/common"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
)

// NewQuerier returns a minting Querier handler.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryConsensusPeers:
			return queryConsensusPeers(ctx, req, k)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])

		}
	}
}

func queryConsensusPeers(ctx sdk.Context, req abci.RequestQuery, k Keeper) ([]byte, error) {
	var params types.QueryConsensusPeersParams

	if err := types.ModuleCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrJSONUnmarshal, "failed to parse params: %s", err)
	}
	consensusPeers, err := k.GetConsensusPeers(ctx, params.ChainId)
	if err != nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "failed to get consensus peers for chainId: %d, Error: %s", params.ChainId, err)
	}
	sink := polycommon.NewZeroCopySink(nil)
	consensusPeers.Serialization(sink)
	return sink.Bytes(), nil
}
