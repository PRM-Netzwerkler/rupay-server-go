package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/KevinGruber2001/rupay-bar-backend/controllers"
	dbCon "github.com/KevinGruber2001/rupay-bar-backend/db/sqlc"
	"github.com/KevinGruber2001/rupay-bar-backend/routes"
	"github.com/KevinGruber2001/rupay-bar-backend/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate"
	_ "github.com/lib/pq"
	"github.com/supertokens/supertokens-golang/recipe/dashboard"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"

	_ "github.com/KevinGruber2001/rupay-bar-backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/golang-migrate/migrate/database/postgres" // Import the postgres database driver
	_ "github.com/golang-migrate/migrate/source/file"       // Import the file source driver
)

var (
	server *gin.Engine
	db     *dbCon.Queries
	ctx    context.Context

	ArticleController            controllers.ArticleController
	ArticleTransactionController controllers.ArticleTransactionController
	ArticleTypeController        controllers.ArticleTypeController
	EventController              controllers.EventController
	TransactionController        controllers.TransactionController

	ArticleRoutes            routes.ArticleRoutes
	ArticleTransactionRoutes routes.ArticleTransactionRoutes
	ArticleTypeRoutes        routes.ArticleTypeRoutes
	EventRoutes              routes.EventRoutes
	TransactionRoutes        routes.TransactionRoutes
)

func runMigrations() {
	m, err := migrate.New(
		"file:///app/db/migration", // your migration folder
		"postgresql://root:secret@postgres:5432/rupay?sslmode=disable")
	if err != nil {
		log.Fatalf("could not create migrate instance: %v", err)
	}
	err = m.Up() // or m.Migrate(0) to apply all migrations
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("could not apply migrations: %v", err)
	}
	log.Println("migrations applied successfully")
}

func init() {

	ctx = context.TODO()
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not loadconfig: %v", err)
	}

	conn, err := sql.Open(config.DbDriver, config.DbSource)
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	db = dbCon.New(conn)

	fmt.Println("PostgreSql connected successfully...")

	// db migrations

	runMigrations()

	// Supertokens Init

	apiBasePath := "/auth"
	websiteBasePath := "/auth"
	st_err := supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			// https://try.supertokens.com is for demo purposes. Replace this with the address of your core instance (sign up on supertokens.com), or self host a core.
			ConnectionURI: "http://localhost:3567",
			// APIKey: <API_KEY(if configured)>,
		},
		AppInfo: supertokens.AppInfo{
			AppName:         "auth_test",
			APIDomain:       "https://localhost:5001",
			WebsiteDomain:   "http://localhost:3000",
			APIBasePath:     &apiBasePath,
			WebsiteBasePath: &websiteBasePath,
		},
		RecipeList: []supertokens.Recipe{
			emailpassword.Init(nil),
			session.Init(nil),
			dashboard.Init(nil),
		},
	})

	if st_err != nil {
		panic(err.Error())
	}

	ArticleController = *controllers.NewArticleController(db, ctx)
	ArticleRoutes = routes.NewRouteArticle(ArticleController)

	ArticleTypeController = *controllers.NewArticleTypeController(db, ctx)
	ArticleTypeRoutes = routes.NewRouteArticleType(ArticleTypeController)

	ArticleTransactionController = *controllers.NewArticleTransactionController(db, ctx)
	ArticleTransactionRoutes = routes.NewRouteArticleTransaction(ArticleTransactionController)

	EventController = *controllers.NewEventController(db, ctx)
	EventRoutes = routes.NewRouteEvent(EventController)

	TransactionController = *controllers.NewTransactionController(db, ctx)
	TransactionRoutes = routes.NewRouteTransaction(TransactionController)

	server = gin.Default()

	// CORS
	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders: append([]string{"content-type"},
			supertokens.GetAllCORSHeaders()...),
		AllowCredentials: true,
	}))

	// Adding the SuperTokens middleware
	// server.Use(func(c *gin.Context) {
	// 	supertokens.Middleware(http.HandlerFunc(
	// 		func(rw http.ResponseWriter, r *http.Request) {
	// 			c.Next()
	// 		})).ServeHTTP(c.Writer, c.Request)
	// 	// we call Abort so that the next handler in the chain is not called, unless we call Next explicitly
	// 	c.Abort()
	// })

}

// @title			Rupay Backend
// @version		1.0
// @description	Backend vor Rupay
// @host			localhost:8888
// @BasePath		/api
func main() {
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	router := server.Group("/api")

	// swagger middleware to serve the API docs
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	ArticleRoutes.ArticleRoute(router)
	ArticleTypeRoutes.ArticleTypeRoute(router)
	ArticleTransactionRoutes.ArticleTransactionRoute(router)
	EventRoutes.EventRoute(router)
	TransactionRoutes.TransactionRoute(router)

	// server.NoRoute(func(ctx *gin.Context) {
	//     ctx.JSON(http.StatusNotFound, gin.H{"status": "failed", "message": fmt.Sprintf("The specified route %s not found", ctx.Request.URL)})
	// })

	log.Fatal(server.Run(":" + config.ServerAddress))
}
