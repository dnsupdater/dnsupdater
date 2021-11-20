package cpanel

import (
	"fmt"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func Present(domain, keyAuth string, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient  *API
		txtRecordData *DNSRecordData
		err           error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	txtRecordData, err = cpanelClient.FindDNSRecord(domain, "TXT")
	if err != nil {
		return err
	}
	if txtRecordData.Domain == "" { // Record not found.
		addResp, err := cpanelClient.AddTXTRecords(domain, keyAuth, 0)
		if err != nil {
			return err
		}
		fmt.Printf("TXT record for %s with keyAuth=%s will be added. Newserial: %s", domain, keyAuth, addResp.Newserial)
		log.Infof("TXT record for %s  with keyAuth=%s will be added. Newserial: %s", domain, keyAuth, addResp.Newserial)

		return nil
	}
	if (txtRecordData.Record == keyAuth) && (txtRecordData.Domain == domain) {
		fmt.Printf("TXT record with same keyAuth=%s already exists for %s Not require any update.", keyAuth, domain)
		log.Infof("TXT record with same keyAuth= %s already exists for %s Not require any update.", keyAuth, domain)
		return nil
	}
	updResp, err := cpanelClient.UpdateTXTRecord(*txtRecordData, keyAuth)
	if err != nil {
		return err
	}
	fmt.Printf("TXT record for %s will be updated with keyAuth=%s. Newserial: %s", domain, keyAuth, updResp.Newserial)
	log.Infof("TXT record for %s will be updated with keyAuth=%s. Newserial: %s", domain, keyAuth, updResp.Newserial)

	return nil
}

func Cleanup(domain string, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient  *API
		txtRecordData *DNSRecordData
		err           error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	txtRecordData, err = cpanelClient.FindDNSRecord(domain, "TXT")
	if err != nil {
		return err
	}

	if txtRecordData.Domain == "" {
		fmt.Printf("TXT Record not found for %s", domain)
		log.Infof("TXT Record not found for %s", domain)
	}

	if txtRecordData.Domain == domain { // Record  found.
		remResp, err := cpanelClient.RemoveTXTRecord(*txtRecordData)
		if err != nil {
			return err
		}
		fmt.Printf("%s TXT Record will be removed. Newserial: %s", domain, remResp.Newserial)
		log.Infof("%s TXT Record will be removed. Newserial: %s", domain, remResp.Newserial)
	}

	return nil
}

func DomainInfo(domain, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient *API
		RecordData   *DNSRecordData
		err          error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	RecordData, err = cpanelClient.FindDNSRecord(domain, "A")
	if err != nil {

		return err
	}

	if RecordData.Domain == "" { // Record not found.
		fmt.Printf("Error: Record not found. This subcommand is only for 'A Type' records.")
		log.Error("Record not found. This subcommand is only for 'A Type' records.")
		return nil
	}
	fmt.Printf("Found 'A type' record: %s / %s", RecordData.Domain, RecordData.Record)
	log.Infof("Found 'A type' record: %s / %s", RecordData.Domain, RecordData.Record)
	return nil
}
