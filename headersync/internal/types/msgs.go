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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Governance message types and routes
const (
	TypeMsgSyncGenesis = "sync_genesis"
	TypeMsgSyncHeaders = "sync_headers"
)

// MsgSend - high level transaction of the coin module
type MsgSyncGenesisParam struct {
	Syncer        sdk.AccAddress
	GenesisHeader []byte
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSyncGenesisParam(syncer sdk.AccAddress, genesisHeader []byte) MsgSyncGenesisParam {
	return MsgSyncGenesisParam{Syncer: syncer, GenesisHeader: genesisHeader}
}

// Route Implements Msg.
func (msg MsgSyncGenesisParam) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgSyncGenesisParam) Type() string { return TypeMsgSyncGenesis }

// ValidateBasic Implements Msg.
func (msg MsgSyncGenesisParam) ValidateBasic() error {
	if msg.Syncer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("address:%s", msg.Syncer.String()))
	}
	if len(msg.GenesisHeader) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing GenesisHeader bytes")
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgSyncGenesisParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Syncer}
}

// GetSignBytes Implements Msg.
func (msg MsgSyncGenesisParam) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgMultiSend - high level transaction of the coin module
type MsgSyncHeadersParam struct {
	Syncer  sdk.AccAddress
	Headers [][]byte
}

// NewMsgMultiSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgSyncHeadersParam(syncer sdk.AccAddress, headers [][]byte) MsgSyncHeadersParam {
	return MsgSyncHeadersParam{Syncer: syncer, Headers: headers}
}

// Route Implements Msg
func (msg MsgSyncHeadersParam) Route() string { return RouterKey }

// Type Implements Msg
func (msg MsgSyncHeadersParam) Type() string { return TypeMsgSyncHeaders }

// ValidateBasic Implements Msg.
func (msg MsgSyncHeadersParam) ValidateBasic() error {
	if msg.Syncer.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, fmt.Sprintf("address:%s", msg.Syncer.String()))
	}
	if len(msg.Headers) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing BlockHeaders bytes")
	}
	return nil
}

// GetSignBytes Implements Msg.
func (msg MsgSyncHeadersParam) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners Implements Msg.
func (msg MsgSyncHeadersParam) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Syncer}
}
