package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/hiennguyen9874/go-boilerplate/config"
	"github.com/hiennguyen9874/go-boilerplate/internal/items"
	"github.com/hiennguyen9874/go-boilerplate/internal/items/presenter"
	"github.com/hiennguyen9874/go-boilerplate/internal/middleware"
	"github.com/hiennguyen9874/go-boilerplate/internal/models"
	"github.com/hiennguyen9874/go-boilerplate/pkg/httpErrors"
	"github.com/hiennguyen9874/go-boilerplate/pkg/logger"
	"github.com/hiennguyen9874/go-boilerplate/pkg/responses"
	"github.com/hiennguyen9874/go-boilerplate/pkg/utils"
)

type itemHandler struct {
	cfg     *config.Config
	itemsUC items.ItemUseCaseI
	logger  logger.Logger
}

func CreateItemHandler(uc items.ItemUseCaseI, cfg *config.Config, logger logger.Logger) items.Handlers {
	return &itemHandler{cfg: cfg, itemsUC: uc, logger: logger}
}

func (h *itemHandler) Create() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		item := new(presenter.ItemCreate)

		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(ctx, item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		newItem, err := h.itemsUC.CreateWithOwner(
			ctx,
			user.Id,
			mapModel(item),
		)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		itemResponse := *mapModelResponse(newItem)
		render.Respond(w, r, responses.CreateSuccessResponse(itemResponse))
	}
}

func (h *itemHandler) Get() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if !user.IsSuperUser && item.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

func (h *itemHandler) GetMulti() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		limit, _ := strconv.Atoi(q.Get("limit"))
		offset, _ := strconv.Atoi(q.Get("offset"))

		ctx := r.Context()

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		var items []*models.Item
		if user.IsSuperUser {
			items, err = h.itemsUC.GetMulti(ctx, limit, offset)
		} else {
			items, err = h.itemsUC.GetMultiByOwnerId(ctx, user.Id, limit, offset)
		}
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}
		render.Respond(w, r, responses.CreateSuccessResponse(mapModelsResponse(items)))
	}
}

func (h *itemHandler) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		item, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if !user.IsSuperUser && item.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		err = h.itemsUC.DeleteWithoutGet(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(item)))
	}
}

func (h *itemHandler) Update() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		item := new(presenter.ItemUpdate)

		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		err = utils.ValidateStruct(r.Context(), item)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrValidation(err)))
			return
		}

		user, err := middleware.GetUserFromCtx(ctx)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		dbItem, err := h.itemsUC.Get(ctx, id)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		if !user.IsSuperUser && dbItem.OwnerId != user.Id {
			render.Render(w, r, responses.CreateErrorResponse(httpErrors.ErrNotEnoughPrivileges(err)))
			return
		}

		values := make(map[string]interface{})
		if item.Title != "" {
			values["title"] = item.Title
		}
		if item.Description != "" {
			values["description"] = item.Description
		}

		updatedItem, err := h.itemsUC.Update(r.Context(), id, values)
		if err != nil {
			render.Render(w, r, responses.CreateErrorResponse(err))
			return
		}

		render.Respond(w, r, responses.CreateSuccessResponse(mapModelResponse(updatedItem)))
	}
}

func mapModel(exp *presenter.ItemCreate) *models.Item {
	return &models.Item{
		Title:       exp.Title,
		Description: exp.Description,
	}
}

func mapModelResponse(exp *models.Item) *presenter.ItemResponse {
	return &presenter.ItemResponse{
		Id:          exp.Id,
		Title:       exp.Title,
		Description: exp.Description,
		OwnerId:     exp.OwnerId,
	}
}

func mapModelsResponse(exp []*models.Item) []*presenter.ItemResponse {
	out := make([]*presenter.ItemResponse, len(exp))
	for i, user := range exp {
		out[i] = mapModelResponse(user)
	}
	return out
}
