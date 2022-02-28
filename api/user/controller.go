package user

import (
	"github.com/gin-gonic/gin"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	
	"time"
	"strings"
	"challenge3/database"
	"challenge3/models"
)

var mySigningKey = "pa$$w0rd"

func Response(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"message": message,
	})
}

func CheckPasswordHash(password string, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}

func GenerateJWT(userAuth *models.User) (string, error) {
	var secretkey = []byte(mySigningKey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = userAuth.Email
	claims["userID"] = userAuth.ID
	claims["role"] = userAuth.Role
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(secretkey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GetPermisstion(role string) string {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	var roleCheck models.Role
	connection.Where("name = ?", role).First(&roleCheck)

	return roleCheck.Permission
}

func GetListUser(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not Log in yet",
		})
		return
	}

	var userList []models.User

	connection.Find(&userList)

	c.HTML(200, "listUser.tmpl", gin.H{
		"userList": userList,
	})
}

func LogIn(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	email := c.PostForm("email")
	password := c.PostForm("password")

	var userAuth models.User

	connection.Where("email = ?", email).First(&userAuth)

	if userAuth.Email == "" {
		Response(c, 200, "Not User")
		return
	}

	if check := CheckPasswordHash(password, userAuth.Password); !check {
		Response(c, 200, "Password is not correct")
		return
	} 
	
	tokenString, err := GenerateJWT(&userAuth)
	if err != nil {
		Response(c, 200, "Cannot generate jwt token")
		return
	}

	c.SetCookie("token", tokenString, 150, "/", "localhost", false, true)
	Response(c, 200, tokenString)
}

func LogOut(c *gin.Context) {
	c.SetCookie("token", "", 150, "/", "localhost", false, true)
	Response(c, 200, "Successful log out")
}

func Register(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	email := c.PostForm("email")
	name := c.PostForm("name")
	password := c.PostForm("password")

	var userCheck models.User
	
	connection.Where("email = ?", email).First(&userCheck)

	if userCheck.Email != "" {
		Response(c, 200, "Email is already existed")
		return
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		Response(c, 200, "Cannot generate hash password")
		return
	}

	var user = models.User{
		Name: name,
		Email: email,
		Password: string(hashPassword),
		Role: "user",
	}

	connection.Create(&user)
	Response(c, 200, "Create user successfully")
}

func CreateUser(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not Log in yet",
		})
		return
	}

	role := c.MustGet("role")
	permission := GetPermisstion(role.(string))

	if canCreate := strings.Contains(permission, "c"); !canCreate && role.(string) != "admin" {
		c.JSON(200, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	password := c.PostForm("password")
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		Response(c, 200, "Cannot generate hash password")
		return
	}

	var newUser = models.User{
		Email: c.PostForm("email"),
		Name: c.PostForm("name"),
		Password: string(hashPassword),
		Role: "user",
	}
	connection.Create(&newUser)

	Response(c, 200, "Create user successfully")
}

func DeleteUser(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not Log in yet",
		})
		return
	}

	role := c.MustGet("role")
	permission := GetPermisstion(role.(string))

	if canDelete := strings.Contains(permission, "d"); !canDelete && role != "admin"{
		c.JSON(200, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	email := c.Param("userEmail")
	var userCheck models.User

	connection.Where("email = ?", email).First(&userCheck)

	if userCheck.Email == "" {
		Response(c, 200, "Do not exist user")
		return
	}

	connection.Delete(&userCheck)
	Response(c, 200, "Delete user successfully")
}

func UpdateUser(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		c.JSON(200, gin.H{
			"message": "Not Log in yet",
		})
		return
	}

	role := c.MustGet("role")
	permission := GetPermisstion(role.(string))

	if canUpdate := strings.Contains(permission, "u"); !canUpdate && role != "admin"{
		c.JSON(200, gin.H{
			"message": "Not Authorized",
		})
		return
	}

	email := c.Param("userEmail")
	var userCheck models.User

	connection.Where("email = ?", email).First(&userCheck)

	if userCheck.Email == "" {
		Response(c, 200, "Do not exist user")
		return
	}

	userCheck.Name = c.PostForm("name")
	userCheck.Password = c.PostForm("password")
	
	connection.Save(&userCheck)
	Response(c, 200, "Update user successfully")

}

func NewRole(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		Response(c, 200, "Not Log in yet")
		return
	}

	role := c.MustGet("role")
	if role != "admin" {
		Response(c, 200, "Not Authorized")
		return
	}

	name := c.PostForm("name")
	permission := c.PostForm("permission")

	var roleCheck models.Role
	connection.Where("name = ?", name).First(&roleCheck)

	if roleCheck.Name != "" {
		Response(c, 200, "This role is available")
		return
	}

	roleCheck.Name = name
	roleCheck.Permission = permission

	connection.Create(&roleCheck)
	Response(c, 200, "Create role successfully")
}

func ChangeRole(c *gin.Context) {
	connection := database.GetDatabase()
	defer database.CloseDatabase(connection)

	if check := c.MustGet("isLogin").(bool); !check {
		Response(c, 200, "Not Log in yet")
		return
	}

	role := c.MustGet("role")
	if role != "admin" {
		Response(c, 200, "Not Authorized")
		return
	}

	email := c.PostForm("email")
	role = c.PostForm("role")

	var (
		userCheck models.User
		roleCheck models.Role
	)

	connection.Where("email = ?", email).First(&userCheck)
	if userCheck.Email == "" {
		Response(c, 200, "Does not exist user")
		return
	}

	connection.Where("name = ?", role).First(&roleCheck)
	if roleCheck.Name == "" {
		Response(c, 200, "Does not exist role")
		return
	}

	userCheck.Role = role.(string)
	connection.Save(&userCheck)
	Response(c, 200, "Change role successfully")

}