# weather-forecast

## Description
`weather-forecast` is a simple http server written in Go that shows local weather forecast for the nearest days in a form of a minimalist website.

## Compilation
In order to compile the source, enter the source code directory and type:

```
go build
```

The compiler should name the binary as `weather` (or `weather.exe` on Windows)
which will be assumed in the next subsection.

## Running
In order to run the server on a UNIX system, type:

```
./weather [-p PORT-NUMBER]
```

To run it on Windows, type:

```
weather [-p PORT-NUMBER]
```

Using the `-p` flag is optional, with the default port numbr being `80`.

## Dependencies
Besides Go's standard library, the application uses a few external APIs to accomplish some of its tasks:

- [MetaWeather](www.metaweather.com), to obtain weather data,
- [ip-api](ip-api.com), to approximate user's geolocation,
- [SeeIP](ip.seeip.org), to easily obtain public IP of the server.
