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

package types // noalias

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	supplyexported "github.com/cosmos/cosmos-sdk/x/supply/exported"
	hs "github.com/polynetwork/cosmos-poly-module/headersync"
	polytype "github.com/polynetwork/poly/core/types"
)

// SupplyKeeper defines the expected supply keeper
type HeaderSyncKeeper interface {
	ProcessHeader(ctx sdk.Context, header *polytype.Header, headerProof []byte, curHeader *polytype.Header) error
	GetConsensusPeers(ctx sdk.Context, chainId uint64) (*hs.ConsensusPeers, error)
}

// SupplyKeeper defines the expected supply keeper
type SupplyKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
	GetModuleAccount(ctx sdk.Context, name string) supplyexported.ModuleAccountI
	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	SetModuleAccount(sdk.Context, supplyexported.ModuleAccountI)

	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	MintCoins(ctx sdk.Context, name string, amt sdk.Coins) error
	BurnCoins(ctx sdk.Context, moduleName string, amt sdk.Coins) error
	SetSupply(ctx sdk.Context, supply supplyexported.SupplyI)
	GetSupply(ctx sdk.Context) (supply supplyexported.SupplyI)
}

type UnlockKeeper interface {
	Unlock(ctx sdk.Context, fromChainId uint64, fromContractAddr sdk.AccAddress, toContractAddr []byte, argsBs []byte) error
	ContainToContractAddr(ctx sdk.Context, toContractAddr []byte, fromChainId uint64) bool
}

type AssetKeeper interface {
	RegisterAsset(ctx sdk.Context, fromChainId uint64, fromContractAddr []byte, toContractAddr []byte, argsBs []byte) error
}
