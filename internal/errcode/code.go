package errcode

const (
	Success ErrCode = 0

	ServerError ErrCode = 10000
	NotFound    ErrCode = 10001

	InvalidInput      ErrCode = 20000
	UnknownClient     ErrCode = 20001
	Unauthorized      ErrCode = 20100
	AuthorizedExpires ErrCode = 20101
	PermissionDenied  ErrCode = 20102
	Unavailable       ErrCode = 20103
	InvalidAccount    ErrCode = 20104

	SystemError      ErrCode = 50000
	DatabaseError    ErrCode = 50001
)
