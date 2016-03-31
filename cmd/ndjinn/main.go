package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

var usage string = `


									oh wow.
`//`

var t0 = time.Now()
var Globalmessage = make(chan string, 3)
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

	// Help Mode and bunk options
	if len(os.Args) > 1 {
			switch os.Args[1] {
				case "-V", "--version": os.Exit(0);
				case "-h", "--help": fmt.Println(usage);os.Exit(0)
				case "-v", "--verbose": break;
				default: fmt.Printf("\nNot a valid command. Try %s\n", os.Args[0] );os.Exit(0)
			}
	}

fmt.Printf("Boot: Started (%s)\n",time.Now().Format("2006-01-02 15:04:05"))

	// Load the configuration file (./config/config.json)
	fmt.Printf("Boot: Loading config.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)
	fmt.Printf(" done.\n")

	// Configure the session cookie store
	fmt.Printf("Boot: Loading cookie store.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	session.Configure(config.Session)
	fmt.Printf(" done.\n")

	fmt.Printf("Boot: Loading PayPal.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	checkout.Configure(config.Checkout)
	fmt.Printf(" done.\n")

	// Connect to database (./database/database.db)

	fmt.Printf("Boot: Loading database.")
	fmt.Printf(".")
	_ = os.Mkdir("backups", 0700)
	fmt.Printf(".")
	fmt.Printf(".")
	database.Connect(config.Database)
	fmt.Printf(" done.\n")

	// Run Maint and Backup
	_, err := maint.RunSchedule()

	// // // //

	database.Close()
	if err != nil {
		fmt.Println(err)
	}
	database.Connect(config.Database)
	fmt.Println("Boot: Database Backup completed.")

	err = os.Mkdir("logs", 0700)
	if err != nil && err.Error() != "mkdir logs: file exists" {
		log.Fatal(err)
	}

	fmt.Println("Boot: Logging requests to ./logs/access.log")
	fmt.Println("Boot: Other logging at ./logs/debug.log")

	// // // //

	// Configure the Google reCAPTCHA prior to loading view plugins
	recaptcha.Configure(config.Recaptcha)

	// Setup the views (./template and ./static)
	fmt.Printf("Boot: Loading views.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	view.Configure(config.View)
	fmt.Printf(" done.\n")

	fmt.Printf("Boot: Loading templates.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	view.LoadTemplates(config.Template.Root, config.Template.Children)
	fmt.Printf(" done.\n")

	fmt.Printf("Boot: Loading view plugins.")
	fmt.Printf(".")
	fmt.Printf(".")
	fmt.Printf(".")
	view.LoadPlugins(
		plugin.TagHelper(config.View),
		plugin.NoEscape(),
		recaptcha.RecaptchaPlugin())

	fmt.Printf(" done.\n")

	// Monitor Uptime
	go func() {
		for {
			amt := time.Duration(rand.Intn(25)+10000)
			time.Sleep(time.Second * amt)
			t1 := time.Now()
			log.Println("Info: Running for " + t1.Sub(t0).String())
		}
	}()



	// Test push notification (for site-wide announcements to all online users)
	go func() {
		for {
				time.Sleep(1 * time.Minute)
				fmt.Println("Gopher")

		}
	}()

	// This is our maintainance every hour. (internal cron)
	go func() {
		c := time.Tick(1 * time.Minute)
		go func() {
			for range c {

				log.Printf("Maint: Starting Maintainance Routine")

				// Database Backup
				_, err := maint.RunSchedule()

				database.Close()
				if err != nil {
					fmt.Println(err)
				}
				database.Connect(config.Database)
				log.Println("Maint: Database Backup completed.")

				_, err = maint.BankTransfer()
				if err != nil {
					fmt.Println(err)
				}
				// Bank Transfer
				log.Println("Maint: Bank Transfer completed.")

				// Bitcoin Sync
				_, err = maint.BitcoinSync()
				if err != nil {
					fmt.Println(err)
				}
				log.Println("Maint: Bitcoin Sync completed.")

			}
		}()

	}()
	fmt.Printf("Boot: Complete (%s)\n",time.Now().Format("2006-01-02 15:04:05"))

if len(os.Args) == 1 {
	if !InitLog() {
		panic("error")
	}
}
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
	f, err := os.OpenFile("logs/debug.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		log.Printf("error opening file: %v", err)
		log.Fatal("Hint: touch ./debug.log, or chown/chmod it so that the "+os.Args[0]+" process can access it.")
		return false
	}
	fmt.Println("Switched to log file. Nothing more to see here.")
	log.SetOutput(f)
	return true
}
