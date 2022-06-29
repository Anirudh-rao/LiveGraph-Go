package main

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pusher/pusher-http-go"
)

var client = pusher.Client{
	AppID:      "1430598",
	Key:        "a986363e1e854c083369",
	Secret:     "3bf0f6eda34b1f884b21",
	Secure:     true,
	Cluster:    "ap2",
	HTTPClient: &http.Client{},
}

type visitData struct {
	Pages int
	Count int
}

func simulate(c echo.Context) error {
	setInterval(func() {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		newVisitsData := visitData{
			Pages: r1.Intn(100),
			Count: r1.Intn(100),
		}
		client.Trigger("visitorsCount", "addNumber", newVisitsData)
	}, 2500, true)

	return c.String(http.StatusOK, "Simulation begun")
}
func setInterval(ourFunc func(), milliseconds int, async bool) chan bool {

	// How often to fire the passed in function in milliseconds
	interval := time.Duration(milliseconds) * time.Millisecond

	// Setup the ticker and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)
	clear := make(chan bool)

	// Put the selection in a go routine so that the for loop is none blocking
	go func() {
		for {
			select {
			case <-ticker.C:
				if async {
					// This won't block
					go ourFunc()
				} else {
					// This will block
					ourFunc()
				}
			case <-clear:
				ticker.Stop()
				return
			}
		}
	}()

	// We return the channel so we can pass in
	// a value to it to clear the interval
	return clear
}
func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	// Middleware

	// Define the HTTP routes
	e.File("/", "public/index.html")
	e.File("/style.css", "public/style.css")
	e.File("/app.js", "public/app.js")
	e.GET("/simulate", simulate)

	// Start server
	e.Logger.Fatal(e.Start(":9000"))
}
