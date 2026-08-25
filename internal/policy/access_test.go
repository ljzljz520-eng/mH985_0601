package policy

import (
	"testing"

	"trainingdesk/internal/model"
)

func TestAuthorization(t *testing.T) {
	identity, err := Normalize(Identity{Name: " lee ", Roles: []Role{RoleReviewer, RoleReviewer}, Stores: []string{"north", "north"}})
	if err != nil || identity.Name != "lee" || len(identity.Roles) != 1 {
		t.Fatalf("identity=%#v err=%v", identity, err)
	}
	r := model.Record{StoreID: "north", Status: model.StatusPending}
	if decision := Authorize(identity, "review", r); !decision.Allowed || !CanView(identity, r) {
		t.Fatalf("decision=%#v", decision)
	}
	if Authorize(identity, "archive", r).Allowed {
		t.Fatal("reviewer should not archive")
	}
}
