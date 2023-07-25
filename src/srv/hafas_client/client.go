package hafasclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type (
	DPlatForm struct {
		Txt string `json:"txt"`
	}
	DStbStop struct {
		PlatPrognose DPlatForm `json:"dPltfS"`
		Cancelled    bool      `json:"dCncl" optional:"true"`
		TimePrognose string    `json:"dTimeR"` // clock time prognosis
		TimePlanned  string    `json:"dTimeS"` // clock time planned
		ProdIndex    int       `json:"dProdX"` // index for lookup of bus name
	}

	DJnyL struct {
		DirTxt   string   `json:"dirTxt"`
		Jid      string   `json:"jid"`
		StbStop  DStbStop `json:"stbStop"`
		StartDay string   `json:"trainStartDate"`
	}
	DRes struct {
		Common DCommon `json:"common"`
		JnyL   []DJnyL `json:"jnyL"`
	}
	DProdEntry struct {
		Name  string `json:"name"`
		Short string `json:"nameS"`
	}
	DCommon struct {
		ProdList []DProdEntry `json:"prodL"`
	}
	DSvcResL struct {
		Res DRes `json:"res"`
	}
	DeparturesRaw struct {
		SvcResL []DSvcResL `json:"svcResL"`
	}
)

type HafasClient struct {
	httpCLient *http.Client
	endpoint   string
}

type Departure struct {
	Destination  string `json:"destination"`
	TimePlanned  string `json:"planned"`
	TimePrognose string `json:"prognosed"`
	Day          string `json:"day"`
	Platform     string `json:"platform"`
	BusName      string `json:"bus"`
	BusNameShort string `json:"busS"`
	JourneyId    string `json:"journey"`
	Cancelled    bool   `json:"cancelled"`
}

func (client *HafasClient) formatStation(station string) map[string]any {
	result := map[string]any{
		"type": "S",
		"lid":  fmt.Sprintf("A=1@L=%s@", station),
	}
	return result
}

func (client *HafasClient) formatJnyFltrl() []map[string]any {
	result := []map[string]any{
		{"type": "PROD", "mode": "INC", "value": "2047"}, // Guessed: Bus??
	}
	return result
}

func (client *HafasClient) formatStationBoard(station string) []map[string]any {
	result := []map[string]any{
		{
			"meth": "StationBoard",
			"req": map[string]any{
				"stbLoc":   client.formatStation(station),
				"type":     "DEP",
				"maxJny":   40,
				"jnyFltrL": client.formatJnyFltrl(),
			},
		},
	}
	return result
}

func (client *HafasClient) formatDepReq(station string) map[string]any {
	result := map[string]any{
		"svcReqL":   client.formatStationBoard(station),
		"ver":       "1.26",
		"lang":      "deu",
		"formatted": false,
		"auth": map[string]any{
			"aid":  "4vV1PaulHallo511icH",
			"type": "AID",
		},
		"client": map[string]any{
			"id":   "AVV_AACHEN",
			"name": "webapp",
			"type": "WEB",
		},
	}
	return result
}

func (draw *DeparturesRaw) toDepartures() ([]Departure, error) {
	if len(draw.SvcResL) > 0 {
		prodList := draw.SvcResL[0].Res.Common.ProdList
		deps := make([]Departure, len(draw.SvcResL[0].Res.JnyL))

		for i, jny := range draw.SvcResL[0].Res.JnyL {
			if jny.StbStop.ProdIndex < len(prodList) && jny.StbStop.ProdIndex >= 0 {
				dep := Departure{
					Destination:  jny.DirTxt,
					TimePlanned:  jny.StbStop.TimePlanned,
					TimePrognose: jny.StbStop.TimePrognose,
					Platform:     jny.StbStop.PlatPrognose.Txt,
					BusName:      prodList[jny.StbStop.ProdIndex].Name,
					BusNameShort: prodList[jny.StbStop.ProdIndex].Short,
					Day:          jny.StartDay,
					JourneyId:    jny.Jid,
					Cancelled:    jny.StbStop.Cancelled,
				}
				deps[i] = dep
			}
		}

		return deps, nil
	}
	return []Departure{}, fmt.Errorf("NO_DEPARTURES")
}

func (client *HafasClient) GetDepartures(station string) ([]Departure, error) {
	// log.Println("GetDepartures")

	data, err := json.Marshal(client.formatDepReq(station))
	if err != nil {
		return []Departure{}, err
	}

	req, err := http.NewRequest("POST", client.endpoint, bytes.NewBuffer(data))
	if err != nil {
		return []Departure{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Accept-Encoding", "gzip, br, deflate")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("user-agent", "BeepBoop")
	req.Header.Set("connection", "keep-alive")

	res, err := client.httpCLient.Do(req)
	if err != nil {
		return []Departure{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Departure{}, err
	}

	var content DeparturesRaw
	if err := json.Unmarshal(body, &content); err != nil {
		return []Departure{}, err
	}

	deps, err := content.toDepartures()
	return deps, err
}

func NewClient() HafasClient {
	return HafasClient{
		endpoint:   "https://auskunft.avv.de/bin/mgate.exe",
		httpCLient: &http.Client{},
	}
}
