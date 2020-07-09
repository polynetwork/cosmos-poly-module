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

package lockproxypip1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/lockproxypip1/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCreateLockProxy:
			return handleMsgCreateLockProxy(ctx, k, msg)
		case types.MsgCreateCoinAndDelegateToProxy:
			return handleMsgCreateCoinAndDelegateToProxy(ctx, k, msg)
		case types.MsgLock:
			return handleMsgLock(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized %s message type: %T", types.ModuleName, msg)
		}
	}
}

func handleMsgCreateLockProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateLockProxy) (*sdk.Result, error) {
	err := k.CreateLockProxy(ctx, msg.Creator)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)
	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgCreateCoinAndDelegateToProxy(ctx sdk.Context, k keeper.Keeper, msg types.MsgCreateCoinAndDelegateToProxy) (*sdk.Result, error) {
	if err := k.CreateCoinAndDelegateToProxy(ctx, msg.Creator, msg.Coin, msg.LockProxyHash, msg.NativeChainId, msg.NativeLockProxyHash, msg.NativeAssetHash); err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

func handleMsgLock(ctx sdk.Context, k keeper.Keeper, msg types.MsgLock) (*sdk.Result, error) {
	err := k.Lock(ctx, msg.LockProxyHash, msg.FromAddress, msg.SourceAssetDenom, msg.ToChainId, msg.ToChainProxyHash, msg.ToChainAssetHash, msg.ToAddressBs, msg.Value, msg.DeductFeeInLock, msg.FeeAmount, msg.FeeAddress)
	if err != nil {
		return nil, err
	}
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
