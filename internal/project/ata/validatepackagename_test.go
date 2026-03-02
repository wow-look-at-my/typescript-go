package ata_test

import (
	"testing"

	"github.com/microsoft/typescript-go/internal/project/ata"
	"github.com/wow-look-at-my/testify/require"
)

func TestValidatePackageName(t *testing.T) {
	t.Parallel()
	t.Run("name cannot be too long", func(t *testing.T) {
		t.Parallel()
		packageName := "a"
		for range 8 {
			packageName += packageName	//nolint:perfsprint
		}
		status, _, _ := ata.ValidatePackageName(packageName)
		require.Equal(t, status, ata.NameTooLong)
	})
	t.Run("package name cannot start with dot", func(t *testing.T) {
		t.Parallel()
		status, _, _ := ata.ValidatePackageName(".foo")
		require.Equal(t, status, ata.NameStartsWithDot)
	})
	t.Run("package name cannot start with underscore", func(t *testing.T) {
		t.Parallel()
		status, _, _ := ata.ValidatePackageName("_foo")
		require.Equal(t, status, ata.NameStartsWithUnderscore)
	})
	t.Run("package non URI safe characters are not supported", func(t *testing.T) {
		t.Parallel()
		status, _, _ := ata.ValidatePackageName("  scope  ")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		status, _, _ = ata.ValidatePackageName("; say ‘Hello from TypeScript!’ #")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		status, _, _ = ata.ValidatePackageName("a/b/c")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
	})
	t.Run("scoped package name is supported", func(t *testing.T) {
		t.Parallel()
		status, _, _ := ata.ValidatePackageName("@scope/bar")
		require.Equal(t, status, ata.NameOk)
	})
	t.Run("scoped name in scoped package name cannot start with dot", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@.scope/bar")
		require.Equal(t, status, ata.NameStartsWithDot)
		require.Equal(t, name, ".scope")
		require.Equal(t, isScopeName, true)
		status, name, isScopeName = ata.ValidatePackageName("@.scope/.bar")
		require.Equal(t, status, ata.NameStartsWithDot)
		require.Equal(t, name, ".scope")
		require.Equal(t, isScopeName, true)
	})
	t.Run("scoped name in scoped package name cannot start with dot", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@_scope/bar")
		require.Equal(t, status, ata.NameStartsWithUnderscore)
		require.Equal(t, name, "_scope")
		require.Equal(t, isScopeName, true)
		status, name, isScopeName = ata.ValidatePackageName("@_scope/_bar")
		require.Equal(t, status, ata.NameStartsWithUnderscore)
		require.Equal(t, name, "_scope")
		require.Equal(t, isScopeName, true)
	})
	t.Run("scope name in scoped package name with non URI safe characters are not supported", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@  scope  /bar")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		require.Equal(t, name, "  scope  ")
		require.Equal(t, isScopeName, true)
		status, name, isScopeName = ata.ValidatePackageName("@; say ‘Hello from TypeScript!’ #/bar")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		require.Equal(t, name, "; say ‘Hello from TypeScript!’ #")
		require.Equal(t, isScopeName, true)
		status, name, isScopeName = ata.ValidatePackageName("@  scope  /  bar  ")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		require.Equal(t, name, "  scope  ")
		require.Equal(t, isScopeName, true)
	})
	t.Run("package name in scoped package name cannot start with dot", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@scope/.bar")
		require.Equal(t, status, ata.NameStartsWithDot)
		require.Equal(t, name, ".bar")
		require.Equal(t, isScopeName, false)
	})
	t.Run("package name in scoped package name cannot start with underscore", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@scope/_bar")
		require.Equal(t, status, ata.NameStartsWithUnderscore)
		require.Equal(t, name, "_bar")
		require.Equal(t, isScopeName, false)
	})
	t.Run("package name in scoped package name with non URI safe characters are not supported", func(t *testing.T) {
		t.Parallel()
		status, name, isScopeName := ata.ValidatePackageName("@scope/  bar  ")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		require.Equal(t, name, "  bar  ")
		require.Equal(t, isScopeName, false)
		status, name, isScopeName = ata.ValidatePackageName("@scope/; say ‘Hello from TypeScript!’ #")
		require.Equal(t, status, ata.NameContainsNonURISafeCharacters)
		require.Equal(t, name, "; say ‘Hello from TypeScript!’ #")
		require.Equal(t, isScopeName, false)
	})
}
