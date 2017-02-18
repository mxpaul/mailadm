package controllers

import (
	//"encoding/json"
	"log"
	"mailadm/models"
	"net/http"

	"github.com/astaxie/beego"
)

// Operations about Users
type DomainController struct {
	beego.Controller
}

func (ctl *DomainController) GetAll() {
	domains, err := models.GetAllDomains()
	if err != nil {
		log.Printf("Domain Getall: maildb error: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "Error, come back later", 502)
		return
	}
	ctl.Data["json"] = domains
	ctl.ServeJSON()
}