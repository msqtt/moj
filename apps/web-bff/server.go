package main

import (
	"log"
	"moj/apps/web-bff/graph"
	"moj/apps/web-bff/middleware"
	"net/http"
	"strconv"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
)

func startHttpServer(resolver *graph.Resolver) {
	port := strconv.FormatInt(int64(resolver.Conf.AppPort), 10)
	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// register middleware
	wrapSrv := middleware.WithClientIp(srv)
	wrapSrv = middleware.WithClientAgent(wrapSrv)
	wrapSrv = middleware.WithAuthToken(wrapSrv)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", wrapSrv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
