//go:build !integration

package workflow

import (
	"os"
	"strings"
	"testing"
)

func TestDailyModelInventoryWorkflowDoesNotPrefetchReflectOnRunnerHost(t *testing.T) {
	lockContent, err := os.ReadFile("../../.github/workflows/daily-model-inventory.lock.yml")
	if err != nil {
		t.Fatalf("failed to read compiled workflow: %v", err)
	}

	lockContentStr := string(lockContent)

	if strings.Contains(lockContentStr, "Fetch Copilot reflect inventory") {
		t.Fatalf("expected compiled workflow to avoid runner-host reflect prefetch step")
	}

	if strings.Contains(lockContentStr, `shell(mkdir -p /tmp/gh-aw/model-inventory && (curl -fsS http://api-proxy:10000/reflect > /tmp/gh-aw/model-inventory/reflect.json || printf "%s" "{\"endpoints\":[],\"error\":\"reflect endpoint unavailable\"}" > /tmp/gh-aw/model-inventory/reflect.json))`) {
		t.Fatalf("expected compiled workflow to avoid the complex Copilot shell allow-tool for /reflect fallback")
	}
}
