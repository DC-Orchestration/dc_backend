package middlewear

import (
    "net/http"
    "fmt"
)

func PermissionMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            userClaims, ok := r.Context().Value(UserContextKey).(*Claims)
            if !ok {
                http.Error(w, "No user in context", http.StatusForbidden)
                return
            }

            userRole := userClaims.Role
            for _, role := range allowedRoles {
                if userRole == role {
                    next.ServeHTTP(w, r)
                    return
                }
            }

            http.Error(w, fmt.Sprintf("Forbidden: requires one of %v", allowedRoles), http.StatusForbidden)
        })
    }
}
