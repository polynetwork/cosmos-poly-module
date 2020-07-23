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

	polycommon "github.com/polynetwork/poly/common"
)

type RegisterAssetTxArgs struct {
	AssetHash       []byte
	NativeAssetHash []byte
}

func (this *RegisterAssetTxArgs) Serialization(sink *polycommon.ZeroCopySink) error {
	sink.WriteVarBytes(this.AssetHash)
	sink.WriteVarBytes(this.NativeAssetHash)
	return nil
}

func (this *RegisterAssetTxArgs) Deserialization(source *polycommon.ZeroCopySource) error {
	assetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize assetHash error")
	}
	nativeAssetHash, eof := source.NextVarBytes()
	if eof {
		return fmt.Errorf("TxArgs deserialize nativeAssetHash error")
	}

	this.AssetHash = assetHash
	this.NativeAssetHash = nativeAssetHash
	return nil
}
