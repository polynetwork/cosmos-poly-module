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
	"testing"
)

func createTestApp(isCheckTx bool) (*simapp.SimApp, sdk.Context) {
	app := simapp.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	return app, ctx
}

func Test_ft_MsgCreateCoins(t *testing.T) {
	app, ctx := createTestApp(true)
	app.SupplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.Coins{sdk.NewCoin("stake", sdk.NewInt(100))}))

	testCases := []struct {
		address       string
		coins         string
		expectSucceed bool
	}{
		{"addr1", "100coin1", true},
		{"addr1", "10000000000000000000000000000000000000000000000000000000000000000000000000coin2", true},
		{"addr1", "1coin3,3coin4,4coin5", true},
		{"addr2", "100coin3,5coin6", false},
		{"addr2", "1coin4,2coin5", false},
		{"addr2", "1coin6,100btcx7", true},
	}
	for _, testCase := range testCases {
		addr := sdk.AccAddress([]byte(testCase.address))
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
		creator := acc.GetAddress()
		require.Equal(t, addr, creator, fmt.Sprintf("expect: %s, got: %s", addr, creator))
		coins, err := sdk.ParseCoins(testCase.coins)
		require.Nil(t, err)
		err = app.FtKeeper.CreateCoins(ctx, creator, coins)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
