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

// Btcx module event types
const (
	AttributeValueCategory = ModuleName

	EventTypeCreateDenom         = "create_denom"
	AttributeKeySourceAssetDenom = "source_asset_denom"
	AttributeKeyFromAssetHash    = "from_asset_hash"
	AttributeKeyRedeemScript     = "redeem_script"

	EventTypeBindAsset      = "bind_asset_hash"
	AttributeKeyCreator     = "creator"
	AttributeKeyToChainId   = "to_chain_id"
	AttributeKeyToAssetHash = "to_asset_hash"

	EventTypeLock           = "lock"
	AttributeKeyFromAddress = "from_address"
	AttributeKeyToAddress   = "to_address"
	AttributeKeyAmount      = "amount"

	EventTypeUnlock = "unlock"
)
