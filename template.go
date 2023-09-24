package email_api

import (
	"bytes"
	"context"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	htmlTemplate "html/template"
)

type ReadTemplateError struct {
	templateName string
	err          error
}

func (readTemplateError ReadTemplateError) Error() string {
	return fmt.Sprintf("Failed to read template %s - %w", readTemplateError.templateName, readTemplateError.err)
}

func readTemplate(ctx context.Context, app *app, fileName string) (string, error) {
	data, err := app.fileStorage.ReadAll(ctx, fileName)
	if err != nil {
		return "", ReadTemplateError{templateName: fileName, err: err}
	}

	return string(data), nil
}

func executeTemplate(template string, data map[string]interface{}) (string, error) {
	// We expect email templates to use title case variables {{ .Title }}
	titleData := make(map[string]interface{})
	for k, v := range data {
		titleK := cases.Title(language.AmericanEnglish)
		titleData[titleK.String(k)] = v
	}

	t, err := htmlTemplate.New("email").
		Option("missingkey=error").
		Parse(template)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, titleData)
	if err != nil {
		return "", err
	}
	return tpl.String(), nil
}
