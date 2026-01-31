package tokens

import (
	"fmt"
	"math"
)

const (
	RoleSystemAdmin = "admin"
	RoleSystemModer = "moderator"
	RoleSystemUser  = "user"
)

var userSystemRoleAllowed = []string{
	RoleSystemAdmin,
	RoleSystemModer,
	RoleSystemUser,
}

var ErrorInvalidUserSystemRole = fmt.Errorf("role must be one of: %s", GetAllSystemUserRoles())

func ValidateUserSystemRole(i string) error {
	for _, role := range userSystemRoleAllowed {
		if role == i {
			return nil
		}
	}

	return fmt.Errorf("'%s', %w", i, ErrorInvalidUserSystemRole)
}

var PriorityUserSystemRoles = map[string]uint8{
	RoleSystemAdmin: math.MaxUint8,
	RoleSystemModer: 2,
	RoleSystemUser:  1,
}

// CompareSystemUserRoles
// res : 1, if first role is higher priority
// res : -1, if second role is higher priority
// res : 0, if roles are equal
func CompareSystemUserRoles(role1, role2 string) (int, error) {
	err := ValidateUserSystemRole(role1)
	if err != nil {
		return -1, err
	}

	err = ValidateUserSystemRole(role2)
	if err != nil {
		return -1, err
	}

	p1, ok1 := PriorityUserSystemRoles[role1]
	p2, ok2 := PriorityUserSystemRoles[role2]

	if !ok1 || !ok2 {
		return -1, nil
	}

	if p1 > p2 {
		return 1, nil
	} else if p1 < p2 {
		return -1, nil
	}
	return 0, nil
}

func GetAllSystemUserRoles() []string {
	return userSystemRoleAllowed
}
