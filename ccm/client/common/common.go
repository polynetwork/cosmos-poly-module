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
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
)

func QueryIfContainContract(cliCtx context.CLIContext, queryRoute string, keystore string, toContractAddr []byte, fromChainId uint64) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryIfContainContract),
		cliCtx.Codec.MustMarshalJSON(types.NewQueryContainToContractParam(keystore, toContractAddr, fromChainId)),
	)
	return res, err
}

func QueryParams(cliCtx context.CLIContext, queryRoute string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.QueryParameters),
		nil,
	)
	return res, err
}

func QueryModuleBalance(cliCtx context.CLIContext, queryRoute string, moduleName string) ([]byte, error) {

	res, _, err := cliCtx.QueryWithData(
		fmt.Sprintf("custom/%s/%s", queryRoute, types.NewQueryModuleBalanceParam(moduleName)),
		nil,
	)
	return res, err
}
