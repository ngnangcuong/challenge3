package post

import (
	"github.com/gin-gonic/gin"

	"flag"
	"fmt"
	"strconv"
	"challenge3/models"
	"challenge3/database"
	repo "challenge3/repository"
)

func GetListPost(c *gin.Context) {
	if isLogin := c.MustGet("isLogin").(bool); !isLogin {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
		return
	}

	flag.Parse()
	_ = flag.Arg(0)
	connection := database.GetDatabase()

	p := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if page <= 0 {
		c.JSON(200, gin.H{
			"message": "Does not exsit page ",
		})
		return
	}

	var postList []models.Post
	var offset = (page - 1) *10
	
	connection.Limit(10).Offset(offset).Find(&postList)
	
	c.JSON(200, postList)
	
}

func CreatePost(c *gin.Context) {
	connection := database.GetDatabase()
	postRepo := repo.NewPostRepo(connection)

	if isLogin := c.MustGet("isLogin").(bool); !isLogin {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
		return
	}

	userID := c.MustGet("userID")
	email := c.MustGet("email")

	content := c.PostForm("content")
	var post = models.Post{
		UserID: uint(userID.(float64)),
		Email: email.(string),
		Content: content,
	}

	postRepo.Create(post)
	c.JSON(200, gin.H{
		"message": "Create post successfully",
	})
}

func UpdatePost(c *gin.Context) {
	connection := database.GetDatabase()
	postRepo := repo.NewPostRepo(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
	}

	role := c.MustGet("role").(string)
	postID := c.Param("postID")
	userID := c.MustGet("userID")
	postCheck, _ := postRepo.Find(postID)

	if postCheck.Content == "" {
		c.JSON(200, gin.H{
			"message": "Does not exist post",
		})
		return
	}

	if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
		c.JSON(401, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	content := c.PostForm("content")
	postRepo.Update(postID, content)

	c.JSON(200, gin.H{
		"message": "Edit post successfully",
	})

}

func DeletePost(c *gin.Context) {
	connection := database.GetDatabase()
	postRepo := repo.NewPostRepo(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
	}
	
	role := c.MustGet("role").(string)
	postID := c.Param("postID")
	userID := c.MustGet("userID")
	postCheck, _ := postRepo.Find(postID)

	if postCheck.Content == ""  {
		c.JSON(200, gin.H{
			"message": "Does not exist post",
		})
		return
	}

	if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
		c.JSON(401, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	postRepo.Delete(postID)
	c.JSON(200, gin.H{
		"message": "Delete post successfully",
	})
}