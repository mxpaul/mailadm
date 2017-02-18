package models

import (
	"fmt"
	"log"
	"regexp"
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

type MailboxAdd struct {
	Localpart string
	Domain    int
	Name      string
	Password  string
	Profile   int `json:",omitempty"`
}

var (
	PgDb              *pg.DB
	mbox_list_query   string
	domain_list_query string
	RegexpLocalpart   *regexp.Regexp
)

func (box *MailboxAdd) Validate() (err error) {
	maxlen := 128
	if box == nil {
		return fmt.Errorf("nil.Validate() is pointless")
	}
	if len(box.Localpart) == 0 {
		return fmt.Errorf("localpart length zero")
	}
	if len(box.Localpart) > maxlen {
		return fmt.Errorf("localpart length %d", len(box.Localpart))
	}
	if len(box.Name) > maxlen {
		return fmt.Errorf("name length %d", len(box.Name))
	}
	if len(box.Password) == 0 {
		return fmt.Errorf("password length zero")
	}
	if len(box.Password) > maxlen {
		return fmt.Errorf("password length %d", len(box.Password))
	}
	if !RegexpLocalpart.MatchString(box.Localpart) {
		return fmt.Errorf("localpart no match regexp %q", RegexpLocalpart)
	}
	if box.Domain <= 0 {
		return fmt.Errorf("domain ID required")
	}
	return
}

func init() {
	timeout, err := beego.AppConfig.Float("maildb::timeout")

	if err != nil {
		log.Fatalf("maildb::timeout in config: %v", err)
	}

	RegexpLocalpart, err = regexp.Compile("^[A-Za-z0-9]+[-_.A-Za-z0-9]*$")
	if err != nil {
		log.Fatalf("Compile localpart regexp: %s", err)
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
	ORDER BY id
`
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

func GetDomainById(id int) (dom Maildomain, err error) {
	domain_select_by_id := `SELECT d.id as Id, d.name as Name, d.default_profile as DefaultProfile 
		FROM t_domain d WHERE id=%d
	`
	res, err := PgDb.Query(&dom, fmt.Sprintf(domain_select_by_id, id))

	if err != nil {
		return dom, fmt.Errorf("maildb: %s", err)
	}
	if res.RowsReturned() == 0 {
		return dom, fmt.Errorf("GetDomainById: not found")
	}

	return
}

func MailboxIdIfExists(box MailboxAdd) (int, error) {
	query := `SELECT u.id as Id FROM t_user u WHERE login=$1::text AND domain=$2::int `

	stmnt, err := PgDb.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("prepare: %s", err)
	}

	var otherbox Mailbox
	res, err := stmnt.Query(&otherbox, box.Localpart, box.Domain)

	if err != nil {
		return 0, fmt.Errorf("maildb: %s", err)
	}
	if res.RowsReturned() > 0 {
		return otherbox.Id, nil
	}

	return 0, nil
}

func CreateMailbox(box MailboxAdd) (err error) {
	query := `INSERT INTO t_user (login,domain,password,profile,fullname) 
		VALUES($1::text, $2::int, $3::text, $4::int, $5::text)
	`
	stmnt, err := PgDb.Prepare(query)
	if err != nil {
		return fmt.Errorf("maildb insert: %s", err)
	}

	_, err = stmnt.Exec(box.Localpart, box.Domain, box.Password, box.Profile, box.Name)

	if err != nil {
		return fmt.Errorf("maildb: %s %v", err, stmnt)
	}

	return
}
