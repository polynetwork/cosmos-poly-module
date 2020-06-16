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

// Minting module event types
const (
	AttributeValueCategory = ModuleName

	EventTypeSyncHeader           = "sync_header"
	AttributeKeyChainId           = "chain_id"
	AttributeKeyHeight            = "height"
	AttributeKeyBlockHash         = "block_hash"
	AttributeKeyNativeChainHeight = "native_chain_height"

	EventTypeCreateCoin                   = "create_coin"
	EventTypeCreateAndDelegateCoinToProxy = "create_and_delegate_coin_to_proxy"
	AttributeKeyToChainId                 = "to_chain_id"
	AttributeKeyToChainProxyHash          = "to_chain_proxy_hash"

	EventTypeBindAsset           = "bind_asset_hash"
	AttributeKeySourceAssetDenom = "source_asset_denom"
	AttributeKeyCreator          = "creator"
	AttributeKeyFromAssetHash    = "from_asset_hash"
	AttributeKeyToChainAssetHash = "to_chain_asset_hash"
	AttributeKeyInitialAmt       = "initial_amt"
	EventTypeLock                = "lock"
	AttributeKeyFromAddress      = "from_address"
	AttributeKeyToAddress        = "to_address"
	AttributeKeyAmount           = "amount"
	AttributeKeyCoins            = "coins"

	EventTypeCreateCrossChainTx = "make_from_cosmos_proof"
	AttributeCrossChainId       = "cross_chainId"
	AttributeKeyTxParamHash     = "make_tx_param_hash"
	AttributeKeyMakeTxParam     = "make_tx_param"

	EventTypeVerifyToCosmosProof                        = "verify_to_cosmos_proof"
	AttributeKeyMerkleValueTxHash                       = "merkle_value.txhash"
	AttributeKeyMerkleValueMakeTxParamTxHash            = "merkle_value.make_tx_param.txhash"
	AttributeKeyMerkleValueMakeTxParamToContractAddress = "merkle_value.make_tx_param.to_contract_address"
	AttributeKeyFromChainId                             = "from_chain_id"

	EventTypeUnlock              = "unlock"
	AttributeKeyFromContractHash = "from_contract_hash"
	AttributeKeyToAssetDenom     = "to_asset_denom"

	EventTypeSetRedeemScript = "set_redeem_script"
	AttributeKeyRedeemKey    = "redeem_key"
	AttributeKeyRedeemScript = "redeem_script"
)
