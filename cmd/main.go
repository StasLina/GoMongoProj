package main

import (
	"LibrarySystem/db"
	"LibrarySystem/routes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getTemplatePath() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), "../templates")
}

// func collectHTMLFiles(root string) ([]string, error) {
// 	var files []string
// 	// Используем filepath.Walk для рекурсивного обхода директории
// 	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
// 		if err != nil {
// 			// Если произошла ошибка, возвращаем ее
// 			return err
// 		}
// 		// Проверяем, является ли файл и имеет ли он расширение .html
// 		if !info.IsDir() && filepath.Ext(path) == ".html" {
// 			files = append(files, path) // Добавляем файл в список
// 		}
// 		return nil // Возвращаем nil, чтобы продолжить обход
// 	})
// 	if err != nil {
// 		return nil, err // Возвращаем ошибку, если она произошла
// 	}
// 	return files, nil // Возвращаем список найденных файлов
// }

func collectHTMLFiles(root string) ([]string, error) {
	var files []string
	// Используем filepath.Walk для рекурсивного обхода директории
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Если произошла ошибка, возвращаем ее
			return err
		}
		// Проверяем, является ли файл и имеет ли он расширение .html
		if !info.IsDir() && filepath.Ext(path) == ".html" {
			// Получаем относительный путь от корневой директории
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err // Возвращаем ошибку, если не удалось получить относительный путь
			}
			files = append(files, "templates/"+relPath) // Добавляем относительный путь в список
		}
		return nil // Возвращаем nil, чтобы продолжить обход
	})
	if err != nil {
		return nil, err // Возвращаем ошибку, если она произошла
	}
	return files, nil // Возвращаем список найденных файлов
}

func main() {
	// Подключение к MongoDB
	db.ConnectDB()

	r := gin.Default()

	// Подключение шаблонов
	r.SetFuncMap(template.FuncMap{
		"toHex": func(id primitive.ObjectID) string {
			return id.Hex()
		},
		"renderContent": func(name string) func(interface{}, *template.Template, io.Writer) error {
			return func(ctx interface{}, tmpl *template.Template, wr io.Writer) error {
				return tmpl.ExecuteTemplate(wr, name, ctx)
			}
		},
	})

	// r.LoadHTMLGlob("templates/**/*.html")
	templatePath := getTemplatePath()
	htmlFiles, err := collectHTMLFiles(templatePath)
	if err != nil {
		panic("Failed to collect HTML templates: " + err.Error())
	}

	fmt.Println("Finded HTML templates:")
	for _, f := range htmlFiles {
		fmt.Println(" -", f)
	}
	// r.LoadHTMLFiles(htmlFiles...)

	// fmt.Printf("pathHtmlTemplate = %s\n", templatePath)
	// r.LoadHTMLGlob(filepath.Join(templatePath, "**/*"))
	r.LoadHTMLGlob("templates/*.html")

	r.GET("/test", func(c *gin.Context) {
		c.HTML(200, "layout.html", gin.H{
			"title": "Test Page",
		})
	})

	db := db.Client.Database("LibrarySystem")
	// Инициализация маршрутов
	routes.SetupRoutes(r, db)

	r.Run(":3000")
}
