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
	TypeMsgCreateLockProxy              = "create_lock_proxy"
	TypeMsgCreateCoinAndDelegateToProxy = "create_delegate_to_proxy"
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
	Creator             sdk.AccAddress
	Coin                sdk.Coin
	LockProxyHash       []byte
	NativeChainId       uint64
	NativeLockProxyHash []byte
	NativeAssetHash     []byte
}

var _ sdk.Msg = MsgCreateCoinAndDelegateToProxy{}

// NewMsgSend - construct arbitrary multi-in, multi-out send msg.
func NewMsgCreateCoinAndDelegateToProxy(creator sdk.AccAddress, coin sdk.Coin, lockProxyHash []byte, nativeChainId uint64, nativeLockProxyHash []byte, nativeAssetHash []byte) MsgCreateCoinAndDelegateToProxy {
	return MsgCreateCoinAndDelegateToProxy{Creator: creator, Coin: coin, LockProxyHash: lockProxyHash, NativeChainId: nativeChainId, NativeLockProxyHash: nativeLockProxyHash, NativeAssetHash: nativeAssetHash}
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
	if msg.NativeChainId == 0 {
		return ErrInvalidChainId(msg.NativeChainId)
	}
	if len(msg.NativeLockProxyHash) == 0 {
		return ErrCreateCoinAndDelegateToProxy("Empty NativeLockProxyHash")
	}
	if len(msg.NativeAssetHash) == 0 {
		return ErrCreateCoinAndDelegateToProxy("Empty NativeAssetHash")
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

type MsgLock struct {
	LockProxyHash    []byte
	FromAddress      sdk.AccAddress
	SourceAssetDenom string
	ToChainId        uint64
	ToChainProxyHash []byte
	ToChainAssetHash []byte
	ToAddressBs      []byte
	Value            sdk.Int
}

func NewMsgLock(lockProxyHash []byte, fromAddress sdk.AccAddress, sourceAssetDenom string, toChainId uint64, toChainProxyHash []byte, toChainAssetHash []byte, toAddress []byte, value sdk.Int) MsgLock {
	return MsgLock{lockProxyHash, fromAddress, sourceAssetDenom, toChainId, toChainProxyHash, toChainAssetHash, toAddress, value}
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
	if msg.SourceAssetDenom == "" {
		return ErrMsgLock("empty MsgBindAssetHash.SourceAssetDenom")
	}
	if msg.ToChainId <= 0 {
		return ErrInvalidChainId(msg.ToChainId)
	}
	if len(msg.ToChainProxyHash) == 0 {
		return ErrMsgLock("empty ToChainProxyHash")
	}
	if len(msg.ToChainAssetHash) == 0 {
		return ErrMsgLock("empty ToChainAssetHash")
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
