package check

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"wait4it/pkg/model"
)

func RunCheck(ctx context.Context, c *model.CheckContext) error {
	cx, err := findCheckModule(c)
	if err != nil {
		return errors.Wrap(err, "can not find the module")
	}

	newCtx, cnl := context.WithTimeout(ctx, time.Duration(c.Config.Timeout)*time.Second)
	defer cnl()

	progress := c.Progress
	if progress == nil {
		progress = func(s string) {}
	}

	progress("Wait4it...")

	if err := ticker(newCtx, cx, progress); err != nil {
		return errors.Wrap(err, "check failed")
	}

	return nil
}

func findCheckModule(c *model.CheckContext) (model.CheckInterface, error) {
	newFunc, ok := cm[c.Config.CheckType]
	if !ok {
		return nil, errors.New("unsupported check type")
	}

	return newFunc(c)
}

func ticker(ctx context.Context, cs model.CheckInterface, progress func(string)) error {
	t := time.NewTicker(1 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			r, err := check(ctx, cs)
			if err != nil {
				return errors.Wrap(err, "check failed")
			}

			if r {
				return nil
			}

			progress(".")
		}
	}
}

func check(ctx context.Context, cs model.CheckInterface) (bool, error) {
	r, eor, err := cs.Check(ctx)
	if err != nil && eor {
		return false, errors.Wrap(err, "failed")
	}

	return r, nil
}
