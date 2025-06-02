package handlers

import (
	"LibrarySystem/services"

	"github.com/gin-gonic/gin"
)

func GetLibraries(c *gin.Context) {
	libraries, err := services.GetAllLibraries()
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(200, "layout.html", gin.H{
		"content_template": "libraries",
		"libraries":        libraries,
	})
}

func ShowAddLibraryForm(c *gin.Context) {
	c.HTML(200, "layout.html", gin.H{
		"content_template": "library_add",
	})
}

func CreateLibrary(c *gin.Context) {
	// Получаем данные из формы
	name := c.PostForm("name")
	address := c.PostForm("address")

	locationNames := c.PostFormArray("locationNames[]")
	locationTypes := c.PostFormArray("locationTypes[]")

	// Проверка на количество филиалов
	if len(locationNames) != len(locationTypes) {
		c.AbortWithStatusJSON(400, gin.H{"error": "Количество названий и типов филиалов не совпадает"})
		return
	}

	// Передаём данные в сервис
	err := services.CreateLibraryService(name, address, locationNames, locationTypes)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	// Перенаправляем
	c.Redirect(302, "/libraries")
}

func EditLibrary(c *gin.Context) {
	id := c.Param("id")
	library, err := services.GetLibraryByID(id)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.HTML(200, "layout.html", gin.H{
		"content_template": "library_edit",
		"Library":          library,
		"libraryID":        id,
	})
}

func UpdateLibrary(c *gin.Context) {
	id := c.Param("id")

	name := c.PostForm("name")
	address := c.PostForm("address")

	LocationIDS := c.PostFormArray("locationIDs[]")
	locationNames := c.PostFormArray("locationNames[]")
	locationTypes := c.PostFormArray("locationTypes[]")

	err := services.UpdateLibraryService(id, name, address, LocationIDS, locationNames, locationTypes)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(302, "/libraries")
}

func DeleteLibrary(c *gin.Context) {
	id := c.Param("id")

	err := services.DeleteLibraryService(id)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.Redirect(302, "/libraries")
}
