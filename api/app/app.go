package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	goji "goji.io"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/spolu/cumulo/api"
	"github.com/spolu/cumulo/lib/db"
	"github.com/spolu/cumulo/lib/env"
	"github.com/spolu/cumulo/lib/errors"
	"github.com/spolu/cumulo/lib/logging"
	"github.com/spolu/cumulo/lib/recoverer"
	"github.com/spolu/cumulo/lib/requestlogger"

	// force initialization of schemas
	_ "github.com/spolu/cumulo/api/model/schemas"
)

// BackgroundContextFromFlags initializes a background context fully loaded
// with everything that could be extracted from the flags.
func BackgroundContextFromFlags(
	envFlag string,
	dsnFlag string,
	prtFlag string,
) (context.Context, error) {
	ctx := context.Background()

	apiEnv := env.Env{
		Environment: env.Test,
		Config:      map[env.ConfigKey]string{},
	}
	if envFlag == "live" {
		apiEnv.Environment = env.Live
	}

	port := fmt.Sprintf("%d", api.DefaultPort)
	if prtFlag != "" {
		port = prtFlag
	}
	apiEnv.Config[api.EnvCfgPort] = port

	ctx = env.With(ctx, &apiEnv)

	apiDB, err := db.NewDBForDSN(ctx,
		dsnFlag,
		"sqlite3://~/.cumulo/cumulo.db")
	if err != nil {
		return nil, errors.Trace(err)
	}
	err = db.CreateDBTables(ctx, "api", apiDB)
	if err != nil {
		return nil, errors.Trace(err)
	}
	ctx = db.WithDB(ctx, "api", apiDB)

	return ctx, nil
}

// Build initializes the app and its web stack.
func Build(
	ctx context.Context,
) (*goji.Mux, error) {
	mux := goji.NewMux()
	mux.Use(requestlogger.Middleware)
	mux.Use(recoverer.Middleware)
	mux.Use(db.Middleware(db.GetDBMap(ctx)))
	mux.Use(env.Middleware(env.Get(ctx)))

	logging.Logf(ctx, "Initializing: environment=%s port=%s",
		env.Get(ctx).Environment, api.GetPort(ctx))

	(&Controller{}).Bind(mux)

	return mux, nil
}

// Serve the goji mux.
func Serve(
	ctx context.Context,
	mux *goji.Mux,
) error {

	s := &http.Server{
		Addr:         fmt.Sprintf(":%s", api.GetPort(ctx)),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      mux,
	}

	logging.Logf(ctx, "Listening: port=%s", api.GetPort(ctx))

	err := gracehttp.Serve(s)
	if err != nil {
		return errors.Trace(err)
	}

	return nil
}
