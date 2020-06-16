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

package simapp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	tmkv "github.com/tendermint/tendermint/libs/kv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func TestGetSimulationLog(t *testing.T) {
	cdc := MakeCodec()

	decoders := make(sdk.StoreDecoderRegistry)
	decoders[auth.StoreKey] = func(cdc *codec.Codec, kvAs, kvBs tmkv.Pair) string { return "10" }

	tests := []struct {
		store       string
		kvPairs     []tmkv.Pair
		expectedLog string
	}{
		{
			"Empty",
			[]tmkv.Pair{{}},
			"",
		},
		{
			auth.StoreKey,
			[]tmkv.Pair{{Key: auth.GlobalAccountNumberKey, Value: cdc.MustMarshalBinaryLengthPrefixed(uint64(10))}},
			"10",
		},
		{
			"OtherStore",
			[]tmkv.Pair{{Key: []byte("key"), Value: []byte("value")}},
			fmt.Sprintf("store A %X => %X\nstore B %X => %X\n", []byte("key"), []byte("value"), []byte("key"), []byte("value")),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.store, func(t *testing.T) {
			require.Equal(t, tt.expectedLog, GetSimulationLog(tt.store, decoders, cdc, tt.kvPairs, tt.kvPairs), tt.store)
		})
	}
}
