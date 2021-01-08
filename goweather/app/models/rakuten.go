package models

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
)

const baseURL = "https://community-open-weather-map.p.rapidapi.com/weather"

type APIClient struct {
	host       string
	key        string
	httpClient *http.Client
}

func New(host, key string) *APIClient {
	apiClient := &APIClient{host, key, &http.Client{}}
	return apiClient
}

func (api APIClient) header() map[string]string {
	return map[string]string{
		"x-rapidapi-host": api.host,
		"x-rapidapi-key":  api.key,
	}
}

func (api *APIClient) getRequest(query map[string]string) (body []byte, err error) {
	baseURL, err := url.Parse(baseURL)
	if err != nil {
		log.Println(err)
		return
	}
	endpoint := baseURL.String()

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println(err)
		return
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	for key, value := range api.header() {
		req.Header.Add(key, value)
	}

	resp, err := api.httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return body, nil
}

type WeatherData struct {
	Coord      Coord     `json:"coord"`
	Weather    []Weather `json:"weather"`
	Base       string    `json:"base"`
	Main       Main      `json:"main"`
	Visibility int       `json:"visibility"`
	Wind       Wind      `json:"wind"`
	Clouds     Clouds    `json:"clouds"`
	Dt         int       `json:"dt"`
	Sys        Sys       `json:"sys"`
	Timezone   int       `json:"timezone"`
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Cod        int       `json:"cod"`
}

type Coord struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
}

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}
type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	TempMin   float64 `json:"temp_min"`
	TempMax   float64 `json:"temp_max"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}
type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}
type Clouds struct {
	All int `json:"all"`
}
type Sys struct {
	Type    int    `json:"type"`
	ID      int    `json:"id"`
	Country string `json:"country"`
	Sunrise int    `json:"sunrise"`
	Sunset  int    `json:"sunset"`
}

func (api *APIClient) GetWeather(place string) (WeatherData, error) {
	query := map[string]string{
		"q": place,
	}
	var weather WeatherData
	resp, err := api.getRequest(query)
	if err != nil {
		log.Printf("Error action=GetWeather_getRequest err=%s", err.Error())
		return weather, err
	}

	if err := json.Unmarshal(resp, &weather); err != nil {
		log.Printf("Responce= %+v", string(resp))
		log.Printf("action=GetWeather err=%s", err.Error())
		weather = WeatherData{}
		return weather, err
	}
	return weather, nil
}

func (api *APIClient) GetIconId(weather WeatherData) string {
	iconId := weather.Weather[0].Icon

	return iconId
}

func (api *APIClient) GetTempAndPres(weather WeatherData) []float64 {
	temp := math.Round((weather.Main.Temp-273.15)*100) / 100
	pres := float64(weather.Main.Pressure)
	mains := []float64{temp, pres}
	return mains
}

func (api *APIClient) GetWeatherDescription(weather WeatherData) string {
	desc := weather.Weather[0].Description
	return desc
}
