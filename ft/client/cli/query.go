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
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/polynetwork/cosmos-poly-module/ft/client/common"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
	"strconv"
)

// GetQueryCmd returns the cli query commands for the minting module.
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
			GetCmdQueryDenomCrossChainInfo(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryDenomInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denom-info [denom]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info incluing the coin creator, coin total supply

Example:
$ %s query %s denom-info btcx
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			res, err := common.QueryDenomInfo(cliCtx, queryRoute, sourceAssetdenom)
			if err != nil {
				return err
			}
			var denomInfo types.DenomInfo
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo of denom:%s is:\n %s\n", sourceAssetdenom, denomInfo.String())
			return nil
		},
	}
}

func GetCmdQueryDenomCrossChainInfo(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "denom-cross-chain-info [denom] [chain_id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query denom cross chain info correlated with a specific chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query a specific denom or coin info correlated with a specific chainId incluing the coin creator,  coin total supply, 
toChainId and the corresponding toAssetHash in hex format

Example:
$ %s query %s denom-cross-chain-info btcx 2
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetdenom := args[0]

			chainId, err := strconv.ParseUint(args[1], 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryDenomCrossChainInfo(cliCtx, queryRoute, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var denomInfo types.DenomCrossChainInfo
			cdc.MustUnmarshalJSON(res, &denomInfo)
			fmt.Printf("denomInfo in detail of denom:%s is:\n %s\n", sourceAssetdenom, denomInfo.String())
			return nil
		},
	}
}
