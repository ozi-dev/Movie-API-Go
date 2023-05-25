package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"selfuelAPI/internal/models"
)

type MovieRepository struct {
	coll *mongo.Collection
}

// NewMovieRepository method is the factory method of MovieRepository struct.
func NewMovieRepository(mc *mongo.Client, dbName, collName string) *MovieRepository {
	return &MovieRepository{
		coll: mc.Database(dbName).Collection(collName),
	}
}

func (mr *MovieRepository) InsertMovie(ctx context.Context, movie *models.Movie) error {
	if _, err := mr.coll.InsertOne(ctx, movie); err != nil {
		return err
	}

	return nil
}

func (mr *MovieRepository) GetMovie(ctx context.Context, id string) (*models.Movie, error) {
	movie := &models.Movie{}

	filter := bson.D{{"id", id}}
	if err := mr.coll.FindOne(ctx, filter).Decode(movie); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, err
	}

	return movie, nil
}
