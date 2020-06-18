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

package ccm

import (
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/keeper"
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamspace = types.DefaultParamspace
	StoreKey          = types.StoreKey
	QuerierRoute      = types.QuerierRoute
	QueryParameters   = types.QueryParameters
	RouterKey         = types.RouterKey

	AttributeValueCategory = types.AttributeValueCategory

	EventTypeCreateCrossChainTx = types.EventTypeCreateCrossChainTx
	AttributeCrossChainId       = types.AttributeCrossChainId
	AttributeKeyTxParamHash     = types.AttributeKeyTxParamHash
	AttributeKeyMakeTxParam     = types.AttributeKeyMakeTxParam

	EventTypeVerifyToCosmosProof                        = types.EventTypeVerifyToCosmosProof
	AttributeKeyMerkleValueTxHash                       = types.AttributeKeyMerkleValueTxHash
	AttributeKeyMerkleValueMakeTxParamTxHash            = types.AttributeKeyMerkleValueMakeTxParamTxHash
	AttributeKeyMerkleValueMakeTxParamToContractAddress = types.AttributeKeyMerkleValueMakeTxParamToContractAddress
	AttributeKeyFromChainId                             = types.AttributeKeyFromChainId
)

var (
	// functions aliases
	RegisterCodec = types.RegisterCodec

	NewKeeper                 = keeper.NewKeeper
	NewQuerier                = keeper.NewQuerier
	NewGenesisState           = types.NewGenesisState
	DefaultGenesisState       = types.DefaultGenesisState
	ValidateGenesis           = types.ValidateGenesis
	NewMsgProcessCrossChainTx = types.NewMsgProcessCrossChainTx
	// key function
	GetCrossChainTxKey = keeper.GetCrossChainTxKey
	GetDoneTxKey       = keeper.GetDoneTxKey
	// variable aliases
	ModuleCdc   = types.ModuleCdc
	OperatorKey = types.OperatorKey
)

type (
	Keeper  = keeper.Keeper
	KeeperI = keeper.KeeperI

	MsgProcessCrossChainTx = types.MsgProcessCrossChainTx
	UnlockKeeper           = types.UnlockKeeper
	GenesisState           = types.GenesisState
	Params                 = types.Params
)
