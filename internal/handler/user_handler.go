package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/dotuan0601/go-backend-template/internal/logger"
	"github.com/dotuan0601/go-backend-template/internal/model"
	"github.com/dotuan0601/go-backend-template/internal/response"
	"github.com/dotuan0601/go-backend-template/internal/service"
)

type UserHandler struct {
	userService service.UserService
	log         *slog.Logger
}

func NewUserHandler(
	userService service.UserService,
	log *slog.Logger,
) *UserHandler {

	return &UserHandler{
		userService: userService,
		log:         log,
	}
}

func (h *UserHandler) CreateUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	h.log.Info(
		"create user called",
	)

	var user model.User

	err := json.NewDecoder(
		r.Body,
	).Decode(&user)

	if err != nil {

		logger.Error(
			r.Context(),
			h.log,
			"invalid create user request",
			err,
		)

		response.Error(
			w,
			response.ErrSystem,
		)

		return
	}

	err = h.userService.CreateUser(
		r.Context(),
		&user,
	)

	if err != nil {

		logger.Error(
			r.Context(),
			h.log,
			"create user failed",
			err,
		)

		response.Error(
			w,
			response.ErrSystem,
		)

		return
	}

	response.Success(
		w,
		user,
	)
}

func (h *UserHandler) GetUsers(
	w http.ResponseWriter,
	r *http.Request,
) {

	h.log.Info(
		"get users called",
	)

	users, err := h.userService.GetUsers(
		r.Context(),
	)

	if err != nil {

		logger.Error(
			r.Context(),
			h.log,
			"get users failed",
			err,
		)

		response.Error(
			w,
			response.ErrSystem,
		)

		return
	}

	response.Success(
		w,
		users,
	)
}

func (h *UserHandler) GetUser(
	w http.ResponseWriter,
	r *http.Request,
) {

	id, err := strconv.ParseInt(
		r.PathValue("id"),
		10,
		64,
	)

	if err != nil {

		logger.Error(
			r.Context(),
			h.log,
			"invalid user id",
			err,
		)

		response.Error(
			w,
			response.ErrSystem,
		)

		return
	}

	user, err := h.userService.GetUser(
		r.Context(),
		id,
	)

	if err != nil {

		logger.Error(
			r.Context(),
			h.log,
			"get user failed",
			err,
		)

		response.Error(
			w,
			response.ErrSystem,
		)

		return
	}

	response.Success(
		w,
		user,
	)
}
