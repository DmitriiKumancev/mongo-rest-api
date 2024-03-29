package integrationtest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/DmitriiKumancev/mongoapi/internal/controllers"
	"github.com/DmitriiKumancev/mongoapi/internal/repository"
	"github.com/DmitriiKumancev/mongoapi/pkg/router"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDbInstance *mongo.Database
)

func TestMain(m *testing.M) {
	log.Println("setup is running")
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	populateDB()
	exitVal := m.Run()
	log.Println("teardown is running")
	_ = testDB.container.Terminate(context.Background())
	os.Exit(exitVal)
}

func InitializeTestRouter() *echo.Echo {
	postgreRepo := repository.New(testDbInstance)

	userController := controllers.New(postgreRepo)

	return router.Initialize(userController)
}

func TestGetPostsWithComments(t *testing.T) {
	apitest.New().
		Handler(InitializeTestRouter()).
		Get("/api/books").
		Header("content-type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.books`, 3)).
		BodyFromFile("responses/books.json").
		End()
}

func TestGetPostsByAuthorForDostoyevski(t *testing.T) {
	userId := "654e618a60034d917aa0ae63"
	apitest.New().
		Handler(InitializeTestRouter()).
		Get(fmt.Sprintf("/api/author/%s/books", userId)).
		Header("content-type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.books`, 2)).
		BodyFromFile("responses/books_dostoyevski.json").
		End()
}

func TestGetPostsByAuthorForMarcusAurelius(t *testing.T) {
	userId := "654e619760034d917aa0ae64"
	apitest.New().
		Handler(InitializeTestRouter()).
		Get(fmt.Sprintf("/api/author/%s/books", userId)).
		Header("content-type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.books`, 1)).
		BodyFromFile("responses/books_marcus.json").
		End()
}

func TestGetPosts_NonExistentAuthor(t *testing.T) {
	userId := "654e619760034d917aa0ae65"
	apitest.New().
		Handler(InitializeTestRouter()).
		Get(fmt.Sprintf("/api/author/%s/books", userId)).
		Header("content-type", "application/json").
		Expect(t).
		Status(http.StatusOK).
		Assert(jsonpath.Len(`$.books`, 0)).
		Assert(jsonpath.Equal("$.books", []interface{}{})).
		End()
}

func TestGetPosts_BadAuthorId(t *testing.T) {
	userId := "tesla"
	apitest.New().
		Handler(InitializeTestRouter()).
		Get(fmt.Sprintf("/api/author/%s/books", userId)).
		Header("content-type", "application/json").
		Expect(t).
		Status(http.StatusBadRequest).
		Assert(jsonpath.Equal("$.error", "the provided hex string is not a valid ObjectID")).
		End()
}

func TestCreatePostSuccess(t *testing.T) {
	apitest.New().
		Handler(InitializeTestRouter()).
		Post("/api/book").
		Header("content-type", "application/json").
		BodyFromFile("requests/create_book_success.json").
		Expect(t).
		Status(http.StatusCreated).
		BodyFromFile("responses/create_book_response.json").
		End()
}

func TestCreatePostAuthorNotExists(t *testing.T) {
	apitest.New().
		Handler(InitializeTestRouter()).
		Post("/api/book").
		Header("content-type", "application/json").
		BodyFromFile("requests/create_book_author_not_exists.json").
		Expect(t).
		Status(http.StatusNotFound).
		Assert(jsonpath.Equal("$.err", "author does not exist")).
		End()
}
