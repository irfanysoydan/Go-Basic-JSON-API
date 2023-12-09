package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

var books = []book{ //bu yapıya slice deniyor. Slice bir dizi gibidir.
	{ID: "1", Title: "Gullivers Travels", Author: "Jonathan Swift", Quantity: 10},
	{ID: "2", Title: "The Catcher in the Rye", Author: "J. D. Salinger", Quantity: 5},
	{ID: "3", Title: "To Kill a Mockingbird", Author: "Harper Lee", Quantity: 20},
}

func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookByID(id string) (*book, error) {

	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "please provide book id"})
		return
	}

	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "book is out of stock"})
		return
	}

	book.Quantity--

	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "please provide book id"})
		return
	}

	book, err := getBookByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	book.Quantity++

	c.IndentedJSON(http.StatusOK, book)
}

func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		return
	}

	books = append(books, newBook)

	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", gin.HandlerFunc(getBooks))
	router.GET("/books/:id", gin.HandlerFunc(bookById))
	router.POST("/books", gin.HandlerFunc(createBook))
	router.PATCH("/checkout", gin.HandlerFunc(checkoutBook))
	router.PATCH("/return", gin.HandlerFunc(returnBook))
	router.Run("localhost:8080")
}
