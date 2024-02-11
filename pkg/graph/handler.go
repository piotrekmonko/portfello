package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/piotrekmonko/portfello/pkg/auth"
	"github.com/piotrekmonko/portfello/pkg/conf"
	"github.com/piotrekmonko/portfello/pkg/dao"
	"net/http"
)

func NewGraphHandler(conf *conf.Config, dbQuerier *dao.DAO, authService *auth.Service) http.Handler {
	graphResolver := &Resolver{
		Conf:        conf,
		Dao:         dbQuerier,
		AuthService: authService,
	}

	graphConfig := Config{
		Resolvers: graphResolver,
	}
	graphConfig.Directives.HasRole = authService.HasRole

	return authService.Middleware(handler.NewDefaultServer(NewExecutableSchema(graphConfig)))
}
