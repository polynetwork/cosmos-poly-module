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
	"reflect"
)

var (
	ErrMarshalSpecificType     = sdkerrors.Register(ModuleName, 1, "ErrMarshalSpecificType")
	ErrUnmarshalSpecificType   = sdkerrors.Register(ModuleName, 2, "ErrUnmarshalSpecificType")
	ErrProcessCrossChainTxType = sdkerrors.Register(ModuleName, 3, "ErrProcessCrossChainTxType")
	ErrVerifyToCosmosTxType    = sdkerrors.Register(ModuleName, 4, "ErrVerifyToCosmosTxType")

	ErrMsgProcessCrossChainTxType = sdkerrors.Register(ModuleName, 5, "ErrMsgProcessCrossChainTxType")
	ErrMsgCreateCrossChainTxType  = sdkerrors.Register(ModuleName, 6, "ErrMsgCreateCrossChainTxType")
	ErrGetModuleBalanceType       = sdkerrors.Register(ModuleName, 7, "ErrGetModuleBalanceType")
)

func ErrMarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(ErrMarshalSpecificType, fmt.Sprintf("Marshal type: %s, Error: %s", reflect.TypeOf(o).String(), err.Error()))
}
func ErrUnmarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(ErrUnmarshalSpecificType, fmt.Sprintf("Umarshal type: %s, Error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrProcessCrossChainTx(reason string) error {
	return sdkerrors.Wrapf(ErrProcessCrossChainTxType, "Reason: %s", reason)
}

func ErrVerifyToCosmosTx(reason string) error {
	return sdkerrors.Wrapf(ErrVerifyToCosmosTxType, "Reason: %s", reason)
}

func ErrMsgProcessCrossChainTx(reason string) error {
	return sdkerrors.Wrapf(ErrMsgProcessCrossChainTxType, "Reason: %s", reason)
}

func ErrMsgCreateCrossChainTx(reason string) error {
	return sdkerrors.Wrapf(ErrMsgCreateCrossChainTxType, "Reason: %s", reason)
}

func ErrGetModuleBalance(reason string) error {
	return sdkerrors.Wrapf(ErrGetModuleBalanceType, "Reason: %s", reason)
}
