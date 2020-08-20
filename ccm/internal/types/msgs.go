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
)

// Governance message types and routes
const (
	TypeMsgProcessCrossChainTx = "process_cross_chain_tx"
	TypeMsgCreateCoins         = "create_coins"
)

type MsgProcessCrossChainTx struct {
	Submitter   sdk.AccAddress // transaction submitter
	FromChainId uint64         // the poly chain id
	Proof       string         // the audit path of cross chain transaction where the root is Header.CrossStateRoot
	Header      string         // the header of height where the cross chain transaction appears
	HeaderProof string         // the audit path of Header where the reliable root is CurHeader.BlockRoot
	CurHeader   string         // any header within current consensus epoch
}

func NewMsgProcessCrossChainTx(submitter sdk.AccAddress, fromChainId uint64, proof, header, headerProof, curHeader string) MsgProcessCrossChainTx {
	return MsgProcessCrossChainTx{submitter, fromChainId, proof, header, headerProof, curHeader}
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
		return ErrMsgProcessCrossChainTx(fmt.Sprintf("MsgCrossChaintx.Proof should not be empty"))
	}
	if len(msg.Header) == 0 {
		return ErrMsgProcessCrossChainTx(fmt.Sprintf("MsgCrossChainTx.Header should not be empty"))
	}
	return nil
}

func (msg MsgProcessCrossChainTx) String() string {
	return fmt.Sprintf(`Process Cross Chain Tx Message:
  Submitter:       		%s
  FromChainId: 			%d
  Proof:    			%s
  Header: 				%s
  HeaderProof: 			%s
  CurHeader:			%s
`, msg.Submitter.String(), msg.FromChainId, msg.Proof, msg.Header, msg.HeaderProof, msg.CurHeader)
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
	if msg.ToChainID == 0 {
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
