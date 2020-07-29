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
	"crypto/sha256"
	"encoding/binary"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/types"
)

const (
	// default paramspace for params keeper
	DefaultParamspace = types.ModuleName
)

var (
	OperatorToLockProxyKey = []byte{0x01}
	BindChainIdPrefix      = []byte{0x02}
	RegistryPrefix         = []byte{0x03}
	BalancePrefix          = []byte{0x04}
)

func GetOperatorToLockProxyKey(operator sdk.AccAddress) []byte {
	return append(OperatorToLockProxyKey, operator...)
}

func GetHashKey(lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) []byte {
	nativeChainIdBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nativeChainIdBytes, nativeChainId)

	lockProxyHashBz := sha256.Sum256(lockProxyHash)
	assetHashBz := sha256.Sum256(assetHash)
	nativeChainIdBz := sha256.Sum256(nativeChainIdBytes)
	nativeLockProxyHashBz := sha256.Sum256(nativeLockProxyHash)
	nativeAssetHashBz := sha256.Sum256(nativeAssetHash)

	return append(append(append(append(lockProxyHashBz[:], assetHashBz[:]...), nativeChainIdBz[:]...), nativeLockProxyHashBz[:]...), nativeAssetHashBz[:]...)
}

func GetRegistryKey(lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) []byte {
	hashKey := GetHashKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	return append(RegistryPrefix, hashKey...)
}

func GetBalanceKey(lockProxyHash []byte, assetHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) []byte {
	hashKey := GetHashKey(lockProxyHash, assetHash, nativeChainId, nativeLockProxyHash, nativeAssetHash)
	return append(BalancePrefix, hashKey...)
}

func GetBindChainIdKey(proxyHash []byte, toChainId uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, toChainId)
	return append(append(BindChainIdPrefix, proxyHash...), b...)
}
