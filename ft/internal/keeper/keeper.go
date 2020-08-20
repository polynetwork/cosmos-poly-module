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
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply/exported"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper of the mint store
type Keeper struct {
	cdc          *codec.Codec
	storeKey     sdk.StoreKey
	authKeeper   types.AccountKeeper
	bankKeeper   types.BankKeeper
	supplyKeeper types.SupplyKeeper
	ccmKeeper    types.CrossChainManager
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc *codec.Codec, key sdk.StoreKey, ak types.AccountKeeper, bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, ccmKeeper types.CrossChainManager) Keeper {

	// ensure mint module account is set
	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the %s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:          cdc,
		storeKey:     key,
		authKeeper:   ak,
		bankKeeper:   bankKeeper,
		supplyKeeper: supplyKeeper,
		ccmKeeper:    ccmKeeper,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) GetModuleAccount(ctx sdk.Context) exported.ModuleAccountI {
	return k.supplyKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) EnsureAccountExist(ctx sdk.Context, addr sdk.AccAddress) error {
	acct := k.authKeeper.GetAccount(ctx, addr)
	if acct == nil {
		return types.ErrAccountNotExist(fmt.Sprintf("account %s does not exist", addr.String()))
	}
	return nil
}

func (k Keeper) CreateCoins(ctx sdk.Context, creator sdk.AccAddress, coins sdk.Coins) error {
	for _, coin := range coins {
		if reason, exist := k.ccmKeeper.ExistDenom(ctx, coin.Denom); exist {
			return types.ErrCreateCoins(fmt.Sprintf("denom:%s already exist, due to reason:%s", coin.Denom, reason))
		}
		k.ccmKeeper.SetDenomCreator(ctx, coin.Denom, creator)
	}
	if err := k.MintCoins(ctx, creator, sdk.NewCoins(coins...)); err != nil {
		return types.ErrCreateCoins(fmt.Sprintf("MintCoins Error: %s", err.Error()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateCoins,
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, coins.String()),
		),
	})
	return nil
}

// MintCoins creates new coins from thin air and adds it to the module account.
// Panics if the name maps to a non-minter module account or if the amount is invalid.
func (k Keeper) MintCoins(ctx sdk.Context, toAcct sdk.AccAddress, amt sdk.Coins) error {
	_, err := k.bankKeeper.AddCoins(ctx, toAcct, amt)
	if err != nil {
		return types.ErrMintCoins(fmt.Sprintf("bankkeeper.AddCoins, Error: %s", err.Error()))
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Inflate(amt)

	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("minted coin:%s to account:%s ", amt.String(), toAcct.String()))

	return nil
}

// BurnCoins burns coins deletes coins from the balance of the module account.
// Panics if the name maps to a non-burner module account or if the amount is invalid.
func (k Keeper) BurnCoins(ctx sdk.Context, fromAcct sdk.AccAddress, amt sdk.Coins) error {

	_, err := k.bankKeeper.SubtractCoins(ctx, fromAcct, amt)
	if err != nil {
		return types.ErrBurnCoins(fmt.Sprintf("bankkeeper.SubtractCoins, Error: %s", err.Error()))
	}

	// update total supply
	supply := k.supplyKeeper.GetSupply(ctx)
	supply = supply.Deflate(amt)
	k.supplyKeeper.SetSupply(ctx, supply)

	logger := k.Logger(ctx)
	logger.Info(fmt.Sprintf("burned coin:%s from account:%s ", amt.String(), fromAcct.String()))
	return nil
}
