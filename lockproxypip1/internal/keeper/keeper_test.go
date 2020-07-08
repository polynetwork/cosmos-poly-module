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
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/polynetwork/cosmos-poly-module/simapp"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"gotest.tools/assert"
	"testing"
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	return app, ctx
}

func Test_lockproxy_CreateLockProxy(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	validLockProxy := sdk.AccAddress([]byte("validLockProxy"))
	err := app.LockProxyKeeper.CreateLockProxy(ctx, validLockProxy)
	require.Nil(t, err)

	testCases := []struct {
		address       string
		coin          string
		createProxy   int // 0: not create, 1: create, 2: will not be created
		expectSucceed bool
	}{
		{"addr1", "100coin1", 0, false},
		{"addr1", "100coin1", 1, true},
		{"addr1", "100coin2", 0, true},
		{"addr2", "3coin3", 2, false},
		{"addr2", "4coin4", 2, false},
		{"addr2", "100coin5", 2, false},
	}
	for _, testCase := range testCases {
		creator := sdk.AccAddress([]byte(testCase.address))
		coin, err := sdk.ParseCoin(testCase.coin)
		require.Nil(t, err)
		if testCase.createProxy == 1 {
			err = app.LockProxyKeeper.CreateLockProxy(ctx, creator)
			require.Nil(t, err)
		}

		err = app.LockProxyKeeper.CreateCoinAndDelegateToProxy(ctx, creator, coin, creator)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
}

func Test_lockproxy_CreateCoinAndDelegateToProxy(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	validLockProxy := sdk.AccAddress([]byte("validLockProxy"))
	err := app.LockProxyKeeper.CreateLockProxy(ctx, validLockProxy)
	require.Nil(t, err)

	testCases := []struct {
		address       string
		coin          string
		lockProxy     []byte
		expectSucceed bool
	}{
		{"addr1", "100coin1", validLockProxy, true},
		{"addr1", "10000000000000000000000000000000000000000000000000000000000000000000000000coin2", validLockProxy, true},
		{"addr1", "3coin3", validLockProxy, true},
		{"addr1", "4coin4", validLockProxy, true},
		{"addr2", "100coin6", sdk.AccAddress([]byte("invalidLockProxy")), false},
	}
	for _, testCase := range testCases {
		addr := sdk.AccAddress([]byte(testCase.address))
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
		creator := acc.GetAddress()
		require.Equal(t, addr, creator, fmt.Sprintf("expect: %s, got: %s", addr, creator))
		coin, err := sdk.ParseCoin(testCase.coin)
		require.Nil(t, err)
		err = app.LockProxyKeeper.CreateCoinAndDelegateToProxy(ctx, creator, coin, testCase.lockProxy)
		if testCase.expectSucceed {
			require.Nil(t, err)
			supply1 := app.SupplyKeeper.GetSupply(ctx).GetTotal()
			assert.Equal(t, true, supply1.AmountOf(coin.Denom).Equal(coin.Amount), "supply amount not equal")
		} else {
			require.Error(t, err)
		}
	}
}

func Test_lockproxy_BindProxyHash(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	testCases := []struct {
		lockProxyCreator   string
		coinCreator        string
		coin               string
		createProxy        int  // 0: not create, 1: create, 2: will not be created
		createCoin         int  // 0: not create, 1: create, 2: will not be created
		createProxySucceed bool // when createProxy == 1, check if expect create proxy succeed or not
		createCoinSucceed  bool // when createCoin == 1, check if expect create coin succeed or not
		toChainId          uint64
		toAssetHash        []byte
		bindSucceed        bool // check if bind succeed
	}{
		{"lp1", "addr1", "100coin1", 0, 0, false, false, 2, []byte{1, 2, 3, 4}, false},
		{"lp1", "addr1", "100coin1", 1, 0, true, false, 2, []byte{1, 2, 3, 5}, false},
		{"lp1", "addr1", "100coin1", 0, 1, true, true, 2, []byte{1, 2, 3, 5}, true},
		{"lp2", "addr2", "100coin2", 1, 1, true, true, 2, []byte{1, 2, 3, 6}, true},
		{"lp5", "addr2", "100coin2", 1, 0, true, false, 2, []byte{1, 2, 3, 9}, true}, // of course, one denom can be bond within two lockproxy, yet, which one is valid depends on the toChain lockproxy contract
	}
	for _, testCase := range testCases {
		proxyCreator := sdk.AccAddress([]byte(testCase.lockProxyCreator))

		if testCase.createProxy == 1 {
			err := app.LockProxyKeeper.CreateLockProxy(ctx, proxyCreator)
			if testCase.createProxySucceed {
				require.Nil(t, err, "expect create lock proxy nil")
			} else {
				require.Error(t, err, "expect create lock proxy error")
			}
		}
		coinCreator := sdk.AccAddress(testCase.coinCreator)
		coin, err := sdk.ParseCoin(testCase.coin)
		require.Nil(t, err, "expect parse coin nil")
		if testCase.createCoin == 1 {

			err := app.LockProxyKeeper.CreateCoinAndDelegateToProxy(ctx, coinCreator, coin, proxyCreator)
			if testCase.createCoinSucceed {
				require.Nil(t, err, "expect create lock proxy nil")
			} else {
				require.Error(t, err, "expect create lock proxy error")
			}
		}

		err = app.LockProxyKeeper.BindAssetHash(ctx, proxyCreator, coin.Denom, testCase.toChainId, testCase.toAssetHash)
		if testCase.bindSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}

	}
}
