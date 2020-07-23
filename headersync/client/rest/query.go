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
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/polynetwork/cosmos-poly-module/headersync/client/common"
	"strconv"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/headersync/current_consensus_peers/{%s}", ChainId),
		queryCurrentCPHandlerFn(cliCtx, queryRoute),
	).Methods("GET")
}

func queryCurrentCPHandlerFn(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)
		chainIdStr := vars[ChainId]
		chainId, err := strconv.ParseUint(chainIdStr, 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, ok := checkResponseQueryCPResponse(w, cliCtx, queryRoute, chainId)
		if !ok {
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryCPResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, chainId uint64) (res []byte, ok bool) {

	res, err := common.QueryConsensusPeers(cliCtx, queryRoute, chainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}
	return res, true
}
