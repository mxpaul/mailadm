package models

import (
	"fmt"
	"log"
	"time"

	"github.com/astaxie/beego"
	"gopkg.in/pg.v5"
)

type Mailbox struct {
	Id        int
	Localpart string
	Domain    string
	Email     string
	Name      string
	Disabled  string
	Profile   string
}

type Maildomain struct {
	Id             int
	Name           string
	Defaultprofile int
}

var (
	PgDb              *pg.DB
	mbox_list_query   string
	domain_list_query string
)

func init() {
	timeout, err := beego.AppConfig.Float("maildb::timeout")
	if err != nil {
		log.Fatalf("maildb::timeout in config: %v", err)
	}
	PgDb = pg.Connect(&pg.Options{
		Addr:     beego.AppConfig.String("maildb::addr"),
		User:     beego.AppConfig.String("maildb::user"),
		Password: beego.AppConfig.String("maildb::password"),
		Database: beego.AppConfig.String("maildb::database"),
	}).WithTimeout(time.Second * time.Duration(timeout))

	mbox_list_query = `
SELECT 
	u.id as Id,
	u.login as Localpart,
	d.name as Domain,
	u.login || '@' || d.name as Email,
	u.fullname as Name,
	u.bool_disabled as Disabled,
	p.name as profile
 FROM 
	t_user u,
	t_domain d,
	t_profile p
 WHERE 
	u.domain=d.id 
	and u.profile=p.id
 ORDER BY Id
`
	domain_list_query = `SELECT
		d.id as Id,
		d.name as Name,
		d.default_profile as DefaultProfile 
	FROM t_domain d
	ORDER BY id`
}

func GetAllMailboxes() (MailboxList []Mailbox, err error) {
	_, err = PgDb.Query(&MailboxList, mbox_list_query)

	if err != nil {
		return nil, fmt.Errorf("select mailboxes: %s", err)
	}

	return
}

func GetAllDomains() (DomainList []Maildomain, err error) {
	_, err = PgDb.Query(&DomainList, domain_list_query)

	if err != nil {
		return nil, fmt.Errorf("select domains: %s", err)
	}

	return
}
