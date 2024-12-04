package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient mantiene la conexión a la base de datos
type MongoDBClient struct {
	Client *mongo.Client
}

// Connect se encarga de establecer la conexión con MongoDB y devolver el cliente
func Connect(mongoURI string) (*MongoDBClient, error) {
	// Configurar el cliente de MongoDB
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Crear un contexto con timeout para evitar bloqueos largos
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Intentar conectar con la base de datos
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error al conectar a MongoDB: %w", err)
	}

	// Comprobar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	log.Println("Conectado a MongoDB")

	// Retornar el cliente envuelto en MongoDBClient
	return &MongoDBClient{Client: client}, nil
}

// Disconnect cierra la conexión con MongoDB
func (db *MongoDBClient) Disconnect() error {
	if db.Client != nil {
		return db.Client.Disconnect(context.Background())
	}
	return nil
}

// Obtener una referencia a una base de datos específica
func (db *MongoDBClient) GetDatabase(dbName string) *mongo.Database {
	return db.Client.Database(dbName)
}

func (client *MongoDBClient) GetCollection(collectionName string) *mongo.Collection {
	return client.Client.Database("yourDatabaseName").Collection(collectionName)
}
