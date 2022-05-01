package utils

import (
	"bytes"
	"html/template"
)

type confirmationTemplateData struct {
	email           string
	confirmationUrl string
}

func NewConfirmationEmailTemplate(path, email, confirmationUrl string) (string, error) {
	data := confirmationTemplateData{email, confirmationUrl}
	return parseConfirmationEmailTemplate(path, data)
}

func parseConfirmationEmailTemplate(path string, data confirmationTemplateData) (string, error) {
	t := template.New("email confirmation")
	t, err := t.Parse(path)
	if err != nil {
		return "", err
	}

	tpl := bytes.Buffer{}
	err = t.Execute(&tpl, data)
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}
