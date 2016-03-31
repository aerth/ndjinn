# ndjinn

**Very much work in progress.**

Notable changes to the `josephspurrier/gowebapp` skeleton would be the directory structure and adding a controller for API requests.

Also I added an experimental maintainance routine and converted the flashes to use Zurb Foundation.

Visitor can sign up, pay membership dues, access a dashboard.

Payments are going through PayPal for now. Future: Will add Bitcoin and maybe another gateway as well.

Dashboard contains forms to manipulate the BoltDB, like adding "listings" or "messages".

## Todo

  * ~~checkout bumps membership level to "paid"~~
  * ~~"paid" members see their own template page dashboard/members~~
  * Use a DB for things


## Overview

load the application settings
start the session
connect to the database
set up the templates
load the routes
attach the middleware

start maintainance routines

starts the web server.


## Structure

The project is organized into the following folders:

~~~

controller	- page logic organized by HTTP methods (GET, POST)
model		- database queries
route		- route information and middleware
components		- packages for templates, MySQL, cryptography, sessions, and json

cmd/ndjinn - build directory, contains binary and dependencies to run.
cmd/ndjinn/static		- location of statically served files like CSS and JS
cmd/ndjinn/template	- HTML templates
cmd/ndjinn/config		- application settings and database schema

~~~

There are a few external packages:

~~~
github.com/gorilla/context				- registry for global request variables
github.com/gorilla/sessions				- cookie and filesystem sessions
github.com/go-sql-driver/mysql 			- MySQL driver
github.com/haisum/recaptcha				- Google reCAPTCHA support
github.com/jmoiron/sqlx 				- MySQL general purpose extensions
github.com/josephspurrier/csrfbanana 	- CSRF protection for gorilla sessions
github.com/julienschmidt/httprouter 	- high performance HTTP request router
github.com/justinas/alice				- middleware chaining
github.com/mattn/go-sqlite3				- SQLite driver
golang.org/x/crypto/bcrypt 				- password hashing algorithm
~~~

The templates are organized into folders:

~~~
about/about.tmpl       - quick info about the app
index/anon.tmpl	       - public home page
index/auth.tmpl	       - home page once you login
login/login.tmpl	   - login page
partial/footer.tmpl	   - footer
partial/menu.tmpl	   - menu at the top of all the pages
register/register.tmpl - register page
base.tmpl		       - base template for all the pages
~~~

## Templates

There are a few template funcs that are available to make working with the templates
and static files easier:

~~~ html
<!-- CSS files with timestamps -->
{{CSS "static/css/normalize3.0.0.min.css"}}
parses to
<link rel="stylesheet" type="text/css" href="/static/css/normalize3.0.0.min.css?1435528339" />

<!-- JS files with timestamps -->
{{JS "static/js/jquery1.11.0.min.js"}}
parses to
<script type="text/javascript" src="/static/js/jquery1.11.0.min.js?1435528404"></script>

<!-- Same page hyperlinks -->
{{LINK "register" "Create a new account."}}
parses to
<a href="/register">Create a new account.</a>

<!-- Output an unescaped variable (not a safe idea, but I find it useful when troubleshooting) -->
{{.SomeVariable | NOESCAPE}}
~~~

There are a few variables you can use in templates as well:

~~~ html
<!-- Use LoginStatus=auth to determine if a user is logged in -->
{{if eq .LoginStatus "auth"}}
You are logged in.
{{else}}
You are not logged in.
{{end}}

<!-- Use BaseURI to print the base URL of the web app -->
<li><a href="{{.BaseURI}}about">About</a></li>

<!-- Use token to output the CSRF token in a form -->
<input type="hidden" name="token" value="{{.token}}">
~~~

It's also easy to add template-specific code before the closing </head> and </body> tags:

~~~ html
<!-- Code is added before the closing </head> tag -->
{{define "head"}}<meta name="robots" content="noindex">{{end}}

...

<!-- Code is added before the closing </body> tag -->
{{define "foot"}}{{JS "//www.google.com/recaptcha/api.js"}}{{end}}
~~~

## JavaScript

You can trigger a flash notification using JavaScript.

