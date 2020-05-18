package rpc

import (
	Err "github.com/dad-go/http/base/error"
)

var (
	dad-goRpcInvalidHash        = responsePacking(Err.INVALID_PARAMS, "invalid hash")
	dad-goRpcInvalidBlock       = responsePacking(Err.INVALID_BLOCK, "invalid block")
	dad-goRpcInvalidTransaction = responsePacking(Err.INVALID_TRANSACTION, "invalid transaction")
	dad-goRpcInvalidParameter   = responsePacking(Err.INVALID_PARAMS, "invalid parameter")

	dad-goRpcUnknownBlock       = responsePacking(Err.UNKNOWN_BLOCK, "unknown block")
	dad-goRpcUnknownTransaction = responsePacking(Err.UNKNOWN_TRANSACTION, "unknown transaction")

	dad-goRpcNil             = responsePacking(Err.INVALID_PARAMS, nil)
	dad-goRpcUnsupported     = responsePacking(Err.INTERNAL_ERROR, "Unsupported")
	dad-goRpcInternalError   = responsePacking(Err.INTERNAL_ERROR, "internal error")
	dad-goRpcIOError         = responsePacking(Err.INTERNAL_ERROR, "internal IO error")
	dad-goRpcAPIError        = responsePacking(Err.INTERNAL_ERROR, "internal API error")
	dad-goRpcSuccess         = responsePacking(Err.SUCCESS, true)
	dad-goRpcFailed          = responsePacking(Err.INTERNAL_ERROR, false)
	dad-goRpcAccountNotFound = responsePacking(Err.INTERNAL_ERROR, "Account not found")

	dad-goRpc = responseSuccess
)

func responseSuccess(result interface{}) map[string]interface{} {
	return responsePacking(Err.SUCCESS, result)
}
func responsePacking(errcode int64, result interface{}) map[string]interface{} {
	resp := map[string]interface{}{
		"error":  errcode,
		"desc":   Err.ErrMap[errcode],
		"result": result,
	}
	return resp
}
