package http

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func parsePaginationParams(ctx *fiber.Ctx) (p, pSize, limit, offset int) {
	page := ctx.Query("page", "1")
	pageSize := ctx.Query("pageSize", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil || pageSizeInt < 1 {
		pageSizeInt = 10
	}

	// -1 because offset starts from 0
	offset = (pageInt - 1) * pageSizeInt
	return pageInt, pageSizeInt, pageSizeInt, offset
}
