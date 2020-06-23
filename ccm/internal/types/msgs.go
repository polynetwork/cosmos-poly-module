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

	sdk "github.com/cosmos/cosmos-sdk/types"

	"encoding/hex"
)

// Governance message types and routes
const (
	TypeMsgProcessCrossChainTx = "process_cross_chain_tx"
	TypeMsgCreateCoins         = "create_coins"
)

type MsgProcessCrossChainTx struct {
	Submitter   sdk.AccAddress
	FromChainId uint64
	Height      uint32
	Proof       string
	Header      []byte
}

func NewMsgProcessCrossChainTx(submitter sdk.AccAddress, fromChainId uint64, height uint32, proof string, header []byte) MsgProcessCrossChainTx {
	return MsgProcessCrossChainTx{submitter, fromChainId, height, proof, header}
}

//nolint
func (msg MsgProcessCrossChainTx) Route() string { return RouterKey }
func (msg MsgProcessCrossChainTx) Type() string  { return TypeMsgProcessCrossChainTx }

// Implements Msg.
func (msg MsgProcessCrossChainTx) ValidateBasic() error {
	if msg.Submitter.Empty() {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "MsgProcessCrossChainTx.Submitter is empty")
	}
	if len(msg.Proof) == 0 {
		// Disable software upgrade proposals as they are currently equivalent
		// to text proposals. Re-enable once a valid software upgrade proposal
		// handler is implemented.
		return ErrMsgProcessCrossChainTx(fmt.Sprintf("MsgCrossChaintx.Proof should not be empty"))
	}
	if len(msg.Header) == 0 {
		return ErrMsgProcessCrossChainTx(fmt.Sprintf("MsgCrossChainTx.Header should not be empty"))
	}
	return nil
}

func (msg MsgProcessCrossChainTx) String() string {
	return fmt.Sprintf(`Bind Proxy Hash Message:
  Submitter:       		%s
  FromChainId: 			%d
  Height:  				%d
  Proof:    			%s
  Header: 				%s
`, msg.Submitter.String(), msg.FromChainId, msg.Height, msg.Proof, hex.EncodeToString(msg.Header))
}

// Implements Msg.
func (msg MsgProcessCrossChainTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgProcessCrossChainTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Submitter}
}

type MsgCreateCrossChainTx struct {
	ToChainID         uint64
	ToContractAddress []byte
	Method            string
	Args              []byte
}

func NewMsgCreateCrossChainTx(toChainId uint64, toContractAddr []byte, method string, args []byte) MsgCreateCrossChainTx {
	return MsgCreateCrossChainTx{ToChainID: toChainId, ToContractAddress: toContractAddr, Method: method, Args: args}
}

//nolint
func (msg MsgCreateCrossChainTx) Route() string { return RouterKey }
func (msg MsgCreateCrossChainTx) Type() string  { return TypeMsgCreateCoins }

// Implements Msg.
func (msg MsgCreateCrossChainTx) ValidateBasic() error {
	if msg.ToChainID <= 0 {
		return ErrMsgCreateCrossChainTx(fmt.Sprintf("invalid chainId: %d", msg.ToChainID))
	}
	if len(msg.ToContractAddress) == 0 {
		return ErrMsgCreateCrossChainTx("ToContractAddress is empty")
	}
	if msg.Method == "" {
		return ErrMsgCreateCrossChainTx("Method is empty")
	}
	if len(msg.Args) == 0 {
		return ErrMsgCreateCrossChainTx("Args is empty")
	}

	return nil
}

func (msg MsgCreateCrossChainTx) String() string {
	return fmt.Sprintf(`Create Coins Message:
  ToChainID:         		%d
  ToContractAddress: 		%x
  Method: 					%s
  Args:						%x
`, msg.ToChainID, msg.ToContractAddress, msg.Method, msg.Args)
}

// Implements Msg.
func (msg MsgCreateCrossChainTx) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// Implements Msg.
func (msg MsgCreateCrossChainTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{}
}
