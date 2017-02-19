package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"mailadm/models"
	"net/http"
	"strconv"

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
	log_prefix := "[MBOX/EDIT]"
	IdString := ctl.GetString(":id")
	Id, err := strconv.ParseUint(IdString, 10, 32)
	if err != nil {
		log.Printf("%s id is not int", log_prefix)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	mailbox, err := models.GetMailUserTupleById(Id)
	if err != nil {
		log.Printf("%s mailbox not found: %s", log_prefix, err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	updatedMailbox, err := ParseMailboxArgsAndUpdateTuple(ctl.Ctx.Input.RequestBody, mailbox)
	if err != nil {
		log.Printf("%s parse args: %s", log_prefix, err)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	if mailbox.Profile != updatedMailbox.Profile {
		_, count, err := models.GetMailProfileTupleById(updatedMailbox.Profile)
		if err != nil {
			log.Printf("%s : %s", log_prefix, err)
			http.Error(ctl.Ctx.ResponseWriter, "Some error", 500)
			return
		}
		if count == 0 {
			log.Printf("%s profile not found: %d", log_prefix, updatedMailbox.Profile)
			http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
			return
		}
	}
	err = models.UpdateMailboxTuple(updatedMailbox)
	if err != nil {
		log.Printf("%s update failed: %s", log_prefix, err)
		http.Error(ctl.Ctx.ResponseWriter, "Some DB error", 500)
		return
	}
	msg := fmt.Sprintf("Success for maildbox id %d (localpart: %s)", Id, updatedMailbox.Login)
	ctl.Data["json"] = msg
	ctl.ServeJSON()
}

func ParseMailboxArgsAndUpdateTuple(jsonString []byte, srcTuple models.MailUserTuple) (tuple models.MailUserTuple, err error) {
	var container interface{}
	err = json.Unmarshal(jsonString, &container)
	if err != nil {
		return tuple, fmt.Errorf("JSON unmarshal: %s", err)
	}
	tuple = srcTuple
	ArgsMap := container.(map[string]interface{})
	for arg_name, arg := range ArgsMap {
		switch arg := arg.(type) {
		case int, int32, int8, int16, uint, uint8, uint16, uint32, uint64, float32, float64:
			if arg_name == "Profile" {
				num := InterfaceToInt(arg)
				if num <= 0 || num > 10000 {
					return tuple, fmt.Errorf("invalid profile id: %d", arg)
				}
				tuple.Profile = num
			} else {
				return tuple, fmt.Errorf("Unknown integer parameter: %d", arg)
			}
		case string:
			if arg_name == "Password" {
				if len(arg) == 0 {
					return tuple, fmt.Errorf("password length zero")
				}
				if len(arg) > 128 {
					return tuple, fmt.Errorf("password length > 128")
				}
				tuple.Password = arg
			} else if arg_name == "Name" {
				if len(arg) > 128 {
					return tuple, fmt.Errorf("name length > 128")
				}
				tuple.Fullname = arg
			} else {
				return tuple, fmt.Errorf("Unknown string parameter: %q", arg_name)
			}
		case bool:
			if arg_name == "Disabled" {
				tuple.Bool_disabled = arg
			} else {
				return tuple, fmt.Errorf("Unknown bool parameter: %q", arg_name)
			}
		default:
			log.Printf("parameter %q of unsupported type", arg_name)
		}
	}
	return
}

func InterfaceToInt(arg interface{}) (num int) {
	switch arg := arg.(type) {
	case int:
		num = int(arg)
	case int8:
		num = int(arg)
	case int16:
		num = int(arg)
	case int32:
		num = int(arg)
	case int64:
		num = int(arg)
	case uint:
		num = int(arg)
	case uint8:
		num = int(arg)
	case uint16:
		num = int(arg)
	case uint32:
		num = int(arg)
	case uint64:
		num = int(arg)
	case float32:
		num = int(arg)
	case float64:
		num = int(arg)
	}
	return
}

// @Title GetAll
// @Description get all Mail Users
// @Success 200 {object} models.Mailbox
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

// @Title GetOne
// @Description Get mailbox data
// @Success 200 {object} models.maildb
// @router /:id [get]
func (ctl *MailboxController) GetOne() {
	log_prefix := "[MBOX/EDIT]"
	IdString := ctl.GetString(":id")
	Id, err := strconv.ParseUint(IdString, 10, 32)
	if err != nil {
		log.Printf("%s id is not int", log_prefix)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	mailbox, count, err := models.GetMailboxFullInfoById(Id)
	if err != nil {
		log.Printf("%s mailbox select: %s", log_prefix, err)
		http.Error(ctl.Ctx.ResponseWriter, "DB error", 500)
		return
	}
	if count == 0 {
		log.Printf("%s mailbox not found id=%d", log_prefix, Id)
		http.Error(ctl.Ctx.ResponseWriter, "Bad arguments", 400)
		return
	}
	ctl.Data["json"] = mailbox
	ctl.ServeJSON()
}
