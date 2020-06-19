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
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/polynetwork/cosmos-poly-module/ccm/client/common"
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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
			GetCmdQueryIfContainContract(queryRoute, cdc),
			GetCmdQueryCcmParams(queryRoute, cdc),
			GetCmdQueryModuleBalance(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

func GetCmdQueryIfContainContract(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "if-contain-contract [module_store_key] [to_contract_addr] [from_chain_id]",
		Args:  cobra.ExactArgs(3),
		Short: "Query if module_store_key module should be targeted to execute `unlock` logic based on ToMerkleValue.MakeTxParam.ToContractAddress and ToMerkleValue.FromChainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`

Example:
$ %s query %s if-contain-contract btcx c330431496364497d7257839737b5e4596f5ac06 2
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			keystore := args[0]
			toContractAddr, _ := hex.DecodeString(args[1])
			fromChainId, _ := strconv.ParseInt(args[2], 10, 64)

			resBs, err := common.QueryIfContainContract(cliCtx, queryRoute, keystore, toContractAddr, uint64(fromChainId))
			if err != nil {
				return err
			}
			var res types.QueryContainToContractRes
			cdc.MustUnmarshalJSON(resBs, &res)
			fmt.Printf("QueryContainToContractAddr res is:\n %s\n", res.String())
			return nil
		},
	}
}

func GetCmdQueryCcmParams(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "parameters",
		Args:  cobra.NoArgs,
		Short: "Query the parameters of ccm module, including chainid used in poly chain network",
		Long: strings.TrimSpace(
			fmt.Sprintf(`

Example:
$ %s query %s parameters
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			paramsBs, err := common.QueryParams(cliCtx, queryRoute)
			if err != nil {
				return err
			}
			var params types.Params
			if err := cdc.UnmarshalJSON(paramsBs, &params); err != nil {
				return err
			}
			fmt.Printf("Paramters res is:\n %s\n", params.String())
			return nil
		},
	}
}

func GetCmdQueryModuleBalance(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use: "module-balance [module_name]",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the module account balance
Example:
$ %s query %s module-balance lockproxy
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			moduleName := args[0]

			res, err := common.QueryModuleBalance(cliCtx, queryRoute, moduleName)
			if err != nil {
				return err
			}
			var balance sdk.Coins
			cdc.MustUnmarshalJSON(res, &balance)
			fmt.Printf("module balance of module: %s is : %s\n", moduleName, balance.String())
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
