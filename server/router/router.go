package router

import (
	"net/http"
	api_customer "server/api/customer"
	api_producer "server/api/producer"
	"server/service"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusOK)
			return
		}

		c.Next()
	}
}

func RouteAll(r *gin.Engine, s *service.Services) {
	api_router := r.Group("api")
	api_router.Use(CORSMiddleware())
	{
		api_router.GET("/ping", func(c *gin.Context) { c.String(200, "pong") })
		producer := api_router.Group("producer")
		{
			producer.GET("/export_data/:id",api_producer.ExportData(s))
			auth_producer := producer.Group("auth")
			{
				auth_producer.POST("/register", api_producer.RedisterProducer(s))
				auth_producer.POST("/signin", api_producer.SigninProducer(s))
			}

			products := producer.Group("products")
			{
				products.DELETE("remove/:product_id", api_producer.RemoveProduct(s))
				products.POST("add", api_producer.AddProduct(s))

			}
			buisness_producer := producer.Group("business")
			{
				buisness_producer.PUT("/:businessId/form", api_producer.HandlePutForm(s))
				buisness_producer.GET("/:businessId/form", api_producer.HandleGetForm(s))
				// buisness.GET("/details/:id", api_customer.GetBuisnessDetailsById(s))
				buisness_producer.GET("/reviews/:buisnessId", api_producer.GetReviews(s))
			}
			promocode := producer.Group("promocode")
			{
				promocode.PUT("/:businessId", api_producer.HandlePutPromocode(s))
			}
			check := producer.Group("check")
			{
				check.POST("generate", api_producer.GenerateCheck(s))
				check.GET("/:id", api_producer.GetCheck(s))
			}

		}
		user := api_router.Group("/user")
		{
			userAuth := user.Group("/auth")
			{
				userAuth.POST("/register", api_customer.HandleClientRegister(s))
				userAuth.POST("/signin", api_customer.HandleClientSignin(s))
			}
			user.POST("review", api_customer.Review(s))
			user.GET("/reviews/:user_id", api_customer.GetUserReviews(s))
			user.GET("/:customerId/promocodes", api_customer.HandleGetPromocodes(s))
		}
		business := api_router.Group("/business")
		{
			business.GET("/", api_customer.QueryBusiness(s))
			business.GET("/:id", api_customer.GetBuisnessDetailsById(s))

		}
	}
}
