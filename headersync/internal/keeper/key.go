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

package keeper

import (
	"encoding/binary"
	"github.com/polynetwork/cosmos-poly-module/headersync/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	BlockHeaderPrefix   = []byte{0x01}
	BlockHashPrefix     = []byte{0x02}
	ConsensusPeerPrefix = []byte{0x03}
	KeyHeightPrefix     = []byte{0x04}

	BlockCurrentHeightKey = []byte("currentHeight")
)

func GetBlockHeaderKey(chainId uint64, blockHash []byte) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(BlockHeaderPrefix, b...), blockHash...)
}

func GetBlockHashKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(BlockHashPrefix, b...), bh...)
}
func GetBlockCurHeightKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(BlockCurrentHeightKey, b...)
}

func GetConsensusPeerKey(chainId uint64, height uint32) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	bh := make([]byte, 4)
	binary.LittleEndian.PutUint32(bh, height)
	return append(append(ConsensusPeerPrefix, b...), bh...)
}

func GetKeyHeightsKey(chainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, chainId)
	return append(append(KeyHeightPrefix, b...), b...)
}
