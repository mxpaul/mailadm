package controllers

import (
	"encoding/json"
	"log"
	"mailadm/models"
	"net/http"

	"github.com/astaxie/beego"
)

// Operations about Users
type MailboxController struct {
	beego.Controller
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (u *MailboxController) Post() {
	var mb models.Mailbox
	json.Unmarshal(u.Ctx.Input.RequestBody, &mb)
	//uid := models.AddUser(user)
	//u.Data["json"] = map[string]string{"uid": uid}
	//u.ServeJSON()
}

// @Title GetAll
// @Description get all Users
// @Success 200 {object} models.User
// @router / [get]
func (mb *MailboxController) GetAll() {
	mbxs, err := models.GetAllMailboxes()
	if err != nil {
		log.Printf("Getall: maildb error: %s", err)
		http.Error(mb.Ctx.ResponseWriter, "Error, come back later", 502)
		return
	}
	mb.Data["json"] = mbxs
	mb.ServeJSON()
}
