package handler

import (
	"strconv"
	"fmt"
	"github.com/labstack/echo"
	"github.com/beewit/beekit/utils"
	"net/http"
)

func ImgCode(c echo.Context) error {
	d := make([]byte, 4)
	s := utils.NewLen(4)
	ss := ""
	d = []byte(s)
	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}
	c.Set("Content-Type", "image/png")
	c.Response().WriteHeader(http.StatusOK)
	utils.NewImage(d, 100, 40).WriteTo(c.Response().Writer)
	fmt.Println(ss)
	c.Response().Flush()
	return nil
}
