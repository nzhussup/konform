package common

import (
	"sort"

	"github.com/nzhussup/konform/internal/schema"
)

type unknownKeyIssue struct {
	Path       string
	Suggestion string
}

func FindUnknownKeyIssues(sc *schema.Schema, doc Document, pathAliases map[string]string) []unknownKeyIssue {
	if sc == nil {
		return nil
	}

	expectedPaths := BuildExpectedLookupPaths(sc, pathAliases)
	actualPaths := FlattenLeafPaths(doc)
	actualSet := sliceToPathSet(actualPaths)

	issues := make([]unknownKeyIssue, 0)

	missingExpected := make([]string, 0, len(expectedPaths))
	for path := range expectedPaths {
		if _, ok := actualSet[path]; ok {
			continue
		}
		missingExpected = append(missingExpected, path)
	}
	sort.Strings(missingExpected)

	actualSetForSuggest := sliceToPathSet(actualPaths)
	for _, missingPath := range missingExpected {
		suggestion, ok := SuggestPath(missingPath, actualSetForSuggest)
		if ok {
			issues = append(issues, unknownKeyIssue{
				Path:       missingPath,
				Suggestion: suggestion,
			})
		}
	}

	return issues
}

func BuildExpectedLookupPaths(sc *schema.Schema, pathAliases map[string]string) map[string]struct{} {
	paths := make(map[string]struct{}, len(sc.Fields))
	for _, field := range sc.Fields {
		if isStructField(field) {
			continue
		}

		lookupPath := ResolveLookupPath(field, pathAliases)
		if lookupPath == "" {
			continue
		}

		paths[lookupPath] = struct{}{}
	}

	return paths
}

func FindUnknownKeys(doc Document, expectedPaths map[string]struct{}) []string {
	leafPaths := FlattenLeafPaths(doc)
	unknown := make([]string, 0)
	for _, path := range leafPaths {
		if _, ok := expectedPaths[path]; ok {
			continue
		}
		unknown = append(unknown, path)
	}
	sort.Strings(unknown)
	return unknown
}

func FlattenLeafPaths(doc Document) []string {
	if doc == nil {
		return nil
	}

	paths := make([]string, 0)
	flattenLeafPaths("", doc, &paths)
	sort.Strings(paths)
	return paths
}

func flattenLeafPaths(prefix string, current any, out *[]string) {
	m, ok := asStringMap(current)
	if !ok {
		if prefix != "" {
			*out = append(*out, prefix)
		}
		return
	}

	if len(m) == 0 {
		if prefix != "" {
			*out = append(*out, prefix)
		}
		return
	}

	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		path := key
		if prefix != "" {
			path = prefix + "." + key
		}
		flattenLeafPaths(path, m[key], out)
	}
}

func sliceToPathSet(paths []string) map[string]struct{} {
	out := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		out[path] = struct{}{}
	}
	return out
}
