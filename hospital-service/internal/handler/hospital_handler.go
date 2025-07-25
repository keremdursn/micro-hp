package handler

import (
	"hospital-service/internal/dto"
	"hospital-service/internal/usecase"
	dt "hospital-shared/dto"
	utilss "hospital-shared/util"

	"github.com/gofiber/fiber/v2"
)

type HospitalHandler struct {
	usecase usecase.HospitalUsecase
}

func NewHospitalHandler(usecase usecase.HospitalUsecase) *HospitalHandler {
	return &HospitalHandler{
		usecase: usecase,
	}
}

// GetHospitalMe godoc
// @Summary     Giriş yapan kullanıcının hastane bilgisini getirir
// @Description Returns the hospital info of the authenticated user
// @Tags        Hospital
// @Produce     json
// @Success     200 {object} dto.HospitalResponse
// @Failure     401 {object} map[string]string
// @Failure     404 {object} map[string]string
// @Router      /api/hospital/me [get]
func (h *HospitalHandler) GetHospitalMe(c *fiber.Ctx) error {
	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	resp, err := h.usecase.GetHospitalByID(user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

func (h *HospitalHandler) CreateHospital(c *fiber.Ctx) error {
	var req dt.CreateHospitalRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.CreateHospital(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateHospitalMe godoc
// @Summary     Giriş yapan kullanıcının hastane bilgisini günceller
// @Description Updates the hospital info of the authenticated user
// @Tags        Hospital
// @Accept      json
// @Produce     json
// @Param       hospital body dto.UpdateHospitalRequest true "Hospital update info"
// @Success     200 {object} dto.HospitalResponse
// @Failure     400 {object} map[string]string
// @Failure     401 {object} map[string]string
// @Router      /api/hospital/me [put]
func (h *HospitalHandler) UpdateHospitalMe(c *fiber.Ctx) error {
	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	req := new(dto.UpdateHospitalRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	resp, err := h.usecase.UpdateHospital(user.HospitalID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
