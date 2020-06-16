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
	ErrInvalidChainIdType      = sdkerrors.Register(ModuleName, 1, "ErrInvalidChainIdType")
	ErrInvalidRedeemScriptType = sdkerrors.Register(ModuleName, 2, "ErrInvalidRedeemScriptType")
	ErrEmptyToAssetHashType    = sdkerrors.Register(ModuleName, 3, "ErrEmptyToAssetHashType")
	ErrCreateDenomType         = sdkerrors.Register(ModuleName, 4, "ErrCreateDenomType")
	ErrBindAssetHashType       = sdkerrors.Register(ModuleName, 5, "ErrBindAssetHashType")
	ErrLockType                = sdkerrors.Register(ModuleName, 6, "ErrLockType")
	ErrUnLockType              = sdkerrors.Register(ModuleName, 7, "ErrUnLockType")
	ErrBurnCoinsType           = sdkerrors.Register(ModuleName, 8, "ErrBurnCoinsType")
	ErrMintCoinsType           = sdkerrors.Register(ModuleName, 9, "ErrMintCoinsType")
)

func ErrInvalidChainId(chainId uint64) error {
	return sdkerrors.Wrapf(ErrInvalidChainIdType, fmt.Sprintf("unknown chainId with id %d", chainId))
}
func ErrInvalidRedeemScript(reason string) error {
	return sdkerrors.Wrapf(ErrInvalidRedeemScriptType, fmt.Sprintf("Reason: %s", reason))
}

func ErrEmptyToAssetHash(toHashStr string) error {
	return sdkerrors.Wrapf(ErrEmptyToAssetHashType, fmt.Sprintf("empty to asset hash %s", toHashStr))
}

func ErrCreateDenom(reason string) error {
	return sdkerrors.Wrapf(ErrCreateDenomType, "Reason: %s", reason)
}

func ErrBindAssetHash(reason string) error {
	return sdkerrors.Wrapf(ErrBindAssetHashType, "Reason: %s", reason)
}

func ErrLock(reason string) error {
	return sdkerrors.Wrapf(ErrLockType, "Reason: %s", reason)
}

func ErrUnLock(reason string) error {
	return sdkerrors.Wrapf(ErrUnLockType, "Reason: %s", reason)
}

func ErrBurnCoins(reason string) error {
	return sdkerrors.Wrapf(ErrBurnCoinsType, "Reason: %s", reason)
}
func ErrMintCoins(reason string) error {
	return sdkerrors.Wrapf(ErrMintCoinsType, "Reason: %s", reason)
}
