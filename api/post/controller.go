package post

import (
	"github.com/gin-gonic/gin"

	"flag"
	"fmt"
	"strconv"
	"challenge3/models"
	"challenge3/database"
)

func GetListPost(c *gin.Context) {
	flag.Parse()
	_ = flag.Arg(0)
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

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

	connection.Find(&postList)
	if len(postList) > (page-1)*10 {
		if len(postList) > page * 10 {
			c.HTML(200, "listPost.tmpl", gin.H{
				"listPost": postList[(page - 1)*10 +1: page*10],
			})
		} else {
			c.HTML(200, "listPost.tmpl", gin.H{
				"listPost": postList[(page - 1)*10 +1:],
			})
		}
		return
	}
	
	c.JSON(200, gin.H{
		"message": "We have less page than you are expected",
	})
}

func CreatePost(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

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

	connection.Create(&post)
	c.JSON(200, gin.H{
		"message": "Create post successfully",
	})
}

func UpdatePost(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
	}

	role := c.MustGet("role").(string)
	postID := c.Param("postID")
	userID := c.MustGet("userID")
	var postCheck models.Post

	connection.First(&postCheck, postID)
	if postCheck.Content == "" {
		c.JSON(200, gin.H{
			"message": "Does not exist post",
		})
		return
	}

	if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
		c.JSON(200, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	content := c.PostForm("content")
	postCheck.Content = content
	connection.Save(&postCheck)

	c.JSON(200, gin.H{
		"message": "Edit post successfully",
	})

}

func DeletePost(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not log in yet",
		})
	}
	
	role := c.MustGet("role").(string)
	postID := c.Param("postID")
	userID := c.MustGet("userID")
	var postCheck models.Post

	connection.First(&postCheck, postID)
	if postCheck.Content == ""  {
		c.JSON(200, gin.H{
			"message": "Does not exist post",
		})
		return
	}

	if postCheck.UserID != uint(userID.(float64)) && role != "admin" {
		c.JSON(200, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	connection.Delete(&postCheck)
	c.JSON(200, gin.H{
		"message": "Delete post successfully",
	})
}