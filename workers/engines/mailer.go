package engines

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"net/smtp"
	"os"
	"path/filepath"
	"sync"

	"github.com/anti-lgbt/medusa/config"
	"github.com/anti-lgbt/medusa/types"
	"gopkg.in/yaml.v2"
)

type Email struct {
	FromAddress string
	FromName    string
	ToAddress   string
	Subject     string
	Reader      io.Reader
}

type MailerEngineWorker struct {
	SMTP_HOST     string
	SMTP_PORT     string
	SMTP_USER     string
	SMTP_PASSWORD string
	SENDER_NAME   string
	SENDER_EMAIL  string
	Config        *types.MailerConfig

	mailerMutex sync.Mutex
}

func NewMailerEngineWorker() *MailerEngineWorker {
	filename, _ := filepath.Abs("config/mailer.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config *types.MailerConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return &MailerEngineWorker{
		SMTP_HOST:     os.Getenv("SMTP_HOST"),
		SMTP_PORT:     os.Getenv("SMTP_PORT"),
		SMTP_USER:     os.Getenv("SMTP_USER"),
		SMTP_PASSWORD: os.Getenv("SMTP_PASSWORD"),
		SENDER_NAME:   os.Getenv("SENDER_NAME"),
		SENDER_EMAIL:  os.Getenv("SENDER_EMAIL"),
		Config:        config,
	}
}

func (w *MailerEngineWorker) Process(payload []byte) error {
	w.mailerMutex.Lock()
	defer w.mailerMutex.Unlock()

	var params *types.EngineMailerPayload
	err := json.Unmarshal(payload, &params)
	if err != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	var event *types.MailerConfigEvent
	for _, e := range w.Config.Events {
		if params.Key == e.Key {
			event = &e
		}
	}

	var templates map[string]types.MailerConfigEventTemplates

	if event == nil {
		config.Logger.Errorf("Error: event not found in config")
		return nil
	}

	templates = event.Templates

	var temp *types.MailerConfigEventTemplates
	if event.Templates != nil {
		for language, t := range templates {
			if params.Language == language {
				temp = &t
			}
		}
	} else {
		config.Logger.Errorf("Error: template not found in config")
		return nil
	}

	if event == nil {
		config.Logger.Errorf("Error: template not found in config")
		return nil
	}

	template_tpl, err := template.ParseFiles("config/mailer/" + temp.TemplatePath)
	if err != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	content_buf := bytes.Buffer{}
	if template_tpl.Execute(&content_buf, params.Record) != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	layout_tpl, err := template.ParseFiles("config/mailer/layout.tpl")
	if err != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	email := Email{
		FromAddress: w.SENDER_EMAIL,
		FromName:    w.SENDER_NAME,
		ToAddress:   params.To,
		Subject:     temp.Subject,
		Reader:      bytes.NewReader(content_buf.Bytes()),
	}

	layout_buf := bytes.Buffer{}
	if layout_tpl.Execute(&layout_buf, email) != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	text, err := ioutil.ReadAll(email.Reader)
	if err != nil {
		config.Logger.Errorf("Error: %v", err)
		return nil
	}

	msg := append(layout_buf.Bytes(), "\r\n"...)
	msg = append(msg, text...)

	auth := smtp.PlainAuth("", w.SMTP_USER, w.SMTP_PASSWORD, w.SMTP_HOST)
	if err := smtp.SendMail(w.SMTP_HOST+":"+w.SMTP_PORT, auth, email.FromAddress, []string{email.ToAddress}, msg); err != nil {
		config.Logger.Errorf("Error: %v", err.Error())
		return nil
	}

	config.Logger.Info("Email Sent Successfully!")

	return nil
}
