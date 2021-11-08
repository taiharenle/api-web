package schema

import (
	"github.com/graphql-go/graphql"
)

// 查询节点入口
var queryFields = graphql.Fields{
	//"test": &test,
}

// 更新节点入口
var mutationFields = graphql.Fields{
	//"test": &test,
}

// 定义查询节点
var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Query",
	Description: "rootQuery",
	Fields:      queryFields,
})

// 定义更新节点
var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Mutation",
	Description: "rootMutation",
	Fields:      mutationFields,
})

// 定义Schema用于http handler处理
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})
