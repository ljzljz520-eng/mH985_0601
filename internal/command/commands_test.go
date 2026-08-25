package command

import "testing"

func TestParseCommand(t *testing.T) {
	req, err := Parse("search store=north limit=10")
	if err != nil || req.Name != "search" || req.String("store") != "north" || req.Int("limit", 0) != 10 || !req.Has("store") {
		t.Fatalf("request=%#v err=%v", req, err)
	}
	if _, err := Parse("search broken"); err == nil {
		t.Fatal("expected malformed parameter error")
	}
}
