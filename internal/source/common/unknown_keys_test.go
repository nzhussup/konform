package common

import (
	"reflect"
	"testing"

	"github.com/nzhussup/konform/internal/schema"
)

func TestFindUnknownKeyIssues(t *testing.T) {
	type nested struct{ Port int }

	makeSchema := func() *schema.Schema {
		var s nested
		var port int
		return &schema.Schema{
			Fields: []schema.Field{
				{
					Path:    "Server",
					KeyName: "server_cfg",
					Type:    reflect.TypeOf(s),
					Value:   reflect.ValueOf(&s).Elem(),
				},
				{
					Path:  "Server.Port",
					Type:  reflect.TypeOf(0),
					Value: reflect.ValueOf(&port).Elem(),
				},
			},
		}
	}

	t.Run("nil schema", func(t *testing.T) {
		got := FindUnknownKeyIssues(nil, Document{"server_cfg": map[string]any{"Port": 8080}}, nil)
		if got != nil {
			t.Fatalf("FindUnknownKeyIssues() = %#v, want nil", got)
		}
	})

	t.Run("missing expected path with suggestion", func(t *testing.T) {
		sc := makeSchema()
		aliases := BuildPathAliases(sc)
		doc := Document{
			"server_cfg": map[string]any{
				"Porrt": 8080,
			},
		}

		got := FindUnknownKeyIssues(sc, doc, aliases)
		if len(got) != 1 {
			t.Fatalf("len(issues) = %d, want 1", len(got))
		}
		if got[0].Path != "server_cfg.Port" {
			t.Fatalf("issue path = %q, want %q", got[0].Path, "server_cfg.Port")
		}
		if got[0].Suggestion != "server_cfg.Porrt" {
			t.Fatalf("issue suggestion = %q, want %q", got[0].Suggestion, "server_cfg.Porrt")
		}
	})

	t.Run("present expected path has no issues", func(t *testing.T) {
		sc := makeSchema()
		aliases := BuildPathAliases(sc)
		doc := Document{
			"server_cfg": map[string]any{
				"Port": 8080,
			},
		}

		got := FindUnknownKeyIssues(sc, doc, aliases)
		if len(got) != 0 {
			t.Fatalf("len(issues) = %d, want 0", len(got))
		}
	})
}

func TestUnknownKeysBuildExpectedLookupPaths(t *testing.T) {
	type nested struct{ Port int }
	var s nested
	var port int
	var timeout int

	sc := &schema.Schema{
		Fields: []schema.Field{
			{
				Path:    "Server",
				KeyName: "server_cfg",
				Type:    reflect.TypeOf(s),
				Value:   reflect.ValueOf(&s).Elem(),
			},
			{
				Path:  "Server.Port",
				Type:  reflect.TypeOf(0),
				Value: reflect.ValueOf(&port).Elem(),
			},
			{
				Path:    "Timeout",
				KeyName: "timeout",
				Type:    reflect.TypeOf(0),
				Value:   reflect.ValueOf(&timeout).Elem(),
			},
		},
	}

	got := BuildExpectedLookupPaths(sc, BuildPathAliases(sc))
	want := map[string]struct{}{
		"server_cfg.Port": {},
		"timeout":         {},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("BuildExpectedLookupPaths() = %#v, want %#v", got, want)
	}
}

func TestUnknownKeysFindUnknownKeys(t *testing.T) {
	doc := Document{
		"server": map[string]any{
			"port": 8080,
			"host": "localhost",
		},
		"mode": "prod",
	}
	expected := map[string]struct{}{
		"server.port": {},
	}

	got := FindUnknownKeys(doc, expected)
	want := []string{"mode", "server.host"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("FindUnknownKeys() = %#v, want %#v", got, want)
	}
}

func TestUnknownKeysFlattenLeafPathsAndSliceToPathSet(t *testing.T) {
	doc := Document{
		"server": map[string]any{
			"port": 8080,
			"tls": map[string]any{
				"enabled": true,
			},
		},
		"empty": map[string]any{},
	}

	paths := FlattenLeafPaths(doc)
	wantPaths := []string{"empty", "server.port", "server.tls.enabled"}
	if !reflect.DeepEqual(paths, wantPaths) {
		t.Fatalf("FlattenLeafPaths() = %#v, want %#v", paths, wantPaths)
	}

	set := sliceToPathSet(paths)
	for _, p := range wantPaths {
		if _, ok := set[p]; !ok {
			t.Fatalf("sliceToPathSet() missing path %q", p)
		}
	}
}
