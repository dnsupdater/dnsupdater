package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dnsupdater/dnsupdater/providers/cpanel"
)

const Version = "0.4"

func help() {
	fmt.Println(`

DNS Updater provides "present" and "cleanup" functions that LEGO expects from the external program provider.

DNS Updater: https://dnsupdater.github.io

Usage:
  # for help
  dnsu help

  # for cPanel provider
  dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> info <domain name> for A record>
  dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> present <domain name> for TXT record> <auth-key>'
  dnsu cpnael --url <cPanel URL> --user <cPanel User> --token <cPanel Token> cleanup <domain name> for TXT record> <auth-key>'

Support Environment Variables:
	DNSU_CPANEL-URL for cPanel URL
	DNSU_CPANEL-USER for cPanel User
	DNSU_CPANEL-TOKEN for cPanel Token

Example:
  # for verify cPanel access   
  dnsu cpnael --url "https://cpanel-hostname:2083" --user cpaneluser --token "RMYKKBIT5TQ1ITFU58VZBQB5TDEYQZN4" info '_acme-challenge.my.example.org.'
 
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
		cpURL, cpUser, cpToken    string
		ecpURL, ecpUser, ecpToken string
	)

	fmt.Printf("\nDNS Updater v%s\n\n", Version)

	ecpURL = os.Getenv("DNSU_CPANEL-URL")
	ecpUser = os.Getenv("DNSU_CPANEL-USER")
	ecpToken = os.Getenv("DNSU_CPANEL-TOKEN")

	cpanelCmd := flag.NewFlagSet("cpanel", flag.ContinueOnError)

	cpanelCmd.StringVar(&cpURL, "url", "", "cPanel url")
	cpanelCmd.StringVar(&cpUser, "user", "", "cPanel user")
	cpanelCmd.StringVar(&cpToken, "token", "", "cPanel token")

	if cpURL == "" {
		cpURL = ecpURL
	}

	if cpUser == "" {
		cpUser = ecpUser
	}

	if cpToken == "" {
		cpToken = ecpToken
	}

	if len(os.Args) < 2 {
		fmt.Println("expected 'help' or 'cpanel' commands")
		os.Exit(1)
	}

	cpanelCmd.Parse(os.Args[2:])
	switch os.Args[1] {
	case "help":
		help()
	case "cpanel":
		if (cpURL == "") || (cpUser == "") || (cpToken == "") {
			fmt.Println("cpanel command must have the flags '--url', '--user' and '--token' or the corresponding environment variables.")
			os.Exit(1)
		}

		subCmd := cpanelCmd.Args()
		if len(subCmd) == 0 {
			fmt.Println("expected 'info', 'present' or 'cleanup' subcommands.")
			os.Exit(1)
		}
		switch subCmd[0] {

		case "info":
			if cpanelCmd.NArg() != 2 {
				fmt.Println("info subcommand must have <domain> parameter.")
				os.Exit(1)
			}
			args := cpanelCmd.Args()

			err := cpanel.DomainInfo(UpdateDomanSuffix(args[1]), cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}

		case "present":
			if cpanelCmd.NArg() != 3 {
				fmt.Println("present subcommand must have two parameters. (<domain> <auth-key>)")
				os.Exit(1)
			}
			args := cpanelCmd.Args()

			err := cpanel.Present(UpdateDomanSuffix(args[1]), args[2], cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}

		case "cleanup":
			if cpanelCmd.NArg() < 2 {
				fmt.Println("cleanup subcommand must have <domain> parameter.")
				os.Exit(1)
			}
			args := cpanelCmd.Args()

			err := cpanel.Cleanup(UpdateDomanSuffix(args[1]), cpURL, cpUser, cpToken)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(2)
			}

		default:
			fmt.Println("expected 'info', 'present' or 'cleanup' subcommands.")
			os.Exit(1)

		}

	default:
		fmt.Println("expected 'help' or 'cpanel' commands.")
		os.Exit(1)
	}
}
