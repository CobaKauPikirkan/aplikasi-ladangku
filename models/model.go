package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct{
	ID 				primitive.ObjectID `json:"_id" bson:"_id"`
	First_Name      *string            `json:"first_name" validate:"required,min=2,max=30"`
	Last_Name       *string            `json:"last_name"  validate:"required,min=2,max=30"`
	Password        *string            `json:"password"   validate:"required,min=6"`
	Email           *string            `json:"email"      validate:"email,required"`
	Phone           *string            `json:"phone"      validate:"required"`
	Created_At      time.Time		   `json:"created_at"`
	Updated_At      time.Time          `json:"updated_at"`
	User_ID         string             `json:"user_id"`
	IsActive		bool				`json:"isActive" bson:"isActive"`
	Ladang			[]Ladang		   `json:"ladang" bson:"ladang"`
	
}

type Ladang struct{
	LadangId 		  primitive.ObjectID 	`bson:"_id"`
	Name              *string			  	`json:"name" bson:"name"`
	Kepadatan_tanaman *uint32		  		`json:"kepadatan_tanaman" bson:"kepadatan_tanaman"`
	Luas_ladang 	  *uint32				`json:"luas_ladang" bson:"luas_ladang"`
	Komoditas		  []Commodity			`json:"komoditas" bson:"komoditas"`
	Todolist		  []Todolist 			`json:"todolist" bson:"todolist"`
	Tanggal_tanam	  time.Time				`json:"tanggal_tanam" bson:"tanggal_tanam"`
	Perkiraan_panen   time.Time				`json:"perkiraan_panen" bson:"perkiraan_panen"`
}

type RequestLadang struct{
	LadangId 		  primitive.ObjectID 	`bson:"_id"`
	Name              *string			  	`json:"name" bson:"name"`
	Kepadatan_tanaman *uint32		  		`json:"kepadatan_tanaman" bson:"kepadatan_tanaman"`
	Luas_ladang 	  *uint32				`json:"luas_ladang" bson:"luas_ladang"`
	Komoditas		  []Commodity			`json:"komoditas" bson:"komoditas"`
	Todolist		  []Todolist 			`json:"todolist" bson:"todolist"`
	Tanggal_tanam	  string				`json:"tanggal_tanam" bson:"tanggal_tanam"`
	Perkiraan_panen   string				`json:"perkiraan_panen" bson:"perkiraan_panen"`
}

type Commodity struct{
	CommodityId		 primitive.ObjectID  	`bson:"_id"`
	Name 			 *string				`json:"name" bson:"name"`
	Panen			 *uint16				`json:"panen" bson:"panen"`
}

type Todolist struct {
	TodolistId 		primitive.ObjectID		`bson:"_id"`
	Date 			time.Time				`json:"date" bson:"date"`
	List			[]List					`json:"list" bson:"list"`
	Todolist_id		string					`json:"todolist_id"bson:"todolist_id"`
}

type RequestTodo struct{
	Date 			string				    `json:"date" bson:"date"`
	List			[]List					`json:"list" bson:"list"`
}

type Check struct{
	IsChecked       bool					`json:"ischecked" bson:"ischecked"`
}

type List struct{
	Todo            *string					`json:"todo" bson:"todo"`
	IsChecked       bool					`json:"ischecked" bson:"ischecked"`
}