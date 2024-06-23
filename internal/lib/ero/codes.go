package ero

const (
	CodeBadRequest           = 400
	CodeUnauthorized         = 401
	CodeExists               = 402
	CodePermissionDenied     = 403
	CodeNotFound             = 404
	CodeUnknownClient        = 499
	CodeInternal             = 500
	CodeUnimplemented        = 501
	CodeTemporaryUnavailable = 502
	CodeCancelled            = 503
	CodeUnknownServer        = 599
)

func ToGrpcCode(eroCode int) uint32 {
	switch eroCode {
	case CodeCancelled:
		return 1
	case CodeUnknownServer:
		return 2
	case CodeUnknownClient:
		return 2
	case CodeBadRequest:
		return 3
	case CodeNotFound:
		return 5
	case CodeExists:
		return 6
	case CodePermissionDenied:
		return 7
	case CodeUnimplemented:
		return 12
	case CodeInternal:
		return 13
	case CodeTemporaryUnavailable:
		return 14
	case CodeUnauthorized:
		return 16
	}
	return 2
}
