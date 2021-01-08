package controllers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"goweather/app/models"

	"goweather/config"
)

var templates = template.Must(template.ParseFiles("app/views/index.html"))

type GetDatas struct {
	Tempelature float64 `json:"tempelature"`
	Pressure    float64 `json:"pressure"`
	Description string  `jso:"description"`
	IconId      string  `json:"iconid"`
	Err         error   `json:"err,omitempty"`
}

func viewIndexHandler(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func cityNameHandler(w http.ResponseWriter, r *http.Request) {
	api := models.New(config.Config.Host, config.Config.Key)
	city_name := r.FormValue("city_name")
	weather, err := api.GetWeather(city_name)
	errors := &GetDatas{
		Err: err,
	}

	if weather.Cod != 0 {
		temps := api.GetTempAndPres(weather)
		tempelature := temps[0]
		pressure := temps[1]
		description := api.GetWeatherDescription(weather)
		iconid := api.GetIconId(weather)

		datas := &GetDatas{
			Tempelature: tempelature,
			Pressure:    pressure,
			Description: description,
			IconId:      iconid,
		}

		err = templates.ExecuteTemplate(w, "index.html", *datas)
		if err != nil {
			log.Printf("cityNameHandler_ExecuteTemplate Error:%+v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = templates.ExecuteTemplate(w, "index.html", *errors)

	}
}

func StartWebServer() error {
	http.HandleFunc("/index/city/", cityNameHandler)
	http.HandleFunc("/index/", viewIndexHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", config.Config.Port), nil)
}
