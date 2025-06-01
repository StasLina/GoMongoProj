// handlers/book_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"LibrarySystem/models"
	"LibrarySystem/repository"
)

type BookHandler struct {
	repo *repository.BookRepository
}

func NewBookHandler(repo *repository.BookRepository) *BookHandler {
	return &BookHandler{repo: repo}
}

func (h *BookHandler) GetAllBooks(c *gin.Context) {
	books, err := h.repo.GetAllBooks()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "layout.html", gin.H{"content_template": "books", "books": books})
}

func (h *BookHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{"content_template": "book_form", "action": "/books/create"})
}

func (h *BookHandler) CreateBook(c *gin.Context) {
	title := c.PostForm("title")
	author := c.PostForm("author")
	pubYearStr := c.PostForm("publicationYear")

	pubYear, err := strconv.Atoi(pubYearStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "layout.html", gin.H{
			"content_template": "book_form",
			"error":            "Неверный год публикации",
			"title":            title,
			"author":           author,
			"year":             pubYearStr,
			"action":           "/books/create",
		})
		return
	}

	book := &models.Book{
		Title:           title,
		Author:          author,
		PublicationYear: pubYear,
	}

	if err := h.repo.CreateBook(book); err != nil {
		c.HTML(http.StatusInternalServerError, "layout.html", gin.H{
			"content_template": "book_form",
			"error":            "Ошибка при создании книги",
			"title":            title,
			"author":           author,
			"year":             pubYearStr,
			"action":           "/books/create",
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/books")
}

func (h *BookHandler) ShowEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	book, err := h.repo.GetBookByID(id)
	if err != nil {
		c.HTML(http.StatusNotFound, "error", gin.H{"error": "Книга не найдена"})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "book_form",
		"book":             book,
		"action":           fmt.Sprintf("/books/%s/update", id.Hex()),
		"id_hex":           id.Hex(),
		"title":            book.Title,
		"author":           book.Author,
		"year":             book.PublicationYear,
	})
}

func (h *BookHandler) UpdateBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	title := c.PostForm("title")
	author := c.PostForm("author")
	pubYearStr := c.PostForm("publicationYear")

	pubYear, err := strconv.Atoi(pubYearStr)
	if err != nil {
		c.HTML(http.StatusBadRequest, "layout.html", gin.H{
			"content_template": "book_form",
			"error":            "Неверный год публикации",
			"title":            title,
			"author":           author,
			"year":             pubYearStr,
			"action":           fmt.Sprintf("/books/%s/update", idStr),
		})
		return
	}

	updated := &models.Book{
		Title:           title,
		Author:          author,
		PublicationYear: pubYear,
	}

	if err := h.repo.UpdateBook(id, updated); err != nil {
		c.HTML(http.StatusInternalServerError, "layout.html", gin.H{
			"content_template": "book_form",
			"error":            "Ошибка при обновлении книги",
			"title":            title,
			"author":           author,
			"year":             pubYearStr,
			"action":           fmt.Sprintf("/books/%s/update", idStr),
		})
		return
	}

	c.Redirect(http.StatusSeeOther, "/books")
}

func (h *BookHandler) DeleteBook(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	if err := h.repo.DeleteBook(id); err != nil {
		c.HTML(http.StatusInternalServerError, "error", gin.H{"error": "Ошибка при удалении книги"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/books")
}
