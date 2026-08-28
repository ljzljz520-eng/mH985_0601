package policy

import (
	"errors"
	"sort"
	"strings"

	"trainingdesk/internal/model"
)

type Role string

const (
	RoleTrainer  Role = "trainer"
	RoleReviewer Role = "reviewer"
	RoleManager  Role = "manager"
	RoleAdmin    Role = "admin"
)

type Identity struct {
	Name   string
	Roles  []Role
	Stores []string
}

type Decision struct {
	Allowed bool
	Reason  string
}

func Authorize(identity Identity, action string, record model.Record) Decision {
	if strings.TrimSpace(identity.Name) == "" {
		return Decision{Reason: "identity is missing"}
	}
	if !storeAllowed(identity, record.StoreID) {
		return Decision{Reason: "identity is outside record store"}
	}
	for _, role := range identity.Roles {
		switch role {
		case RoleAdmin:
			return Decision{Allowed: true, Reason: "administrator"}
		case RoleManager:
			if action == "archive" || action == "publish" || action == "change" {
				return Decision{Allowed: true, Reason: "manager policy"}
			}
		case RoleReviewer:
			if action == "review" && record.Status == model.StatusPending {
				return Decision{Allowed: true, Reason: "reviewer policy"}
			}
		case RoleTrainer:
			if action == "register" || action == "change" {
				return Decision{Allowed: true, Reason: "trainer policy"}
			}
		}
	}
	return Decision{Reason: "role does not permit action"}
}

func storeAllowed(identity Identity, storeID string) bool {
	if len(identity.Stores) == 0 {
		return false
	}
	for _, allowed := range identity.Stores {
		if allowed == "*" || allowed == storeID {
			return true
		}
	}
	return false
}

func Normalize(identity Identity) (Identity, error) {
	if strings.TrimSpace(identity.Name) == "" {
		return Identity{}, errors.New("identity name is required")
	}
	roles := make(map[Role]bool)
	for _, role := range identity.Roles {
		switch role {
		case RoleTrainer, RoleReviewer, RoleManager, RoleAdmin:
			roles[role] = true
		default:
			return Identity{}, errors.New("unknown role")
		}
	}
	stores := make(map[string]bool)
	for _, storeID := range identity.Stores {
		if strings.TrimSpace(storeID) != "" {
			stores[storeID] = true
		}
	}
	normalized := Identity{Name: strings.TrimSpace(identity.Name)}
	for role := range roles {
		normalized.Roles = append(normalized.Roles, role)
	}
	for storeID := range stores {
		normalized.Stores = append(normalized.Stores, storeID)
	}
	sort.Slice(normalized.Roles, func(i, j int) bool { return normalized.Roles[i] < normalized.Roles[j] })
	sort.Strings(normalized.Stores)
	return normalized, nil
}

func CanView(identity Identity, record model.Record) bool {
	if !storeAllowed(identity, record.StoreID) {
		return false
	}
	if record.Status == model.StatusDraft {
		for _, role := range identity.Roles {
			if role == RoleTrainer || role == RoleReviewer || role == RoleManager || role == RoleAdmin {
				return true
			}
		}
		return false
	}
	return true
}
