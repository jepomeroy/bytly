package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"bytly/model"
	"bytly/utils"
)

type errMessage struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func createBytly(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		msg := errMessage{}
		msg.Message = "error parsing a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	b := model.Bytly{}

	if err := json.Unmarshal(body, &b); err != nil {
		msg := errMessage{}
		msg.Message = "error parsing a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	if b.Random {
		b.Bytly = utils.RandomUrl(8)
	}
	b.Clicked = 0

	bytly, err := model.CreateBytly(b)

	if err != nil {
		msg := errMessage{}

		msg.Message = "error creating a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	return c.JSON(http.StatusOK, bytly)
}

func getBytlies(c echo.Context) error {
	bytlies, err := model.GetAllBytlies()

	if err != nil {
		msg := errMessage{}

		msg.Message = "error retrieving all bytlies"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusInternalServerError, msg)
	}

	return c.JSON(http.StatusOK, bytlies)
}

func deleteBytly(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	model.DeleteBytly(uint64(id))
	return c.NoContent(http.StatusNoContent)
}

func getBytly(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	bytly, err := model.GetBytlyById(uint64(id))

	if err != nil {
		msg := errMessage{}

		msg.Message = "error retrieving bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusInternalServerError, msg)
	}

	return c.JSON(http.StatusOK, bytly)
}

func redirectBytly(c echo.Context) error {
	shortcut := c.Param("bytly")
	bytly, err := model.GetBytlyByShortcut(shortcut)

	if err != nil {
		msg := errMessage{}

		msg.Message = "error retrieving bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusInternalServerError, msg)
	}

	bytly.Clicked += 1
	model.UpdateBytly(bytly)

	return c.Redirect(http.StatusTemporaryRedirect, bytly.Redirect)
}

func updateBytly(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)

	if err != nil {
		msg := errMessage{}
		msg.Message = "error parsing a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	b := model.Bytly{}

	if err := json.Unmarshal(body, &b); err != nil {
		msg := errMessage{}
		msg.Message = "error parsing a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	bytly, err := model.UpdateBytly(b)

	if err != nil {
		msg := errMessage{}

		msg.Message = "error creating a bytly"
		msg.Error = fmt.Sprint("{}", err)
		return c.JSON(http.StatusBadRequest, msg)
	}

	return c.JSON(http.StatusOK, bytly)
}

func Setup() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.GET("/r/:bytly", redirectBytly)
	e.DELETE("/bytly/:id", deleteBytly)
	e.GET("/bytly", getBytlies)
	e.GET("/bytly/:id", getBytly)
	e.PATCH("/bytly", updateBytly)
	e.POST("/bytly", createBytly)

	e.Logger.Fatal(e.Start(":5000"))
}
