package json

import (
	"encoding/json"

	"github.com/nzhussup/konform/internal/schema"
	"github.com/nzhussup/konform/internal/source/common"
)

type FileSource struct {
	path           string
	callerDir      string
	suggestionMode common.UnknownKeySuggestionMode
}

func NewFileSource(path string, callerDir string, suggestionMode common.UnknownKeySuggestionMode) FileSource {
	return FileSource{path: path, callerDir: callerDir, suggestionMode: suggestionMode}
}

func (s FileSource) Load(sc *schema.Schema) error {
	return common.LoadFileWithMode(sc, s.path, s.callerDir, "json", func(data []byte) (common.Document, error) {
		var doc common.Document
		if err := json.Unmarshal(data, &doc); err != nil {
			return nil, err
		}
		return doc, nil
	}, s.suggestionMode)
}
