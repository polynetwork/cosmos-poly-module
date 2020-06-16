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

package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/polynetwork/cosmos-poly-module/btcx/client/common"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

// GetQueryCmd returns the cli query commands for the btcx module.
func GetQueryCmd(queryRoute string, cdc *codec.Codec) *cobra.Command {
	ccQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	ccQueryCmd.AddCommand(
		flags.GetCommands(
			GetCmdQueryDenomInfo(queryRoute, cdc),
			GetCmdQueryDenomInfoWithChainId(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

func GetCmdQueryDenomInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denom-info [denom]",
		Args:  cobra.ExactArgs(1),
		Short: "Query denom info",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info incluing the coin creator,  coin total supply, the 
redeem script and redeem script hash

Example:
$ %s query %s denom-info btcx
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]

			res, err := common.QueryDenomInfo(cliCtx, queryRoute, denom)
			if err != nil {
				return err
			}
			var denomInfo types.DenomInfo
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo of denom:%s is:\n %s\n", denom, denomInfo.String())
			return nil
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryDenomInfoWithChainId(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denom-cc-info [denom] [chain_id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query denom info correlated with a specific chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info correlated with a specific chainId incluing the coin creator,  coin total supply, the 
redeem script and redeem script hash, toChainId and the corresponding toAssetHash in hex format

Example:
$ %s query %s denom-cc-info btcx 2
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			denom := args[0]
			toChainId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryDenomCrossChainInfo(cliCtx, queryRoute, denom, toChainId)
			if err != nil {
				return err
			}

			var denomCCInfo types.DenomCrossChainInfo
			cdc.MustUnmarshalJSON(res, &denomCCInfo)
			fmt.Printf("denom cross chain Info of denom:%s for chainId:%d is:\n %s\n", denom, toChainId, denomCCInfo.String())
			return nil
		},
	}
}