~~~ javascript
flashError("You must type in a username.");

flashSuccess("Record created!");

flashNotice("There seems to be a piece missing.");

flashWarning("Something does not seem right...");
~~~

## Controllers

The controller files all share the same package name. This cuts down on the
number of packages when you are mapping the routes. It also forces you to use
a good naming convention for each of the funcs so you know where each of the
funcs are located and what type of HTTP request they each are mapped to.

### These are a few things you can do with controllers.

Access a gorilla session:

~~~ go
// Get the current session
sess := session.Instance(r)
...
// Close the session after you are finished making changes
sess.Save(r, w)
~~~

Trigger 1 of 4 different types of flash messages on the next page load (no other code needed):

~~~ go
sess.AddFlash(view.Flash{"Sorry, no brute force :-)", view.FlashNotice})
sess.Save(r, w) // Ensure you save the session after making a change to it
~~~

Validate form fields are not empty:

~~~ go
// Ensure a user submitted all the required form fields
if validate, missingField := view.Validate(r, []string{"email", "password"}); !validate {
	sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
	sess.Save(r, w)
	LoginGET(w, r)
	return
}
~~~

Render a template:

~~~ go
// Create a new view
v := view.New(r)

// Set the template name
v.Name = "login/login"

// Assign a variable that is accessible in the form
v.Vars["token"] = csrfbanana.Token(w, r, sess)

// Refill any form fields from a POST operation
view.Repopulate([]string{"email"}, r.Form, v.Vars)

// Render the template
v.Render(w)
~~~

Return the flash messages during an Ajax request:

~~~ go
// Get session
sess := session.Instance(r)

// Set the flash message
sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
sess.Save(r, w)

// Display the flash messages as JSON
v := view.New(r)
v.SendFlashes(w)
~~~

Handle the database query:

~~~ go
// Get database result
result, err := model.UserByEmail(email)

// Determine if password is correct
if err == sql.ErrNoRows {
	// User does not exist
} else if err != nil {
	// Display error message
} else if passhash.MatchString(result.Password, password) {
	// Password matches!
} else {
	// Password does not match
}
~~~

Send an email:

~~~ go
// Email a user
err := email.SendEmail(email.ReadConfig().From, "This is the subject", "This is the body!")
if err != nil {
	log.Println(err)
	sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
	sess.Save(r, w)
	return
}
~~~

Validate a form if the Google reCAPTCHA is enabled in the config:

~~~ go
// Validate with Google reCAPTCHA
if !recaptcha.Verified(r) {
	sess.AddFlash(view.Flash{"reCAPTCHA invalid!", view.FlashError})
	sess.Save(r, w)
	RegisterGET(w, r)
	return
}
~~~

## Database

It's a good idea to abstract the database layer out so if you need to make
changes, you don't have to look through business logic to find the queries. All
the queries are stored in the models folder.

The user.go file at the root of the model directory is a compliation of all the queries for each database type.

Connect to the database (only once needed in your application):

~~~ go
// Connect to database
database.Connect(config.Database)
~~~

Read from the database:

~~~ go
result := User{}
err := database.DB.Get(&result, "SELECT id, password, status_id, nickname FROM user WHERE email = ? LIMIT 1", email)
return result, err
~~~

Write to the database:

~~~ go
_, err := database.DB.Exec("INSERT INTO user (nickname, last_name, email, password) VALUES (?,?,?,?)", nickname, last_name, email, password)
return err
~~~

## Middleware

There are a few pieces of middleware included. The package called csrfbanana
protects against Cross-Site Request Forgery attacks and prevents double submits.
The package httprouterwrapper provides helper functions to make funcs compatible
with httprouter. The package logrequest will log every request made against the
website to the console. The package pprofhandler enables pprof so it will work
with httprouter. In route.go, all the individual routes use alice to make
chaining very easy.

## Configuration

To make the web app a little more flexible, you can make changes to different
components in one place through the config.json file. If you want to add any
of your own settings, you can add them to config.json and update the structs
in gowebapp.go and the individual files so you can reference them in your code.

To enable HTTPS, set UseHTTPS to true, create a folder called tls in the root,
and then place the certificate and key files in that folder.
