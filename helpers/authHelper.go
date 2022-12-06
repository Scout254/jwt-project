package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CheckUserType(c *gin.Context, role string)(err error){
	userType := c.GetString("user_type")
	err = nil

	if userType != role{
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(c *gin.Context, userId string)(err error){
	userType := c.GetString("user_type")
	uid := c.GetString("uid")
	err = nil

	if userType == "USER" && uid != userId{
		err = errors.New("unauthorized to access this resource")
		return err
	}

	err = CheckUserType(c ,userType)
	return err
}

func UpdateAllTokens(signedToken  string , signedRefreshToken string , userId string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{"token",signedToken})
	updateObj = append(updateObj, bson.E{ "refresh_token",signedToken})

	Updated_at,_ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{"updated_at",Updated_at})

	upsert := true
	filter := bson.M{"user_id":userId}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_ , err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{ "$set",updateObj},
		},
		&opt,
	)
	defer cancel()
	if err != nil{
		log.Panic(err)
		return
	}
	
}

func ValidateToken(signedToken string)(claims *signedDetails, msg string){
	token , err := jwt.ParseWithClaims(
		signedToken,
		&signedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil{
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*signedDetails)
	if !ok{
		msg = fmt.Sprintf("the token is invalid")
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix(){
		msg = fmt.Sprintf("token is expired")
		msg = err.Error()
		return
	}
	return claims , msg
}