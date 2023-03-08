package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Silarsone/image-loader/internal/config"
	"github.com/Silarsone/image-loader/internal/constants"
	"github.com/Silarsone/image-loader/internal/middleware"
	"github.com/Silarsone/image-loader/internal/model"
	"github.com/Silarsone/image-loader/internal/response"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strconv"
)

//go:generate mockgen -source ./server.go -destination ../mock/server.go -package mock

type User struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Login       string   `json:"login"`
	Password    string   `json:"password"`
	Description string   `json:"description"`
	ImageUrls   []string `json:"imageUrls,omitempty"`
}

type controller interface {
	AddUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, id int64) (model.User, error)
	UpdateUser(ctx context.Context, modelUser model.User) error
	DeleteUser(ctx context.Context, id int64) error
	Authorize(ctx context.Context, login, password string) (string, error)
	AddFile(ctx context.Context, image model.Image) error
}

type Server struct {
	listenURI string
	logger    *logrus.Logger
	r         chi.Router
	c         controller
	cfg       *config.Config
}

func NewServer(listenURI string, logger *logrus.Logger, c controller, cfg *config.Config) *Server {
	return &Server{
		listenURI: listenURI,
		logger:    logger,
		r:         chi.NewRouter(),
		c:         c,
		cfg:       cfg,
	}
}

func (s *Server) RegisterRoutes() {
	s.r.Use(middleware.Logger(s.logger))

	s.r.Get("/user/auth", s.HandleAuthorize)
	s.r.Post("/user/add", s.HandleAddUser)

	s.r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(s.cfg.JWTKeyword, s.logger))

		r.Get("/user/{userID}", s.HandleGetUser)
		r.Put("/user/update", s.HandleUpdateUser)
		r.Delete("/user/{userID}", s.HandleDeleteUser)
		r.Post("/image/add", s.HandleAddFile)
	})
}

func (s *Server) StartServer() {
	srv := http.Server{
		Addr:    s.listenURI,
		Handler: s.r,
	}

	s.logger.Info("server is running!")
	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Fatal(err)
	}
}

// HandleAuthorize issues a JWT
//
//	@Summary      Authorize
//	@Description  Issue JWT
//	@Tags         auth
//	@Accept       json
//	@Produce      json
//	@Param        user    body     User  true  "authorize user"
//	@Success      200  {array}   response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      404  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /user/auth [get]
func (s *Server) HandleAuthorize(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	token, err := s.c.Authorize(r.Context(), user.Login, user.Password)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := response.ParseResponse(token, false)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}
}

// HandleAddUser adds a new user
//
//	@Summary      AddUser
//	@Description  add a new user
//	@Tags         user
//	@Accept       json
//	@Produce      json
//	@Param        user    body     User  true  "authorize user"
//	@Success      200  {array}   response.Response
//	@Failure      400  {object}  response.Response
//	@Failure      404  {object}  response.Response
//	@Failure      500  {object}  response.Response
//	@Router       /user/add [post]
func (s *Server) HandleAddUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	err = s.c.AddUser(r.Context(), user.toModel())
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	user, err := s.c.GetUser(context.Background(), id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	b, err := json.Marshal(&user)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(b)
	if err != nil {
		s.logger.Error(err)
	}
}

func (s *Server) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "userID")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	err = s.c.DeleteUser(r.Context(), id)
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			s.logger.Error(err)
		}
	}(r.Body)

	err = s.c.UpdateUser(r.Context(), user.toModel())
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleAddFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("fileKey")
	if err != nil {
		s.handleError(err, http.StatusBadRequest, w)
	}

	defer file.Close()
	userID, err := userIDFromCtx(r.Context())
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	err = s.c.AddFile(r.Context(), model.Image{
		UserID:    userID,
		Name:      header.Filename,
		Data:      file,
		Extension: ".jpg",
	})
	if err != nil {
		s.handleError(err, http.StatusInternalServerError, w)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleError(err error, status int, w http.ResponseWriter) {
	s.logger.Error(err)
	w.WriteHeader(status)

	b, err := response.ParseResponse(err.Error(), true)
	if err != nil {
		s.logger.Error(err)
	}

	_, err = w.Write(b)
	if err != nil {
		s.logger.Error(err)
	}
}

func (u User) toModel() model.User {
	return model.User{
		ID:          u.ID,
		Name:        u.Name,
		Description: u.Description,
		Login:       u.Login,
		Password:    u.Password,
	}
}

func userIDFromCtx(ctx context.Context) (int, error) {
	idAny := ctx.Value(constants.IdCtxKey)

	id, ok := idAny.(int)
	if !ok {
		return 0, fmt.Errorf("couldn't cast user id from context")
	}

	return id, nil
}
