package processes

import (
	"chatbot"
	"html/template"
	"strings"
)

const (
	TPL_QRCODE          = `@{{.Name}}, please do not send qrcode`
	TPL_SERVICE_HOTLINE = `@{{.Name}}, our 24 hours hotline is {{.Phone}}`
)

type Handle func(record *chatbot.ReviewRecord) ([]byte, error)

func HandleQrcode(record *chatbot.ReviewRecord) ([]byte, error) {
	if len(record.ImgUrl) > 0 {
		tmpl, err := template.New("qrcode").Parse(TPL_QRCODE)
		if err != nil {
			return nil, err
		}

		data := struct {
			Name string
		}{
			Name: record.SenderNickname,
		}

		var output strings.Builder
		err = tmpl.Execute(&output, data)
		if err != nil {
			return nil, err
		}

		return []byte(output.String()), nil
	}
	return nil, nil
}

func HandleHotLine(record *chatbot.ReviewRecord) ([]byte, error) {
	if strings.Contains(record.Text, "hotline") {
		tmpl, err := template.New("hotline").Parse(TPL_SERVICE_HOTLINE)
		if err != nil {
			return nil, err
		}

		data := struct {
			Name  string
			Phone string
		}{
			Name:  record.SenderNickname,
			Phone: "86-010-11116666",
		}

		var output strings.Builder
		err = tmpl.Execute(&output, data)
		if err != nil {
			return nil, err
		}

		return []byte(output.String()), nil
	}
	return nil, nil
}
