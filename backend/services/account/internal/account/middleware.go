package account

import (
	"net/http"
	"strings"
)

func (h *Handler) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
		if !ok || token == "" {
			writeError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		c, err := h.tokens.parse(token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "missing or invalid token")
			return
		}

		if c.Role != RoleAdmin {
			writeError(w, http.StatusForbidden, "administrator role required")
			return
		}

		next(w, r)
	}
}
