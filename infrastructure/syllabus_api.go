package infrastructure

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"syllabus-scraper/domain"
)

type AuraResponse struct {
	Actions []struct {
		State       string `json:"state"`
		ReturnValue struct {
			ReturnValue struct {
				Result []SyllabusData `json:"result"`
			} `json:"returnValue"`
		} `json:"returnValue"`
	} `json:"actions"`
}

type SyllabusData struct {
	ID            string `json:"Id"`
	CourseName    string `json:"R_SlCourseName__c"`
	PersonalName  string `json:"R_SlPersonalName__c"`
	WeekDayPeriod string `json:"R_SlWeekDayPeriod__c"`
	CampusInfo    string `json:"R_SlCampusInfo__c"`
}

type SyllabusAPI struct {
	client *http.Client
}

func NewSyllabusAPI() *SyllabusAPI {

	jar, _ := cookiejar.New(nil)

	return &SyllabusAPI{
		client: &http.Client{
			Jar: jar,
		},
	}
}

func (s *SyllabusAPI) GetSyllabus(query domain.SyllabusQuery) ([]domain.Syllabus, error) {

	baseURL := "https://syllabus.ritsumei.ac.jp/syllabus/s/"
	apiURL := "https://syllabus.ritsumei.ac.jp/syllabus/s/sfsites/aura"

	req, _ := http.NewRequest("GET", baseURL, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	auraContext := `{
		"mode":"PROD",
		"fwuid":"",
		"app":"siteforce:communityApp",
		"loaded":{},
		"dn":[],
		"globals":{},
		"uad":true
	}`

	message := map[string]interface{}{
		"actions": []map[string]interface{}{
			{
				"id":                "1;a",
				"descriptor":        "aura://ApexActionController/ACTION$execute",
				"callingDescriptor": "UNKNOWN",
				"params": map[string]interface{}{
					"namespace": "",
					"classname": "R_SyllabusPublicPageController",
					"method":    "getSyllabusRecords",
					"params": map[string]interface{}{
						"action": map[string]interface{}{
							"lang":               "ja",
							"keyword":            nil,
							"faculty":            query.Faculty,
							"year":               query.Year,
							"term":               query.Term,
							"week":               query.Week,
							"period":             query.Period,
							"professionalCareer": nil,
							"limits":             100,
						},
					},
					"cacheable":      false,
					"isContinuation": false,
				},
			},
		},
	}

	messageBytes, _ := json.Marshal(message)

	form := url.Values{
		"message":      {string(messageBytes)},
		"aura.context": {auraContext},
		"aura.pageURI": {"/syllabus/s/?language=ja"},
		"aura.token":   {"null"},
	}

	postReq, _ := http.NewRequest(
		"POST",
		apiURL,
		bytes.NewBufferString(form.Encode()),
	)

	postReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	postReq.Header.Set("User-Agent", "Mozilla/5.0")

	postResp, err := s.client.Do(postReq)
	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(postResp.Body)
	postResp.Body.Close()

	var auraResp AuraResponse

	err = json.Unmarshal(body, &auraResp)
	if err != nil {
		return nil, err
	}

	results := auraResp.Actions[0].ReturnValue.ReturnValue.Result

	var syllabuses []domain.Syllabus

	for _, r := range results {

		syllabuses = append(syllabuses, domain.Syllabus{
			ID:            r.ID,
			CourseName:    r.CourseName,
			PersonalName:  r.PersonalName,
			WeekDayPeriod: r.WeekDayPeriod,
			CampusInfo:    r.CampusInfo,
		})
	}

	return syllabuses, nil
}
