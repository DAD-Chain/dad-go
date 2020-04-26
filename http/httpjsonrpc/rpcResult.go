package httpjsonrpc

var (
	dad-goRpcInvalidHash = responsePacking("invalid hash")
	dad-goRpcInvalidBlock = responsePacking("invalid block")
	dad-goRpcInvalidTransaction = responsePacking("invalid transaction")
	dad-goRpcInvalidParameter = responsePacking("invalid parameter")

	dad-goRpcUnknownBlock = responsePacking("unknown block")
	dad-goRpcUnknownTransaction = responsePacking("unknown transaction")

	dad-goRpcNil = responsePacking(nil)
	dad-goRpcUnsupported = responsePacking("Unsupported")
	dad-goRpcInternalError = responsePacking("internal error")
	dad-goRpcIOError = responsePacking("internal IO error")
	dad-goRpcAPIError = responsePacking("internal API error")
	dad-goRpcSuccess = responsePacking(true)
	dad-goRpcFailed = responsePacking(false)
	dad-goRpcAccountNotFound = responsePacking(("Account not found"))

	dad-goRpc = responsePacking
)
