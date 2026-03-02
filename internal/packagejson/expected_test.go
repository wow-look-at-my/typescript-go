package packagejson_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/json"
	"github.com/microsoft/typescript-go/internal/packagejson"
	"github.com/wow-look-at-my/testify/require"
)

func TestExpected(t *testing.T) {
	t.Parallel()

	type packageJson struct {
		Name	packagejson.Expected[string]	`json:"name"`
		Version	packagejson.Expected[string]	`json:"version"`
		Exports	packagejson.Expected[any]	`json:"exports"`
		Main	packagejson.Expected[string]	`json:"main"`
	}

	var p packageJson

	jsonString := `{
		"name": "test",
		"version": 2,
		"exports": null
	}`

	err := json.Unmarshal([]byte(jsonString), &p)
	require.NoError(t, err)

	require.Equal(t, p.Name.Valid, true)
	require.Equal(t, p.Name.Value, "test")

	require.Equal(t, p.Version.Valid, false)
	require.Equal(t, p.Version.Value, "")

	require.True(t, p.Exports.Null)
	require.Equal(t, p.Exports.Valid, false)

	require.Equal(t, p.Main.Valid, false)
	require.Equal(t, p.Main.Null, false)
	require.Equal(t, p.Main.Value, "")
}
