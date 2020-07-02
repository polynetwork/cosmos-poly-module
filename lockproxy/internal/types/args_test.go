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

package types_test

import (
	"encoding/hex"
	"fmt"
	"github.com/polynetwork/cosmos-poly-module/lockproxy/internal/types"
	polycommon "github.com/polynetwork/poly/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestTxArgs_Deserialization(t *testing.T) {
	argsBs, _ := hex.DecodeString("06657468786363148af0a30541c146b214f5ac59c7999eacae0592d600407a10f35a0000000000000000000000000000000000000000000000000000")
	txArgs := new(types.TxArgs)
	err := txArgs.Deserialization(polycommon.NewZeroCopySource(argsBs), 32)
	require.Nil(t, err)
	require.Equal(t, "657468786363", hex.EncodeToString(txArgs.ToAssetHash), "to asset hash not correct")
	require.Equal(t, "8af0a30541c146b214f5ac59c7999eacae0592d6", hex.EncodeToString(txArgs.ToAddress), "to address not correct")
	require.Equal(t, "100000000000000", txArgs.Amount.String(), "to amount not correct")
}

func TestTxArgs_Serialization(t *testing.T) {
	txArgs := new(types.TxArgs)
	txArgs.ToAssetHash, _ = hex.DecodeString("0000000000000000000000000000000000000000")
	txArgs.ToAddress, _ = hex.DecodeString("344cfc3b8635f72f14200aaf2168d9f75df86fd3")
	txArgs.Amount, _ = new(big.Int).SetString("100000000000000000000000000000000000004", 10)
	sink := polycommon.NewZeroCopySink(nil)
	err := txArgs.Serialization(sink, 32)
	require.Nil(t, err)
	fmt.Printf("sink.Bytes are %x\n", sink.Bytes())
}
