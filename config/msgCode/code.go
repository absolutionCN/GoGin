package msgCode

const (
	SUCCESS                        = 1000
	ERROR                          = 5000
	UN_LOGGING                     = 5001
	ERROR_NOT_EXIST_AUTH           = 5002
	ERROR_AUTH_TOKEN_CHECK_FAIL    = 5003
	ERROR_AUTH_TOKEN_CHECK_TIMEOUT = 5004
	INVALID_PARAMS                 = 10400
	ERROR_EXIST_TOKEN              = 10001
	ERROR_NOT_EXIST_TOKEN          = 10002
	ERROR_NOT_EXIST_TOKENID        = 10003
	ERROR_SYNC_YAPI                = 20001
	ERROR_SYNC_RUNNING             = 20002
	ERROR_NOT_EXIST_GROUP          = 30001
	ERROR_EXIST_AUTH               = 30002
	ERROR_EXIST_TAG                = 30003
	ERROR_NOT_EXIST_TAG            = 30004
	ERROR_NOT_EXIST_ARTICLE        = 30005
)

const (
	TaskSuccess = iota
	TaskError
	TaskRunning
)
