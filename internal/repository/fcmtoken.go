package repository

import (
	"context"
	"fmt"
	"log"
	"main/internal/model"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Fcmtoken interface {
	CreateFCMTokenUser(c echo.Context, ctx context.Context, userid string) (string, error)
}

type fcmtoken struct {
	Db *mongo.Collection
}

func NewFcmToken(db *mongo.Collection) *fcmtoken {
	return &fcmtoken{
		db,
	}
}

func (r *fcmtoken) CreateFCMTokenUser(c echo.Context, ctx context.Context, userid string) (string, error) {
	fcmtoken := c.Request().Header.Get("FcmToken")

	filter := bson.M{"user_id": userid}

	var existingDocument model.FCMToken

	err := r.Db.FindOne(ctx, filter).Decode(&existingDocument)

	if err == mongo.ErrNoDocuments {
		// Dokumen tidak ditemukan, maka Anda bisa membuatnya
		newFcm := model.FCMToken{
			UserId:    userid,
			Fcm:       fcmtoken,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		_, err := r.Db.InsertOne(ctx, newFcm)
		if err != nil {
			log.Fatal(err)
		}
		return "Document created", err
	} else if err != nil {
		fmt.Println("error sni ", err.Error())
		return "Error created document", err
	} else {
		// Dokumen sudah ada, namun perlu pengecekan apakah fcm berbeda
		if existingDocument.Fcm != fcmtoken {
			// if not the same, do update
			update := bson.M{
				"$set": bson.M{"fcm": fcmtoken},
			}
			_, err := r.Db.UpdateOne(ctx, filter, update)
			if err != nil {
				return "Error updating document", err
			}
			return "Document updated", nil
		}
		return "Document already", nil
	}
}
