package customers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ActorsUseCase interface {
	CreateCustomer(c *gin.Context)
}

type actorsUseCase struct {
	actorsRepo ActorsRepository
}

func NewActorsUseCase(actorsRepo ActorsRepository) ActorsUseCase {
	return &actorsUseCase{
		actorsRepo: actorsRepo,
	}
}

func (uc *actorsUseCase) CreateCustomer(c *gin.Context) {
	var actors Actor

	token_key := c.GetHeader("token_key")
	var flagVerified int
	checker, err := uc.actorsRepo.FindByUsername(token_key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if checker.Username == "superadmin" {
		flagVerified = 1
	} else {
		if checker != nil && checker.RoleID == "2" {
			flagVerified = 2
		}
	}

	if flagVerified == 1 || flagVerified == 2 {
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

		if actors.RoleID == "3" {
			if err := uc.actorsRepo.Create(&actors); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Can't create except Customer Actor"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": actors})
}
