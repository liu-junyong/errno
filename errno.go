package errno

import (
	"fmt"
)

const (
	InternalErrorLimit = 10007
)

const (
	ErrCodeBadRequest          = 400
	ErrCodeForbidden           = 403
	ErrCodeValidateErr         = 422
	ErrCodeTooManayRequest     = 249
	ErrCodeInternalServerError = 500
	ErrCodeInvalidHTTPMethod   = 405
)

// api、service 错误码尽量收敛到此
var (
	errMap = make(map[int32]*ErrNo, 0)
	OK     = RegisterErrNo(0, "success")

	// 10000 ~ 10999 服务内部错误
	Unkown              = RegisterErrNo(10000, "Server Internal Error")
	DBNotFound          = RegisterErrNo(10001, "Record Not Found")
	RPCFailed           = RegisterErrNo(10002, "RPC Failed")
	RPCResConvertFailed = RegisterErrNo(10003, "RPC Res Convert failed")
	JSONMarshalFailed   = RegisterErrNo(10004, "JSON Marshal Failed")
	JSONUnmarshalFailed = RegisterErrNo(10005, "JSON Unmarshal failed")
	DBError             = RegisterErrNo(10006, "DB Error")
	PackError           = RegisterErrNo(10007, "Pack Error")
	GetLockFailed       = RegisterErrNo(10008, "Get Lock Failed")

	// 11000 ~ 11999 通用业务错误
	ParamWrong     = RegisterErrNo(11001, "invalid param")
	SessionExpired = RegisterErrNo(11002, "no_login")
	NoPermission   = RegisterErrNo(11003, "no_permission")
	DupOper        = RegisterErrNo(11004, "dup operation")
	GenIDFailed    = RegisterErrNo(11005, "gen id failed")
	MustLogin      = RegisterErrNo(11006, "Must Login")

	// 12000 ~ 12111
	InvalidSymbol     = RegisterErrNo(12000, "invalid symbol")
	SymbolLoadFailed  = RegisterErrNo(12001, "unable to get snapshot E3020")
	AccountLoadFailed = RegisterErrNo(12001, "unable to get snapshot E3021")

	// HTTP Error不注册
	HTTPAuthMissing             = New(400, "Authorization header is missing")
	HTTPTokenBearerMissing      = New(400, "Bearer is missing")
	HTTPTokenBearerBase64       = New(400, "Bearer is not properly encoded in base64")
	HTTPInvalidToken            = New(400, "invalid token")
	HTTPConfirmedDeviceNotMatch = New(400, "confirmed device not match")

	HTTPTokenDecodeError = New(400, "token decode failed")
)

// ErrNo is the core object of the package.
type ErrNo struct {
	StatusCode    int32
	StatusMessage string
}

// Error is a function of error interface.
func (en *ErrNo) Error() string {
	return fmt.Sprintf("StatusCode: %d, StatusMessage: %s",
		en.StatusCode, en.StatusMessage)
}

// New returns an instance of ErrNo object with the given arguments.
func New(statusCode int32, statusMessage string) *ErrNo {
	return &ErrNo{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	}
}

// CopyWithPrompt returns a new ErrNo instance with the prompt argument.
func (en *ErrNo) CopyWithPrompt(prompt string) *ErrNo {
	return &ErrNo{
		StatusCode:    en.StatusCode,
		StatusMessage: en.StatusMessage,
	}
}

// RegisterErrNo creates an instance of ErrNo and stores it in errMap before returns it.
func RegisterErrNo(statusCode int32, statusMessage string) *ErrNo {
	err := New(statusCode, statusMessage)
	errMap[statusCode] = err
	return err
}

// GetErrNo gets the ErrNo instance in errMap and returns it.
func GetErrNo(statusCode int32) *ErrNo {
	if errNo, ok := errMap[statusCode]; ok {
		return errNo
	}
	return Unkown
}

// GetApiErrNo: 屏蔽一些内部错误, 避免抛给用户
func GetApiErrNo(statusCode int32) *ErrNo {
	if statusCode == 0 {
		return nil
	}
	if statusCode <= int32(InternalErrorLimit) {
		return Unkown
	}
	if errNo, ok := errMap[statusCode]; ok {
		return errNo
	}
	return Unkown
}

// GuessErr returns a ErrNo instance by judging the err argument.
func GuessErr(err interface{}) *ErrNo {
	if err == nil {
		return OK
	}
	switch e := err.(type) {
	case *ErrNo:
		return e
	default:
		return Unkown
	}
}
