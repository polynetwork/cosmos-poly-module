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
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/polynetwork/cosmos-poly-module/ft/internal/types"
	"github.com/polynetwork/cosmos-poly-module/lockproxy"
)

func (k Keeper) CreateCoinAndDelegateToProxy(ctx sdk.Context, creator sdk.AccAddress, coin sdk.Coin, lockproxyHash []byte) error {

	if reason, exist := k.ExistDenom(ctx, coin.Denom); exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("denom:%s already exist, due to reason:%s", coin.Denom, reason))
	}
	if exist := k.lockProxyKeeper.EnsureLockProxyExist(ctx, lockproxyHash); !exist {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("lockproxy with hash: %s not created", lockproxyHash))

	}
	//k.SetOperator(ctx, denom, creator)
	k.ccmKeeper.SetDenomCreator(ctx, coin.Denom, creator)

	if err := k.supplyKeeper.MintCoins(ctx, lockproxy.ModuleName, sdk.NewCoins(coin)); err != nil {
		return types.ErrCreateCoinAndDelegateToProxy(fmt.Sprintf("supplyKeeper.MintCoins Error: %s", err.Error()))
	}
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCreateAndDelegateCoinToProxy,
			sdk.NewAttribute(types.AttributeKeySourceAssetDenom, coin.Denom),
			sdk.NewAttribute(types.AttributeKeyCreator, creator.String()),
			sdk.NewAttribute(types.AttributeKeyAmount, coin.Amount.String()),
		),
	})
	k.Logger(ctx).Info(fmt.Sprintf("creator:%s initialized coin: %s ", creator.String(), coin.String()))
	return nil
}
