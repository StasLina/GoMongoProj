package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"LibrarySystem/models"
	"LibrarySystem/repository"
)

type InterlibraryRequestHandler struct {
	repo *repository.InterlibraryRequestRepository
}

func NewInterlibraryRequestHandler(repo *repository.InterlibraryRequestRepository) *InterlibraryRequestHandler {
	return &InterlibraryRequestHandler{repo: repo}
}

func (h *InterlibraryRequestHandler) GetAll(c *gin.Context) {
	requests, err := h.repo.GetAllInterlibraryRequests()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch requests"})
		return
	}
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "interlibrary_requests",
		"requests":         requests})
}

func (h *InterlibraryRequestHandler) GetCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "interlibrary_request_form",
		"isEdit":           false})
}

func (h *InterlibraryRequestHandler) Create(c *gin.Context) {
	det, _ := h.parseBookDetails(c)
	req := models.InterlibraryRequest{
		SourceLocation:      parseLocationItem(c, "source"),
		DestinationLocation: parseLocationItem(c, "destination"),
		RequestDate:         primitive.NewDateTimeFromTime(time.Now()),
		Status:              c.PostForm("status"),
		Details:             det,
	}

	result, err := h.repo.CreateInterlibraryRequest(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/requests/%s", result.InsertedID.(primitive.ObjectID).Hex()))
}

func (h *InterlibraryRequestHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	req, err := h.repo.GetInterlibraryRequestByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "interlibrary_request_form",
		"request":          req})
}

func (h *InterlibraryRequestHandler) GetEditForm(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	req, err := h.repo.GetInterlibraryRequestByID(id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.HTML(http.StatusOK, "layout.html", gin.H{
		"content_template": "interlibrary_request_form",
		"isEdit":           true, "request": req})
}

func (h *InterlibraryRequestHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	det, _ := h.parseBookDetails(c)
	req := models.InterlibraryRequest{
		SourceLocation:      parseLocationItem(c, "source"),
		DestinationLocation: parseLocationItem(c, "destination"),
		RequestDate:         primitive.NewDateTimeFromTime(time.Now()),
		Status:              c.PostForm("status"),
		Details:             det,
	}

	fmt.Print("Details length: ", len(req.Details), "\n")
	_, err := h.repo.UpdateInterlibraryRequest(id, req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to update request"})
		return
	}
	c.Redirect(http.StatusFound, "/requests")
}

func (h *InterlibraryRequestHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, _ := primitive.ObjectIDFromHex(idStr)
	_, err := h.repo.DeleteInterlibraryRequest(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete request"})
		return
	}
	c.Redirect(http.StatusFound, "/requests")
}

// --- Вспомогательные функции ---

func parseLocationItem(c *gin.Context, prefix string) models.LocationItem {
	locID := c.PostForm(prefix + "_location_id")
	libID := c.PostForm(prefix + "_library_id")

	var locationID, libraryID primitive.ObjectID

	if locID == "" {
		locationID = repository.GenerateObjectId()
	} else {
		locationID, _ = primitive.ObjectIDFromHex(locID)
	}

	if libID == "" {
		libraryID = repository.GenerateObjectId()
	} else {
		libraryID, _ = primitive.ObjectIDFromHex(libID)
	}

	return models.LocationItem{
		LocationID: locationID,
		LibraryID:  libraryID,
	}
}

// func parseBookDetails(c *gin.Context) []models.BookDetail {
// 	detailsCountStr := c.PostForm("details_count")
// 	detailsCount, _ := strconv.Atoi(detailsCountStr)
// 	details := make([]models.BookDetail, 0, detailsCount)

// 	for i := 0; i < detailsCount; i++ {
// 		bookIDStr := c.PostForm(fmt.Sprintf("detail[%d][book_id]", i))
// 		title := c.PostForm(fmt.Sprintf("detail[%d][title]", i))
// 		author := c.PostForm(fmt.Sprintf("detail[%d][author]", i))
// 		quantityStr := c.PostForm(fmt.Sprintf("detail[%d][quantity]", i))
// 		quantity, _ := strconv.Atoi(quantityStr)

// 		var bookID primitive.ObjectID
// 		if bookIDStr == "" {
// 			bookID = repository.GenerateObjectId()
// 		} else {
// 			bookID, _ = primitive.ObjectIDFromHex(bookIDStr)
// 		}

// 		details = append(details, models.BookDetail{
// 			BookID:   bookID,
// 			Title:    title,
// 			Author:   author,
// 			Quantity: quantity,
// 		})
// 	}

// 	fmt.Print(len(details))
// 	return details
// }

func (h *InterlibraryRequestHandler) parseBookDetails(c *gin.Context) ([]models.BookDetail, error) {
	var details []models.BookDetail
	i := 0

	for {
		bookIDStr := c.PostForm(fmt.Sprintf("detail[%d][book_id]", i))
		title := c.PostForm(fmt.Sprintf("detail[%d][title]", i))
		author := c.PostForm(fmt.Sprintf("detail[%d][author]", i))
		quantityStr := c.PostForm(fmt.Sprintf("detail[%d][quantity]", i))

		// Если хотя бы одно поле отсутствует — выходим из цикла
		if title == "" && author == "" && quantityStr == "" {
			break
		}

		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity at index %d: %v", i, err)
		}

		var bookID primitive.ObjectID
		if bookIDStr == "" {
			bookID = h.repo.GenerateObjectId()
		} else {
			var err error
			bookID, err = primitive.ObjectIDFromHex(bookIDStr)
			if err != nil {
				return nil, fmt.Errorf("invalid book_id at index %d: %v", i, err)
			}
		}

		details = append(details, models.BookDetail{
			BookID:   bookID,
			Title:    title,
			Author:   author,
			Quantity: quantity,
		})

		i++
	}

	return details, nil
}
