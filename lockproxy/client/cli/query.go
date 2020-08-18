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
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"strings"

	"github.com/spf13/cobra"

	"encoding/hex"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/client/common"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
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
			GetCmdQueryProxyByOperator(queryRoute, cdc),
			GetCmdQueryProxyHash(queryRoute, cdc),
			GetCmdQueryAssetHash(queryRoute, cdc),
		)...,
	)

	return ccQueryCmd
}

// GetCmdQueryValidatorOutstandingRewards implements the query validator outstanding rewards command.
func GetCmdQueryProxyByOperator(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxy-hash-by-operator [operator_address]",
		Args:  cobra.ExactArgs(1),
		Short: "Query lockproxy hex string by the operator/creator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the created lock proxy contract hash created by the operator

Example:
$ %s query %s proxy-hash-by-operator cosmos1c0n2e6kuzp03pqm3av9q2v0fqn6ql3z5c5ddw7
`,
				version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			operatorAddr, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}
			res, err := common.QueryProxyByOperator(cliCtx, queryRoute, operatorAddr)

			if err != nil {
				return err
			}
			var proxyHash []byte
			cdc.MustUnmarshalJSON(res, &proxyHash)
			fmt.Printf("creator:%s with lock proxy hash:%x \n", operatorAddr, proxyHash)
			return nil
		},
	}
}

// GetCmdQueryParams implements a command to return the current minting
// parameters.
func GetCmdQueryProxyHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "proxy-hash [lock_proxy_hash/proxy_creator_address] [chainId]",
		Short: "Query the proxy hash deployed in another chain with chainId",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query proxy contract hash bond with self in another blockchain 
with chainId

Example:
$ %s query %s proxy-hash e931a4f7020caaacf3ce942567625ebbc0a0ab35 2
Or
$ %s query %s proxy-hash cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf 2
`,
				version.ClientName, types.ModuleName, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			lockProxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				lockProxyBs, err1 := hex.DecodeString(args[0])
				if err1 != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, fmt.Sprintf("lockproxy: %s or operator decord Error: %s", err, err1))
				}
				lockProxy = append(lockProxy, lockProxyBs...)
			}
			chainIdStr := args[1]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryProxyHash(cliCtx, queryRoute, lockProxy, chainId)
			if err != nil {
				return err
			}
			var proxyHash []byte
			cdc.MustUnmarshalJSON(res, &proxyHash)
			fmt.Printf("toChain proxy_hash: %s\n", hex.EncodeToString(proxyHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}

// GetCmdQueryInflation implements a command to return the current minting
// inflation value.
func GetCmdQueryAssetHash(queryRoute string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "asset-hash [lock_proxy_hash/operator] [sourceassetdenom] [chainId]",
		Short: "Query the asset hash in chainId chain corresponding with soureAssetDenom",
		Args:  cobra.ExactArgs(3),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the currently synced height of chainId blockchain

Example:
$ %s query %s proxy-hash e931a4f7020caaacf3ce942567625ebbc0a0ab35 stake 2
Or
$ %s query %s proxy-hash cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf stake 2
`,
				version.ClientName, types.ModuleName, version.ClientName, types.ModuleName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			lockProxy, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				lockProxyBs, err1 := hex.DecodeString(args[0])
				if err1 != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, fmt.Sprintf("lockproxy: %s or operator decord Error: %s", err, err1))
				}
				lockProxy = append(lockProxy, lockProxyBs...)
			}

			sourceAssetdenom := args[1]

			chainIdStr := args[2]

			chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
			if err != nil {
				return err
			}
			res, err := common.QueryAssetHash(cliCtx, queryRoute, lockProxy, sourceAssetdenom, chainId)
			if err != nil {
				return err
			}
			var assetHash []byte
			cdc.MustUnmarshalJSON(res, &assetHash)
			fmt.Printf("asset_hash: %s\n", hex.EncodeToString(assetHash))
			//return cliCtx.PrintOutput(hex.EncodeToString(proxyHash))
			return nil
		},
	}
}
