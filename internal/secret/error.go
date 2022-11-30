package secret

import "errors"

var ErrSecretNotFound = errors.New("secret not found")
var ErrUserHaveNotSecret = errors.New("not secret for this user")
