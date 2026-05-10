package db

import (
	"context"
	"errors"
	"time"

	"car-lock-system/backend/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingRepository struct {
	collection *mongo.Collection
}

func NewBookingRepository() *BookingRepository {
	return &BookingRepository{
		collection: DB.Collection("bookings"),
	}
}

func (br *BookingRepository) CreateBooking(ctx context.Context, booking *models.Booking) error {
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	result, err := br.collection.InsertOne(ctx, booking)
	if err != nil {
		return err
	}
	booking.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (br *BookingRepository) GetBookingByID(ctx context.Context, id string) (*models.Booking, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New("invalid booking id")
	}

	var booking models.Booking
	err = br.collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&booking)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("booking not found")
		}
		return nil, err
	}
	return &booking, nil
}

func (br *BookingRepository) UpdateBookingStatus(ctx context.Context, id string, status string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid booking id")
	}

	result, err := br.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("booking not found")
	}
	return nil
}

func (br *BookingRepository) UpdateBookingDetails(ctx context.Context, id, from, to, note string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.New("invalid booking id")
	}

	update := bson.M{"updated_at": time.Now()}
	if from != "" {
		update["from"] = from
	}
	if to != "" {
		update["to"] = to
	}
	if note != "" {
		update["note"] = note
	}

	result, err := br.collection.UpdateOne(
		ctx,
		bson.M{"_id": objID},
		bson.M{"$set": update},
	)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("booking not found")
	}
	return nil
}

func (br *BookingRepository) GetBookingsByUser(ctx context.Context, username string) ([]models.Booking, error) {
	cursor, err := br.collection.Find(ctx, bson.M{"username": username})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []models.Booking
	if err = cursor.All(ctx, &bookings); err != nil {
		return nil, err
	}
	return bookings, nil
}
