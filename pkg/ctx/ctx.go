package ctx

import "context"

type contextKey string

const (
	ContextId contextKey = "context_id"
)

func NewCtx(id string) context.Context {
	return context.WithValue(context.TODO(), ContextId, id)
}

func InheritCtx(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ContextId, id)
}

func SetCtxValue(ctx context.Context, key string, value any) context.Context {
	return context.WithValue(ctx, contextKey(key), value)
}

func GetCtxValue(ctx context.Context, key string) any {
	return ctx.Value(contextKey(key))
}
