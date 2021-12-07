package task

import (
	"context"
	"os"

	"cuelang.org/go/cue"
	"github.com/rs/zerolog/log"
	"go.dagger.io/dagger/compiler"
	"go.dagger.io/dagger/plancontext"
	"go.dagger.io/dagger/solver"
)

func init() {
	Register("SecretFile", func() Task { return &secretFileTask{} })
}

type secretFileTask struct {
}

func (c secretFileTask) Run(ctx context.Context, pctx *plancontext.Context, _ solver.Solver, v *compiler.Value) (*compiler.Value, error) {
	lg := log.Ctx(ctx)

	var secretFile struct {
		Path string
	}

	if err := v.Decode(&secretFile); err != nil {
		return nil, err
	}

	lg.Debug().Str("path", secretFile.Path).Msg("loading secret")

	plaintext, err := os.ReadFile(secretFile.Path)
	if err != nil {
		return nil, err
	}

	secret := pctx.Secrets.New(string(plaintext))
	out := compiler.NewValue()
	if err := out.FillPath(cue.ParsePath("contents"), secret.MarshalCUE()); err != nil {
		return nil, err
	}
	return out, nil
}