package controller

import (
	"context"
	"strconv"

	"log"
	"net/http"
	"time"

	"github.com/CobaKauPikirkan/aplikasi-ladangku/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddTodoList() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == "" {
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

		ladang_id := c.Query("ladang")
		if ladang_id == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var todolist models.Todolist

		todolist.TodolistId = primitive.NewObjectID()
		var request models.RequestTodo
		err = c.BindJSON(&request)
		if err != nil {
			c.IndentedJSON(http.StatusNotAcceptable, err.Error())
			return
		}

		layout := "2006-01-02"
		todolist.Date, _ = time.Parse(layout, request.Date)
		todolist.List = request.List
		todolist.Todolist_id = todolist.TodolistId.Hex()

		tes := ladang_id
		filter := bson.D{primitive.E{Key: "_id", Value: userId}}
		update := bson.M{"$push": bson.M{"ladang." + tes + ".todolist": todolist}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}

		c.IndentedJSON(200, "Successfully Added to the ladang")
	}
}

func GetAllTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == "" {
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

		ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		defer cancel()
		
		match:=	bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}}}
		unwind1:=bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$ladang"}}}}
		unwind2:=bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$ladang.todolist"}}}}
		project:= bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0},{Key: "todo", Value: "$ladang.todolist"}}}}
		
		var listing []bson.M
		currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{match,unwind1,unwind2, project})
		if err != nil {
			log.Println(err)
			return
		}
		ctx.Done()

		if err = currentresults.All(ctx, &listing); err != nil {
			log.Println(err)
			return
		}

		var slice []interface{}
		
		for _, result := range listing {
			slice = append(slice, result)
		}

		c.IndentedJSON(200, slice)

		defer currentresults.Close(ctx)
	}
	
}

func GetOneTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == "" {
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

		todo := c.Query("todo")
		if todo == "" {
			c.Header("content-type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid code"})
			c.Abort()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 *time.Second)
		defer cancel()
		
		match:=	bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: userId}}}}
		unwind1:=bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$ladang"}}}}
		unwind2:=bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$ladang.todolist"}}}}
		project:= bson.D{{Key: "$project", Value: bson.D{{Key: "_id", Value: 0},{Key: "todo", Value: "$ladang.todolist"}}}}
		
		var listing []bson.M
		currentresults, err := userCollection.Aggregate(ctx, mongo.Pipeline{match,unwind1,unwind2, project})
		if err != nil {
			log.Println(err)
			return
		}
		ctx.Done()

		if err = currentresults.All(ctx, &listing); err != nil {
			log.Println(err)
			return
		}

		var slice []interface{}
		
		for _, result := range listing {
			slice = append(slice, result)
		}

		no, err :=strconv.Atoi(todo)
		if err != nil {
			log.Println(err)
			return
		}
		c.IndentedJSON(200, slice[no])

		defer currentresults.Close(ctx)
	}
}



func EditTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}
		ladang := c.Query("ladang")
		if ladang == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}
		
		todo := c.Query("todo")
		if todo == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		list := c.Query("list")
		if list == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}
		var edittodolist models.Check
		
		if err := c.BindJSON(&edittodolist); err != nil{
			c.IndentedJSON(http.StatusBadRequest, err.Error())
		}

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		
		edittodolist.IsChecked = true
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "ladang."+ladang+".todolist."+todo+".list."+list+".ischecked",Value: edittodolist.IsChecked}}}}
		
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, "Something went wrong")
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "successfuly update")
	}
}


func DeleteTodo() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("userid")
		if user_id == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}

		usert_id, err := primitive.ObjectIDFromHex(user_id)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, "internal server error")
		}
		ladang := c.Query("ladang")
		if ladang == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}
		
		todoid := c.Query("todoid")
		if todoid == ""{
			c.Header("content-type", "application/json")
			c.JSON(http.StatusInternalServerError, gin.H{"error":"invalid"})
			c.Abort()
			return
		}
		
		
		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()
		
		
		filter := bson.D{primitive.E{Key: "_id", Value: usert_id}}
		update := bson.M{ "$pull": bson.M{"ladang."+ladang+".todolist":  bson.M{"todolist_id":todoid}}}
		
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.IndentedJSON(500, err)
			return
		}
		defer cancel()
		ctx.Done()
		c.IndentedJSON(200, "successfuly updated")
	}
}
