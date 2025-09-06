package user

import (
	"strconv"
	"time"

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
// @Summary List all users
// @Description Get all registered users
// @Tags users
// @Produce json
// @Success 200 {array} UserResponse
// @Failure 500 {object} map[string]string
// @Router /api/v1/users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	allUsers, err := h.service.GetAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(allUsers)
}

// GetUser godoc
// @Summary Create a new user
// @Description Create user with username and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserCreateRequest true "User data"
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

	userResp, err := h.service.GetUser(c.Context(), types.UserId(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err})
	}

	return c.JSON(userResp)
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

	signedToken, err := auth.GenerateJWTToken(user.ID, time.Now().Add(time.Hour*72).Unix())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create token"})
	}

	return c.Status(fiber.StatusOK).JSON(UserLoginResponse{Token: signedToken})
}

// CreateUser godoc
// @Summary Get a user by ID
// @Description Fetch a user given their ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/v1/register [get]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var createReq UserCreateRequest
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
