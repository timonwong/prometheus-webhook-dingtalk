package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"
	"unicode"
)

var cfg = struct {
	fs *flag.FlagSet

	listenAddress    string
	dingTalkProfiles dingTalkProfilesFlag
	requestTimeout   time.Duration
}{}

func init() {
	cfg.fs = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	cfg.fs.Usage = usage

	cfg.fs.StringVar(&cfg.listenAddress, "web.listen-address", ":8060",
		"Address to listen on for web interface.",
	)
	cfg.fs.Var(&cfg.dingTalkProfiles, "ding.profile",
		"Custom DingTalk profile (can specify multiple times, <profile>=<dingtalk-url>).",
	)
	cfg.fs.DurationVar(&cfg.requestTimeout, "ding.timeout", 5*time.Second,
		"Timeout for invoking DingTalk webhook.",
	)
}

func parse(args []string) error {
	err := cfg.fs.Parse(args)
	if err != nil || len(cfg.fs.Args()) != 0 {
		if err != flag.ErrHelp {
			fmt.Fprintf(os.Stderr, "Invalid command line arguments. Help: %s -h", os.Args[0])
		}
		if err == nil {
			err = fmt.Errorf("bon-flag argument on command line: %q", cfg.fs.Args()[0])
		}
		return err
	}

	return nil
}

var helpTmpl = strings.TrimSpace(`
usage: prometheus-webhook-dingtalk [<args>]
{{ range $cat, $flags := . }}{{ if ne $cat "." }} == {{ $cat | upper }} =={{ end }}
  {{ range $flags }}
   -{{ .Name }} {{ .DefValue | quote }}
      {{ .Usage | wrap 80 6 }}
  {{ end }}
{{ end }}
`)

func usage() {
	t := template.New("usage")
	t = t.Funcs(template.FuncMap{
		"wrap": func(width, indent int, s string) (ns string) {
			width = width - indent
			length := indent
			for _, w := range strings.SplitAfter(s, " ") {
				if length+len(w) > width {
					ns += "\n" + strings.Repeat(" ", indent)
					length = 0
				}
				ns += w
				length += len(w)
			}
			return strings.TrimSpace(ns)
		},
		"quote": func(s string) string {
			if len(s) == 0 || s == "false" || s == "true" || unicode.IsDigit(rune(s[0])) {
				return s
			}
			return fmt.Sprintf("%q", s)
		},
		"upper": strings.ToUpper,
	})
	t = template.Must(t.Parse(helpTmpl))

	groups := make(map[string][]*flag.Flag)

	// Bucket flags into groups based on the first of their dot-separated levels.
	cfg.fs.VisitAll(func(fl *flag.Flag) {
		parts := strings.SplitN(fl.Name, ".", 2)
		if len(parts) == 1 {
			groups["."] = append(groups["."], fl)
		} else {
			name := parts[0]
			groups[name] = append(groups[name], fl)
		}
	})
	for cat, fl := range groups {
		if len(fl) < 2 && cat != "." {
			groups["."] = append(groups["."], fl...)
			delete(groups, cat)
		}
	}

	if err := t.Execute(os.Stdout, groups); err != nil {
		panic(fmt.Errorf("error executing usage template: %s", err))
	}
}

type dingTalkProfilesFlag struct {
	profiles map[string]string
}

func (c *dingTalkProfilesFlag) Set(opt string) error {
	if c.profiles == nil {
		c.profiles = make(map[string]string)
	}

	parts := strings.SplitN(opt, "=", 2)
	profile, webhookURL := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])

	if profile == "" {
		return errors.New("profile part cannot be empty")
	}
	if webhookURL == "" {
		return errors.New("webhook-url part cannot be emtpy")
	}

	c.profiles[profile] = webhookURL
	return nil
}

func (c *dingTalkProfilesFlag) String() string {
	return ""
}
