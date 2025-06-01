package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"LibrarySystem/models"
	"LibrarySystem/repository"
)

type HardBookHandler struct {
	repo *repository.HardBookRepository
}

func NewHardBookHandler(repo *repository.HardBookRepository) *HardBookHandler {
	return &HardBookHandler{repo: repo}
}

func (h *HardBookHandler) GetAll(c *gin.Context) {
	books, err := h.repo.GetAllHardBooks()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{"books": books, "content_template": "hard_book"})
}

func (h *HardBookHandler) GetForm(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "hard_book_form",
		"title":            "Добавить книгу",
		"action":           "/hard_books",
		"method":           "POST",
		"submitText":       "Создать",
		"book":             models.HardBook{},
	})
}

func (h *HardBookHandler) Create(c *gin.Context) {
	inventoryID := c.PostForm("inventoryId")
	acqDateStr := c.PostForm("acquisitionDate")

	acqDate, err := time.Parse("2006-01-02", acqDateStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		return
	}

	var invID primitive.ObjectID
	if inventoryID == "" {
		invID = primitive.NewObjectID()
	} else {
		invID, _ = primitive.ObjectIDFromHex(inventoryID)
	}

	book := &models.HardBook{
		InventoryID:     invID,
		AcquisitionDate: primitive.NewDateTimeFromTime(acqDate),
	}

	if err := h.repo.CreateHardBook(book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/hard_books")
}

func (h *HardBookHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	book, err := h.repo.GetHardBookByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var str = fmt.Sprintf("/hard_books/%s/update", id)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "hard_book_form",
		"title":            "Изменить книгу",
		"action":           str,
		"method":           "POST",
		"submitText":       "Обновить",
		"book":             book,
		"toHex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
	})
}

func (h *HardBookHandler) Update(c *gin.Context) {
	id := c.Param("id")
	inventoryID := c.PostForm("inventoryId")
	acqDateStr := c.PostForm("acquisitionDate")

	acqDate, err := time.Parse("2006-01-02", acqDateStr)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Неверный формат даты"})
		return
	}

	var invID primitive.ObjectID
	if inventoryID == "" {
		invID = primitive.NewObjectID()
	} else {
		invID, _ = primitive.ObjectIDFromHex(inventoryID)
	}

	book := &models.HardBook{
		InventoryID:     invID,
		AcquisitionDate: primitive.NewDateTimeFromTime(acqDate),
	}

	if err := h.repo.UpdateHardBook(id, book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/hard_books")
}

func (h *HardBookHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.repo.DeleteHardBook(id); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusSeeOther, "/hard_books")
}
