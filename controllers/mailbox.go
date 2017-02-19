package controllers

import (
	"encoding/json"
	"fmt"
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
func (ctl *MailboxController) Create() {
	var mailbox models.MailboxAdd
	err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &mailbox)
	if err != nil {
		log.Printf("arg error: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	err = mailbox.Validate()
	if err != nil {
		log.Printf("MailboxAdd: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	domain, err := models.GetDomainById(mailbox.Domain)
	if err != nil {
		log.Printf("domain %d load: %s", mailbox.Domain, err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	// For now we just ignore profilesent in request and use domain default profile
	//if mailbox.Profile == 0 {
	mailbox.Profile = domain.Defaultprofile
	//}
	existingBoxId, err := models.MailboxIdIfExists(mailbox)
	if err != nil {
		log.Printf("exist check: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "DB Error", 500)
		return
	}
	if existingBoxId > 0 {
		log.Printf("Mailbox %s@%s exists", mailbox.Localpart, domain.Name)
		http.Error(ctl.Ctx.ResponseWriter, "Mailbox exists", 500)
		return
	}
	err = models.CreateMailbox(mailbox)
	if err != nil {
		log.Printf("MailboxAdd: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "DB Error", 500)
		return
	}
	msg := fmt.Sprintf("Mailbox %s@%s created", mailbox.Localpart, domain.Name)
	log.Print(msg)
	ctl.Data["json"] = msg
	ctl.ServeJSON()
}

// @Title Edit
// @Description Put mailbox edited data
// @Success 200 {object} models.maildb
// @router /:id [put]
func (ctl *MailboxController) Edit() {
	Id := ctl.GetString(":id")
	log.Printf("PUT REQUEST FOR ID %d", Id)
	var mailbox models.MailboxEdit
	err := json.Unmarshal(ctl.Ctx.Input.RequestBody, &mailbox)
	if err != nil {
		log.Printf("Mailbox EDIT arg error: %s", err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	ctl.Data["json"] = mailbox
	ctl.ServeJSON()
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
