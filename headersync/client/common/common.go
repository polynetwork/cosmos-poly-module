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

package common

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryHeader(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryHeader),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryHeaderParams(chainId, height)),
	)
	return res, err
}

func QueryCurrentHeaderHeight(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryCurrentHeight),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryCurrentHeightParams(chainId)),
	)

	return res, err
}
func QueryKeyHeights(cliCtx context.CLIContext, queryRoute string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryKeyHeights),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryKeyHeightsParams(chainId)),
	)

	return res, err
}

func QueryKeyHeight(cliCtx context.CLIContext, queryRoute string, chainId uint64, height uint32) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryKeyHeight),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryKeyHeightParams(chainId, height)),
	)

	return res, err
}
