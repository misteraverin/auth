package errdomain

type ConstError string

func (e ConstError) Error() string {
	return string(e)
}

const ErrNotCorrectBasicAuth = ConstError("not correct basic auth")
const ErrWrongLoginOrPassword = ConstError("wrong login or password")
const ErrUserIsNotExist = ConstError("user is not exist")

const ErrTokenExpired = ConstError("token expired")
const ErrTokenInvalid = ConstError("token invalid")
