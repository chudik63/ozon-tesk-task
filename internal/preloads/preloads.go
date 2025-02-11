package preloads

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
)

func GetPreloads(ctx context.Context) (preloads []string) {
	defer func() {
		if r := recover(); r != nil {
			preloads = []string{}
		}
	}()

	opCtx := graphql.GetOperationContext(ctx)

	return GetNestedPreloads(
		opCtx,
		graphql.CollectFieldsCtx(ctx, nil),
		"",
	)
}

func GetNestedPreloads(ctx *graphql.OperationContext, fields []graphql.CollectedField, prefix string) (preloads []string) {
	for _, column := range fields {
		prefixColumn := GetPreloadString(prefix, column.Name)
		preloads = append(preloads, prefixColumn)
		preloads = append(preloads, GetNestedPreloads(ctx, graphql.CollectFields(ctx, column.Selections, nil), prefixColumn)...)
	}
	return
}

func GetPreloadString(prefix, name string) string {
	if len(prefix) > 0 {
		return prefix + "." + name
	}
	return name
}
