package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"LibrarySystem/models"
	"LibrarySystem/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReaderHandler struct {
	repo *repository.ReaderRepository
}

func NewReaderHandler(repo *repository.ReaderRepository) *ReaderHandler {
	return &ReaderHandler{repo: repo}
}

func (h *ReaderHandler) GetAllReaders(c *gin.Context) {
	readers, err := h.repo.GetAllReaders()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "reader",
		"Title":            "Читатели",
		"Readers":          readers,
		"toHex": func(oid primitive.ObjectID) string {
			return oid.Hex()
		},
	})
}

func (h *ReaderHandler) GetReaderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	reader, err := h.repo.GetReaderByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	c.JSON(http.StatusOK, reader)
}

func (h *ReaderHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "reader_form",
		"Title":            "Добавить читателя",
		"Action":           "/readers",
		"Method":           "POST",
	})
}

func (h *ReaderHandler) CreateReader(c *gin.Context) {
	reader := parseReaderForm(c, h.repo)
	if err := h.repo.CreateReader(reader); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/readers")
}

func (h *ReaderHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	reader, err := h.repo.GetReaderByID(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Reader not found"})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "reader_form",
		"Title":            "Редактировать читателя",
		"Action":           fmt.Sprintf("/readers/%s/update", idStr),
		"Method":           "POST",
		"Reader":           reader,
		"toHex": func(oid primitive.ObjectID) string {
			return oid.Hex()
		},
	})
}

func (h *ReaderHandler) UpdateReader(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	reader := parseReaderForm(c, h.repo)
	if err := h.repo.UpdateReader(id, reader); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/readers")
}

func (h *ReaderHandler) DeleteReader(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	if err := h.repo.DeleteReader(id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusSeeOther, "/readers")
}

func parseReaderForm(c *gin.Context, repo *repository.ReaderRepository) *models.Reader {
	fullName := c.PostForm("fullName")

	categoryID := c.PostForm("categoryId")
	categoryName := c.PostForm("categoryName")
	maxBorrowDays, _ := strconv.Atoi(c.PostForm("maxBorrowDays"))

	// Attributes
	attributeNames := c.PostFormArray("attributeName")
	attributeValues := c.PostFormArray("attributeValue")
	var attributes []models.Attribute
	for i := 0; i < len(attributeNames); i++ {
		if attributeNames[i] != "" && attributeValues[i] != "" {
			attributes = append(attributes, models.Attribute{
				Name:  attributeNames[i],
				Value: attributeValues[i],
			})
		}
	}

	// Subscriptions
	locationIDs := c.PostFormArray("locationId")
	startDates := c.PostFormArray("startDate")
	endDates := c.PostFormArray("endDate")
	maxBooksArr := c.PostFormArray("maxBooks")
	var subscriptions []models.Subscription
	for i := 0; i < len(locationIDs); i++ {
		locID := locationIDs[i]
		startDate, _ := time.Parse("2006-01-02", startDates[i])
		endDate, _ := time.Parse("2006-01-02", endDates[i])
		maxBooks, _ := strconv.Atoi(maxBooksArr[i])

		var objID primitive.ObjectID
		if locID == "" {
			objID = repo.GenerateObjectId()
		} else {
			objID, _ = primitive.ObjectIDFromHex(locID)
		}

		subscriptions = append(subscriptions, models.Subscription{
			LocationID: objID,
			StartDate:  primitive.NewDateTimeFromTime(startDate),
			EndDate:    primitive.NewDateTimeFromTime(endDate),
			MaxBooks:   maxBooks,
		})
	}

	return &models.Reader{
		FullName: fullName,
		Category: models.Category{
			CategoryID:    mustParseOID(categoryID, repo),
			Name:          categoryName,
			MaxBorrowDays: maxBorrowDays,
		},
		Attributes:    attributes,
		Subscriptions: subscriptions,
	}
}

func mustParseOID(idStr string, repo *repository.ReaderRepository) primitive.ObjectID {
	if idStr == "" {
		return repo.GenerateObjectId()
	}
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return repo.GenerateObjectId()
	}
	return id
}
