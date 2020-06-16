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
	ErrDeserializeHeaderFailType    = sdkerrors.Register(ModuleName, 1, "ErrDeserializeHeaderFailType")
	ErrSerializeHeaderFailType      = sdkerrors.Register(ModuleName, 2, "ErrSerializeHeaderFailType")
	ErrHeaderEmptyType              = sdkerrors.Register(ModuleName, 3, "ErrHeaderEmptyType")
	ErrFindKeyHeightFailType        = sdkerrors.Register(ModuleName, 4, "ErrFindKeyHeightFailType")
	ErrGetConsensusPeersFailType    = sdkerrors.Register(ModuleName, 5, "ErrGetConsensusPeersFailType")
	ErrBookKeeperNumErrType         = sdkerrors.Register(ModuleName, 6, "ErrBookKeeperNumErrType")
	ErrInvalidPublicKeyType         = sdkerrors.Register(ModuleName, 7, "ErrInvalidPublicKeyType")
	ErrUnmarshalSpecificType        = sdkerrors.Register(ModuleName, 8, "ErrUnmarshalSpecificType")
	ErrVerifyMultiSigFailType       = sdkerrors.Register(ModuleName, 9, "ErrVerifyMultiSignatureFailType")
	ErrDeserializeConsensusPeerType = sdkerrors.Register(ModuleName, 10, "ErrDeserializeConsensusPeerType")
	ErrGenesisHeaderExistType       = sdkerrors.Register(ModuleName, 11, "Genesis Header already exist")
	ErrSyncBlockHeaderType          = sdkerrors.Register(ModuleName, 12, "ErrSyncBlockHeaderType")
)

func ErrSyncBlockHeader(operation string, chainId uint64, height uint32, err error) error {
	return sdkerrors.Wrapf(ErrSyncBlockHeaderType, "operation: %s, chainId: %d, height: %d,  Error :%s", operation, chainId, height, err.Error())
}
func ErrDeserializeHeader(err error) error {
	return sdkerrors.Wrap(ErrDeserializeHeaderFailType, fmt.Sprintf("Header deserialization Error:%s", err.Error()))
}

func ErrSerializeHeader(err error) error {
	return sdkerrors.Wrap(ErrSerializeHeaderFailType, fmt.Sprintf("Header serialization Error:%s", err.Error()))
}
func ErrHeaderEmpty(headerhash []byte) error {
	return sdkerrors.Wrap(ErrHeaderEmptyType, fmt.Sprintf("Header of headerHash: %x is empty", headerhash))
}
func ErrDeserializeConsensusPeer(err error) error {
	return sdkerrors.Wrap(ErrDeserializeConsensusPeerType, fmt.Sprintf("ConsensusPeer deserialization Error:%s", err.Error()))
}

func ErrMarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(ErrUnmarshalSpecificType, fmt.Sprintf("Marshal type: %s, Error: %s", reflect.TypeOf(o).String(), err.Error()))
}
func ErrUnmarshalSpecificTypeFail(o interface{}, err error) error {
	return sdkerrors.Wrap(ErrUnmarshalSpecificType, fmt.Sprintf("Umarshal type: %s, Error: %s", reflect.TypeOf(o).String(), err.Error()))
}

func ErrFindKeyHeight(height uint32, chainId uint64) error {
	return sdkerrors.Wrap(ErrFindKeyHeightFailType, fmt.Sprintf("Can not find key height with height: %d and chainId: %d", height, chainId))
}

func ErrGetConsensusPeers(height uint32, chainId uint64) error {
	return sdkerrors.Wrap(ErrGetConsensusPeersFailType, fmt.Sprintf("For chainId: %d, height: %d, Get consensus peers empty error", height, chainId))
}

func ErrBookKeeperNum(headerBookKeeperNum int, consensusNodeNum int) error {
	return sdkerrors.Wrap(ErrBookKeeperNumErrType, fmt.Sprintf("Header Bookkeepers number: %d must more than 2/3 consensus node number: %d", headerBookKeeperNum, consensusNodeNum))
}

func ErrInvalidPublicKey(pubkey string) error {
	return sdkerrors.Wrap(ErrInvalidPublicKeyType, fmt.Sprintf("Invalid pubkey: %s", pubkey))
}

func ErrVerifyMultiSigFail(err error, height uint32) error {
	return sdkerrors.Wrap(ErrVerifyMultiSigFailType, fmt.Sprintf("Verify multi signature Error: %s of height: %d", err.Error(), height))
}
