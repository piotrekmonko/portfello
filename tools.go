//go:build tools
// +build tools

package tools

import (
	_ "github.com/99designs/gqlgen"
	_ "github.com/google/wire/cmd/wire"
	_ "github.com/spf13/cobra"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)
