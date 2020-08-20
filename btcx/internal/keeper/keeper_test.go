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
	supply "github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/polynetwork/cosmos-poly-module/btcx"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/btcx/internal/types"
	"github.com/polynetwork/cosmos-poly-module/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"testing"
)

// returns context and an app with updated btcx keeper
func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	return app, ctx
}

func btcx_initSupply(t *testing.T, app *simapp.SimApp, ctx sdk.Context) sdk.Coins {
	coinsStr := "1000000000btca,1000000000000000000btcb"
	initTokens, err := sdk.ParseCoins(coinsStr)
	require.Equal(t, nil, err, "Parse Coins error should be nil")
	totalSupply := sdk.NewCoins(initTokens...)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(totalSupply))

	total := app.SupplyKeeper.GetSupply(ctx).GetTotal()

	require.Equal(t, totalSupply, total, "supply should be initTokens")
	return total
}

func Test_btcx_MsgCreateDenom(t *testing.T) {
	app, ctx := createTestApp(true)
	btcx_initSupply(t, app, ctx)

	testCases := []struct {
		address       string
		denom         string
		redeemScrit   string
		expectSucceed bool
	}{
		{"addr1", "btcx1", "1234", true},
		{"addr1", "btcx2", "123456", true},
		{"addr1", "btcx2", "12345678", false},
		{"addr2", "btcx3", "12345678", true}, // which is valid depends on which denom is bonded in poly chain, btcx2 or btcx3
		{"addr2", "btcx4", "123456789", false},
		{"addr2", "btcx5", "1234567890", true},
	}
	for _, testCase := range testCases {
		addr := sdk.AccAddress([]byte(testCase.address))
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
		creator := acc.GetAddress()
		require.Equal(t, addr, creator, fmt.Sprintf("expect: %s, got: %s", addr, creator))
		err := app.BtcxKeeper.CreateDenom(ctx, creator, testCase.denom, testCase.redeemScrit)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}

	querier := keeper.NewQuerier(app.BtcxKeeper)
	for _, testCase := range testCases {
		if testCase.expectSucceed {
			query := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", btcx.StoreKey, types.QueryDenomInfo),
				Data: app.Codec().MustMarshalJSON(types.NewQueryDenomInfoParam(testCase.denom)),
			}
			denomInfoBs, err := querier(ctx, []string{types.QueryDenomInfo}, query)
			require.NoError(t, err)
			var denomInfo types.DenomInfo
			btcx.ModuleCdc.MustUnmarshalJSON(denomInfoBs, &denomInfo)
			require.Equal(t, sdk.AccAddress([]byte(testCase.address)).String(), denomInfo.Creator, "creator is not correct")
			require.Equal(t, testCase.denom, denomInfo.Denom, "denom is not correct")
			require.True(t, denomInfo.TotalSupply.Equal(sdk.ZeroInt()), "total supply is not correct")
			require.Equal(t, testCase.redeemScrit, denomInfo.RedeemScipt, "redeem script is not correct")
		}
	}
}

func Test_btcx_MsgBindAssetHash(t *testing.T) {
	app, ctx := createTestApp(true)
	btcx_initSupply(t, app, ctx)

	addr := sdk.AccAddress([]byte("addr"))
	acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	creator := acc.GetAddress()
	redeemScript := "12345678"
	require.Equal(t, addr, creator, fmt.Sprintf("expect: %s, got: %s", addr, creator))
	err := app.BtcxKeeper.CreateDenom(ctx, creator, "btcx1", redeemScript)
	require.Nil(t, err)
	invalidCreator := app.AccountKeeper.NewAccountWithAddress(ctx, sdk.AccAddress([]byte("invalidCreator"))).GetAddress()
	testCases := []struct {
		creator       sdk.AccAddress
		denom         string
		toChainId     uint64
		toAssetHash   []byte
		expectSucceed bool
	}{
		{creator, "btcx1", 2, []byte{1, 2, 3, 4}, true},
		{invalidCreator, "btcx1", 2, []byte{1, 2, 3, 5}, false},
		{creator, "btcx1", 3, []byte{1, 2, 3, 6}, true},
	}

	for _, testCase := range testCases {
		err := app.BtcxKeeper.BindAssetHash(ctx, testCase.creator, testCase.denom, testCase.toChainId, testCase.toAssetHash)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
	querier := keeper.NewQuerier(app.BtcxKeeper)
	for _, testCase := range testCases {
		if testCase.expectSucceed {
			query := abci.RequestQuery{
				Path: fmt.Sprintf("custom/%s/%s", btcx.StoreKey, types.QueryDenomCrossChainInfo),
				Data: app.Codec().MustMarshalJSON(types.NewQueryDenomCrossChainInfoParam(testCase.denom, testCase.toChainId)),
			}
			denomInfoBs, err := querier(ctx, []string{types.QueryDenomCrossChainInfo}, query)
			require.NoError(t, err)
			var denomInfo types.DenomCrossChainInfo
			btcx.ModuleCdc.MustUnmarshalJSON(denomInfoBs, &denomInfo)
			require.Equal(t, testCase.creator.String(), denomInfo.Creator, "creator is not correct")
			require.Equal(t, testCase.denom, denomInfo.Denom, "denom is not correct")
			require.Equal(t, hex.EncodeToString([]byte(testCase.denom)), denomInfo.AssetHash, "denom is not correct")
			require.True(t, denomInfo.TotalSupply.Equal(sdk.ZeroInt()), "total supply is not correct")
			require.Equal(t, redeemScript, denomInfo.RedeemScipt, "redeem script is not correct")
			require.Equal(t, testCase.toChainId, denomInfo.ToChainId, "redeem script is not correct")
			require.Equal(t, hex.EncodeToString(testCase.toAssetHash), denomInfo.ToAssetHash, "redeem script is not correct")
		}
	}
}

func Test_btcx_MsgLock(t *testing.T) {
	app, ctx := createTestApp(true)
	btcx_initSupply(t, app, ctx)

	total := app.SupplyKeeper.GetSupply(ctx).GetTotal()
	btcx1CoinStr := "100btcx1"
	btcx1Coin, err := sdk.ParseCoin(btcx1CoinStr)
	require.Nil(t, err)

	denom := btcx1Coin.Denom
	creator := sdk.AccAddress([]byte("creator"))

	err = app.BtcxKeeper.CreateDenom(ctx, creator, denom, "12345678")
	require.Nil(t, err)

	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(total.Add(btcx1Coin)))
	require.Equal(t, btcx1Coin.Amount, app.SupplyKeeper.GetSupply(ctx).GetTotal().AmountOf("btcx1"), "btcx1 amount should be 100")

	balance, err := app.BankKeeper.AddCoins(ctx, creator, sdk.Coins{btcx1Coin})
	require.Nil(t, err)
	require.Equal(t, sdk.Coins{btcx1Coin}, balance, "create balance is not equal to 100btcx1")

	err = app.BtcxKeeper.BindAssetHash(ctx, creator, denom, 2, []byte{1, 2, 3, 4})
	require.Nil(t, err)
	err = app.BtcxKeeper.BindAssetHash(ctx, creator, denom, 3, []byte{1, 2, 3, 5})
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
		err := app.BtcxKeeper.Lock(ctx, testCase.from, testCase.denom, testCase.toChainId, testCase.toAddr, testCase.amount)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}

	balance = app.BankKeeper.GetCoins(ctx, creator)
	require.Equal(t, "97btcx1", balance.String(), "balnace of creator is not balanced")
}
