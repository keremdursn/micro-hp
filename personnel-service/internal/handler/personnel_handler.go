package handler

import (
	"strconv"

	utilss "hospital-shared/util"
	"personnel-service/internal/dto"
	"personnel-service/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type PersonnelHandler struct {
	usecase usecase.PersonnelUsecase
}

func NewPersonnelHandler(usecase usecase.PersonnelUsecase) *PersonnelHandler {
	return &PersonnelHandler{
		usecase: usecase,
	}
}

// ListAllJobGroups godoc
// @Summary     Tüm meslek gruplarını listeler
// @Description Returns all job groups
// @Tags        Personnel
// @Produce     json
// @Success     201 {array} dto.JobGroupLookup
// @Router      /api/personnel/job-groups [get]
func (h *PersonnelHandler) ListAllJobGroups(c *fiber.Ctx) error {
	resp, err := h.usecase.ListAllJobGroups()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// ListTitleByJobGroup godoc
// @Summary     Seçili meslek grubuna ait unvanları listeler
// @Description Returns all titles for a given job group
// @Tags        Personnel
// @Produce     json
// @Param       job_group_id query int true "Job Group ID"
// @Success     201 {array} dto.TitleLookup
// @Failure     400 {object} map[string]string
// @Router      /api/personnel/titles [get]
func (h *PersonnelHandler) ListTitleByJobGroup(c *fiber.Ctx) error {

	jobGroupID, err := strconv.ParseUint(c.Query("job_group_id", "0"), 10, 64)
	if err != nil || jobGroupID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid job_group_id"})
	}

	resp, err := h.usecase.ListTitleByJobGroup(uint(jobGroupID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(resp)
}

// AddStaff godoc
// @Summary     Personel ekler
// @Description Adds a new staff member
// @Tags        Personnel
// @Accept      json
// @Produce     json
// @Param       staff body dto.AddStaffRequest true "Staff info"
// @Success     201 {object} dto.StaffResponse
// @Failure     400 {object} map[string]string
// @Router      /api/personnel/staff [post]
func (h *PersonnelHandler) AddStaff(c *fiber.Ctx) error {
	var req dto.AddStaffRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.AddStaff(&req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// UpdateStaff godoc
// @Summary     Personel günceller
// @Description Updates a staff member
// @Tags        Personnel
// @Accept      json
// @Produce     json
// @Param       id path int true "Staff ID"
// @Param       staff body dto.UpdateStaffRequest true "Staff info"
// @Success     201 {object} dto.StaffResponse
// @Failure     400 {object} map[string]string
// @Router      /api/personnel/staff/{id} [put]
func (h *PersonnelHandler) UpdateStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid staff id"})
	}

	req := new(dto.UpdateStaffRequest)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.UpdateStaff(uint(id), req, user.HospitalID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

// DeleteStaff godoc
// @Summary     Personel siler
// @Description Deletes a staff member
// @Tags        Personnel
// @Produce     json
// @Param       id path int true "Staff ID"
// @Success     200 {object} map[string]string
// @Failure     400 {object} map[string]string
// @Router      /api/personnel/staff/{id} [delete]
func (h *PersonnelHandler) DeleteStaff(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid staff id"})
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.usecase.DeleteStaff(uint(id), user.HospitalID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Staff deleted"})
}

// ListStaff godoc
// @Summary     Personelleri listeler (filtreli ve sayfalı)
// @Description Lists staff with filters and pagination
// @Tags        Personnel
// @Produce     json
// @Param       page query int false "Page number"
// @Param       size query int false "Page size"
// @Param       first_name query string false "First name"
// @Param       last_name query string false "Last name"
// @Param       tc query string false "TC Kimlik No"
// @Param       job_group_id query int false "Job Group ID"
// @Param       title_id query int false "Title ID"
// @Success     201 {object} dto.StaffListResponse
// @Failure     400 {object} map[string]string
// @Router      /api/personnel/staff [get]
func (h *PersonnelHandler) ListStaff(c *fiber.Ctx) error {

	// 1. Query parametrelerinden filtreleri ve sayfa/size al
	page, _ := strconv.Atoi(c.Query("page", "1"))
	size, _ := strconv.Atoi(c.Query("size", "10"))
	firstName := c.Query("first_name", "")
	lastName := c.Query("last_name", "")
	tc := c.Query("tc", "")

	// job_group_id ve title_id parse edilir
	var jobGroupID *uint
	if v := c.Query("job_group_id", ""); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			jid := uint(id)
			jobGroupID = &jid
		}
	}
	var titleID *uint
	if v := c.Query("title_id", ""); v != "" {
		if id, err := strconv.ParseUint(v, 10, 64); err == nil && id > 0 {
			tid := uint(id)
			titleID = &tid
		}
	}

	user := utilss.GetUserInfo(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	resp, err := h.usecase.ListStaff(user.HospitalID, dto.StaffListFilter{
		FirstName:  firstName,
		LastName:   lastName,
		TC:         tc,
		JobGroupID: jobGroupID,
		TitleID:    titleID,
	}, page, size)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(resp)
}
