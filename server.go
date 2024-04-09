package main

import (
	"database/sql"
	"fmt"
	"go-echo-restApi/database"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Post struct {
	Id int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

func main() {
	app := echo.New()

	user, password, host, port, database := database.DbSourceName()

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, database))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	
	app.Use(middleware.CORS())

	app.GET("/post", func(c echo.Context) error {
		row, err := db.Query("SELECT id, title, content FROM Post")
		if err != nil {
			return c.String( http.StatusBadRequest ,err.Error())
		}

		result := Posts{}
		for row.Next(){
			post := Post{}
			err :=row.Scan(&post.Id, &post.Title, &post.Content)
			if err != nil {
				return c.String( http.StatusBadRequest,err.Error())
			}
			result.Posts = append(result.Posts, post)
		}

		return c.JSON(http.StatusOK, result)
	})

	app.GET("/post/:id", func(c echo.Context) error {
		id := c.Param("id")
		row := db.QueryRow("SELECT id, title, content FROM Post WHERE id = ?", id)

		if row.Err() != nil {
			return c.String(http.StatusBadRequest, "Invalid Query")
		} 

		post := Post{}
		row.Scan(&post.Id, &post.Title, &post.Content)

		return c.JSON(http.StatusOK, post)
	})

	app.POST("/post", func(c echo.Context) error {
		post := Post{}

		err := c.Bind(&post)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		_, err = db.Exec("INSERT INTO Post (title, content) VALUES (? , ?)", post.Title, post.Content)
		if err!= nil {
			c.String( http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusCreated, map[string]string{"message": "success"})
	})

	app.PUT("/post/:id", func(c echo.Context) error {
		id := c.Param("id")
		post := Post{}
		err := c.Bind(&post)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		_, err = db.Exec("UPDATE Post set title = ?, content = ? where id = ?", post.Title, post.Content, id )
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "data updated successfully"})
	})

	app.DELETE("/post/:id", func(c echo.Context) error {
		id := c.Param("id")
		_, err := db.Exec("DELETE FROM Post WHERE id = ?", id)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		
		return c.JSON(http.StatusAccepted, map[string]string{"message": "data deleted successfully"})

	})

	app.Logger.Fatal(app.Start(":4000"))
}