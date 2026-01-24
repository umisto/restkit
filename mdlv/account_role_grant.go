package mdlv

import (
	"net/http"

	"github.com/netbill/logium"
	"github.com/netbill/restkit/ape"
	"github.com/netbill/restkit/ape/problems"
	"github.com/netbill/restkit/tokens"
	"github.com/netbill/restkit/tokens/roles"
)

func AccountRoleGrant(
	log logium.Logger,
	ctxKey int,
	allowedRoles map[string]bool,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, ok := ctx.Value(ctxKey).(tokens.AccountJwtData)
			if !ok {
				log.Errorf("missing AuthorizationHeader header")
				ape.RenderErr(w, problems.Unauthorized("Missing AuthorizationHeader header"))

				return
			}

			if err := roles.ValidateUserSystemRole(user.Role); err != nil {
				log.WithError(err).Errorf("account role not valid")
				ape.RenderErr(w, problems.Unauthorized("account role not valid"))

				return
			}

			if !allowedRoles[user.Role] {
				log.Errorf("account role not allowedRoles")
				ape.RenderErr(w, problems.Forbidden("account role not allowedRoles"))

				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
