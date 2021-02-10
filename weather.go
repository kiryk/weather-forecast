package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

var port = flag.String("p", "80", "port number")

// Formats for API queries
var locationFmt = "http://www.metaweather.com/api/location/search/?lattlong=%s"
var weatherFmt = "http://www.metaweather.com/api/location/%d/"
var publicIpFmt = "https://ip.seeip.org/json"
var geolocationFmt = "http://ip-api.com/json/%s"

// Weather forecast page template
var pageTpl *template.Template

// Used for seeip.org responses
type SeeIpResponse struct {
	Ip string `json:"ip"`
}

// Used for ip-api responses
type IPApiResponse struct {
	Status  string  `json:"status"`
	Message string  `json:"message"`
	Latt    float32 `json:"lat"`
	Long    float32 `json:"lon"`
}

// Full forecast data, i.e. full MetaWeather API response
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

// Location data in MetaWeather responses
type Location struct {
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Lattlong     string `json:"latt_long"`
	Woeid        int    `json:"woeid"`
//Distance     int    `json:"distance"`
}

// Source information in MetaWeather responses
type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// A single day forecast in MetaWeather responses
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

// Some values from MetaWeather use imperial units, this function
// scales them to metric
func imperialToMetric(f *Forecast) {
	const KmPerMile = 1.609344

	f.WindSpeed *= KmPerMile
	f.Visibility *= KmPerMile
}

// Convert a YYYY-MM-DD date string into a nicer one, i.e.:
// - if the day is today, change it to "Today"
// - if the day is tomorrow, change it to "Tomorrow"
// - otherwise format it to be like "Monday, 2 Jan 2006"
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

// Fetch a structure from a JSON API. The `data` value has to be passed
// as a pointer to the destination structure.
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

// Fetch approximate geolocation using given ip.
func fetchLattlong(ip string) (string, error) {
	var resp IPApiResponse

	query := fmt.Sprintf(geolocationFmt, ip)
	if err := fetchStruct(query, &resp); err != nil {
		return "", err
	}

	if resp.Status != "success" {
		return "", fmt.Errorf("%s geolocation failed: %s", ip, resp.Message)
	}

	return fmt.Sprintf("%g,%g", resp.Latt, resp.Long), nil
}

// Fetch _our_ public IP
func fetchPublicIp() (string, error) {
	var resp SeeIpResponse

	if err := fetchStruct(publicIpFmt, &resp); err != nil {
		return "", fmt.Errorf("coudln't get public ip: %e", err)
	}

	return resp.Ip, nil
}

// Generate a full page containing the forecast
func showWeather(w http.ResponseWriter, r *http.Request) {
	var locations []Location
	var weather Weather
	var lattlong string

	query := r.URL.Query()

	if lattlongs, ok := query["lattlong"]; ok {
		lattlong = lattlongs[0]
	} else {
		var err error
		var ip string

		ip = r.RemoteAddr[:strings.IndexByte(r.RemoteAddr, ':')]
		if lattlong, err = fetchLattlong(ip); err != nil {
			// we failed, but perhaps the IP is just local

			if ip, err = fetchPublicIp(); err != nil {
				http.Error(w, "couldn't get users public IP",
			             http.StatusInternalServerError)
				log.Println(err)
				return
			}
			if lattlong, err = fetchLattlong(ip); err != nil {
				http.Error(w, "couldn't get users geolocation",
				           http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
	}

	locationQuery := fmt.Sprintf(locationFmt, lattlong)
	if err := fetchStruct(locationQuery, &locations); err != nil {
		http.Error(w, "could not fetch users woeid",
		           http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if len(locations) < 1 {
		http.Error(w, "user's location was not found",
		           http.StatusInternalServerError)
		log.Println(fmt.Errorf("user's location was not found"))
		return
	}

	weatherQuery := fmt.Sprintf(weatherFmt, locations[0].Woeid)
	if err := fetchStruct(weatherQuery, &weather); err != nil {
		http.Error(w, "user's location was not found",
		           http.StatusInternalServerError)
		log.Println(err)
		return
	}

	for i := range weather.Forecasts {
		imperialToMetric(&weather.Forecasts[i]);
		dateToReadable(&weather.Forecasts[i].ApplicableDate);
	}

	if err := pageTpl.Execute(w, weather); err != nil {
		http.Error(w, "user's location was not found",
		           http.StatusInternalServerError)
		log.Println(err)
		return
	}
}

func main() {
	var err error

	flag.Parse()

	pageTpl, err = template.New("page.tpl").ParseFiles("page.tpl")
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", showWeather)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
