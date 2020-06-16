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

package keeper_test

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/polynetwork/cosmos-poly-module/btcx"
	"github.com/polynetwork/cosmos-poly-module/ft"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ft_crossed_independently_MsgCreateDenom(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	testCases := []struct {
		address       string
		denom         string
		expectSucceed bool
	}{
		{"addr1", "coin1", true},
		{"addr1", "coin2", true},
		{"addr1", "coin1", false},
	}
	for _, testCase := range testCases {
		creator := sdk.AccAddress([]byte(testCase.address))
		err := app.FtKeeper.CreateDenom(ctx, creator, testCase.denom)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func Test_ft_crossed_independently_MsgBindAssetHash(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	creator1 := sdk.AccAddress([]byte("addr1"))
	creator2 := sdk.AccAddress([]byte("addr2"))
	err := app.FtKeeper.CreateDenom(ctx, creator1, "coin1")
	require.Nil(t, err, "create denom error")
	err = app.FtKeeper.CreateDenom(ctx, creator2, "coin2")
	require.Nil(t, err, "create denom error")

	testCases := []struct {
		operator      sdk.AccAddress
		denom         string
		toChainId     uint64
		toAssetHash   []byte
		expectSucceed bool
	}{
		{creator1, "coin1", 2, []byte{1, 2, 3, 4}, true},
		{creator2, "coin1", 2, []byte{1, 2, 3, 4, 5, 6}, false},
		{creator2, "coin2", 3, []byte{1, 2}, true},
		{creator2, "coin2", 2, []byte{1, 2, 3, 4}, true},
		{creator2, "coin2", 2, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, true},
	}
	querier := keeper.NewQuerier(app.FtKeeper)
	for _, testCase := range testCases {
		err := app.FtKeeper.BindAssetHash(ctx, testCase.operator, testCase.denom, testCase.toChainId, testCase.toAssetHash)
		if testCase.expectSucceed {
			require.Nil(t, err)

			query := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", ft.StoreKey, types.QueryDenomInfo),
				Data: app.Codec().MustMarshalJSON(types.NewQueryDenomInfoParam(testCase.denom)),
			}
			denomInfoBs, err := querier(ctx, []string{types.QueryDenomInfo}, query)
			require.NoError(t, err)
			var denomInfo types.DenomInfo
			btcx.ModuleCdc.MustUnmarshalJSON(denomInfoBs, &denomInfo)
			require.Equal(t, testCase.operator.String(), denomInfo.Creator, "creator is not correct")
			require.Equal(t, testCase.denom, denomInfo.Denom, "denom is not correct")
			require.Equal(t, hex.EncodeToString([]byte(testCase.denom)), denomInfo.AssetHash, "asset hash is not correct")
			require.True(t, denomInfo.TotalSupply.Equal(sdk.ZeroInt()), "total supply is not correct")

			query = abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", ft.StoreKey, types.QueryDenomCrossChainInfo),
				Data: app.Codec().MustMarshalJSON(types.NewQueryDenomCrossChainInfoParam(testCase.denom, testCase.toChainId)),
			}
			denomInfoBs, err = querier(ctx, []string{types.QueryDenomCrossChainInfo}, query)
			require.NoError(t, err)
			var denomCrossChainInfo types.DenomCrossChainInfo
			btcx.ModuleCdc.MustUnmarshalJSON(denomInfoBs, &denomCrossChainInfo)
			require.Equal(t, testCase.operator.String(), denomCrossChainInfo.Creator, "creator is not correct")
			require.Equal(t, testCase.denom, denomCrossChainInfo.Denom, "denom is not correct")
			require.Equal(t, hex.EncodeToString([]byte(testCase.denom)), denomCrossChainInfo.AssetHash, "asset hash is not correct")
			require.True(t, denomCrossChainInfo.TotalSupply.Equal(sdk.ZeroInt()), "total supply is not correct")
			require.Equal(t, testCase.toChainId, denomCrossChainInfo.ToChainId, "to chain id is not correct")
			require.Equal(t, hex.EncodeToString(testCase.toAssetHash), denomCrossChainInfo.ToAssetHash, "to asset hash is not correct")
		} else {
			require.Error(t, err)
		}
	}
}

func Test_ft_crossed_independently_MsgLock(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	total := app.SupplyKeeper.GetSupply(ctx).GetTotal()
	coin1CoinStr := "100coin1"
	coin1Coin, err := sdk.ParseCoin(coin1CoinStr)
	require.Nil(t, err)

	denom := coin1Coin.Denom
	creator := sdk.AccAddress([]byte("creator"))

	err = app.FtKeeper.CreateDenom(ctx, creator, denom)
	require.Nil(t, err)

	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(total.Add(coin1Coin)))
	require.Equal(t, coin1Coin.Amount, app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf("coin1"), "btcx1 amount should be 100")

	balance, err := app.BankKeeper.AddCoins(ctx, creator, sdk.Coins{coin1Coin})
	require.Nil(t, err)
	require.Equal(t, sdk.Coins{coin1Coin}, balance, "create balance is not equal to 100btcx1")

	err = app.FtKeeper.BindAssetHash(ctx, creator, denom, 2, []byte{1, 2, 3, 4})
	require.Nil(t, err)
	err = app.FtKeeper.BindAssetHash(ctx, creator, denom, 3, []byte{1, 2, 3, 5})
	require.Nil(t, err)

	testCases := []struct {
		from          sdk.AccAddress
		denom         string
		toChainId     uint64
		toAddr        []byte
		amount        sdk.Int
		expectSucceed bool
	}{
		{creator, denom, 2, []byte{1, 2}, sdk.NewInt(1), true},
		{creator, denom, 2, []byte{1, 3}, sdk.NewInt(100), false},
		{creator, denom, 3, []byte{1, 4}, sdk.NewInt(2), true},
		{creator, denom, 4, []byte{1, 5}, sdk.NewInt(3), false},
	}

	for _, testCase := range testCases {
		err := app.FtKeeper.Lock(ctx, testCase.from, testCase.denom, testCase.toChainId, testCase.toAddr, testCase.amount)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}

	balance = app.BankKeeper.GetCoins(ctx, creator)
	require.Equal(t, "97coin1", balance.String(), "balnace of creator is not balanced")

}
