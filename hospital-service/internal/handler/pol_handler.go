package handler

import (
	"strconv"

	"hospital-service/internal/config"
	"hospital-service/internal/dto"
	"hospital-service/internal/usecase"
	"hospital-shared/jwt"

	"github.com/gofiber/fiber/v2"
)

type PolyclinicHandler struct {
	polyclinicUsecase usecase.PolyclinicUsecase
	config            *config.Config
}

func NewPolyclinicHandler(polyclinicUsecase usecase.PolyclinicUsecase, cfg *config.Config) *PolyclinicHandler {
	return &PolyclinicHandler{
		polyclinicUsecase: polyclinicUsecase,
		config:            cfg,
	}
}

// ListAllPolyclinics godoc
// @Summary     Tüm poliklinikleri listeler
// @Description Returns all polyclinics
// @Tags        Polyclinic
// @Produce     json
// @Success     200 {array} dto.PolyclinicLookup
// @Router      /api/polyclinic/all [get]
func (h *PolyclinicHandler) ListAllPolyclinics(c *fiber.Ctx) error {
	resp, err := h.polyclinicUsecase.ListAllPolyclinics()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// AddHospitalPolyclinic godoc
// @Summary     Hastaneye poliklinik ekler
// @Description Adds a polyclinic to the hospital
// @Tags        Polyclinic
// @Accept      json
// @Produce     json
// @Param       polyclinic body dto.AddHospitalPolyclinicRequest true "Polyclinic info"
// @Success     201 {object} dto.HospitalPolyclinicResponse
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/add [post]
func (h *PolyclinicHandler) AddHospitalPolyclinic(c *fiber.Ctx) error {
	req := new(dto.AddHospitalPolyclinicRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := jwt.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.polyclinicUsecase.AddPolyclinicToHospital(req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListHospitalPolyclinic godoc
// @Summary     Hastanenin polikliniklerini listeler (sayfalı)
// @Description Lists hospital's polyclinics with pagination
// @Tags        Polyclinic
// @Produce     json
// @Param       page query int false "Page number"
// @Param       size query int false "Page size"
// @Success     200 {object} dto.HospitalPolyclinicListResponse
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/list [get]
func (h *PolyclinicHandler) ListHospitalPolyclinic(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	user := jwt.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.polyclinicUsecase.ListHospitalPolyclinic(user.HospitalID, page, size)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// RemoveHospitalPolyclinic godoc
// @Summary     Hastaneden poliklinik siler
// @Description Removes a polyclinic from the hospital
// @Tags        Polyclinic
// @Produce     json
// @Param       id path int true "Hospital Polyclinic ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/polyclinic/remove/{id} [delete]
func (h *PolyclinicHandler) RemoveHospitalPolyclinic(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid id"})
	}

	user := jwt.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	err = h.polyclinicUsecase.RemoveHospitalPolyclinic(uint(id), user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Polyclinic removed from hospital"})
}

func (h *PolyclinicHandler) GetHospitalPolyclinic(c *fiber.Ctx) error {
	hpID, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid ID"})
	}

	hp, err := h.polyclinicUsecase.GetHospitalPolyclinic(uint(hpID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "not found"})
	}

	return c.JSON(hp)
}
