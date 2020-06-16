//nolint
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
	"fmt"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrInvalidChainIdType = sdkerrors.Register(ModuleName, 1, "ErrInvalidChainIdType")

	ErrMsgBindAssetHashType             = sdkerrors.Register(ModuleName, 2, "ErrMsgBindAssetHashType")
	ErrBindAssetHashType                = sdkerrors.Register(ModuleName, 3, "ErrBindAssetHashType")
	ErrAccountNotExistType              = sdkerrors.Register(ModuleName, 4, "ErrAccountNotExistType")
	ErrCreateCoinsType                  = sdkerrors.Register(ModuleName, 5, "ErrCreateCoinsType")
	ErrCreateCoinAndDelegateToProxyType = sdkerrors.Register(ModuleName, 6, "ErrCreateCoinAndDelegateToProxyType")
	ErrMsgCreateDenomType               = sdkerrors.Register(ModuleName, 7, "ErrMsgCreateDenomType")
	ErrCreateDenomType                  = sdkerrors.Register(ModuleName, 8, "ErrCreateDenomType")
	ErrBurnCoinsType                    = sdkerrors.Register(ModuleName, 9, "ErrBurnCoinsType")
	ErrMintCoinsType                    = sdkerrors.Register(ModuleName, 10, "ErrMintCoinsType")
	ErrMsgLockType                      = sdkerrors.Register(ModuleName, 11, "ErrMsgLockType")
	ErrLockType                         = sdkerrors.Register(ModuleName, 12, "ErrLockType")
	ErrUnLockType                       = sdkerrors.Register(ModuleName, 13, "ErrUnLockType")
)

func ErrInvalidChainId(chainId uint64) error {
	return sdkerrors.Wrapf(ErrInvalidChainIdType, fmt.Sprintf("Reason: unknown chainId with id %d", chainId))
}

func ErrMsgBindAssetHash(reason string) error {
	return sdkerrors.Wrapf(ErrMsgBindAssetHashType, fmt.Sprintf("Reason %s", reason))
}

func ErrBindAssetHash(reason string) error {
	return sdkerrors.Wrapf(ErrBindAssetHashType, fmt.Sprintf("Reason %s", reason))
}

func ErrAccountNotExist(reason string) error {
	return sdkerrors.Wrapf(ErrAccountNotExistType, fmt.Sprintf("Reason: %s", reason))
}

func ErrCreateCoins(reason string) error {
	return sdkerrors.Wrapf(ErrCreateCoinsType, fmt.Sprintf("Reason: %s", reason))
}

func ErrMsgCreateDenom(reason string) error {
	return sdkerrors.Wrapf(ErrMsgCreateDenomType, fmt.Sprintf("Reason: %s", reason))
}
func ErrCreateDenom(reason string) error {
	return sdkerrors.Wrapf(ErrCreateDenomType, fmt.Sprintf("Reason: %s", reason))
}

func ErrCreateCoinAndDelegateToProxy(reason string) error {
	return sdkerrors.Wrapf(ErrCreateCoinAndDelegateToProxyType, fmt.Sprintf("Reason: %s", reason))
}

func ErrMsgLock(reason string) error {
	return sdkerrors.Wrapf(ErrMsgLockType, fmt.Sprintf("Reason: %s", reason))
}

func ErrLock(reason string) error {
	return sdkerrors.Wrapf(ErrLockType, fmt.Sprintf("Reason: %s", reason))
}
func ErrUnLock(reason string) error {
	return sdkerrors.Wrapf(ErrUnLockType, fmt.Sprintf("Reason: %s", reason))
}

func ErrBurnCoins(reason string) error {
	return sdkerrors.Wrapf(ErrBurnCoinsType, "Reason: %s", reason)
}
func ErrMintCoins(reason string) error {
	return sdkerrors.Wrapf(ErrMintCoinsType, "Reason: %s", reason)
}
