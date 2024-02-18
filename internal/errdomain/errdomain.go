package errdomain

import "errors"

var ErrWrongLoginOrPassword = errors.New("wrong login or password")
var ErrUserIsNotExist = errors.New("user is not exist")

var ErrTokenExpired = errors.New("token expired")
var ErrTokenInvalid = errors.New("token invalid")
