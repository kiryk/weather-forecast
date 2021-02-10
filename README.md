# weather-forecast

## Description
`weather-forecast` is a simple HTTP server written in Go that shows local weather forecast for the nearest days in a form of a minimalist website.

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

Using the `-p` flag is optional, with the default port number being `80`.

## Dependencies
Besides Go's standard library, the application uses a few external APIs to accomplish some of its tasks:

- [MetaWeather](https://www.metaweather.com), to obtain weather data,
- [ip-api](https://ip-api.com), to approximate user's geolocation,
- [SeeIP](https://seeip.org), to easily obtain public IP of the server.

## What could be changed
First of all, the server works over HTTP instead of HTTPS. This is however intended as the app was meant to be simple.

Secondly, the server fetches a little bit more data than is presented to the user, so more information can be easily added to the website by making small changes in the `page.tpl` file. It hasn't been done only to keep the layout clear.
