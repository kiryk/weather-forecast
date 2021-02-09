package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	//"io/ioutil"
	"log"
	"net/http"
	//"os"
	//"path"
	//"sort"
	//"strings"
	"time"
)

var port = flag.String("p", "80", "port number")

var locationFmt = "http://www.metaweather.com/api/location/search/?lattlong=%s"
var weatherFmt = "http://www.metaweather.com/api/location/%d/"

var consentTpl *template.Template
var weatherTpl *template.Template

type Location struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Lattlong     string `json:"latt_long"`
	Woeid        int    `json:"woeid"`
//Distance     int    `json:"distance"`
}

type Forecast struct {
	Id               int     `json:"id"`
	ApplicableDate   string  `json:"applicable_date"`
	WeatherStateName string  `json:"weather_state_name"`
	WeatherStateAbbr string  `json:"weather_state_abbr"`
	WindSpeed        float32 `json:"wind_speed"`
	WindDirDegrees   float32 `json:"wind_direction"`
	WindDirCompass   string  `json:"wind_direction_compass"`
	MinTempCelsius   float32 `json:"min_temp"`
	MaxTempCelsius   float32 `json:"max_temp"`
	TheTempCelsius   float32 `json:"the_temp"`
	AirPressureHpa   float32 `json:"air_pressure"`
	HumidityPercent  float32 `json:"humidity"`
	Visibility       float32 `json:"visibility"`
	QualityPercent   float32 `json:"predictability"`
}

type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Weather struct {
	Location

	Time      string     `json:"time"`
	Timezone  string     `json:"timezone_name"`
	Sunrise   string     `json:"sun_rise"`
	Sunset    string     `json:"sun_set"`

	Parent    Location   `json:"parent"`
	Forecasts []Forecast `json:"consolidated_weather"`
	Sources   []Source   `json:"sources"`
}

func imperialToMetric(f *Forecast) {
	const KmPerMile = 1.609344

	f.WindSpeed *= KmPerMile
	f.Visibility *= KmPerMile
}

func dateToReadable(date *string) {
	wt, _ := time.Parse("2006-01-02", *date)
	wy, wm, wd := wt.Date()
	ty, tm, td := time.Now().Date();

	switch {
	case wy == ty && wm == tm && wd == td:
		*date = "Today"
	case wy == ty && wm == tm && wd == td+1:
		*date = "Tomorrow"
	default:
		*date = wt.Format("Monday, 2 Jan 2006")
	}
}

func fetchStruct(query string, data interface{}) error {
	client := http.Client{}

	res, err := client.Get(query)
	if err != nil {
		return fmt.Errorf("failed to fetch: %s", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status from API server: %s", res.Status)
	}

	if err := json.NewDecoder(res.Body).Decode(data); err != nil {
		return fmt.Errorf("could not parse API response: %e", err)
	}
	return nil
}

func askCoords(w http.ResponseWriter, r *http.Request) {
	if err := consentTpl.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}
}

func showWeather(w http.ResponseWriter, r *http.Request) {
	var locations []Location
	var weather Weather
	var lattlong string

	query := r.URL.Query()

	if lattlongs, ok := query["lattlong"]; !ok {
		askCoords(w, r)
		return
	} else {
		lattlong = lattlongs[0]
	}

	locationQuery := fmt.Sprintf(locationFmt, lattlong)
	if err := fetchStruct(locationQuery, &locations); err != nil {
		log.Println(err)
		return
	}

	if len(locations) < 1 {
		askCoords(w, r)
		return
	}

	weatherQuery := fmt.Sprintf(weatherFmt, locations[0].Woeid)
	if err := fetchStruct(weatherQuery, &weather); err != nil {
		log.Println(err)
		return
	}

	for i := range weather.Forecasts {
		imperialToMetric(&weather.Forecasts[i]);
		dateToReadable(&weather.Forecasts[i].ApplicableDate);
	}

	if err := weatherTpl.Execute(w, weather); err != nil {
		log.Println(err)
		return
	}
}

func main() {
	var err error

	flag.Parse()

	/*Funcs(template.FuncMap{"title": strings.Title, "cutext": cutext})*/
	consentTpl, err = template.New("consent.html").ParseFiles("consent.html")
	if err != nil {
		log.Fatal(err)
	}

	weatherTpl, err = template.New("weather.html").ParseFiles("weather.html")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", showWeather)
	http.HandleFunc("/theme.css", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/theme.css")
	})

	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
