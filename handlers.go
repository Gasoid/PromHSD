package main

import (
	"errors"
	"fmt"
	"net/http"
	"promhsd/db"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type readJsonPayload struct {
	Id      string             `json:"id"`
	Time    time.Time          `json:"time"`
	Name    string             `json:"name" binding:"required"`
	Entries []entryJsonPayload `json:"entries" binding:"required"`
}

type createJsonPayload struct {
	Name    string             `json:"name" binding:"required"`
	Entries []entryJsonPayload `json:"entries" binding:"required"`
}

type entryJsonPayload struct {
	Targets string `json:"targets" binding:"required"`
	Labels  string `json:"labels" binding:"required"`
}

type updateJsonPayload struct {
	Name    string             `json:"name" binding:"required"`
	Entries []entryJsonPayload `json:"entries" binding:"required"`
}

func (p *createJsonPayload) validate() (*db.Target, error) {
	err := &db.ValidationError{Text: "Validation failed"}
	if p.Name == "" {
		err.Text = "Field name is empty"
		return nil, err
	}
	t := db.NewTarget()
	t.Name = p.Name
	for _, e := range p.Entries {
		targets := strings.Split(e.Targets, ",")
		entry := db.NewEntry()
		entry.Targets = targets
		labels := strings.Split(e.Labels, ",")
		for _, l := range labels {
			kv := strings.Split(l, "=")
			if len(kv) != 2 {
				err.Text = "Labels are invalid"
				return nil, err
			}
			entry.Labels[kv[0]] = kv[1]
		}
		t.Entries = append(t.Entries, *entry)
	}
	return t, nil
}

func (p *updateJsonPayload) validate() (*db.Target, error) {
	err := &db.ValidationError{Text: "Validation failed"}
	if p.Name == "" {
		err.Text = "Field name is empty"
		return nil, err
	}
	t := db.NewTarget()
	t.Name = p.Name
	for _, e := range p.Entries {
		targets := strings.Split(e.Targets, ",")
		entry := db.NewEntry()
		entry.Targets = targets
		labels := strings.Split(e.Labels, ",")
		for _, l := range labels {
			kv := strings.Split(l, "=")
			if len(kv) != 2 {
				err.Text = "Labels are invalid"
				return nil, err
			}
			entry.Labels[kv[0]] = kv[1]
		}
		t.Entries = append(t.Entries, *entry)
	}
	return t, nil
}

func convertToJson(t *db.Target) *readJsonPayload {
	r := &readJsonPayload{Name: t.Name, Id: t.ID.String(), Time: t.Time, Entries: make([]entryJsonPayload, 0, len(t.Entries))}
	for _, entry := range t.Entries {
		labels := make([]string, 0, len(entry.Labels))
		for k, v := range entry.Labels {
			labels = append(labels, fmt.Sprintf("%s=%s", k, v))
		}
		e := entryJsonPayload{Targets: strings.Join(entry.Targets, ","), Labels: strings.Join(labels, ",")}
		r.Entries = append(r.Entries, e)
	}
	return r
}

// sourcesHandler godoc
// @Summary      getTargetsHandler
// @Description  returns targets
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  []db.Target{}
// @Router       /targets/ [get]
func getTargetsHandler(c *gin.Context) {
	targets := []db.Target{}
	err := dbService.List(&targets)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"targets": targets})
}

// sourcesHandler godoc
// @Summary      createTargetHandler
// @Description  creates target, returns id
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  "id"
// @Router       /target/ [post]
// @Param        payload  body createJsonPayload  true  "name"
func createTargetHandler(c *gin.Context) {
	payload := createJsonPayload{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "make sure name, targets and labels are sent"})
		return
	}
	t, err := payload.validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err_fields": err.Error()})
		return
	}
	err = dbService.Create(t)
	if err != nil {
		if errors.As(err, &db.ErrValidation) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
	c.JSON(http.StatusOK, gin.H{"id": t.ID})
}

// sourcesHandler godoc
// @Summary      getTargetHandler
// @Description  returns target
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  readJsonPayload
// @Param        id  path     string  true  "target id"
// @Router       /target/{id} [get]
func getTargetHandler(c *gin.Context) {
	t := db.NewTarget()
	t.ID = db.ID(c.Param("id"))
	err := dbService.Get(t)
	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		if errors.As(err, &db.ErrValidation) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return

	}
	c.JSON(http.StatusOK, gin.H{"target": convertToJson(t)})
}

// sourcesHandler godoc
// @Summary      updateTargetHandler
// @Description  returns id
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  "id"
// @Param        id  path     string  true  "target id"
// @Router       /target/{id} [post]
// @Param        payload  body updateJsonPayload  true  "name"
func updateTargetHandler(c *gin.Context) {
	payload := updateJsonPayload{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "make sure name, targets and labels are sent"})
		return
	}
	t, err := payload.validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err_fields": err.Error()})
		return
	}
	t.ID = db.ID(c.Param("id"))
	err = dbService.Update(t)
	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		if errors.As(err, &db.ErrValidation) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": t.ID})
}

// sourcesHandler godoc
// @Summary      removeTargetHandler
// @Description  delete item by id
// @Produce      json
// @Accept       json
// @Success      200  {array}  string  "id"
// @Param        id  path     string  true  "target id"
// @Router       /target/{id} [delete]
func removeTargetHandler(c *gin.Context) {
	t := db.NewTarget()
	t.ID = db.ID(c.Param("id"))
	err := dbService.Delete(t)
	if err != nil {
		if errors.As(err, &db.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		if errors.As(err, &db.ErrValidation) {
			c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

func prometheusHandler(c *gin.Context) {
	t := db.NewTarget()
	t.ID = db.ID(c.Param("id"))
	err := dbService.Get(t)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"err": err.Error()})
		return
	}

	c.JSON(http.StatusOK, t.Entries)
}

func healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{})
}
