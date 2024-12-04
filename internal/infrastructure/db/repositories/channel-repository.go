package dbrepositories

import (
	"albion-killbot/internal/entities"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChannelRepository maneja las operaciones de acceso a la base de datos relacionadas con los canales.
type ChannelRepository struct {
	DB *mongo.Collection
}
type MongoDBClient struct {
	Client *mongo.Client
}

// NewChannelRepository crea un nuevo repositorio para Channels.
func NewChannelRepository(db *mongo.Collection) *ChannelRepository {
	return &ChannelRepository{DB: db}
}

// FindByChannelID busca un canal por su ID de canal.
func (cr *ChannelRepository) FindChannels() ([]*entities.Channel, error) {
	var channels []*entities.Channel

	// Realizamos la búsqueda en la colección, obteniendo el cursor
	cursor, err := cr.DB.Find(nil, bson.M{}) // Aquí obtienes el cursor
	if err != nil {
		// Si ocurre un error al buscar los canales
		return nil, fmt.Errorf("error al buscar canales: %w", err)
	}
	defer cursor.Close(nil) // Asegúrate de cerrar el cursor al final

	// Iteramos sobre el cursor para llenar el slice de canales
	for cursor.Next(nil) {
		var channel entities.Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, fmt.Errorf("error al decodificar canal: %w", err)
		}
		// Agregamos el canal decodificado al slice
		channels = append(channels, &channel)
	}

	// Verificamos si ocurrió un error durante la iteración
	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("error al iterar sobre los canales: %w", err)
	}

	return channels, nil
}

// FindByChannelID busca un canal por su ID de canal.
func (cr *ChannelRepository) FindByChannelID(channelID string) (*entities.Channel, error) {
	var channel entities.Channel
	err := cr.DB.FindOne(nil, bson.M{"channel_id": channelID}).Decode(&channel)
	if err != nil {
		// Si no se encuentra, retorna nil y el error
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("error al buscar el canal por ID: %w", err)
	}
	return &channel, nil
}

// Save guarda un nuevo canal en la base de datos.
func (cr *ChannelRepository) Save(channel *entities.Channel) (*entities.Channel, error) {
	fmt.Println(cr.DB)
	if cr.DB == nil {
		return nil, fmt.Errorf("la conexión a la base de datos no está inicializada")
	}
	_, err := cr.DB.InsertOne(nil, channel)
	if err != nil {
		return nil, fmt.Errorf("error al guardar el canal: %w", err)
	}
	return channel, nil
}

// Update actualiza un canal en la base de datos.
func (cr *ChannelRepository) Update(channel *entities.Channel) (*entities.Channel, error) {
	// Aquí se usa `ReplaceOne` para reemplazar el documento completo basado en el `channel_id`
	_, err := cr.DB.ReplaceOne(nil, bson.M{"channel_id": channel.ChannelID}, channel)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar el canal: %w", err)
	}
	return channel, nil
}

// FindOneAndUpdate encuentra un canal y lo actualiza. Es útil si necesitas cambiar solo ciertos campos.
func (cr *ChannelRepository) FindOneAndUpdate(filter bson.M, update bson.M) (*entities.Channel, error) {
	var channel entities.Channel
	updateOpts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	err := cr.DB.FindOneAndUpdate(nil, filter, update, updateOpts).Decode(&channel)
	if err != nil {
		return nil, fmt.Errorf("error al actualizar el canal: %w", err)
	}
	return &channel, nil
}
