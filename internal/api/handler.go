package api

import (
	"context"
	"log"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"github.com/streamersonglist/dragonfly-flex/internal/dragonfly"
)

type CheckRoleResponse struct {
	Body string
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	node, err := dragonfly.NewNode()
	if err != nil {
		log.Panicf("failed to create node: %s", err.Error())
	}

	humaApi := humago.New(mux, huma.DefaultConfig("Dragonfly API", "1.0.0"))

	huma.Get(humaApi, "/check/role", func(ctx context.Context, input *struct{}) (*CheckRoleResponse, error) {
		role, err := node.CheckRole(ctx)
		if err != nil {
			return nil, err
		}

		resp := &CheckRoleResponse{
			Body: *role,
		}

		return resp, nil
	})

	return mux
}
