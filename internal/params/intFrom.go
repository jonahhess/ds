package params

import "context"

func IntFrom(ctx context.Context, name string) (int, bool) {
	val, ok := ctx.Value(ctxKey(name)).(int)
	return val, ok
}
