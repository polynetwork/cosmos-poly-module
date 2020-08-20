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

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgCreateLockProxy              = "create_lock_proxy"
	TypeMsgCreateCoinAndDelegateToProxy = "create_delegate_to_proxy"
	TypeMsgBindProxyHash                = "bind_proxy_hash"
	TypeMsgBindAssetHash                = "bind_asset_hash"
	TypeMsgLock                         = "lock"
)

// MsgSend - high level transaction of the coin module
type MsgCreateLockProxy struct {
	Creator sdk.AccAddress
}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgCreateLockProxy(creator sdk.AccAddress) MsgCreateLockProxy {
	return MsgCreateLockProxy{creator}
}

// Route Implements Msg.
func (msg MsgCreateLockProxy) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateLockProxy) Type() string { return TypeMsgCreateLockProxy }

// ValidateBasic Implements Msg.
func (msg MsgCreateLockProxy) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgCreateLockProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// GetSignBytes Implements Msg.
func (msg MsgCreateLockProxy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// MsgSend - high level transaction of the coin module
type MsgCreateCoinAndDelegateToProxy struct {
	Creator       sdk.AccAddress
	Coin          sdk.Coin
	LockProxyHash []byte
}

var _ sdk.Msg = MsgCreateCoinAndDelegateToProxy{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgCreateCoinAndDelegateToProxy(creator sdk.AccAddress, coin sdk.Coin, lockProxyHash []byte) MsgCreateCoinAndDelegateToProxy {
	return MsgCreateCoinAndDelegateToProxy{Creator: creator, Coin: coin, LockProxyHash: lockProxyHash}
}

// Route Implements Msg.
func (msg MsgCreateCoinAndDelegateToProxy) Route() string { return RouterKey }

// Type Implements Msg.
func (msg MsgCreateCoinAndDelegateToProxy) Type() string { return TypeMsgCreateCoinAndDelegateToProxy }

// ValidateBasic Implements Msg.
func (msg MsgCreateCoinAndDelegateToProxy) ValidateBasic() error {
	if msg.Creator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if !msg.Coin.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, msg.Coin.String())
	}
	return nil
}

// GetSigners Implements Msg.
func (msg MsgCreateCoinAndDelegateToProxy) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Creator}
}

// GetSignBytes Implements Msg.
func (msg MsgCreateCoinAndDelegateToProxy) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

type MsgBindProxyHash struct {
	Operator         sdk.AccAddress
	ToChainId        uint64
	ToChainProxyHash []byte
}

func NewMsgBindProxyHash(operator sdk.AccAddress, toChainId uint64, toChainProxyHash []byte) MsgBindProxyHash {
	return MsgBindProxyHash{operator, toChainId, toChainProxyHash}
}

//nolint
func (msg MsgBindProxyHash) Route() string { return RouterKey }
func (msg MsgBindProxyHash) Type() string  { return TypeMsgBindProxyHash }

// Implements Msg.
func (msg MsgBindProxyHash) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if msg.ToChainId == 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToChainProxyHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrMsgBindProxyHash("Empty MsgBindProxyHash.ToChainProxyHash")
	}

	return nil
}

func (msg MsgBindProxyHash) String() string {
	return fmt.Sprintf(`MsgBindProxyHash:
  Operator:       		%s(%x)
  ToChainId:			%d
  ToChainProxyHash:     %s
`, msg.Operator.String(), msg.Operator.Bytes(), msg.ToChainId, hex.EncodeToString(msg.ToChainProxyHash))
}

// Implements Msg.
func (msg MsgBindProxyHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindProxyHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

type MsgBindAssetHash struct {
	Operator         sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAssetHash      []byte
}

func NewMsgBindAssetHash(operator sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAssetHash []byte) MsgBindAssetHash {
	return MsgBindAssetHash{operator, sourceAssetDenom, toChainId, toAssetHash}
}

//nolint
func (msg MsgBindAssetHash) Route() string { return RouterKey }
func (msg MsgBindAssetHash) Type() string  { return TypeMsgBindAssetHash }

// Implements Msg.
func (msg MsgBindAssetHash) ValidateBasic() error {
	if msg.Operator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if err := sdk.ValidateDenom(msg.SourceAssetDenom); err != nil {
		return ErrMsgBindAssetHash(fmt.Sprintf("MsgBindAssetHash.SourceAssetDenom: %s is invalid, err: %v", msg.SourceAssetDenom, err))
	} else if _, err := sdk.ParseCoin("10" + msg.SourceAssetDenom); err != nil {
		return ErrMsgBindAssetHash(fmt.Sprintf("Invalid denom: %s", msg.SourceAssetDenom))
	}
	if msg.ToChainId == 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToAssetHash) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrMsgBindAssetHash("Empty MsgBindAssetHash.ToAssetHash")
	}
	return nil
}

func (msg MsgBindAssetHash) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Signer:         	%s
  SourceAssetDenom: %s
  ToChainId:  		%d
  ToAssetHash:      %s
`, msg.Operator.String(), msg.SourceAssetDenom, msg.ToChainId, hex.EncodeToString(msg.ToAssetHash))
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgBindAssetHash) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Operator}
}

type MsgLock struct {
	LockProxyHash    []byte
	FromAddress      sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToAddressBs      []byte
	Value            sdk.Int
}

func NewMsgLock(lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toAddress []byte, value sdk.Int) MsgLock {
	return MsgLock{lockProxyHash, fromAddress, sourceAssetDenom, toChainId, toAddress, value}
}

//nolint
func (msg MsgLock) Route() string { return RouterKey }
func (msg MsgLock) Type() string  { return TypeMsgLock }

// Implements Msg.
func (msg MsgLock) ValidateBasic() error {
	if len(msg.LockProxyHash) == 0 {
		return ErrMsgLock("empty MsgLock.LockProxyHash")
	}
	if msg.FromAddress.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if err := sdk.ValidateDenom(msg.SourceAssetDenom); err != nil {
		return ErrMsgLock(fmt.Sprintf("MsgBindAssetHash.SourceAssetDenom: %s is invalid, err: %v", msg.SourceAssetDenom, err))
	}
	if msg.ToChainId == 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToAddressBs) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrMsgLock("empty MsgLock.ToAssetHash")
	}
	if msg.Value.IsNegative() {
		return ErrMsgLock(fmt.Sprintf("MsgLock.Value: %s should not be negative", msg.Value.String()))
	}
	return nil
}

func (msg MsgLock) String() string {
	return fmt.Sprintf(`MsgLock:
  LockProxyHash: 		%x
  FromAddress:          %s
  SourceAssetDenom:     %s
  ToChainId:            %d
  ToAddress:            %x
  Value:                %s
`, msg.LockProxyHash, msg.FromAddress.String(), msg.SourceAssetDenom, msg.ToChainId, msg.ToAddressBs, msg.Value.String())
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
