package graphql_test

import (
	"context"
	"testing"

	ggql "github.com/graph-gophers/graphql-go"
	"github.com/graphql-go/graphql"
	pgql "github.com/playlyfe/go-graphql"
)

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "RootQueryType",
				Fields: graphql.Fields{
					"hello": &graphql.Field{
						Type: graphql.String,
						Resolve: func(p graphql.ResolveParams) (interface{}, error) {
							return "world", nil
						},
					},
				},
			}),
	},
)

func BenchmarkGoGraphQLMaster(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: "{hello}",
		})
	}
}

var schema2 = `
    type RootQueryType {
        hello: String
    }
  `
var resolvers = map[string]interface{}{
	"RootQueryType/hello": func(params *pgql.ResolveParams) (interface{}, error) {
		return "world", nil
	},
}
var executor, _ = pgql.NewExecutor(schema2, "RootQueryType", "", resolvers)

func BenchmarkPlaylyfeGraphQLMaster(b *testing.B) {
	for i := 0; i < b.N; i++ {
		context := map[string]interface{}{}
		variables := map[string]interface{}{}
		executor.Execute(context, "{hello}", variables, "")
	}
}

type helloWorldResolver1 struct{}

func (r *helloWorldResolver1) Hello() string {
	return "world"
}

var schema3 = ggql.MustParseSchema(`
schema {
  query: Query
}
type Query {
  hello: String!
}
`, &helloWorldResolver1{})

func BenchmarkGophersGraphQLMaster(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := context.Background()
		variables := map[string]interface{}{}
		schema3.Exec(ctx, "{hello}", "", variables)
	}
}
