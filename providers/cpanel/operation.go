package cpanel

import (
	"fmt"

	"github.com/pkg/errors"
)

func Present(domain, keyAuth string, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient  *API
		txtRecordData *TXTRecordData
		err           error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	txtRecordData, err = cpanelClient.FindTXTRecord(domain, "TXT")
	if err != nil {
		return err
	}

	if txtRecordData.Domain == "" { // Record not found.
		addResp, err := cpanelClient.AddTXTRecords(domain, keyAuth, 0)
		if err != nil {
			return err
		}

		fmt.Printf("TXT record will be added. Newserial: %s ", addResp.Newserial)

		return nil
	}

	if (txtRecordData.Record == keyAuth) && (txtRecordData.Domain == domain) {
		fmt.Printf("TXT record with same keyAuth already exists. Not require any update.")
		return nil
	}

	updResp, err := cpanelClient.UpdateTXTRecord(*txtRecordData, keyAuth)
	if err != nil {
		fmt.Printf("err: %s\n", err.Error())
	} else {
		fmt.Printf("TXT record will be updated. Newserial: %s ", updResp.Newserial)
	}

	return nil
}

func Cleanup(domain string, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient  *API
		txtRecordData *TXTRecordData
		err           error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	txtRecordData, err = cpanelClient.FindTXTRecord(domain, "TXT")
	if err != nil {
		return err
	}

	if txtRecordData.Domain == "" {
		fmt.Printf("Record not found for %s", domain)
	}

	if txtRecordData.Domain == domain { // Record  found.
		remResp, err := cpanelClient.RemoveTXTRecord(*txtRecordData)
		if err != nil {
			return err
		}
		fmt.Printf("Record will be removed. Newserial: %s ", remResp.Newserial)
	}

	return nil
}

func DomainInfo(domain, cpURL, cpUser, cpToken string) error {
	var (
		cpanelClient  *API
		txtRecordData *TXTRecordData
		err           error
	)

	if cpanelClient, err = NewClient(
		cpURL, cpUser, cpToken); err != nil {
		return errors.Wrapf(err, "Couldn't create remote cPanel API client")
	}

	txtRecordData, err = cpanelClient.FindTXTRecord(domain, "A")
	if err != nil {
		return err
	}

	if txtRecordData.Domain == "" { // Record not found.
		fmt.Printf("Record not found. This subcommand is only for 'A Type' records.")
		return nil
	}

	fmt.Printf("Found 'A type' record:\n name: %s, address: %s\n", txtRecordData.Domain, txtRecordData.Record)

	return nil
}
