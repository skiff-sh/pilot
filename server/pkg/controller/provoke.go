package controller

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v3"
	cmap "github.com/orcaman/concurrent-map/v2"
	"github.com/skiff-sh/pilot/api/go/pilot"
	"github.com/skiff-sh/pilot/pkg/behavior/behaviortype"
	"google.golang.org/grpc/codes"
)

func NewProvokeHandler(store cmap.ConcurrentMap[string, behaviortype.Interface]) fiber.Handler {
	return func(c fiber.Ctx) error {
		ctx := c.UserContext()
		name := c.Params("name")
		if name == "" {
			return c.Status(http.StatusBadRequest).SendString("name is required")
		}

		beh, ok := store.Get(name)
		if !ok {
			return c.Status(http.StatusNotFound).SendString(fmt.Sprintf("name %s not found", name))
		}

		resp, err := beh.Provoke(ctx)
		if err != nil {
			return err
		}

		if resp.Status != nil && resp.Status.Code() != codes.OK {
			return resp.Status.Err()
		}

		return c.JSON(&pilot.ProvokeBehavior_Response{Body: resp.Body.ToProto()})
	}
}
