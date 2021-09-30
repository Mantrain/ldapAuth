// Package LdapAuth a ldap authentication plugin.
package ldapAuth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/go-ldap/ldap/v3"
)

const (
	defaultRealm        = "traefik"
	authorizationHeader = "Authorization"
	contentType         = "Content-Type"
)

// Config the plugin configuration.
type Config struct {
	Enabled      bool   `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Debug        bool   `json:"debug,omitempty" yaml:"debug,omitempty"`
	Host         string `json:"host,omitempty" yaml:"host,omitempty"`
	Port         uint16 `json:"port,omitempty" yaml:"port,omitempty"`
	BaseDn       string `json:"baseDn,omitempty" yaml:"baseDn,omitempty"`
	UserUniqueId string `json:"useruniqueid,omitempty" yaml:"useruniqueid,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Enabled:      true,
		Debug:        false,
		Host:         "example.com",
		Port:         389,
		BaseDn:       "dc=example,dc=org",
		UserUniqueId: "uid", // Usually uid or sAMAccountname
	}
}

// LdapAuth Struct plugin.
type LdapAuth struct {
	next   http.Handler
	name   string
	config *Config
}

// New created a new LdapAuth plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	log.Println("Starting", name, "Middleware...")
	if config.Debug {
		log.Println("Enabled       =>", config.Enabled)
		log.Println("Host          =>", config.Host)
		log.Println("Port          =>", config.Port)
		log.Println("BaseDn        =>", config.BaseDn)
		log.Println("UserUniqueId  =>", config.UserUniqueId)
	}

	return &LdapAuth{
		name:   name,
		next:   next,
		config: config,
	}, nil
}

func (la *LdapAuth) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	if !la.config.Enabled {
		log.Printf("%s Disabled! Passing request...", la.name)
		la.next.ServeHTTP(rw, req)
		return
	}

	user, password, ok := req.BasicAuth()

	if !ok {
		// No valid 'Authentication: Basic xxxx' header found in request
		la.RequireAuth(rw, req)
		return
	}

	isValidUser := la.ldapCheckUser(user, password)

	if !isValidUser {
		log.Printf("Authentication failed")
		la.RequireAuth(rw, req)
		return
	} else {
		log.Printf("Authentication succeeded")
	}

	// Sanitize Some Headers Infos
	req.URL.User = url.User(user)
	req.Header["LDAP-User"] = []string{user}
	req.Header.Del("Authorization")

	la.next.ServeHTTP(rw, req)
}

func (la *LdapAuth) ldapCheckUser(user, password string) bool {
	log.Printf("user=%s password=%s", user, password)
	if user == "john" {
		return true
	}
	return false
}

func (la *LdapAuth) RequireAuth(w http.ResponseWriter, req *http.Request) {
	w.Header().Set(contentType, "text/plan")
	w.Header().Set("WWW-Authenticate", `Basic realm="`+defaultRealm+`"`)
	w.WriteHeader(401)
	w.Write([]byte(fmt.Sprintf("%d %s bbb\n", http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))))
}
