package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/aerth/ndjinn/components/checkout"
	"github.com/aerth/ndjinn/components/database"
	"github.com/aerth/ndjinn/components/email"
	"github.com/aerth/ndjinn/components/jsonconfig"
	"github.com/aerth/ndjinn/components/maint"
	"github.com/aerth/ndjinn/components/recaptcha"
	"github.com/aerth/ndjinn/components/server"
	"github.com/aerth/ndjinn/components/session"
	"github.com/aerth/ndjinn/components/view"
	"github.com/aerth/ndjinn/components/view/plugin"
	"github.com/aerth/ndjinn/route"
)

var t0 = time.Now()

// Version
func Version() string {
	return "ndjinn v0.2"
}

func init() {
	// Verbose logging with file name and line number
	// log.SetFlags(log.Lshortfile)

	// Normal logging
	log.SetFlags(log.LstdFlags)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {

	// Print Version Info
	fmt.Println(Version())

	// Help Mode
	if len(os.Args) > 1 {
		fmt.Println("Usage:	")
		os.Exit(0)
	}

	// Switch to log file
	/*
		if !InitLog() {
			panic("Log Error")
		}
	*/

	// Load the configuration file (./config/config.json)
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Configure the session cookie store
	session.Configure(config.Session)
	checkout.Configure(config.Checkout)

	// Connect to database (./database.db)
	database.Connect(config.Database)

	// Run Maint and Backup
	_, err := maint.RunSchedule()

	// // // //

	database.Close()
	if err != nil {
		fmt.Println(err)
	}
	database.Connect(config.Database)
	log.Println("Database Backup completed.")

	// // // //

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views (./template and ./static)
	view.Configure(config.View)
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		recaptcha.RecaptchaPlugin())

	// Monitor Uptime
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			t1 := time.Now()
			log.Println("Running for " + t1.Sub(t0).String())
			//	fmt.Println(database.View("user", "test@example.com", &user))
		}
	}()

	// This is our maintainance every hour. (internal cron)
	go func() {
		c := time.Tick(1 * time.Minute)
		go func() {
			for range c {

				log.Printf("Starting Maintainance Routine")

				// Database Backup
				_, err := maint.RunSchedule()

				database.Close()
				if err != nil {
					fmt.Println(err)
				}
				database.Connect(config.Database)
				log.Println("Database Backup completed.")

				_, err = maint.BankTransfer()
				if err != nil {
					fmt.Println(err)
				}
				// Bank Transfer
				log.Println("Bank Transfer completed.")

				// Bitcoin Sync
				_, err = maint.BitcoinSync()
				if err != nil {
					fmt.Println(err)
				}
				log.Println("Bitcoin Sync completed.")

			}
		}()

	}()

	// Start the listener
	server.Run(route.LoadHTTP(), route.LoadHTTPS(), config.Server)
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database  database.DatabaseInfo   `json:"Database"`
	Email     email.SMTPInfo          `json:"Email"`
	Recaptcha recaptcha.RecaptchaInfo `json:"Recaptcha"`
	Server    server.Server           `json:"Server"`
	Session   session.Session         `json:"Session"`
	Template  view.Template           `json:"Template"`
	View      view.View               `json:"View"`
	Checkout  checkout.CheckoutInfo   `json:"Checkout"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

// InitLog switches the log engine to a file, rather than stdout
func InitLog() bool {
	f, err := os.OpenFile("./debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Printf("error opening file: %v", err)
		log.Fatal("Hint: touch ./debug.log, or chown/chmod it so that the deltanine process can access it.")
		return false
	}
	log.Println("Powered on. Switched to log file.")
	log.SetOutput(f)
	log.Println("Powered on. Switched to log file.")
	return true
}
