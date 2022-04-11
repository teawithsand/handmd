package fsal

import "errors"

var ErrReadOnlyFS = errors.New("handmd/util/fsal: read only fs")
var ErrOperationNotSupported = errors.New("handmd/util/fsal: given fs or file does not support this operation")
