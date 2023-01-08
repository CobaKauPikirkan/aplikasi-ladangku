package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CobaKauPikirkan/aplikasi-ladangku/database"
	"github.com/CobaKauPikirkan/aplikasi-ladangku/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.UserData(database.Client, "User")
var commodityCollection *mongo.Collection = database.Commodity(database.Client, "Commodity")
var Validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

func VerifyPassword(userPassword string, givenPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenPassword),[]byte(userPassword))
	valid := true
	msg := ""

	if err != nil {
		msg = "user or password Incorrect"
		valid = false
	}

	return valid, msg
}

func SignUp() gin.HandlerFunc {
	return func (c *gin.Context)  {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var user models.User

		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		validationErr := Validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr})
			return
		}
		//check duplication
		count, err := userCollection.CountDocuments(ctx, bson.M{"email":user.Email})
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user is Already exist"})
		}

		count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "This phone is Already exist"})
			return
		}

		time.LoadLocation("Indonesia")

		//password 
		password := HashPassword(*user.Password)
		user.Password = &password

		user.ID = primitive.NewObjectID()
		user.Created_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_At, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.User_ID = user.ID.Hex()

		user.Ladang = make([]models.Ladang, 0)

		result, inserterr := userCollection.InsertOne(ctx, user) 
		if inserterr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"user didnt created"})
			return
		}
		defer cancel()

		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user or password incorrect"})
			return
		}

		isValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()

		if !isValid {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			fmt.Println(msg)
			return
		}

		c.JSON(http.StatusFound, foundUser)
	}
}

func AddCommodity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		var commodity models.Commodity

		if err := c.BindJSON(&commodity); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		commodity.CommodityId = primitive.NewObjectID()
		_, anyerr := commodityCollection.InsertOne(ctx, commodity)
		if anyerr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": anyerr})
			return
		}
		defer cancel()
		c.JSON(200, "succesfully create product")
	}
}

func GetCommodityAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		var commodityList []models.Commodity
		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		cursor, err := commodityCollection.Find(ctx, bson.D{{}})
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "something went wrong, please try againt after a moment")
			return
		}

		err = cursor.All(ctx, &commodityList)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		defer cursor.Close(ctx)

		if err := cursor.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid")
			return
		}

		c.IndentedJSON(200, commodityList)
	}
}

func SearchCommodityByQuery() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchCommodity []models.Commodity
		queryParam := c.Query("name")
		if queryParam == ""{
			log.Println("query is Empty")
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalisd search index"})
			c.Abort()
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		searchquerydb, err := commodityCollection.Find(ctx, bson.M{"name": bson.M{"$regex":queryParam}})
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}

		err = searchquerydb.All(ctx, &searchCommodity)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(404, "invalid")
			return
		}
		defer searchquerydb.Close(ctx)

		if err := searchquerydb.Err(); err != nil{
			log.Println(err)
			c.IndentedJSON(400, "invalid request")
			return
		}

		defer cancel()
		c.IndentedJSON(200, searchCommodity)
	}
}
func CommodityById() gin.HandlerFunc {
	return func(c *gin.Context) {
		var searchCommodity models.Commodity
		queryParam := c.Query("id")
		if queryParam == ""{
			log.Println("query is Empty")
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalisd search index"})
			c.Abort()
			return
		}

		id, err := primitive.ObjectIDFromHex(queryParam)
		if err != nil {
			log.Println(err)
			return
		} 

		var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

	 	err = commodityCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&searchCommodity)
		if err != nil {
			c.IndentedJSON(404, "something went wrong in fetching the dbquery")
			return
		}

		// err = searchquerydb.All(ctx, &searchCommodity)
		// if err != nil {
		// 	log.Println(err)
		// 	c.IndentedJSON(404, "invalid")
		// 	return
		// }
		// defer searchquerydb.Close(ctx)

		// if err := searchquerydb.Err(); err != nil{
		// 	log.Println(err)
		// 	c.IndentedJSON(400, "invalid request")
		// 	return
		// }

		defer cancel()
		c.IndentedJSON(200, searchCommodity)
	}
}