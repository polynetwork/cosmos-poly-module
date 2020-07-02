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
	"github.com/polynetwork/poly/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeer_Serialization(t *testing.T) {
	paramSerialize := new(Peer)
	paramSerialize.Index = 1
	paramSerialize.PeerPubkey = "abcdefg"
	sink := common.NewZeroCopySink(nil)
	paramSerialize.Serialization(sink)

	paramDeserialize := new(Peer)
	err := paramDeserialize.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, paramDeserialize, paramSerialize)
}

func TestConsensusPeers_Serialization(t *testing.T) {
	paramSerialize := new(ConsensusPeers)
	paramSerialize.Height = 1
	paramSerialize.ChainID = 0
	peer1 := &Peer{Index: 1, PeerPubkey: "abcd"}
	peer2 := &Peer{Index: 2, PeerPubkey: "efgh"}
	paramSerialize.PeerMap = make(map[string]*Peer)
	paramSerialize.PeerMap[peer1.PeerPubkey] = peer1
	paramSerialize.PeerMap[peer2.PeerPubkey] = peer2
	sink := common.NewZeroCopySink(nil)
	paramSerialize.Serialization(sink)

	paramDeserialize := new(ConsensusPeers)
	err := paramDeserialize.Deserialization(common.NewZeroCopySource(sink.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, paramDeserialize, paramSerialize)

	s1 := common.NewZeroCopySink(nil)
	p := *paramSerialize
	p.Serialization(s1)
	p2 := new(ConsensusPeers)
	err = p2.Deserialization(common.NewZeroCopySource(s1.Bytes()))
	assert.Nil(t, err)
	assert.Equal(t, p, *p2)
}
