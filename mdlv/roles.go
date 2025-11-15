package mdlv

import (
	"net/http"

	"github.com/chains-lab/ape"
	"github.com/chains-lab/ape/problems"
	"github.com/chains-lab/restkit/roles"
	"github.com/chains-lab/restkit/token"
)

func SystemRoleGrant(ctxKey interface{}, allowedRoles map[string]bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(ctxKey).(token.UserData)
			if !ok {
				ape.RenderErr(w,
					problems.Unauthorized("Missing AuthorizationHeader header"),
				)

				return
			}

			if err := roles.ValidateUserSystemRole(user.Role); err != nil {
				ape.RenderErr(w,
					problems.Unauthorized("SystemUser role not valid"),
				)

				return
			}

			if !allowedRoles[user.Role] {
				ape.RenderErr(w,
					problems.Forbidden("SystemUser role not allowedRoles"),
				)

				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
