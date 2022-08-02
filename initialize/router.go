package initialize

import (
	"net/http"

	"github.com/akazwz/gin-api/api/auth"
	"github.com/akazwz/gin-api/api/file"
	"github.com/akazwz/gin-api/api/posts"
	"github.com/akazwz/gin-api/api/s3/r2"
	"github.com/akazwz/gin-api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	//  cors 跨域
	r.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	// 404 not found
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Not Found",
		})
	})

	//Teapot  418
	r.GET("teapot", func(c *gin.Context) {
		c.JSON(http.StatusTeapot, gin.H{
			"message": "I'm a teapot",
			"story": "This code was defined in 1998 " +
				"as one of the traditional IETF April Fools' jokes," +
				" in RFC 2324, Hyper Text Coffee Pot Control Protocol," +
				" and is not expected to be implemented by actual HTTP servers." +
				" However, known implementations do exist.",
		})
	})

	// 文件路由组
	fileGroup := r.Group("/file").Use(middleware.LimitByRequest(3))
	{
		// 简单上传
		fileGroup.POST("", file.UploadFile)
		// 分块上传
		fileGroup.POST("/chunk", file.UploadChunk)
		fileGroup.POST("/chunk/merge", file.MergeChunk)
		fileGroup.GET("/chunk/state", file.ChunkState)
	}
	// auth 路由组
	authGroup := r.Group("/auth").Use(middleware.LimitByRequest(3))
	{
		authGroup.POST("/signup", auth.SignupByUsernamePwd)
		authGroup.POST("/login", auth.LoginByUsernamePwd)
		/* me jwt auth  */
		authGroup.GET("/me", middleware.JWTAuth(), auth.Me)
	}

	// s3
	s3Group := r.Group("/s3").Use(middleware.LimitByRequest(3))
	{
		// 直传
		s3Group.POST("/r2/upload", r2.Upload)
		// https://docs.aws.amazon.com/amazonglacier/latest/dev/uploading-an-archive-mpu-using-rest.html
		s3Group.POST("/r2/upload/:key", r2.CreateMultipartUpload)
		s3Group.PUT("/r2/upload/:key", r2.UploadPart)
		s3Group.POST("/r2/upload/complete", r2.CompleteMultipartUpload)
		s3Group.DELETE("/r2/upload/:key", r2.AbortMultipartUpload)
		s3Group.GET("/r2/upload/:key", r2.ListParts)
		s3Group.GET("/r2/upload", r2.ListMultipartUploads)
	}

	postsGroup := r.Group("/posts").Use(middleware.LimitByRequest(3))
	{
		postsGroup.GET("/:id", posts.GetPostById)
		postsGroup.POST("", middleware.JWTAuth(), posts.CreatePost)
		postsGroup.GET("", posts.FindPosts)
	}

	return r
}
