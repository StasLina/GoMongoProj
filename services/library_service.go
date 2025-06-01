package services

import (
	"LibrarySystem/models"
	"LibrarySystem/repository"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllLibraries() ([]models.Library, error) {
	return repository.GetAllLibraries()
}

func CreateLibraryService(name, address string, locationNames, locationTypes []string) error {
	// Проверки
	if name == "" || address == "" {
		return errors.New("название и адрес обязательны")
	}

	if len(locationNames) != len(locationTypes) {
		return errors.New("количество названий и типов филиалов не совпадает")
	}

	// Создаём список филиалов
	var locations []models.Location
	for i := range locationNames {
		locations = append(locations, models.Location{
			LocationID: repository.GenerateObjectId(),
			Name:       locationNames[i],
			Type:       locationTypes[i],
		})
	}

	// Создаём новую библиотеку
	newLibrary := models.Library{
		Name:      name,
		Address:   address,
		Locations: locations,
	}

	// Сохраняем в БД
	return repository.CreateLibrary(newLibrary)
}

func GetLibraryByID(id string) (models.Library, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Library{}, err
	}

	library, err := repository.GetLibraryByID(objectID)
	if err != nil {
		return models.Library{}, err
	}

	return library, nil
}

func UpdateLibraryService(id string, name, address string, locationIDs, locationNames, locationTypes []string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	fmt.Printf("locationIDs: %v\n", locationIDs)
	fmt.Printf("locationNames: %v\n", locationNames)
	fmt.Printf("locationTypes: %v\n", locationTypes)

	fmt.Printf("len(locationIDs): %d\n", len(locationIDs))
	fmt.Printf("len(locationNames): %d\n", len(locationNames))
	fmt.Printf("len(locationTypes): %d\n", len(locationTypes))

	if len(locationIDs) != len(locationNames) || len(locationNames) != len(locationTypes) {
		return errors.New("количество полей не совпадает")
	}

	var locations []models.Location

	for i := range locationIDs {
		var locID primitive.ObjectID

		if i >= len(locationIDs) || locationIDs[i] == "" {
			locID = primitive.NewObjectID()
		} else {
			parsedID, err := primitive.ObjectIDFromHex(locationIDs[i])
			if err != nil {
				return err
			}
			locID = parsedID
		}

		locations = append(locations, models.Location{
			LocationID: locID,
			Name:       locationNames[i],
			Type:       locationTypes[i],
		})
	}

	library := models.Library{
		ID:        objectID,
		Name:      name,
		Address:   address,
		Locations: locations,
	}

	return repository.UpdateLibrary(objectID, library)
}

func DeleteLibraryService(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	return repository.DeleteLibrary(objectID)
}
