package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func FindDocuments(collectionName string, filter interface{}) ([]bson.M, error) {
	// Obtener la colección de la base de datos
	collection := DBClient.Database("albion-bot").Collection(collectionName)

	// Obtener los documentos
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		log.Printf("Error al buscar documentos: %v", err)
		return nil, err
	}

	// Convertir los resultados en un slice de mapas
	var results []bson.M
	if err = cursor.All(context.Background(), &results); err != nil {
		log.Printf("Error al leer documentos: %v", err)
		return nil, err
	}

	fmt.Printf("Documentos encontrados: %v", results)
	return results, nil
}

func InsertDocument(collectionName string, document interface{}) error {
	// Obtener la colección de la base de datos
	collection := DBClient.Database("albion-bot").Collection(collectionName)

	// Insertar el documento
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Printf("Error al insertar documento: %v", err)
		return err
	}

	fmt.Println("Documento insertado correctamente")
	return nil
}
