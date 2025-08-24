package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// return page, pageSize
func getPaginationParams(ctx *fiber.Ctx) (int, int) {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("pageSize", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 || pageSizeInt > 100 {
		pageSizeInt = 20
	}

	return pageInt, pageSizeInt
}

func paginationToLimitOffset(page, pageSize int) (int, int) {
	// -1 because offset starts from 0
	offset := (page - 1) * pageSize
	return pageSize, offset
}
