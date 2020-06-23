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
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
)

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryProxyByOperator(cliCtx context.CLIContext, queryRoute string, operator sdk.AccAddress) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProxyByOperator),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryProxyByOperatorParam(operator)),
	)
	return res, err
}

// QueryDelegatorTotalRewards queries delegator total rewards.
func QueryProxyHash(cliCtx context.CLIContext, queryRoute string, lockProxyHash []byte, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryProxyHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryProxyHashParam(lockProxyHash, chainId)),
	)
	return res, err
}

func QueryAssetHash(cliCtx context.CLIContext, queryRoute string, lockProxy []byte, sourceAssetDenom string, chainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryAssetHash),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryAssetHashParam(lockProxy, sourceAssetDenom, chainId)),
	)
	return res, err
}
