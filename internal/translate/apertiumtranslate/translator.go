package apertiumtranslate

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/eeeXun/gtt/internal/translate/core"
)

const (
	textURL = "https://www.apertium.org/apy/translate?langpair=%s|%s&q=%s"
)

type ApertiumTranslate struct {
	*core.Language
	*core.TTSLock
	core.EngineName
}

func NewApertiumTranslate() *ApertiumTranslate {
	return &ApertiumTranslate{
		Language:   core.NewLanguage(),
		TTSLock:    core.NewTTSLock(),
		EngineName: core.NewEngineName("ApertiumTranslate"),
	}
}

func (t *ApertiumTranslate) GetAllLang() []string {
	return lang
}

func (t *ApertiumTranslate) Translate(message string) (translation *core.Translation, err error) {
	translation = new(core.Translation)
	var data map[string]interface{}

	urlStr := fmt.Sprintf(
		textURL,
		langCode[t.GetSrcLang()],
		langCode[t.GetDstLang()],
		url.QueryEscape(message),
	)
	res, err := http.Get(urlStr)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	if len(data) <= 0 {
		return nil, errors.New("Translation not found")
	}

	switch res.StatusCode {
	case 200:
		translation.TEXT = fmt.Sprintf("%v",
			data["responseData"].(map[string]interface{})["translatedText"])
	default:
		return nil, errors.New(
			fmt.Sprintf("%s does not support translate from %s to %s.\nSee available pair on %s",
				t.GetEngineName(),
				t.GetSrcLang(),
				t.GetDstLang(),
				"https://www.apertium.org/",
			))
	}

	return translation, nil
}

func (t *ApertiumTranslate) PlayTTS(lang, message string) error {
	defer t.ReleaseLock()

	return errors.New(t.GetEngineName() + " does not support text to speech")
}
