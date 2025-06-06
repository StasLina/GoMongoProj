package routes

import (
	"LibrarySystem/handlers"
	"LibrarySystem/repository"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupBooks(r *gin.Engine, db *mongo.Database) {
	repo := repository.NewBookRepository(db)
	handler := handlers.NewBookHandler(repo)

	r.GET("/books", func(c *gin.Context) { handler.GetAllBooks(c) })
	r.GET("/books/new", func(c *gin.Context) { handler.ShowCreateForm(c) })
	r.POST("/books/create", func(c *gin.Context) { handler.CreateBook(c) })
	r.GET("/books/:id/update", func(c *gin.Context) { handler.ShowEditForm(c) })
	r.POST("/books/:id/update", func(c *gin.Context) { handler.UpdateBook(c) })
	r.GET("/books/:id/delete", func(c *gin.Context) { handler.DeleteBook(c) })
}

func setupBooksRepo(r *gin.Engine, db *mongo.Database) {
	repo := repository.NewBookInventoryRepository(db)
	handler := handlers.NewBookInventoryHandler(repo)
	r.GET("/book_inventory", handler.ListBooks)
	r.GET("/book_inventory/new", handler.ShowCreateForm)
	r.POST("/book_inventory/new", handler.CreateBook)
	r.GET("/book_inventory/edit/:id", handler.ShowEditForm)
	r.POST("/book_inventory/edit/:id", handler.UpdateBook)
	r.GET("/book_inventory/delete/:id", handler.DeleteBook)
}

func setupHardBooksRoutes(r *gin.Engine, db *mongo.Database) {
	hardBookRepo := repository.NewHardBookRepository(db)
	handler := handlers.NewHardBookHandler(hardBookRepo)

	// Роуты
	r.GET("/hard_books", handler.GetAll)
	r.GET("/hard_books/new", handler.GetForm)
	r.POST("/hard_books", handler.Create)
	r.GET("/hard_books/:id", handler.GetByID)
	r.POST("/hard_books/:id/update", handler.Update)
	r.GET("/hard_books/:id/delete", handler.Delete)

}

func setupReaderRoutes(r *gin.Engine, db *mongo.Database) {
	readerRepo := repository.NewReaderRepository(db)
	readerHandler := handlers.NewReaderHandler(readerRepo)

	// Routes
	r.GET("/readers", readerHandler.GetAllReaders)
	r.GET("/readers/:id", readerHandler.GetReaderByID)
	r.GET("/readers/new", readerHandler.ShowCreateForm)
	r.POST("/readers", readerHandler.CreateReader)
	r.GET("/readers/:id/edit", readerHandler.ShowEditForm)
	r.POST("/readers/:id/update", readerHandler.UpdateReader)
	r.GET("/readers/:id/delete", readerHandler.DeleteReader)
}

func setupBookIssues(r *gin.Engine, db *mongo.Database) {
	repo := repository.NewBookIssueRepository(db)
	handler := handlers.NewBookIssueHandler(repo)
	// Маршруты
	r.GET("/bookissues", handler.GetAll)
	r.GET("/bookissues/new", handler.ShowCreateForm)
	r.POST("/bookissues/create", handler.CreateBookIssue)
	r.GET("/bookissues/:id/update", handler.GetByID)
	r.POST("/bookissues/:id/update", handler.UpdateBookIssue)
	r.GET("/bookissues/:id/delete", handler.DeleteBookIssue)
}

func setupInterLibraryRequest(r *gin.Engine, db *mongo.Database) {
	reqRepo := repository.NewInterlibraryRequestRepository(db)
	reqHandler := handlers.NewInterlibraryRequestHandler(reqRepo)

	r.GET("/requests", reqHandler.GetAll)
	r.GET("/requests/new", reqHandler.GetCreateForm)
	r.POST("/requests/create", reqHandler.Create)
	r.GET("/requests/:id", reqHandler.GetByID)
	r.GET("/requests/:id/edit", reqHandler.GetEditForm)
	r.POST("/requests/:id/update", reqHandler.Update)
	r.POST("/requests/:id/delete", func(c *gin.Context) {
		if c.PostForm("_method") == "DELETE" {
			reqHandler.Delete(c)
		} else {
			c.AbortWithStatusJSON(400, gin.H{"error": "Unknown method"})
		}
	})
}
func SetupRoutes(r *gin.Engine, db *mongo.Database) {
	r.GET("/", func(c *gin.Context) { c.HTML(200, "index.html", nil) })
	r.GET("/libraries", handlers.GetLibraries)
	r.GET("/libraries/new", handlers.ShowAddLibraryForm)
	r.POST("/libraries", handlers.CreateLibrary)
	r.GET("/libraries/:id/edit", handlers.EditLibrary)
	r.POST("/libraries/:id/update", handlers.UpdateLibrary)
	r.GET("/libraries/:id/delete", handlers.DeleteLibrary)

	setupBooks(r, db)
	setupBooksRepo(r, db)
	setupHardBooksRoutes(r, db)
	setupReaderRoutes(r, db)
	setupBookIssues(r, db)
	setupInterLibraryRequest(r, db)
}
