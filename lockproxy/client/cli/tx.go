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
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
	"math/big"
	"strconv"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s module send transaction subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}
	txCmd.AddCommand(flags.PostCommands(
		SendCreateLockProxyTxCmd(cdc),
		SendCreateCoinAndDelegateToProxyTxCmd(cdc),
		SendBindProxyHashTxCmd(cdc),
		SendBindAssetHashTxCmd(cdc),
		SendLockTxCmd(cdc),
	)...)
	return txCmd
}

func SendCreateLockProxyTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-lock-proxy [creator]",
		Short: "Create lockproxy contract by creator, needs creator's signature",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s create-lock-proxy cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf 
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			msg := types.NewMsgCreateLockProxy(cliCtx.GetFromAddress())
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendCreateCoinAndDelegateToProxyTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-coin-delegate [creator] [coin] [lock_proxy_hash/lock_proxy_creator]",
		Short: "Create coin by creator, and immediately delegate to the lock proxy module account",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s create-coin-delegate cosmos1lzk4nch5v2snduup2uujpud9j6gqeunqarx2d9 1000mst1, e931a4f7020caaacf3ce942567625ebbc0a0ab35
Or 
$ $ %s tx %s create-coin-delegate cosmos1lzk4nch5v2snduup2uujpud9j6gqeunqarx2d9 1000mst1, cosmos1ayc6faczpj42eu7wjsjkwcj7h0q2p2e4vrlkzf
`,
				version.ClientName, types.ModuleName, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			creator, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			coin, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			lockProxy, err := hex.DecodeString(args[2])
			if err != nil {
				lockProxyBs, err1 := sdk.AccAddressFromBech32(args[2])
				if err1 != nil {
					return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, fmt.Sprintf("lockproxy: %s or operator from hex or from Bech32 Error: %s", err, err1))
				}
				lockProxy = append(lockProxy, lockProxyBs...)
			}
			msg := types.NewMsgCreateCoinAndDelegateToProxy(creator, coin, sdk.AccAddress(lockProxy))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindProxyHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind-proxy-hash [to_chain_id] [to_chain_proxy_hash]",
		Short: "bindproxyhash by the operator",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s bind-proxy-hash 3 11223344556677889900
Or
$ %s tx %s bind-proxy-hash 3 0x11223344556677889900
`,
				version.ClientName, types.ModuleName, version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			toChainIdStr := args[0]
			toChainProxyHashStr := args[1]

			targetChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}
			if toChainProxyHashStr[0:2] == "0x" {
				toChainProxyHashStr = toChainProxyHashStr[2:]
			}
			toChainProxyHash, err := hex.DecodeString(toChainProxyHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindProxyHash(cliCtx.GetFromAddress(), targetChainId, toChainProxyHash)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendBindAssetHashTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind-asset-hash [source_asset_denom] [to_chainId] [to_asset_hash]",
		Short: "bind asset hash by the operator, ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s bind-asset-hash ont 3 00000000000000000001 100000
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			sourceAssetDenom := args[0]

			toChainIdStr := args[1]
			toChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			toAssetHashStr := args[2]
			if toAssetHashStr[0:2] == "0x" {
				toAssetHashStr = toAssetHashStr[2:]
			}
			toAssetHash, err := hex.DecodeString(toAssetHashStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'targetProxyHash' error:%v", err)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgBindAssetHash(cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAssetHash)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}

func SendLockTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lock [lock_proxy_hash] [source_asset_denom] [to_chain_id] [to_address] [amount]",
		Short: "lock amount of source_asset_denom and aim to release amount in to_chain_id chain to to_address",
		Long: strings.TrimSpace(
			fmt.Sprintf(`
Example:
$ %s tx %s lock 12341234 ont 3 616f2a4a38396ff203ea01e6c070ae421bb8ce2d 123 
`,
				version.ClientName, types.ModuleName,
			),
		),
		Args: cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			lockProxyHash, err := hex.DecodeString(args[0])
			if err != nil {
				return err
			}

			sourceAssetDenom := args[1]

			toChainIdStr := args[2]
			toChainId, err := strconv.ParseUint(toChainIdStr, 10, 64)
			if err != nil {
				return err
			}

			toAddressStr := args[3]
			toAddress, err := hex.DecodeString(toAddressStr)
			if err != nil {
				return fmt.Errorf("decode hex string 'toAddress' error:%v", err)
			}

			valueBigInt, ok := big.NewInt(0).SetString(args[4], 10)
			if !ok {
				return fmt.Errorf("read value as big int from args[3] failed")
			}
			value := sdk.NewIntFromBigInt(valueBigInt)

			// build and sign the transaction, then broadcast to Tendermint
			msg := types.NewMsgLock(lockProxyHash, cliCtx.GetFromAddress(), sourceAssetDenom, toChainId, toAddress, value)
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	return cmd
}
