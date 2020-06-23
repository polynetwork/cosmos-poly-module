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
	"encoding/hex"
	"fmt"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/polynetwork/cosmos-poly-module/ccm/client/common"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, queryRoute string) {
	r.HandleFunc(
		fmt.Sprintf("/ccm/if_contain_contract/{%s}/{%s}/{%s}", ModuleStoreKey, ToContract, FromChainId),
		queryIfContainContract(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/ccm/parameters"),
		queryParams(cliCtx, queryRoute),
	).Methods("GET")

	r.HandleFunc(
		fmt.Sprintf("/ccm/module_balance/{%s}", ModuleName),
		queryModuleBalance(cliCtx, queryRoute),
	).Methods("GET")
}

func queryIfContainContract(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		// query for rewards from a particular delegator
		vars := mux.Vars(r)

		toContract, err := hex.DecodeString(vars[ToContract])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		fromChainId, err := strconv.ParseUint(vars[FromChainId], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		res, ok := checkResponseQueryIfContainContractResponse(w, cliCtx, queryRoute, vars[ModuleStoreKey], toContract, fromChainId)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryIfContainContractResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string, keyStore string, toContract []byte, fromChainId uint64) (res []byte, ok bool) {

	res, err := common.QueryIfContainContract(cliCtx, queryRoute, keyStore, toContract, fromChainId)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}

func queryParams(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		res, ok := checkResponseQueryParamsResponse(w, cliCtx, queryRoute)
		if !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkResponseQueryParamsResponse(
	w http.ResponseWriter, cliCtx context.CLIContext, queryRoute string) (res []byte, ok bool) {

	res, err := common.QueryParams(cliCtx, queryRoute)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return nil, false
	}

	return res, true
}

func queryModuleBalance(cliCtx context.CLIContext, queryRoute string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}
		vars := mux.Vars(r)
		res, err := common.QueryModuleBalance(cliCtx, queryRoute, vars[ModuleName])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
