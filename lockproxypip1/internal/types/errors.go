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
	ErrInvalidChainIdType               = sdkerrors.Register(ModuleName, 1, "ErrInvalidChainIdType")
	ErrMsgBindAssetHashType             = sdkerrors.Register(ModuleName, 2, "ErrMsgBindAssetHashType")
	ErrMsgLockType                      = sdkerrors.Register(ModuleName, 3, "ErrMsgLockType")
	ErrAccountNotExistType              = sdkerrors.Register(ModuleName, 4, "ErrAccountNotExistType")
	ErrCreateLockProxyType              = sdkerrors.Register(ModuleName, 5, "ErrCreateLockProxyType")
	ErrBindProxyHashType                = sdkerrors.Register(ModuleName, 6, "ErrBindProxyHashType")
	ErrBindAssetHashType                = sdkerrors.Register(ModuleName, 7, "ErrBindAssetHashType")
	ErrLockType                         = sdkerrors.Register(ModuleName, 8, "ErrLockType")
	ErrUnLockType                       = sdkerrors.Register(ModuleName, 9, "ErrUnLockType")
	ErrMsgBindProxyHashType             = sdkerrors.Register(ModuleName, 10, "ErrMsgBindProxyHashType")
	ErrCreateCoinAndDelegateToProxyType = sdkerrors.Register(ModuleName, 11, "ErrCreateCoinAndDelegateToProxyType")
	ErrRegistryAlreadyExistsType        = sdkerrors.Register(ModuleName, 12, "ErrRegistryAlreadyExistsType")
	ErrRegisterAssetType                = sdkerrors.Register(ModuleName, 13, "ErrRegisterAssetType")
)

func ErrInvalidChainId(chainId uint64) error {
	return sdkerrors.Wrapf(ErrInvalidChainIdType, fmt.Sprintf("Reason: unknown chainId with id %d", chainId))
}

func ErrMsgBindAssetHash(reason string) error {
	return sdkerrors.Wrapf(ErrMsgBindAssetHashType, fmt.Sprintf("Reason: %s", reason))
}
func ErrMsgLock(reason string) error {
	return sdkerrors.Wrapf(ErrMsgLockType, fmt.Sprintf("Reason: %s", reason))
}

func ErrAccountNotExist(reason string) error {
	return sdkerrors.Wrapf(ErrAccountNotExistType, fmt.Sprintf("Reason: %s", reason))
}

func ErrCreateLockProxy(reason string) error {
	return sdkerrors.Wrapf(ErrCreateLockProxyType, fmt.Sprintf("Reason: %s", reason))
}

func ErrBindProxyHash(reason string) error {
	return sdkerrors.Wrapf(ErrBindProxyHashType, fmt.Sprintf("Reason: %s", reason))
}

func ErrBindAssetHash(reason string) error {
	return sdkerrors.Wrapf(ErrBindAssetHashType, fmt.Sprintf("Reason: %s", reason))
}

func ErrLock(reason string) error {
	return sdkerrors.Wrapf(ErrLockType, fmt.Sprintf("Reason: %s", reason))
}
func ErrUnLock(reason string) error {
	return sdkerrors.Wrapf(ErrUnLockType, fmt.Sprintf("Reason: %s", reason))
}

func ErrMsgBindProxyHash(reason string) error {
	return sdkerrors.Wrapf(ErrMsgBindProxyHashType, fmt.Sprintf("Reason: %s", reason))
}

func ErrCreateCoinAndDelegateToProxy(reason string) error {
	return sdkerrors.Wrapf(ErrCreateCoinAndDelegateToProxyType, fmt.Sprintf("Reason: %s", reason))
}

func ErrRegistryAlreadyExists(reason string) error {
	return sdkerrors.Wrapf(ErrRegistryAlreadyExistsType, fmt.Sprintf("Reason: %s", reason))
}

func ErrRegisterAsset(reason string) error {
	return sdkerrors.Wrapf(ErrRegisterAssetType, fmt.Sprintf("Reason: %s", reason))
}
