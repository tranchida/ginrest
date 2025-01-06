package main

import (
	"net/http"
	"tranchida/ginrest/pkg/message"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//store, err := message.NewSQLiteStore("message.db")
	store, err := message.NewMemoryMessageStore()
	if err != nil {
		panic(err)
	}
	handler := messageHandler{store: store}

	e.GET("/", homepage)
	e.GET("/messages", handler.list)
	e.POST("/messages", handler.add)
	e.GET("/messages/:id", handler.get)
	e.PUT("/messages/:id", handler.update)
	e.DELETE("/messages/:id", handler.remove)

	if err := e.Start(":8080"); err != nil {
		panic(err)
	}
}

func homepage(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

type messageHandler struct {
	store message.MessageStore
}

func (mh *messageHandler) list(c echo.Context) error {
	messages, err := mh.store.List()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, messages)
}

func (mh *messageHandler) add(c echo.Context) error {
	var msg message.Message
	if err := c.Bind(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	if err := mh.store.Add(msg.Id, msg); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, msg)
}

func (mh *messageHandler) get(c echo.Context) error {
	id := c.Param("id")
	msg, err := mh.store.Get(id)
	if err != nil {
		if err == message.ErrMessageNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "message not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, msg)
}

func (mh *messageHandler) update(c echo.Context) error {
	id := c.Param("id")
	var msg message.Message
	if err := c.Bind(&msg); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	msg.Id = id
	if err := mh.store.Update(id, msg); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, msg)
}

func (mh *messageHandler) remove(c echo.Context) error {
	id := c.Param("id")
	if err := mh.store.Remove(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
