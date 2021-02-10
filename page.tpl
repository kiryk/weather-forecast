<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8" />
		<style>
			body {
				max-width: 500pt;
				margin: auto;
				font-family: sans-serif;
				font-size: 12pt;
				font-weight: lighter;
			}

			details {
				cursor: pointer;
			}

			h1 {
				text-align: center;
				font-size: 400%;
				color: gold;
				weight: bold;
			}

			h2 {
				text-align: center;
				font-size: 150%;
				color: cornflowerblue;
				weight: bold;
			}

			.weather {
				margin-bottom: 10pt;
				box-shadow: 0 0 3pt lightskyblue;
				border: solid 1px lightgray;
				border-radius: 8pt;
				overflow: hidden;
			}

			.weather img {
				height: 12pt;
			}

			.weather p, .pair, .triplet {
				margin: 0;
				padding: 10pt;
			}

			.pair p {
				display: inline-block;
				margin: 0pt;
				padding: 0pt;
				width: 50%
			}

			.triplet p {
				display: inline-block;
				margin: 0pt;
				padding: 0pt;
				width: 33.33%
			}

			.weather > *:nth-child(odd) {
				background-color: aliceblue;
			}

			.weather > *:nth-child(even) {
				background-color: white;
			}

			.weather > *:first-child {
				font-weight: bold;
				background-color: lightskyblue;
			}
		</style>
	</head>
	<body>
		<h1>Weather</h1>
		<h2>{{.LocationType}} of {{.Title}}</h2>
		{{range .Forecasts}}
			<div class="weather">
				<p>{{.ApplicableDate}}</p>
				<p>
					{{.WeatherStateName}}
					<img src="https://www.metaweather.com/static/img/weather/{{.WeatherStateAbbr}}.svg">
				</p>
				<div class="pair">
					<p>Min: {{.MinTempCelsius}}°C</p><p>Max: {{.MaxTempCelsius}}°C</p>
				</div>
				<div class="pair">
					<p>Pressure: {{.AirPressureHpa}} Hpa</p><p>Wind: {{.WindSpeed}} km/h</p>
				</div>
				<div class="pair">
					<p>Humidity: {{.HumidityPercent}}%</p><p>Visibility: {{.Visibility}} km</p>
				</div>
			</div>
		{{end}}
		<details>
 			<summary>Details about this forecast</summary>
			<p>
				The information in the forecast comes from the following sources:
				<ul>
					{{range .Sources}}<li><a href="{{.URL}}">{{.Title}}</a></li>{{end}}
				</ul>
				and has been collected and shared by <a href="https://www.metaweather.com">MetaWeather</a> API.
			</p>
			<p>Your geolocation has been obtained thanks to <a href="https://ip-api.com">ip-api</a> and <a href="https://seeip.org">SeeIP</a>.</p>
		</details>
	</body>
</html>
