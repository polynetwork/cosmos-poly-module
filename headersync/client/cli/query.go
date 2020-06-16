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
	polycommon "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/common"
	polytype "github.com/polynetwork/cosmos-poly-module/headersync/poly-utils/core/types"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/polynetwork/cosmos-poly-module/headersync/client/common"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
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
			GetCmdQueryHeader(queryRoute, cdc),
			GetCmdQueryCurrentHeight(queryRoute, cdc),
			GetCmdQueryKeyHeights(queryRoute, cdc),
			GetCmdQueryKeyHeight(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryHeader(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "header [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query header of chainId of height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query block header for a specific height 
already synced from another blockchain, normally, relayer-chain (with chainId=0), into current chain 

Example:
$ %s query %s header 0 1
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]
			heightStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			height, err := strconv.ParseUint(heightStr, 10, 32)
			if err != nil {
				return err
			}

			res, err := common.QueryHeader(cliCtx, queryRoute, uint64(chainId), uint32(height))
			if err != nil {
				return err
			}
			var header polytype.Header
			if err := header.Deserialization(polycommon.NewZeroCopySource(res)); err != nil {
				fmt.Printf("Query PolyChain header Deserialization")
			}
			fmt.Printf("header of height:%d is:\n %s\n", header.Height, header.String())
			return nil
		},
	}
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryCurrentHeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "current-height [chainId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query %s current-height 0
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryCurrentHeaderHeight(cliCtx, queryRoute, chainId)
			if err != nil {
				return err
			}
			var height uint32
			cdc.MustUnmarshalJSON(res, &height)
			fmt.Printf("current synced header height of chainid:%d is: %d\n", chainId, height)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

func GetCmdQueryKeyHeights(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "key-heights [chainId]",
		Args:  cobra.ExactArgs(1),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query %s key-heights 0
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			chainIdStr := args[0]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}

			res, err := common.QueryKeyHeights(cliCtx, queryRoute, chainId)
			if err != nil {
				return err
			}
			var keyHeights types.KeyHeights
			cdc.MustUnmarshalJSON(res, &keyHeights)
			fmt.Printf("Key heights of chainid: %d is: %+v\n", chainId, keyHeights)
			return nil
			//return cliCtx.PrintOutput(MCHeader{header})
		},
	}
}

func GetCmdQueryKeyHeight(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "key-height [chainId] [height]",
		Args:  cobra.ExactArgs(2),
		Short: "Query block height",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query %s key-height 0 100
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			chainId, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}
			height, err := strconv.Atoi(args[1])
			if err != nil {
				return err
			}
			res, err := common.QueryKeyHeight(cliCtx, queryRoute, chainId, uint32(height))
			if err != nil {
				return err
			}
			var keyHeight uint32
			cdc.MustUnmarshalJSON(res, &keyHeight)
			fmt.Printf("Key height of chainid: %d is: %d\n", chainId, keyHeight)
			return nil
		},
	}
}
