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
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgBindAssetHash = "bind_asset_hash"
	TypeMsgLock          = "lock"
	TypeMsgCreateDenom   = "create_coin"
)

type MsgCreateDenom struct {
	Creator      sdk.AccAddress
	Denom        string
	RedeemScript string
}

func NewMsgCreateDenom(creator sdk.AccAddress, denom string, redeemScript string) MsgCreateDenom {
	return MsgCreateDenom{Creator: creator, Denom: denom, RedeemScript: redeemScript}
}

//nolint
func (msg MsgCreateDenom) Route() string { return RouterKey }
func (msg MsgCreateDenom) Type() string  { return TypeMsgCreateDenom }

// Implements Msg.
func (msg MsgCreateDenom) ValidateBasic() error {
	if msg.Creator.Empty() {
		return ErrCreateDenom("MsgCreateDenom.Creator is empty")
	}
	if err := sdk.ValidateDenom(msg.Denom); err != nil {
		return ErrCreateDenom(fmt.Sprintf("MsgCreateDenom.Denom: %s is illegal, Error:%v", msg.Denom, err))
	}
	if _, err := hex.DecodeString(msg.RedeemScript); err != nil {
		return ErrCreateDenom(fmt.Sprintf("MsgCreateDenom.RedeemScript: %s is not hex string format, Error:%v", msg.RedeemScript, err))
	}
	return nil
}

func (msg MsgCreateDenom) String() string {
	return fmt.Sprintf(`MsgCreateDenom:
  Creator:         %s
  Denom: 		   %s
  RedeemScript:    %s
`, msg.Creator.String(), msg.Denom, msg.RedeemScript)
}

// Implements Msg.
func (msg MsgCreateDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateDenom) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgBindAssetHash struct {
	Creator          sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAssetHash      []byte
}

func NewMsgBindAssetHash(creator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) MsgBindAssetHash {
	return MsgBindAssetHash{creator, sourceAssetDenom, toChainId, toAssetHash}
}

//nolint
func (msg MsgBindAssetHash) Route() string { return RouterKey }
func (msg MsgBindAssetHash) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindAssetHash) ValidateBasic() error {
	if msg.Creator.Empty() {
		return ErrBindAssetHash(fmt.Sprintf("Empty address:%s", msg.Creator.String()))
	}
	if err := sdk.ValidateDenom(msg.SourceAssetDenom); err != nil {
		return ErrCreateDenom(fmt.Sprintf("MsgBindAssetHash.SourceAssetDenom: %s is invalid, err: %v", msg.SourceAssetDenom, err))
	}
	if msg.ToChainId == 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToAssetHash) == 0 {
		return ErrEmptyToAssetHash(hex.EncodeToString(msg.ToAssetHash))
	}
	return nil
}

func (msg MsgBindAssetHash) String() string {
	return fmt.Sprintf(`MsgBindAssetHash:
  Creator:          %s
  SourceAssetDenom: %s
  TargetChainId:    %d
  TargetAssetHash:  %x
`, msg.Creator.String(), msg.SourceAssetDenom, msg.ToChainId, msg.ToAssetHash)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

type MsgLock struct {
	FromAddress      sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAddressBs      []byte
	Value            sdk.Int
}

func NewMsgLock(fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddress []byte, value sdk.Int) MsgLock {
	return MsgLock{fromAddress, sourceAssetDenom, toChainId, toAddress, value}
}

//nolint
func (msg MsgLock) Route() string { return RouterKey }
func (msg MsgLock) Type() string  { return TypeMsgLock }

// Implements Msg.
func (msg MsgLock) ValidateBasic() error {
	if msg.FromAddress.Empty() {
		return ErrLock(fmt.Sprintf("MsgLock.FromAddress is empty, address:%s", msg.FromAddress.String()))
	}
	if err := sdk.ValidateDenom(msg.SourceAssetDenom); err != nil {
		return ErrLock(fmt.Sprintf("MsgLock.SourceAssetDenom: %s is invalid, err: %v", msg.SourceAssetDenom, err))
	}
	if msg.ToChainId == 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToAddressBs) == 0 {
		return ErrEmptyToAssetHash(hex.EncodeToString(msg.ToAddressBs))
	}
	if msg.Value.IsNegative() {
		return errors.New("bind asset param limit should be positive")
	}
	return nil
}

func (msg MsgLock) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  FromAddress:         %s
  SourceAssetDenom: %s
  ToChainId:  %d
  ToAddress:     %s
  Value: %s
`, msg.FromAddress.String(), msg.SourceAssetDenom, msg.ToChainId, hex.EncodeToString(msg.ToAddressBs), msg.Value.String())
}

// Implements Msg.
func (msg MsgLock) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgLock) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.FromAddress}
}
