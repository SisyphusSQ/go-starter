package vo

import (
	"github.com/SisyphusSQ/go-starter/v2/utils"
)

const MaxPageSize = 100

func ValidateBaseList(page, pageSize int) error {
	if page <= 0 || pageSize <= 0 || pageSize > MaxPageSize {
		return utils.ErrBadParamInput
	}

	return nil
}
