package mdlv

import (
	"net/http"

	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
	"github.com/netbill/restkit/auth"
	"github.com/netbill/restkit/auth/roles"
)

func (s Service) RoleGrant(allowedRoles map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(s.ctxKey).(auth.AccountData)
			if !ok {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			if err := roles.ValidateUserSystemRole(user.Role); err != nil {
				ape.RenderErr(w,
					problems.Unauthorized("account role not valid"),
				)

				return
			}

			if !allowedRoles[user.Role] {
				ape.RenderErr(w,
					problems.Forbidden("account role not allowedRoles"),
				)

				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
