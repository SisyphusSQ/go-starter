package vo

import (
	"go-starter/utils"
)

func ValidateBaseList(page, pageSize int) error {
	if page == 0 || pageSize == 0 {
		return utils.ErrBadParamInput
	}

	return nil
}
