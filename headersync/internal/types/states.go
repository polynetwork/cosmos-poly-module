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
	"sort"
	"strconv"
)

type Peer struct {
	Index      uint32
	PeerPubkey string
}

func (this *Peer) Serialization(sink *polycommon.ZeroCopySink) {
	sink.WriteUint32(this.Index)
	sink.WriteVarBytes([]byte(this.PeerPubkey))
}

func (this *Peer) Deserialization(source *polycommon.ZeroCopySource) error {
	index, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("utils.DecodeVarUint, deserialize index error")
	}
	peerPubkey, eof := source.NextString()
	if eof {
		return fmt.Errorf("utils.DecodeString, deserialize peerPubkey error")
	}
	this.Index = uint32(index)
	this.PeerPubkey = peerPubkey
	return nil
}

type ConsensusPeers struct {
	ChainID uint64
	Height  uint32
	PeerMap map[string]*Peer
}

func (this *ConsensusPeers) Serialization(sink *polycommon.ZeroCopySink) {
	sink.WriteUint64(this.ChainID)
	sink.WriteUint32(this.Height)
	sink.WriteVarUint(uint64(len(this.PeerMap)))
	var peerList []*Peer
	for _, v := range this.PeerMap {
		peerList = append(peerList, v)
	}
	sort.SliceStable(peerList, func(i, j int) bool {
		return peerList[i].PeerPubkey > peerList[j].PeerPubkey
	})
	for _, v := range peerList {
		v.Serialization(sink)
	}
}

func (this *ConsensusPeers) Deserialization(source *polycommon.ZeroCopySource) error {
	chainID, eof := source.NextUint64()
	if eof {
		return fmt.Errorf("utils.DecodeVarUint, deserialize chainID error")
	}
	height, eof := source.NextUint32()
	if eof {
		return fmt.Errorf("utils.DecodeVarUint, deserialize height error")
	}
	n, eof := source.NextVarUint()
	if eof {
		return fmt.Errorf("utils.DecodeVarUint, deserialize HeightList length error")
	}
	peerMap := make(map[string]*Peer)
	for i := 0; uint64(i) < n; i++ {
		peer := new(Peer)
		if err := peer.Deserialization(source); err != nil {
			return fmt.Errorf("deserialize peer error: %v", err)
		}
		peerMap[peer.PeerPubkey] = peer
	}
	this.ChainID = chainID
	this.Height = uint32(height)
	this.PeerMap = peerMap
	return nil
}

func (this *ConsensusPeers) String() string {
	var peerList []*Peer
	for _, v := range this.PeerMap {
		peerList = append(peerList, v)
	}
	sort.SliceStable(peerList, func(i, j int) bool {
		return peerList[i].PeerPubkey > peerList[j].PeerPubkey
	})
	var peerMapStr string
	for _, v := range peerList {
		peerMapStr += "\t\t\t\t\t\t" + strconv.FormatUint(uint64(v.Index), 10) + ":" + v.PeerPubkey + "\n"
	}
	return fmt.Sprintf(`
	ChainID          : %d
	Height           : %d
	PeerMap		     : 
%s	
`, this.ChainID, this.Height, fmt.Sprintf("%s", peerMapStr))
}
