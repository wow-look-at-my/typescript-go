package packagejson_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/wow-look-at-my/testify/require"
)

func TestExports(t *testing.T) {
	t.Parallel()

	t.Run("UnmarshalJSONV2", func(t *testing.T) {
		t.Parallel()
		testExports(t, func(in []byte, out any) error { return json.Unmarshal(in, out) })
	})
}

func testExports(t *testing.T, unmarshal func([]byte, any) error) {
	type Exports struct {
		Imports	packagejson.ExportsOrImports	`json:"imports"`
		Exports	packagejson.ExportsOrImports	`json:"exports"`
	}

	var e Exports

	jsonString := `{
		"imports": {
			"#foo": {
				"import": "./foo.ts"
			}
		},
		"exports": {
			".": {
				"import": "./test.ts",
				"default": "./test.ts"
			},
			"./test": [
				"./test1.ts",
				"./test2.ts",
				null
			],
			"./null": null
		}
	}`

	err := unmarshal([]byte(jsonString), &e)
	require.NoError(t, err)

	require.True(t, e.Exports.IsSubpaths())
	require.Equal(t, e.Exports.AsObject().Size(), 3)
	require.True(t, e.Exports.AsObject().GetOrZero(".").IsConditions())
	require.True(t, e.Exports.AsObject().GetOrZero(".").AsObject().GetOrZero("import").Type == packagejson.JSONValueTypeString)
	require.Equal(t, e.Exports.AsObject().GetOrZero("./test").AsArray()[2].Type, packagejson.JSONValueTypeNull)
	require.True(t, e.Exports.AsObject().GetOrZero("./null").Type == packagejson.JSONValueTypeNull)

	require.True(t, e.Imports.IsImports())
	require.Equal(t, e.Imports.AsObject().Size(), 1)
	require.True(t, e.Imports.AsObject().GetOrZero("#foo").IsConditions())
	require.True(t, e.Imports.AsObject().GetOrZero("#foo").AsObject().GetOrZero("import").Type == packagejson.JSONValueTypeString)
}
