package customers

import "github.com/gin-gonic/gin"

type ActorsController struct {
	actorsUseCase ActorsUseCase
}

func NewActorsController(actorsUseCase ActorsUseCase) *ActorsController {
	return &ActorsController{
		actorsUseCase: actorsUseCase,
	}
}

func (ctrl *ActorsController) CreateCustomer(c *gin.Context) {
	ctrl.actorsUseCase.CreateCustomer(c)
}
