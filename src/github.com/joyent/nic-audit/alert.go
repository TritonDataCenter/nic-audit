/*
 * Copyright (c) 2017, Joyent, Inc. All rights reserved.
 *
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/.
 */
package main

import (
	"container/list"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

import (
	"github.com/jordan-wright/email"
	"github.com/joyent/triton-go/compute"
)

var alertLogger *log.Logger

func init() {
	alertLogger = log.New(os.Stdout, "", log.LstdFlags)
}

// Alert describes the properties of an offending network match such that
// it can be acted upon by a system administrator.
type Alert struct {
	Instance     compute.Instance
	Account      Account
	NicGroupName string
	NicGroupIds  []string
}

// processAlerts iterates an aggregated list of alerts containing
// offending network details and triggers an alert action for each
// alert.
func processAlerts(alerts list.List, client compute.ComputeClient,
	config Configuration) {

	aggregate := "Instances with offending network combinations have been found.\n"

	for e := alerts.Front(); e != nil; e = e.Next() {
		var alert Alert = e.Value.(Alert)
		account := alert.Account

		// Always write the alert being processed right away so that we
		// know the current item being processed
		logAlert(alert)

		aggregate += "\n======================================================\n"
		aggregate += " Offending network match detected\n"
		aggregate += "======================================================\n"
		aggregate += fmt.Sprintf("  Account: %v\n", alert.Account.AccountName)
		aggregate += fmt.Sprintf("  Account Description: %v\n", alert.Account.Description)
		aggregate += fmt.Sprintf("  Triton URL: %v\n", alert.Account.TritonUrl)
		aggregate += fmt.Sprintf("  Match Group: %v\n", alert.NicGroupName)
		aggregate += fmt.Sprintf("  Networks Matched: %v\n", alert.NicGroupIds)
		aggregate += fmt.Sprintf("  Instance ID: %v\n", alert.Instance.ID)
		aggregate += fmt.Sprintf("  Instance Name: %v\n", alert.Instance.Name)
		aggregate += fmt.Sprintf("  Instance IPs: %v\n", alert.Instance.IPs)
		aggregate += fmt.Sprintf("  Instance Firewall Enabled: %v\n",
			alert.Instance.FirewallEnabled)
		aggregate += fmt.Sprintf("  Instance Networks: %v\n", alert.Instance.Networks)

		if len(account.NetworksToRemove) > 0 {
			networksRemoved, removeErr := removeNICsBasedOnNetworks(
				account.NetworksToRemove, alert.Instance, client,
				config.PrivateNetworkBlocks)
			if removeErr != nil {
				log.Printf("Error removing network for instance [%v]: %v\n",
					alert.Instance.ID, removeErr)
			} else {
				aggregate += fmt.Sprintf("  Instance Networks Removed: %v\n",
					networksRemoved)
			}
		}

		if len(config.EmailAlerts.SmtpServer) > 0 {
			emailAlerts(config.EmailAlerts, aggregate)
		} else {
			log.Println("Alert email is disable because no SMTP server " +
				"has been set")
		}
	}
}

// logAlert outputs a message to STDOUT reporting the specified alert.
func logAlert(alert Alert) {
	alertLogger.Printf("%v: %v (%v) %v\n", alert.NicGroupName,
		alert.Instance.Name, alert.Instance.ID, alert.Instance.IPs)
}

// emailAlerts emails the contents of the specified body text to the
// specified email recipients. Typically the body would contain an aggregation
// of all of the email alters triggered per account.
func emailAlerts(emailAlertConfig EmailAlerts, body string) {
	mail := email.NewEmail()
	mail.From = fmt.Sprintf("%v <%v>", emailAlertConfig.FromName,
		emailAlertConfig.From)

	mail.To = emailAlertConfig.To
	mail.Bcc = emailAlertConfig.BCC
	mail.Cc = emailAlertConfig.CC
	mail.Subject = emailAlertConfig.Subject
	mail.Headers.Add("Content-Type", "text/html; charset=UTF-8")
	mail.Headers.Add("MIME-version", "1.0")
	mail.Headers.Add("Content-Transfer-Encoding", "quoted-printable")

	mail.HTML = []byte(fmt.Sprintf("<html><body><pre>%v</pre></body></html>", body))
	mail.Text = []byte(body)

	mail.Send(emailAlertConfig.SmtpServer, smtp.PlainAuth("", "", "",
		emailAlertConfig.SmtpServer))
}
