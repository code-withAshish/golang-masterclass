package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io"
	"masterclass/database"
	"masterclass/models"
	"net/http"
	"os"
)

type Response struct {
	Message string `json:"message"`
	Name    string `json:"username"`
}

type QueryParams struct {
	Query  string `json:"query" query:"query"`
	Ashish string `json:"ashish" query:"ashish"`
	Nehal  string `json:"nehal" query:"nehal"`
}

// echo is the express js of golang
func main() {
	// making a new instance of echo
	app := echo.New()

	db := database.ConnectToDB()
	defer db.Close(context.Background())
	//registering middlewares
	app.Use(middleware.Logger())
	app.Use(middleware.Recover())

	app.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, &Response{
			Message: "Hello World",
			Name:    "Ashish",
		})
	})

	app.GET("/second", func(ctx echo.Context) error {
		var queries QueryParams
		err := ctx.Bind(&queries)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{})
		}
		return ctx.JSON(http.StatusOK, queries)
	})

	app.GET("/getNote/:id", func(ctx echo.Context) error {
		noteId := ctx.Param("id")
		fmt.Println(noteId)
		var fetchedNote models.Notes
		err := db.QueryRow(context.Background(), "select * from notes where id = $1", noteId).Scan(&fetchedNote.ID, &fetchedNote.Title, &fetchedNote.Description, &fetchedNote.Content, &fetchedNote.UserID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &Response{Message: err.Error()})
		}
		return ctx.JSON(http.StatusOK, fetchedNote)
	})

	app.POST("/createNote", func(ctx echo.Context) error {
		var note models.Notes
		err := ctx.Bind(&note)
		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{Message: "Form data extraction failed succesfully!!!"})
		}
		_, err = db.Query(context.Background(), "insert into notes(title, description, content, user_id) values ($1,$2,$3,$4)", note.Title, note.Description, note.Content, note.UserID)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &Response{Message: err.Error()})
		}
		return ctx.JSON(http.StatusOK, &Response{Message: "Note Created"})
	})

	app.GET("/getAllNotesOfUser/:id", func(ctx echo.Context) error {
		userId := ctx.Param("id")
		fmt.Println(userId)
		rows, err := db.Query(context.Background(), "select * from notes where user_id = $1", userId)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &Response{Message: err.Error()})
		}
		collectRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[models.Notes])
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, &Response{Message: err.Error()})
		}
		return ctx.JSON(http.StatusOK, collectRows)
	})
	app.POST("/form", func(ctx echo.Context) error {
		form, err := ctx.MultipartForm()

		if err != nil {
			return ctx.JSON(http.StatusBadRequest, &Response{
				Message: err.Error(),
				Name:    "error",
			})
		}

		files := form.File["files"]

		for _, file := range files {
			// Source
			src, err := file.Open()
			if err != nil {
				return err
			}
			defer src.Close()

			// Destination
			dst, err := os.Create(file.Filename)
			if err != nil {
				return err
			}
			defer dst.Close()

			// Copy
			if _, err = io.Copy(dst, src); err != nil {
				return err
			}

		}
		return ctx.JSON(http.StatusOK, &Response{Message: "file received"})
	})
	app.Logger.Fatal(app.Start(":3000"))
}
