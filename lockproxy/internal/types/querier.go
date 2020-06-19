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
	QueryProxyHash       = "proxy_hash"
	QueryAssetHash       = "asset_hash"
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
type QueryProxyHashParam struct {
	LockProxyHash []byte
	ChainId       uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryProxyHashParam(lockProxyHash []byte, chainId uint64) QueryProxyHashParam {
	return QueryProxyHashParam{lockProxyHash, chainId}
}

type QueryAssetHashParam struct {
	LockProxyHash    []byte
	SourceAssetDenom string
	ChainId          uint64
}

func NewQueryAssetHashParam(lockProxyHash []byte, sourceAssetDenom string, chainId uint64) QueryAssetHashParam {
	return QueryAssetHashParam{LockProxyHash: lockProxyHash, SourceAssetDenom: sourceAssetDenom, ChainId: chainId}
}
