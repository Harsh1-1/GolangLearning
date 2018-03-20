package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type AttackDetails struct {
	Rate     uint64
	Duration int
	Target   string
}

func attack(r uint64, d int, t string) {
	dur := time.Duration(d)
	rate := r // per second
	duration := dur * time.Second
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    t,
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration) {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Max Latencies: %s\n", metrics.Latencies.Mean)
	fmt.Printf("Rate: %f", metrics.Rate)
	fmt.Printf("Duration: %f", metrics.Duration)
	fmt.Printf("Bytes in %f", metrics.BytesIn.Total)
	fmt.Println(metrics.Latencies)
	fmt.Print(metrics)
	// return metrics
	// metrics.Latencies

}

func main() {

	tmpl := template.Must(template.ParseFiles("attack.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			tmpl.Execute(w, nil)
			return
		}
		rate, _ := strconv.Atoi(r.FormValue("rate"))
		duration, _ := strconv.Atoi(r.FormValue("duration"))
		details := AttackDetails{
			Rate:     uint64(rate),
			Duration: duration,
			Target:   r.FormValue("target"),
		}

		// do something with details

		// fmt.Println(details.Email)
		// fmt.Fprintf(w, "Hey rate is %d and duration is: %d", details.Rate, details.Duration)

		go attack(details.Rate, details.Duration, details.Target)

		_ = details

		tmpl.Execute(w, struct{ Success bool }{true})
	})

	http.ListenAndServe(":8000", nil)
	// attack(100, 4, "https://google.co.in")
}
