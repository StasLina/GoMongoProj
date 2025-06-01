package handlers

import (
	"context"
	"net/http"
	"strconv"

	"LibrarySystem/models"
	"LibrarySystem/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookInventoryHandler struct {
	repo *repository.BookInventoryRepository
}

func NewBookInventoryHandler(repo *repository.BookInventoryRepository) *BookInventoryHandler {
	return &BookInventoryHandler{repo: repo}
}

func (h *BookInventoryHandler) ListBooks(c *gin.Context) {
	ctx := c.Request.Context()
	books, err := h.repo.GetAllBooksInventory(ctx)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	data := map[string]interface{}{
		"content_template": "book_inventory",
		"books":            books,
	}

	c.HTML(200, "layout.html", data)
}

func (h *BookInventoryHandler) ShowCreateForm(c *gin.Context) {
	data := map[string]interface{}{
		"content_template": "book_inventory_create",
	}
	c.HTML(200, "layout.html", data)
}

func (h *BookInventoryHandler) CreateBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	qtyStr := c.PostForm("quantity")

	qty, _ := strconv.Atoi(qtyStr)

	bookIDStr := c.PostForm("bookId")
	locIDStr := c.PostForm("locationId")

	var bookID, locationID primitive.ObjectID
	if bookIDStr == "" {
		bookID = h.repo.GenerateObjectId()
	} else {
		id, _ := primitive.ObjectIDFromHex(bookIDStr)
		bookID = id
	}

	if locIDStr == "" {
		locationID = h.repo.GenerateObjectId()
	} else {
		id, _ := primitive.ObjectIDFromHex(locIDStr)
		locationID = id
	}

	book := &models.BookInventory{
		ID:         h.repo.GenerateObjectId(),
		BookID:     bookID,
		LocationID: locationID,
		Title:      title,
		Author:     author,
		Quantity:   qty,
	}

	if err := h.repo.CreateBookInventory(context.Background(), book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/book_inventory")
}

func (h *BookInventoryHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	book, err := h.repo.GetBookInventoryByID(context.Background(), id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	data := map[string]interface{}{
		"book":             book,
		"content_template": "book_inventory_edit", // Имя шаблона с content
	}
	c.HTML(200, "layout.html", data)
}

func (h *BookInventoryHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	title := c.PostForm("title")
	author := c.PostForm("author")
	qtyStr := c.PostForm("quantity")
	qty, _ := strconv.Atoi(qtyStr)

	bookIDStr := c.PostForm("bookId")
	locIDStr := c.PostForm("locationId")

	var bookID, locationID primitive.ObjectID
	if bookIDStr == "" {
		bookID = h.repo.GenerateObjectId()
	} else {
		id, _ := primitive.ObjectIDFromHex(bookIDStr)
		bookID = id
	}

	if locIDStr == "" {
		locationID = h.repo.GenerateObjectId()
	} else {
		id, _ := primitive.ObjectIDFromHex(locIDStr)
		locationID = id
	}

	book := &models.BookInventory{
		ID:         id,
		BookID:     bookID,
		LocationID: locationID,
		Title:      title,
		Author:     author,
		Quantity:   qty,
	}

	if err := h.repo.UpdateBookInventory(context.Background(), book); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/book_inventory")
}

func (h *BookInventoryHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	if err := h.repo.DeleteBookInventory(context.Background(), id); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusSeeOther, "/book_inventory")
}
