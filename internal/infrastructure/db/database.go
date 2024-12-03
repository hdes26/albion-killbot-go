// infrastructure/database/database.go

package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DBClient es una variable global que mantendrá la conexión a MongoDB
var DBClient *mongo.Client

// Connect se encarga de establecer la conexión con MongoDB
func Connect(mongoURI string) error {
	// Configurar el cliente de MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Intentar conectar con la base de datos
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	// Comprobar la conexión
	err = client.Ping(context.Background(), nil)
	if err != nil {
		return err
	}

	// Asignar el cliente a la variable global DBClient
	DBClient = client
	log.Println("Conectado a MongoDB")

	return nil
}

// Disconnect cierra la conexión con MongoDB
func Disconnect() error {
	if DBClient != nil {
		return DBClient.Disconnect(context.Background())
	}
	return nil
}
