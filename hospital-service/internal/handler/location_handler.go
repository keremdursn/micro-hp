package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"hospital-service/internal/usecase"
)

type LocationHandler struct {
	usecase usecase.LocationUsecase
}

func NewLocationHandler(usecase usecase.LocationUsecase) *LocationHandler {
	return &LocationHandler{usecase: usecase}
}

// ListCities godoc
// @Summary     Tüm şehirleri listeler
// @Description Returns all cities (il)
// @Tags        Location
// @Produce     json
// @Success     200 {array} dto.CityLookup
// @Router      /api/cities [get]
func (h *LocationHandler) ListCities(c *fiber.Ctx) error {
	resp, err := h.usecase.ListAllCities()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}

// ListDistrictsByCity godoc
// @Summary     Seçili ile ait ilçeleri listeler
// @Description Returns all districts (ilçe) for a given city
// @Tags        Location
// @Produce     json
// @Param       city_id query int true "City ID"
// @Success     200 {array} dto.DistrictLookup
// @Failure     400 {object} map[string]string
// @Router      /api/districts [get]
func (h *LocationHandler) ListDistrictsByCity(c *fiber.Ctx) error {
	cityID, err := strconv.ParseUint(c.Query("city_id", "0"), 10, 64)
	if err != nil || cityID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid city_id"})
	}
	resp, err := h.usecase.ListDistrictsByCity(uint(cityID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
