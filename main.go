// This is supposed to be a webserver responding to nearly all request with 418 I'm a teapod, but logging those requests to a file.
// It should return a redirect to my github page for the robots.txt path.

package main

import (
	"net/http"
	"os"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
)

type logger struct{ *logrus.Logger }

func (l *logger) Say(msg string) {
	l.Info(msg)
}
func (l *logger) Sayf(fmt string, args ...interface{}) {
	l.Infof(fmt, args)
}
func (l *logger) SayWithField(msg string, k string, v interface{}) {
	l.WithField(k, v).Info(msg)
}
func (l *logger) SayWithFields(msg string, fields map[string]interface{}) {
	l.WithFields(fields).Info(msg)
}

func NewLogger() *logger {

	logLevel := logrus.InfoLevel
	log := logrus.New()
	log.SetLevel(logLevel)
	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "logs/access.log",
		MaxSize:    50,  // megabytes
		MaxBackups: 5,   // amouts
		MaxAge:     365, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		logrus.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	log.SetOutput(colorable.NewColorableStdout())
	log.SetFormatter(&logrus.TextFormatter{
		PadLevelText:    true,
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.AddHook(rotateFileHook)

	return &logger{log}
}

func main() {
	log := NewLogger()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		requesterIP := r.RemoteAddr
		log.Printf(
			"%s\t\t%s\t\t%s\t\t%v",
			r.Method,
			r.RequestURI,
			requesterIP,
			time.Now(),
		)

		log.WithFields(logrus.Fields{
			"method":     r.Method,
			"requester":  requesterIP,
			"requestURI": r.RequestURI,
			"time":       time.Now(),
		}).Info("I'm a teapot!")

		http.Error(w, http.StatusText(http.StatusTeapot), http.StatusTeapot)
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		requesterIP := r.RemoteAddr

		log.WithFields(logrus.Fields{
			"method":     r.Method,
			"requester":  requesterIP,
			"requestURI": r.RequestURI,
			"time":       time.Now(),
		}).Info("Redirecting to github!")

		http.Redirect(w, r, "https://tobias-oberdoerfer.matrx.me/", http.StatusMovedPermanently)
	})

	logFile, err := os.OpenFile("access-log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	server := &http.Server{
		Addr:         ":8080",
		Handler:      http.DefaultServeMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("Starting server on port 8080")
	log.Fatal(server.ListenAndServe())
}
