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

	"github.com/stretchr/testify/require"
	"testing"
)

func Test_ft_crossed_via_lockproxy_CreateCoinAndDelegateToProxy(t *testing.T) {
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
		coins, err := sdk.ParseCoin(testCase.coin)
		require.Nil(t, err)
		err = app.FtKeeper.CreateCoinAndDelegateToProxy(ctx, creator, coins, testCase.lockProxy)
		if testCase.expectSucceed {
			require.Nil(t, err)
		} else {
			require.Error(t, err)
		}
	}
}
