package api

import (
	"context"
	"fmt"

	"github.com/spolu/cumulo/lib/env"
	"github.com/spolu/cumulo/lib/logging"
)

const (
	// EnvCfgPort is the port on which to run the mint.
	EnvCfgPort env.ConfigKey = "port"
)

// GetPort retrieves the current mint port from the given contest.
func GetPort(
	ctx context.Context,
) string {
	return env.Get(ctx).Config[EnvCfgPort]
}

// Logf shells out to logging.Logf adding the mint host as prefix.
func Logf(
	ctx context.Context,
	format string,
	v ...interface{},
) {
	logging.Logf(ctx, fmt.Sprintf("[api:%s] ", GetPort(ctx))+format, v...)
}
