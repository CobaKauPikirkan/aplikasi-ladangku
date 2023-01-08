package controller

import (
	"context"

	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CobaKauPikirkan/aplikasi-ladangku/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddLadang() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		ladang, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}

		var ladangs models.Ladang
		var request models.RequestLadang
		ladangs.LadangId = primitive.NewObjectID()
		err = c.BindJSON(&request)
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		layout := "2006-01-02"
		
		ladangs.Name = request.Name
		ladangs.Kepadatan_tanaman = request.Kepadatan_tanaman
		ladangs.Luas_ladang = request.Luas_ladang

		ladangs.Tanggal_tanam, _ = time.Parse(layout, request.Tanggal_tanam)
		ladangs.Perkiraan_panen, _ = time.Parse(layout, request.Perkiraan_panen)
		ladangs.Komoditas = make([]models.Commodity, 0)
		ladangs.Todolist = make([]models.Todolist, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		match_filter := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: ladang}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$ladang"}}}}
		group := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$ladang_id"}, {Key: "count",Value: bson.D{primitive.E{Key: "$sum",Value: 1 }}}}}}
	
		pointcursor, err := commodityCollection.Aggregate(ctx, mongo.Pipeline{match_filter, unwind, group})
		if err != nil {
			c.IndentedJSON(500, "internal server error")
			return
		}
		
		var ladangInfo []bson.M

		err = pointcursor.All(ctx, &ladangInfo)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		filter := bson.D{primitive.E{Key: "_id", Value: ladang}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "ladang", Value: ladangs}}}}
		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		c.IndentedJSON(200, result)

		defer cancel()
		ctx.Done()
	}
}

func AddCommodityToLadang() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}

		commodity_id := c.Query("commodity")
		if commodity_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}
		commodity, err := primitive.ObjectIDFromHex(commodity_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}

		ladang_id := c.Query("ladang")
		if ladang_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		ctx, cancel:= context.WithTimeout(context.Background(), 100 *time.Second)
		defer cancel()


		var commoditycart []models.Commodity

		searchfromdb, err := commodityCollection.Find(ctx, bson.M{"_id": commodity})
		if err != nil {
			log.Println(err)
			return 
		}
		
		err = searchfromdb.All(ctx, &commoditycart)
		if err != nil {
			log.Println(err)
			return 
		}

		fmt.Println("ini komoditas ",commoditycart)

		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.M{"$push": bson.M{"ladang."+ ladang_id +".komoditas": bson.M{"$each": commoditycart}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully Added to the ladang")
	}
}

func GetLadangAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}

		var ladangs models.User

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

	 	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&ladangs)
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}

		defer cancel()
		c.IndentedJSON(200, ladangs.Ladang)
	}
}
func GetLadangByarray() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}
		ladang := c.Query("ladang")
		if ladang == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		var ladangs models.User

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

	 	err = userCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: userId}}).Decode(&ladangs)
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}
		no, err :=strconv.Atoi(ladang)
		if err != nil {
			log.Println(err)
			return
		}
		defer cancel()
		c.IndentedJSON(200, ladangs.Ladang[no])
	}
}

func DeleteCommodity() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		userId, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server erro")
			return
		}

		ladang := c.Query("ladang")
		if ladang == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		ladangempty := make([]models.Commodity, 0)
		filter3 := bson.D{primitive.E{Key: "_id", Value: userId}}
		update3 := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "ladang."+ladang+".komoditas", Value: ladangempty}}}}

		_, err  = userCollection.UpdateOne(ctx, filter3, update3)
		if err != nil {
			log.Println(err)
			return 
		}

		defer cancel()
		c.IndentedJSON(200, "succesfully delete commodity")
	}
}