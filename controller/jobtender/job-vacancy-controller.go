package jobtender

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yusufwira/lern-golang-gin/entity/job_tender"
	"gorm.io/gorm"
)

type JobVacancyController struct {
	JobVacancyRepo *job_tender.JobVacancyRepo
}

func GetJobVacancyController(Db *gorm.DB) *JobVacancyController {
	return &JobVacancyController{
		JobVacancyRepo: job_tender.GetJobVacancyRepo(Db)}
}

func (c *JobVacancyController) GetDetailJob(ctx *gin.Context) {
	id := ctx.Param("id")
	idJob, _ := strconv.Atoi(id)
	model, err := c.JobVacancyRepo.FindJobByID(idJob)
	if err == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusOK,
			"info":   "Success",
			"Data":   model})
		return
	} else {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"info":   "Failed",
			"Data":   err})
		return
	}

}
