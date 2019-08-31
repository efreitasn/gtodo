package template

import (
	"context"

	"github.com/efreitasn/gtodo/internal/data/todo"
	"github.com/efreitasn/gtodo/pkg/flash"
)

type ctxKey string

var dataCtxKey ctxKey = "templateData"

// Data is the data used in the rendered templates.
type Data struct {
	FlashMessage *flash.Message
	Mode         string
	Auth         bool
	TodosDone    []todo.Todo
	TodosNotDone []todo.Todo
}

// DataFromContext gets data from a context.
func DataFromContext(ctx context.Context) *Data {
	if data, ok := ctx.Value(dataCtxKey).(*Data); ok {
		return data
	}

	return nil
}

// ContextWithData adds data to a context.
func ContextWithData(ctx context.Context, p *Data) context.Context {
	return context.WithValue(ctx, dataCtxKey, p)
}
