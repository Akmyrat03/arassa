package errlst

import "errors"

const (
	ErrInvalidLangID = "Invalid or missing language ID"
	ErrInvalidPage   = "Invalid page number"
	ErrInvalidLimit  = "Invalid limit"
)

var (
	ErrBadRequest = errors.New("bad request")
)
