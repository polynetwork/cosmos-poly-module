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

package rest

import (
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/polynetwork/cosmos-poly-module/ccm/internal/types"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/ccm/process_crosschain_tx", ProcessCrossChainTxRequestHandlerFn(cliCtx)).Methods("POST")

}

// SendReq defines the properties of a send request's body.
type ProcessCrossChainTxReq struct {
	BaseReq     rest.BaseReq `json:"base_req" yaml:"base_req"`
	FromChainId uint64       `json:"from_chain_id" yaml:"from_chain_id"`
	Proof       string       `json:"proof" yaml:"proof"`
	Header      string       `json:"header" yaml:"header"`
	HeaderProof string       `json:"header_proof" yaml:"header_proof"`
	CurHeader   string       `json:"cur_header" yaml:"cur_header"`
}

// SendRequestHandlerFn - http request handler to send coins to a address.
func ProcessCrossChainTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req ProcessCrossChainTxReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}
		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		msg := types.NewMsgProcessCrossChainTx(fromAddr, req.FromChainId, req.Proof, req.Header, req.HeaderProof, req.CurHeader)
		utils.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
