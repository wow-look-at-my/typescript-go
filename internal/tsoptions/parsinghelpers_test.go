package tsoptions

import (
	"reflect"
	"strings"
	"testing"

	"github.com/microsoft/typescript-go/internal/core"
	"github.com/wow-look-at-my/testify/assert"
)

func TestParseCompilerOptionNoMissingFields(t *testing.T) {
	t.Parallel()
	var missingKeys []string
	for _, field := range reflect.VisibleFields(reflect.TypeFor[core.CompilerOptions]()) {
		if !field.IsExported() {
			continue
		}

		keyName := field.Name
		// use the JSON key from the tag, if present
		// e.g. `json:"dog[,anythingelse]"` --> dog
		if jsonTag, ok := field.Tag.Lookup("json"); ok {
			keyName = strings.SplitN(jsonTag, ",", 2)[0]
		}
		val := reflect.Zero(field.Type).Interface()
		co := core.CompilerOptions{}
		found := parseCompilerOptions(keyName, val, &co)
		if !found {
			missingKeys = append(missingKeys, keyName)
		}
	}
	assert.LessOrEqual(t, len(missingKeys), 0)

}
