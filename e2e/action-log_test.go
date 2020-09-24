package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cucumber/godog"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type HistoryResponse struct {
	ID          int    `json:"id"`
	UserID      int    `json:"userId"`
	Project     string `json:"project"`
	Description string `json:"description"`
}

type TestData struct {
	responseStatusCode int
	responseBody       []byte
	responseJson       string
}

func (t *TestData) StepDefinitioninition1(httpMethod, url string) error {
	client := &http.Client{}

	req, err := http.NewRequest(httpMethod, url, nil)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	t.responseStatusCode = resp.StatusCode
	t.responseBody, err = ioutil.ReadAll(resp.Body)

	return nil
}

func (t *TestData) StepDefinitioninition2(code int) error {
	if t.responseStatusCode != code {
		return fmt.Errorf("неожиданный код состояния: %d != %d", t.responseStatusCode, code)
	}

	return nil
}

func (t *TestData) StepDefinitioninition3(text string) error {
	expectedText := strings.TrimSpace(string(t.responseBody))

	if expectedText != text {
		return fmt.Errorf("неожиданный текст: %s != %s", expectedText, text)
	}

	return nil
}

func (t *TestData) StepDefinitioninition4(
	httpMethod, url,
	contentType string,
	data *godog.DocString,
) (err error) {
	var r *http.Response

	switch httpMethod {
	case http.MethodPost, http.MethodGet:
		replacer := strings.NewReplacer("\n", "", "\t", "")
		cleanJson := replacer.Replace(data.Content)
		r, err = http.Post(url, contentType, bytes.NewReader([]byte(cleanJson)))
	default:
		err = fmt.Errorf("неизвестный метод: %s", httpMethod)
	}

	if err != nil {
		return
	}

	t.responseStatusCode = r.StatusCode
	t.responseBody, err = ioutil.ReadAll(r.Body)
	t.responseJson = strings.TrimSpace(string(t.responseBody))

	return
}

func (t *TestData) StepDefinitioninition5(body *godog.DocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanJson := replacer.Replace(body.Content)

	var expected, actual HistoryResponse

	if err := json.Unmarshal([]byte(cleanJson), &actual); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(t.responseJson), &expected); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("структуры должны совпадать, %v vs. %v", expected, actual)
	}

	return nil
}

func (t *TestData) StepDefinitioninition6(body *godog.DocString) error {
	replacer := strings.NewReplacer("\n", "", "\t", "")
	cleanJson := replacer.Replace(body.Content)

	var expected, actual []HistoryResponse

	if err := json.Unmarshal([]byte(cleanJson), &actual); err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(t.responseJson), &expected); err != nil {
		return err
	}

	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("структуры должны совпадать, %v vs. %v", expected, actual)
	}

	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	test := new(TestData)

	ctx.Step(`^1. Я отправляю "([^"]*)" запрос на "([^"]*)"$`, test.StepDefinitioninition1)
	ctx.Step(`^Код ответа должен быть (\d+)$`, test.StepDefinitioninition2)
	ctx.Step(`^В ответ должен быть текст "([^"]*)"$`, test.StepDefinitioninition3)

	ctx.Step(`^2. Я отправляю "([^"]*)" запрос на "([^"]*)" с заголовком "([^"]*)" и данными:$`, test.StepDefinitioninition4)
	ctx.Step(`^Код ответа должен быть (\d+)$`, test.StepDefinitioninition2)
	ctx.Step(`^В ответ должен получить следующие данные:$`, test.StepDefinitioninition5)

	ctx.Step(`^3. Я отправляю "([^"]*)" запрос на "([^"]*)" с заголовком "([^"]*)" и данными:$`, test.StepDefinitioninition4)
	ctx.Step(`^Код ответа должен быть (\d+)$`, test.StepDefinitioninition2)

	ctx.Step(`^4. Я отправляю "([^"]*)" запрос на "([^"]*)" с заголовком "([^"]*)" и данными:$`, test.StepDefinitioninition4)
	ctx.Step(`^Код ответа должен быть (\d+)$`, test.StepDefinitioninition2)
	ctx.Step(`^В ответе должен получить следующий список:$`, test.StepDefinitioninition6)
}
