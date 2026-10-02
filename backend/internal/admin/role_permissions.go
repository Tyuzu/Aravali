package admin

// Built-in role grants allowed to each actor type.
var assignableRoleMatrix = map[string][]string{
	"admin":     {"user", "farmer", "worker", "moderator", "admin"},
	"moderator": {"user", "farmer", "worker", "admin"},
}

// AssignableRolesFor returns the list of roles a user may grant based on their current roles.
// An admin may grant any role, while a moderator may grant any role except moderator.
func AssignableRolesFor(actorRoles []string) []string {
	seen := map[string]struct{}{}
	ordered := make([]string, 0, len(actorRoles))

	for _, role := range actorRoles {
		normalized := NormalizeRoleName(role)
		if normalized == "" {
			continue
		}
		if _, ok := seen[normalized]; ok {
			continue
		}
		seen[normalized] = struct{}{}
		ordered = append(ordered, normalized)
	}

	for _, actorRole := range ordered {
		if allowed, ok := assignableRoleMatrix[actorRole]; ok {
			result := make([]string, 0, len(allowed))
			for _, role := range allowed {
				if role == "" {
					continue
				}
				result = append(result, role)
			}
			return result
		}
	}

	return nil
}

// CanAssignRole returns true when the actor is permitted to grant the requested role.
func CanAssignRole(actorRoles []string, requestedRole string) bool {
	normalizedRequested := NormalizeRoleName(requestedRole)
	if normalizedRequested == "" {
		return false
	}

	for _, role := range AssignableRolesFor(actorRoles) {
		if role == normalizedRequested {
			return true
		}
	}

	return false
}
