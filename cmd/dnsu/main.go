package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dnsupdater/dnsupdater/providers/cpanel"
	log "github.com/sirupsen/logrus"
)

const Version = "0.7.2"

func help() {
	fmt.Println(`

DNS Updater provides "present" and "cleanup" functions that LEGO expects from the external program provider.

DNS Updater: https://dnsupdater.github.io

Usage:
  # for help
  dnsu help

  # for cPanel provider
  dnsu cpanel --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] info <domain name for A record>
  dnsu cpanel --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] present <domain name> for TXT record> <auth-key>'
  dnsu cpanel --url <cPanel URL> --user <cPanel User> --token <cPanel Token> [--logoutput <log file name>] cleanup <domain name> for TXT record> <auth-key>'

Supported Environment Variables:
    DNSU_LOG-OUTPUT for log file
	DNSU_CPANEL-URL for cPanel URL
	DNSU_CPANEL-USER for cPanel User
	DNSU_CPANEL-TOKEN for cPanel Token

Example:
  # for verify cPanel access   
  dnsu cpanel --url "https://cpanel-hostname:2083" --user cpaneluser --token "RMYKKBIT5TQ1ITFU58VZBQB5TDEYQZN4" info '_acme-challenge.my.example.org.'
 
`)
}

func UpdateDomanSuffix(domain string) string {
	if strings.HasSuffix(domain, ".") {
		return domain
	}
	return domain + "."
}

func main() {

	var (
		logOutput, cpURL, cpUser, cpToken     string
		eLogOutput, ecpURL, ecpUser, ecpToken string
		mdl                                   string
	)

	fmt.Printf("\nDNS Updater v%s\n\n", Version)

	log.SetOutput(os.Stdout)

	eLogOutput = os.Getenv("DNSU_LOG-OUTPUT")
	ecpURL = os.Getenv("DNSU_CPANEL-URL")
	ecpUser = os.Getenv("DNSU_CPANEL-USER")
	ecpToken = os.Getenv("DNSU_CPANEL-TOKEN")

	cpanelCmd := flag.NewFlagSet("cpanel", flag.ContinueOnError)

	cpanelCmd.StringVar(&logOutput, "logoutput", "", "log otput file")
	cpanelCmd.StringVar(&cpURL, "url", "", "cPanel url")
	cpanelCmd.StringVar(&cpUser, "user", "", "cPanel user")
	cpanelCmd.StringVar(&cpToken, "token", "", "cPanel token")

	if logOutput == "" {
		logOutput = eLogOutput
	}

	if cpURL == "" {
		cpURL = ecpURL
	}

	if cpUser == "" {
		cpUser = ecpUser
	}

	if cpToken == "" {
		cpToken = ecpToken
	}

	logFile := "dnsu.log"

	switch logOutput {
	case "":
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
	case "stdout":
		log.SetOutput(os.Stdout)

	case "stderr":
		log.SetOutput(os.Stderr)

	default:
		logFile = logOutput
		file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			fmt.Printf("Error: %s", err.Error())
			log.Fatal(err)
		}
		defer file.Close()
		log.SetOutput(file)
	}

	log.SetLevel(log.InfoLevel)
	// log.SetFormatter(&log.JSONFormatter{})

	if len(os.Args) < 2 {
		fmt.Println("Error: expected 'help' or 'cpanel' commands")
		log.Fatal("expected 'help' or 'cpanel' commands")
	}

	cpanelCmd.Parse(os.Args[2:])
	mdl = os.Args[1]

	switch {
	case mdl == "help":
		help()
	case mdl == "cpanel":
		log.Println("Geldi----cpanel-")
		if (cpURL == "") || (cpUser == "") || (cpToken == "") {
			fmt.Println("Error: cpanel command must have the flags '--url', '--user' and '--token' or the corresponding environment variables.")
			log.Fatal("cpanel command must have the flags '--url', '--user' and '--token' or the corresponding environment variables.")
		}

		subCmd := cpanelCmd.Args()
		if len(subCmd) == 0 {
			fmt.Println("Error: expected 'info', 'present' or 'cleanup' subcommands.")
			log.Fatal("expected 'info', 'present' or 'cleanup' subcommands.")
		}

		switch subCmd[0] {

		case "info":
			if cpanelCmd.NArg() != 2 {
				fmt.Println("Error: info subcommand must have <domain> parameter.")
				log.Fatal("info subcommand must have <domain> parameter.")

			}
			args := cpanelCmd.Args()

			err := cpanel.DomainInfo(UpdateDomanSuffix(args[1]), cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				log.Fatal(err.Error())
			}

		case "present":
			if cpanelCmd.NArg() != 3 {
				fmt.Println("Error: present subcommand must have two parameters. (<domain> <auth-key>)")
				log.Fatal("present subcommand must have two parameters. (<domain> <auth-key>)")
			}
			args := cpanelCmd.Args()

			err := cpanel.Present(UpdateDomanSuffix(args[1]), args[2], cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				log.Fatal(err.Error())
			}

		case "cleanup":
			if cpanelCmd.NArg() < 2 {
				fmt.Println("Error: cleanup subcommand must have <domain> parameter.")
				log.Fatal("cleanup subcommand must have <domain> parameter.")
			}
			args := cpanelCmd.Args()
			err := cpanel.Cleanup(UpdateDomanSuffix(args[1]), cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Printf("Error: %s", err.Error())
				log.Fatal(err.Error())
			}
		default:
			fmt.Println("Error: expected 'info', 'present' or 'cleanup' subcommands.")
			log.Fatal("expected 'info', 'present' or 'cleanup' subcommands.")
		}

	default:
		fmt.Println("Error: expected 'help' or 'cpanel' commands.")
		log.Fatal("expected 'help' or 'cpanel' commands.")
	}
}
