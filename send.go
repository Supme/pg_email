package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Supme/emailwriters"
	"github.com/microo8/plgo"
	"log"
	"mime"
	"net/smtp"
	"time"
)

//Add to postgresql.conf parameters:
//email.serverhost = 'smtp.server.tld'
//email.serverport = '465'
//email.fromname = 'Postgres email sender'
//email.fromemail = 'postgres@email.tld'
//email.username = 'smtp_user'
//email.password ='smtp_password'

const (
	boundaryAlternative      = "===============ALTERNATIVE=="
	boundaryAlternativeBegin = "--" + boundaryAlternative + "\r\n"
	boundaryAlternativeEnd   = "--" + boundaryAlternative + "--\r\n"
)

//SendEmail sending email
//function arguments (to_email, to_name, subject, text_html, text_plain) return error or blank string then ok
func SendEmail(toEmail, toName, subject, textHTML, textPlain string) string {
	//ntcLog := plgo.NewNoticeLogger("", log.Ltime|log.Lshortfile)
	errLog := plgo.NewErrorLogger("", log.Ltime|log.Lshortfile)
	defer func() {
		if r := recover(); r != nil {
			errLog.Println("recover after panic: ", r)
		}
	}()

	conf, err := getConfigValues("email.serverhost", "email.serverport", "email.fromemail", "email.fromname", "email.username", "email.password")
	if err != nil {
		errLog.Printf("Get config %s", err)
	}

	auth := smtp.PlainAuth("", conf["email.username"], conf["email.password"], conf["email.serverhost"])

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         conf["email.serverhost"],
	}

	conn, err := tls.Dial("tcp", conf["email.serverhost"]+":"+conf["email.serverport"], tlsconfig)
	if err != nil {
		return err.Error()
	}

	client, err := smtp.NewClient(conn, conf["email.serverhost"])
	if err != nil {
		return err.Error()
	}

	if err = client.Auth(auth); err != nil {
		return err.Error()
	}

	if err = client.Mail(conf["email.fromemail"]); err != nil {
		return err.Error()
	}

	if err = client.Rcpt(toEmail); err != nil {
		return err.Error()
	}

	w, err := client.Data()
	if err != nil {
		return err.Error()
	}

	_, err = w.Write([]byte("MIME-Version: 1.0\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("From: " + mime.BEncoding.Encode("utf-8", conf["email.fromname"]) + " <" + conf["email.fromemail"] + ">" + "\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("To: " + mime.BEncoding.Encode("utf-8", toName) + " <" + toEmail + ">" + "\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Content-Type: multipart/alternative;\r\n\tboundary=\"" + boundaryAlternative + "\"\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Subject: " + mime.BEncoding.Encode("utf-8", subject) + "\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	// Plain text
	_, err = w.Write([]byte(boundaryAlternativeBegin))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("MIME-Version: 1.0\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Content-Type: text/plain;\r\n\tcharset=\"utf-8\"\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	err = EmailWriters.Base64Writer(w, textPlain)
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	// HTML
	_, err = w.Write([]byte(boundaryAlternativeBegin))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("MIME-Version: 1.0\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Content-Type: text/html;\r\n\tcharset=\"utf-8\"\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	err = EmailWriters.Base64Writer(w, textHTML)
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}

	_, err = w.Write([]byte(boundaryAlternativeEnd))
	if err != nil {
		return err.Error()
	}
	_, err = w.Write([]byte("\r\n"))
	if err != nil {
		return err.Error()
	}

	err = w.Close()
	if err != nil {
		return err.Error()
	}
	err = client.Quit()
	if err != nil {
		return err.Error()
	}

	return ""
}

func getConfigValues(keys ...string) (map[string]string, error) {
	db, err := plgo.Open() //open the connection to DB
	if err != nil {
		return nil, fmt.Errorf("cannot open DB: %s", err)
	}
	defer db.Close()

	val := map[string]string{}
	for i := range keys {
		stmt, err := db.Prepare(fmt.Sprintf("SHOW %s", keys[i]), nil)
		if err != nil {
			return nil, fmt.Errorf("prepare query: %s", err)
		}
		row, err := stmt.QueryRow()
		if err != nil {
			return nil, fmt.Errorf("query row: %s", err)
		}
		var v string
		err = row.Scan(&v)
		if err != nil {
			return nil, fmt.Errorf("row scan: %s", err)
		}
		val[keys[i]] = v
	}
	return val, nil
}
