package main

import (
	"core/internal/core"
	"core/internal/textinput"
	"encoding/json"
	"strings"
	"testing"
)

// integration test
func TestMain(t *testing.T) {
	t.Run("textinput", func(t *testing.T) {
		text := `
			type User struct {
				Name string "json:\"name\""
			}
		`
		structName := "User"

		loader := textinput.NewDynamicStructLoader(text, structName)

		t.Cleanup(func() {
			if err := loader.Done(); err != nil {
				t.Logf("failed to clean up loader: %v", err)
			}
		})

		if err := loader.Load(); err != nil {
			t.Fatalf("failed to load plugin: %v", err)
		}

		out, err := loader.NewInstance()
		if err != nil {
			t.Fatalf("failed to load struct: %v", err)
		}

		jsonText := core.Stringify(out)

		t.Logf("created instance: \n%s", jsonText)

		// check that the jsonText doesn't contain "error"
		if strings.Contains(jsonText, "error") {
			t.Fatalf("failed to create instance, got error in json: %s", jsonText)
		}

		type User struct {
			Name string "json:\"name\""
		}

		var user User

		if err := json.Unmarshal([]byte(jsonText), &user); err != nil {
			t.Errorf("json text is not backward compatible: %s, got error: %s", jsonText, err)
		}
	})
}
