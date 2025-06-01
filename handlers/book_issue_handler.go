package handlers

import (
	"fmt"
	"net/http"
	"time"

	"LibrarySystem/models"
	"LibrarySystem/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookIssueHandler struct {
	repo *repository.BookIssueRepository
}

func NewBookIssueHandler(repo *repository.BookIssueRepository) *BookIssueHandler {
	return &BookIssueHandler{repo: repo}
}

// GetAllBookIssues — вывод списка выдач книг
func (h *BookIssueHandler) GetAll(c *gin.Context) {
	issues, err := h.repo.GetAllBookIssues(c.Request.Context())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Не удалось загрузить записи"})
		return
	}

	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "book_issues",
		"title":            "Выдачи книг",
		"issues":           issues,
		"toHex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
	})
}

// GetByID — получение конкретной записи по ID
func (h *BookIssueHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	issue, err := h.repo.GetBookIssueByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	action := fmt.Sprintf("/bookissues/%s/update", idStr)
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "book_issue_form",
		"title":            "Изменить запись",
		"action":           action,
		"method":           "POST",
		"submitText":       "Обновить",
		"issue":            issue,
		"toHex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
	})
}

// ShowCreateForm — отображает форму создания новой записи
func (h *BookIssueHandler) ShowCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "book_issue_form",
		"title":            "Новая выдача книги",
		"action":           "/bookissues/create",
		"method":           "POST",
		"submitText":       "Создать",
		"issue": models.BookIssue{
			IssueDate:  time.Now(),
			DueDate:    time.Now().AddDate(0, 1, 0), // +1 месяц
			ReturnDate: nil,
		},
		"toHex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
	})
}

// CreateBookIssue — создаёт новую запись
func (h *BookIssueHandler) CreateBookIssue(c *gin.Context) {
	if c.Request.Method == "POST" {
		c.Request.ParseForm()

		hardBookIDStr := c.PostForm("hardBookId")
		hardBookID, _ := primitive.ObjectIDFromHex(hardBookIDStr)
		if hardBookID.IsZero() {
			hardBookID = h.repo.GenerateObjectId()
		}

		bookIssueUserIDStr := c.PostForm("bookIssueUserId")
		bookIssueUserID, _ := primitive.ObjectIDFromHex(bookIssueUserIDStr)
		if bookIssueUserID.IsZero() {
			bookIssueUserID = h.repo.GenerateObjectId()
		}

		locationIDStr := c.PostForm("locationId")
		locationID, _ := primitive.ObjectIDFromHex(locationIDStr)

		issueDateStr := c.PostForm("issueDate")
		issueDate, _ := time.Parse("2006-01-02", issueDateStr)

		dueDateStr := c.PostForm("dueDate")
		dueDate, _ := time.Parse("2006-01-02", dueDateStr)

		returnDateStr := c.PostForm("returnDate")
		var returnDate *time.Time
		if returnDateStr != "" {
			rd, _ := time.Parse("2006-01-02", returnDateStr)
			returnDate = &rd
		}

		issue := models.BookIssue{
			ID: h.repo.GenerateObjectId(),
			Book: models.BookItem{
				HardBookID: hardBookID,
				Title:      c.PostForm("bookTitle"),
				Author:     c.PostForm("bookAuthor"),
			},
			BookIssueUser: models.BookIssueUser{
				BookIssuesID: bookIssueUserID,
				FullName:     c.PostForm("fullName"),
			},
			IssueDate:  issueDate,
			DueDate:    dueDate,
			ReturnDate: returnDate,
			LocationID: locationID,
			Status:     c.PostForm("status"),
		}

		err := h.repo.CreateBookIssue(c.Request.Context(), issue)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании записи"})
			return
		}
		c.Redirect(http.StatusSeeOther, "/bookissues")
	}
}

// UpdateBookIssue — обновляет существующую запись
func (h *BookIssueHandler) UpdateBookIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	c.Request.ParseForm()

	hardBookIDStr := c.PostForm("hardBookId")
	hardBookID, _ := primitive.ObjectIDFromHex(hardBookIDStr)
	if hardBookID.IsZero() {
		hardBookID = h.repo.GenerateObjectId()
	}

	bookIssueUserIDStr := c.PostForm("bookIssueUserId")
	bookIssueUserID, _ := primitive.ObjectIDFromHex(bookIssueUserIDStr)
	if bookIssueUserID.IsZero() {
		bookIssueUserID = h.repo.GenerateObjectId()
	}

	locationIDStr := c.PostForm("locationId")
	locationID, _ := primitive.ObjectIDFromHex(locationIDStr)

	issueDateStr := c.PostForm("issueDate")
	issueDate, _ := time.Parse("2006-01-02", issueDateStr)

	dueDateStr := c.PostForm("dueDate")
	dueDate, _ := time.Parse("2006-01-02", dueDateStr)

	returnDateStr := c.PostForm("returnDate")
	var returnDate *time.Time
	if returnDateStr != "" {
		rd, _ := time.Parse("2006-01-02", returnDateStr)
		returnDate = &rd
	}

	updated := models.BookIssue{
		Book: models.BookItem{
			HardBookID: hardBookID,
			Title:      c.PostForm("bookTitle"),
			Author:     c.PostForm("bookAuthor"),
		},
		BookIssueUser: models.BookIssueUser{
			BookIssuesID: bookIssueUserID,
			FullName:     c.PostForm("fullName"),
		},
		IssueDate:  issueDate,
		DueDate:    dueDate,
		ReturnDate: returnDate,
		LocationID: locationID,
		Status:     c.PostForm("status"),
	}

	err := h.repo.UpdateBookIssue(c.Request.Context(), id, updated)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении записи"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/bookissues")
}

// DeleteBookIssue — удаляет запись
func (h *BookIssueHandler) DeleteBookIssue(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)

	err := h.repo.DeleteBookIssue(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Redirect(http.StatusSeeOther, "/bookissues")
}
