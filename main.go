package main

import (
    "github.com/miekg/dns"
    "log"
    "math/rand"
    "strings"
    "os"
    "fmt"
    "time"
    "flag"
    "regexp"
)

const (
    DEFAULT_PORT = 53
    DEFAULT_TIMEOUT = time.Minute
    DEFAULT_PROTO = "tcp"
)

var (
        LOG *log.Logger
        PROTO string
        TIMEOUT time.Duration
        TCP_REGEX *regexp.Regexp

        DEFAULT_SERVERS = []string {    "8.8.8.8",
                                        "8.8.4.4",
                                        "156.154.70.1",
                                        "156.154.71.1",
                                        "208.67.222.222",
                                        "208.67.220.220",
                                        "198.153.192.1",
                                        "198.153.194.1" }


                                    )

func init () {
    LOG = log.New (os.Stderr, "[DNS PROXY] ", log.LstdFlags)

    var proto, pattern string
    var timeout int

    flag.StringVar (&proto, "p", "", "Which proto use for lookup domain")
    flag.IntVar (&timeout, "t", 0, "How many seconds to timeout")
    flag.StringVar (&pattern, "r", "", "Regex pattern for match domain to use tcp proto")
    flag.Parse ()

    switch proto {
    case "tcp" :
        PROTO = "tcp"
    case "udp" :
        PROTO = "udp"
    default:
        PROTO = DEFAULT_PROTO
    }

    if timeout > 5 {
        TIMEOUT = time.Duration(timeout) * time.Second
    } else {
        TIMEOUT = DEFAULT_TIMEOUT
    }

    if len(pattern) > 0 {
        if re, err := regexp.Compile (pattern); err != nil {
            LOG.Fatalf ("Compiling pattern [%s] was %s", pattern, err)
        } else {
            TCP_REGEX = re
        }
    }


    LOG.Printf ("Timeout duration: %s", TIMEOUT)
    if TCP_REGEX != nil {
        LOG.Printf ("Compiling tcp regex pattern [%s]", TCP_REGEX)
    }
}

type Proxy struct {
    Servers []string
}

func ( p Proxy ) ServeDNS (w dns.ResponseWriter, r *dns.Msg) {
    c := dns.NewClient ()

    c.Net = PROTO
    if TCP_REGEX != nil {
        for _, q := range r.Question {
            if TCP_REGEX.MatchString (q.Name) {
                LOG.Printf ("Tcp proto regex match: %s", q.Name)
                c.Net = "tcp"
                break
            }
        }
    }

    c.ReadTimeout = TIMEOUT
    c.WriteTimeout = TIMEOUT

    if rs, err := c.Exchange (r, p.Server()); err == nil {
        w.Write (rs)
    } else {
        dns.Refused (w, r)
        LOG.Printf ("%s %s", w.RemoteAddr (), err)
    }
}

func ( p Proxy ) Server () string {
    sl := len(p.Servers)
    if sl > 0 {
        i := rand.Intn (sl)
        s := p.Servers[i]
        if strings.Index (s, ":") == -1 {
            s = fmt.Sprintf ("%s:%d", s, DEFAULT_PORT)
        }
        return s
    }
    return "8.8.8.8:53"
}

func main () {

    var servers []string
    args := flag.Args ()

    if len(args) > 0 {
        servers = append (servers, args...)
    } else {
        servers = DEFAULT_SERVERS
    }

    LOG.Printf ("Servers: %s", servers)

    proxyer := Proxy{servers}

    if err := dns.ListenAndServe ("127.0.0.1:53", "udp", proxyer); err != nil {
        LOG.Fatal (err)
    }
}
