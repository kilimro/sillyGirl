package core

import "testing"

func TestAdminPanelsRouteReplacesLegacyListRoutes(t *testing.T) {
	legacy := map[string]bool{
		"/api/admin/qinglong/panels": false,
		"/api/admin/daidai/panels":   false,
	}
	foundUnified := false
	for _, route := range apiRouteSnapshot() {
		if route.Method != GET {
			continue
		}
		if route.Path == "/api/admin/panels" {
			foundUnified = true
		}
		if _, exists := legacy[route.Path]; exists {
			legacy[route.Path] = true
		}
	}
	if !foundUnified {
		t.Fatal("unified admin panels route is not registered")
	}
	for path, registered := range legacy {
		if registered {
			t.Fatalf("legacy panel list route is still registered: %s", path)
		}
	}
}
