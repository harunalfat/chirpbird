package helpers

import (
	"github.com/harunalfat/chirpbird/backend/entities"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IsExistsInEntityArray(data entities.Entities, searchFor string) bool {
	for i := 0; i < data.GetLength(); i++ {
		dataID := data.GetID(i)
		if dataID == searchFor {
			return true
		}
	}

	return false
}

// will return random ObjectID if invalid hex is provided
func ObjectIDFromHex(hex string) primitive.ObjectID {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		id = primitive.NewObjectID()
	}

	return id
}

func ObjectIDsFromHexes(hexes []string) []primitive.ObjectID {
	var result []primitive.ObjectID
	for _, hex := range hexes {
		id := ObjectIDFromHex(hex)
		result = append(result, id)
	}

	return result
}

func HexesFromObjectIDs(objectIDs []primitive.ObjectID) []string {
	var hexes []string
	for _, id := range objectIDs {
		hexes = append(hexes, id.Hex())
	}
	return hexes
}
