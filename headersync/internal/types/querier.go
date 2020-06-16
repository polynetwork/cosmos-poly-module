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

const (
	QueryHeader        = "header"
	QueryCurrentHeight = "current_height"
	QueryKeyHeights    = "key_heights"
	QueryKeyHeight     = "key_height"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryHeaderParams struct {
	ChainId uint64
	Height  uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryHeaderParams(chainId uint64, height uint32) QueryHeaderParams {
	return QueryHeaderParams{ChainId: chainId, Height: height}
}

type QueryCurrentHeightParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryCurrentHeightParams(chainId uint64) QueryCurrentHeightParams {
	return QueryCurrentHeightParams{ChainId: chainId}
}

type QueryKeyHeightsParams struct {
	ChainId uint64
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryKeyHeightsParams(chainId uint64) QueryKeyHeightsParams {
	return QueryKeyHeightsParams{ChainId: chainId}
}

type QueryKeyHeightParams struct {
	ChainId uint64
	Height  uint32
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryKeyHeightParams(chainId uint64, height uint32) QueryKeyHeightParams {
	return QueryKeyHeightParams{ChainId: chainId, Height: height}
}
