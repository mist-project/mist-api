package api

import (
	"net/http"

	pb "mistapi/src/protos/v1/gen"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// A completely separate router for administrator routes
func appserverRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", list)
	// r.Get("/{id}", getUser)
	// r.Post("/", createAppserver)
	return r
}

type Appserver struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	IsOwner bool   `json:"is_owner"`
}

// list godoc
// @Summary      List Appservers
// @Description  Get a list of appservers
// @Tags         appserver
// @Accept       json
// @Produce      json
// @Success      200  {array}  Appserver
// @Router       /api/v1/appserver [get]
func list(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := setupContext()
	defer cancel()

	grpcC := GetGRPCConnFromContext(r)
	c := Client{Conn: grpcC}
	response, err := c.GetServerClient().ListAppservers(
		ctx, &pb.ListAppserversRequest{},
	)

	if err != nil {
		// TODO: add better error handling
		render.JSON(w, r, err)
		return
	}

	res := make([]Appserver, 0, len(response.Appservers))

	for _, a := range response.Appservers {
		res = append(res, Appserver{
			ID:      a.Id,
			Name:    a.Name,
			IsOwner: a.IsOwner,
		})
	}

	render.JSON(w, r, res)
}
