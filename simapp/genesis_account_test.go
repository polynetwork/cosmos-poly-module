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

package simapp_test

import (
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/simapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
)

func TestSimGenesisAccountValidate(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	vestingStart := time.Now().UTC()

	coins := sdk.NewCoins(sdk.NewInt64Coin("test", 1000))
	baseAcc := authtypes.NewBaseAccount(addr, nil, pubkey, 0, 0)
	require.NoError(t, baseAcc.SetCoins(coins))

	testCases := []struct {
		name    string
		sga     simapp.SimGenesisAccount
		wantErr bool
	}{
		{
			"valid basic account",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
			},
			false,
		},
		{
			"invalid basic account with mismatching address/pubkey",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(addr, nil, secp256k1.GenPrivKey().PubKey(), 0, 0),
			},
			true,
		},
		{
			"valid basic account with module name",
			simapp.SimGenesisAccount{
				BaseAccount: authtypes.NewBaseAccount(sdk.AccAddress(crypto.AddressHash([]byte("testmod"))), nil, nil, 0, 0),
				ModuleName:  "testmod",
			},
			false,
		},
		{
			"valid basic account with invalid module name/pubkey pair",
			simapp.SimGenesisAccount{
				BaseAccount: baseAcc,
				ModuleName:  "testmod",
			},
			true,
		},
		{
			"valid basic account with valid vesting attributes",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			false,
		},
		{
			"valid basic account with invalid vesting end time",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins,
				StartTime:       vestingStart.Add(2 * time.Hour).Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			true,
		},
		{
			"valid basic account with invalid original vesting coins",
			simapp.SimGenesisAccount{
				BaseAccount:     baseAcc,
				OriginalVesting: coins.Add(coins...),
				StartTime:       vestingStart.Unix(),
				EndTime:         vestingStart.Add(1 * time.Hour).Unix(),
			},
			true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.wantErr, tc.sga.Validate() != nil)
		})
	}
}
