package handler

import (
	"strconv"

	"auth-service/internal/dto"
	"auth-service/internal/usecase"
	utilss "hospital-shared/util"

	"github.com/gofiber/fiber/v2"
)

type SubUserHandler struct {
	usecase usecase.SubUserUsecase
}

func NewSubUserHandler(usecase usecase.SubUserUsecase) *SubUserHandler {
	return &SubUserHandler{
		usecase: usecase,
	}
}

// CreateSubUser godoc
// @Summary     Alt kullanıcı ekler
// @Description Adds a new subuser (yetkili/çalışan)
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Param       subuser body dto.CreateSubUserRequest true "SubUser info"
// @Success     201 {object} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser [post]
func (h *SubUserHandler) CreateSubUser(c *fiber.Ctx) error {
	req := new(dto.CreateSubUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.CreateSubUser(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListSubUsers godoc
// @Summary     Alt kullanıcıları listeler
// @Description Lists all subusers for the hospital
// @Tags        SubUser
// @Produce     json
// @Success     200 {array} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser [get]
func (h *SubUserHandler) ListSubUsers(c *fiber.Ctx) error {
	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.ListSubUsers(user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// UpdateSubUser godoc
// @Summary     Alt kullanıcıyı günceller
// @Description Updates a subuser
// @Tags        SubUser
// @Accept      json
// @Produce     json
// @Param       id path int true "SubUser ID"
// @Param       subuser body dto.UpdateSubUserRequest true "SubUser info"
// @Success     200 {object} dto.SubUserResponse
// @Failure     400 {object} map[string]string
// @Router      /api/subuser/{id} [put]
func (h *SubUserHandler) UpdateSubUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	req := new(dto.UpdateSubUserRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.UpdateSubUser(uint(id), req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteSubUser godoc
// @Summary     Alt kullanıcıyı siler
// @Description Deletes a subuser
// @Tags        SubUser
// @Produce     json
// @Param       id path int true "SubUser ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/subuser/{id} [delete]
func (h *SubUserHandler) DeleteSubUser(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.DeleteSubUser(uint(id), user.HospitalID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User deleted"})
}
