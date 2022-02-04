package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"strconv"
)

func ParseID(c *gin.Context) (uint, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, err
	}

	if id < 0 {
		return 0, errors.Errorf("id(%d) is not uint", id)
	}
	return uint(id), nil
}
