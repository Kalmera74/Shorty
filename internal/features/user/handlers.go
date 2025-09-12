package user

import (
	"strconv"
	"time"

	"github.com/Kalmera74/Shorty/internal/features/shortener"
	"github.com/Kalmera74/Shorty/internal/types"
	"github.com/Kalmera74/Shorty/pkg/auth"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

type UserHandler struct {
	service IUserService
}

func NewUserHandler(service IUserService) *UserHandler {
	return &UserHandler{service: service}
}

// GetAllUsers godoc
// @Summary List all users with pagination
// @Description Get all registered users (paginated)
// @Tags users
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param pageSize query int false "Number of items per page" default(10)
// @Success 200 {object} PaginatedUsersResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	allUsers, total, err := h.service.GetAllUsers(c.Context(), page, pageSize)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Map to response DTO
	users := make([]UserResponse, 0, len(allUsers))
	for _, user := range allUsers {
		users = append(users, UserResponse{
			Id:       uint(user.ID),
			UserName: user.UserName,
			Email:    user.Email,
		})
	}

	// Build paginated response
	return c.JSON(PaginatedUsersResponse{
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: (total + pageSize - 1) / pageSize,
		Data:       users,
	})
}

// GetUser godoc
// @Summary Create a new user
// @Description Create user with username and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserResponse true "User data"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [post]
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	userModel, err := h.service.GetUser(c.Context(), types.UserId(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}

	var shortsResponse []shortener.ShortResponse
	if len(userModel.Shorts) > 0 {
		shortsResponse = make([]shortener.ShortResponse, 0, len(userModel.Shorts))
		for _, shortModel := range userModel.Shorts {
			shortsResponse = append(shortsResponse, shortener.ShortResponse{
				Id:          shortModel.ID,
				OriginalUrl: shortModel.OriginalUrl,
				ShortUrl:    shortModel.ShortUrl,
			})
		}
	}

	userRespond := UserResponse{
		Id:       uint(userModel.ID),
		UserName: userModel.UserName,
		Email:    userModel.Email,
		Shorts:   shortsResponse,
	}
	return c.JSON(userRespond)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password, returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body UserLoginRequest true "User login credentials"
// @Success      200 {object} UserLoginResponse "JWT token successfully generated"
// @Failure      400 {object} map[string]string "Invalid request body"
// @Failure      401 {object} map[string]string "Invalid credentials"
// @Failure      500 {object} map[string]string "Could not create token"
// @Router       /api/v1/login [post]
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req UserLoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	err := validate.Struct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	user, err := h.service.VerifyCredentials(c.Context(), req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	signedToken, err := auth.GenerateJWTToken(user.ID, user.Role, time.Now().Add(time.Hour*72).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})
	}
	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    signedToken,
		Expires:  time.Now().Add(72 * time.Hour),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{Token: signedToken})
}

// CreateUser godoc
// @Summary      Register a new user
// @Description  Creates a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body UserRegisterRequest true "User registration data"
// @Success      201 {object} UserResponse
// @Failure      400 {object} map[string]string "Invalid request"
// @Failure      500 {object} map[string]string "Failed to create user"
// @Router       /api/v1/register [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var createReq UserRegisterRequest
	if err := c.BodyParser(&createReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	err := validate.Struct(createReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	createdUser, err := h.service.CreateUser(c.Context(), createReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update username or email of a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UserUpdateRequest true "Updated user data"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var updateReq UserUpdateRequest
	if err := c.BodyParser(&updateReq); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	reqErr := validate.Struct(updateReq)
	if reqErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": reqErr.Error()})
	}

	if err := h.service.UpdateUser(c.Context(), types.UserId(id), updateReq); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Remove a user by ID
// @Tags users
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	if err := h.service.DeleteUser(c.Context(), types.UserId(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
