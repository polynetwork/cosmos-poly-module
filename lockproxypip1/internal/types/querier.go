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

package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	QueryProxyByOperator = "query_proxy_by_operator"
	QueryRegistry        = "registry"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryProxyByOperatorParam struct {
	Operator sdk.AccAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryProxyByOperatorParam(operator sdk.AccAddress) QueryProxyByOperatorParam {
	return QueryProxyByOperatorParam{Operator: operator}
}

// QueryBalanceParams defines the params for querying an account balance.
type QueryRegistryParam struct {
	LockProxyHash       []byte
	AssetHash           []byte
	NativeChainId       uint64
	NativeLockProxyHash []byte
	NativeAssetHash     []byte
}

// NewQueryRegistryParam creates a new instance of QueryRegistryParam.
func NewQueryRegistryParam(lockProxyHash []byte, assetHash []byte, chainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) QueryRegistryParam {
	return QueryRegistryParam{lockProxyHash, assetHash, chainId, nativeLockProxyHash, nativeAssetHash}
}
