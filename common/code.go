package common

type Code int

const (
	CodeSuccess = 1000 + iota
	CodeServerError
	CodeInvalidPassword
	CodeUserExist
	CodeVoteTimeExpired
)

var MsgMap = map[Code]string{
	CodeSuccess:         "success",
	CodeServerError:     "server error",
	CodeInvalidPassword: "invalid username or password",
	CodeUserExist:       "user already exist",
	CodeVoteTimeExpired: "vote time expired",
}

func ToMsg(code Code) string {
	msg, ok := MsgMap[code]
	if !ok {
		return MsgMap[CodeServerError]
	}
	return msg
}
