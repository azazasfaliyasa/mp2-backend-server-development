package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type Actors struct {
	gorm.Model
	Username  string `json:"username"`
	Password  string `json:"password"`
	Tokey_key string `json:"token_key"`
	Role_id   string `json:"role_id"`
	Flag_act  string `json:"flag_act"`
	Flag_ver  string `json:"flag_ver"`
}

var db *gorm.DB
var err error

func initDB() {
	dsn := "root:mysql@tcp(localhost:3306)/milestone1?parseTime=true"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
}

func createCustomer(c *gin.Context) {
	var actors Actors
	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var flag_verified int
	if checker.Username == "superadmin" {
		flag_verified = 1
	} else if checker.Role_id == "2" {
		flag_verified = 2
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Baca data JSON dari body permintaan
		if err := c.BindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Membuat objek hash dari algoritma SHA-256
		hash := sha256.New()
		// Mengupdate hash dengan data yang ingin di-hash
		hash.Write([]byte(actors.Password))
		// Mengambil nilai hash sebagai array byte
		hashBytes := hash.Sum(nil)
		// Mengubah array byte menjadi representasi heksadesimal
		hashString := hex.EncodeToString(hashBytes)

		actors.Password = hashString

		if actors.Role_id == "3" {
			err := db.Select("username", "password", "role_id", "flag_act").Create(&actors).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create except Customer Actor"})
			return
		}

		// Tampilkan respons berhasil
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

}

func createAdmin(c *gin.Context) {
	var actors Actors

	var checker Actors

	token_key := c.GetHeader("token_key")
	var flag_verified int
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if checker.Username == "superadmin" {
		flag_verified = 1
	} else if checker.Role_id == "2" {
		flag_verified = 2
	}

	if flag_verified == 2 {
		// Baca data JSON dari body permintaan
		if err := c.BindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Membuat objek hash dari algoritma SHA-256
		hash := sha256.New()
		// Mengupdate hash dengan data yang ingin di-hash
		hash.Write([]byte(actors.Password))
		// Mengambil nilai hash sebagai array byte
		hashBytes := hash.Sum(nil)
		// Mengubah array byte menjadi representasi heksadesimal
		hashString := hex.EncodeToString(hashBytes)

		actors.Password = hashString

		if actors.Role_id == "2" {
			err := db.Select("username", "password", "role_id", "flag_act").Create(&actors).Error
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create except admin"})
			return
		}

		// Tampilkan respons berhasil
		c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

}

func getWaitingApproved(c *gin.Context) {
	var actors []Actors
	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if checker.Username == "superadmin" {
		// Dapatkan semua data user dari database dengan kondisi WHERE flag_ver = nil
		if err := db.Where("flag_ver", nil).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data user
		c.JSON(http.StatusOK, gin.H{"users": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": "Permission Denied"})
	}
}

func getCustomer(c *gin.Context) {
	var actors []Actors
	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var flag_verified int
	if checker.Username == "superadmin" {
		flag_verified = 1
	} else if checker.Role_id == "2" {
		flag_verified = 2
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan nilai halaman dan ukuran halaman dari query string untuk pagination
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "10")

		// Konversi nilai halaman dan ukuran halaman ke tipe data yang sesuai
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = 10
		}

		// Dapatkan semua data admin dari database dengan kondisi WHERE role_id = 2
		var totalRecords int64
		db.Model(&Actors{}).Where("role_id = ?", "2").Count(&totalRecords)

		// Hitung offset berdasarkan halaman dan ukuran halaman
		offset := (page - 1) * size
		if offset < 0 {
			offset = 0
		}

		// Query data admin dengan pagination
		if err := db.Where("role_id = ?", "3").Offset(offset).Limit(size).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data admin dan informasi pagination
		c.JSON(http.StatusOK, gin.H{
			"users": actors,
			"page":  page,
			"size":  size,
			"total": totalRecords,
			"pages": int(math.Ceil(float64(totalRecords) / float64(size))),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

}

func getActorsById(c *gin.Context) {
	var actors Actors
	userID := c.Param("id")

	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var flag_verified int
	if checker.Username == "superadmin" {
		flag_verified = 1
	} else if checker.Role_id == "2" {
		flag_verified = 2
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan data user dari database berdasarkan ID
		if err := db.First(&actors, userID).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data user
		c.JSON(http.StatusOK, gin.H{"user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
	}

}

func getAdmin(c *gin.Context) {
	var actors []Actors
	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var flag_verified int
	if checker.Username == "superadmin" {
		flag_verified = 1
	} else if checker.Role_id == "2" {
		flag_verified = 2
	}

	if flag_verified == 1 || flag_verified == 2 {
		// Dapatkan nilai halaman dan ukuran halaman dari query string untuk pagination
		pageStr := c.DefaultQuery("page", "1")
		sizeStr := c.DefaultQuery("size", "10")

		// Konversi nilai halaman dan ukuran halaman ke tipe data yang sesuai
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			page = 1
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil || size < 1 {
			size = 10
		}

		// Dapatkan semua data admin dari database dengan kondisi WHERE role_id = 2
		var totalRecords int64
		db.Model(&Actors{}).Where("role_id = ?", "2").Count(&totalRecords)

		// Hitung offset berdasarkan halaman dan ukuran halaman
		offset := (page - 1) * size
		if offset < 0 {
			offset = 0
		}

		// Query data admin dengan pagination
		if err := db.Where("role_id = ?", "2").Offset(offset).Limit(size).Find(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan data admin dan informasi pagination
		c.JSON(http.StatusOK, gin.H{
			"users": actors,
			"page":  page,
			"size":  size,
			"total": totalRecords,
			"pages": int(math.Ceil(float64(totalRecords) / float64(size))),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}
}

func updateAdmin(c *gin.Context) {
	var actors Actors
	userID := c.Param("id")

	var checker Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if checker.Username == "superadmin" {
		// Dapatkan data user dari database berdasarkan ID
		if err := db.First(&actors, userID).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Baca data JSON dari body permintaan
		if err := c.ShouldBindJSON(&actors); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Simpan perubahan ke database
		if err := db.Select("username", "password", "flag_act", "flag_ver").Save(&actors).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// Tampilkan respons berhasil
		c.JSON(http.StatusOK, gin.H{"message": "User updated successfully", "user": actors})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
	}
}

func deleteCustomer(c *gin.Context) {
	var actor Actors
	var deleter Actors

	token_key := c.GetHeader("token_key")
	if err := db.Where("token_key", token_key).First(&deleter).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var flag_verified int
	if deleter.Username == "superadmin" {
		flag_verified = 1
	} else if deleter.Role_id == "2" {
		flag_verified = 2
	}

	actorID := c.Param("id")
	// Dapatkan data user dari database berdasarkan ID
	if err := db.First(&actor, actorID).Error; err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Hapus data user dari database
	if flag_verified == 1 {
		if err := db.Delete(&actor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if flag_verified == 2 && actor.Role_id == "3" {
		if err := db.Delete(&actor).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Permission Denied"})
		return
	}

	// Tampilkan respons berhasil
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func loginAuth(c *gin.Context) {
	var checker Actors

	username, password, status := c.Request.BasicAuth()

	// Membuat objek hash dari algoritma SHA-256
	hash := sha256.New()
	// Mengupdate hash dengan data yang ingin di-hash
	hash.Write([]byte(password))
	// Mengambil nilai hash sebagai array byte
	hashBytes := hash.Sum(nil)
	// Mengubah array byte menjadi representasi heksadesimal
	hashString := hex.EncodeToString(hashBytes)

	password = hashString

	if !status {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := db.Where("username", username).First(&checker).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": username})
		return
	}
	if checker.Password == password {
		// Inisialisasi klaim-klaim yang ingin Anda sertakan dalam token
		claims := jwt.MapClaims{
			"sub":  checker.ID,
			"name": checker.Username,
			"iat":  time.Now().Unix(),
			"exp":  time.Now().Add(time.Hour * 1).Unix(),
		}
		// Tandatangani token dengan secret key
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte("secret-key"))
		if err != nil {
			// Errors handler
		}

		if err := db.Model(&checker).Update("token_key", signedToken).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login has been successful", "token_key": signedToken})
		return
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wrong Password"})
		return
	}
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	//actorsRepo := customers.NewActorsRepository(db)
	//actorsUseCase := customers.NewActorsUseCase(actorsRepo)
	//actorsController := customers.NewActorsController(actorsUseCase)

	//r.POST("/customers", actorsController.CreateCustomer)

	r.POST("/customers", createCustomer)
	r.POST("/admin", createAdmin)
	r.GET("/approved", getWaitingApproved)
	r.GET("/customers", getCustomer)
	r.GET("/customers/:id", getActorsById)
	r.GET("/admin", getAdmin)
	r.GET("/admin:id", getActorsById)
	r.PUT("/admin/:id", updateAdmin)
	r.DELETE("/customers/:id", deleteCustomer)
	r.POST("/login", loginAuth)

	return r
}

func main() {
	initDB()
	r := setupRouter()

	// Port server
	err := r.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
