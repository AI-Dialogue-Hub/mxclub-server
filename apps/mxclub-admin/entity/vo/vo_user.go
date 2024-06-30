package vo

import (
	"mxclub/domain/user/entity/enum"
)

type UserVO struct {
	Name     string        `json:"name"`
	Role     enum.RoleType `json:"Role"`
	JwtToken string        `json:"JwtToken"`
}
