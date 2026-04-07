package env

import (
	"fmt"
	"os"

	"github.com/nzhussup/conform/internal/decode"
	"github.com/nzhussup/conform/internal/errs"
	"github.com/nzhussup/conform/internal/schema"
)

func Load(sc *schema.Schema) error {
	if sc == nil {
		return errs.InvalidSchemaNil
	}

	for _, field := range sc.Fields {
		envName := field.EnvName
		if envName == "" {
			continue
		}

		raw, ok := os.LookupEnv(envName)
		if !ok {
			continue
		}

		if err := decode.SetFieldValue(field, raw); err != nil {
			ctx := fmt.Sprintf("env %q -> %s", envName, field.Path)
			return errs.WrapDecode(errs.Decode, ctx, err)
		}
	}

	return nil
}
