package cmd

import (
	"fmt"
	"github.com/akira/go-puppetdb"
	"log"
	"math"
	"sort"
	"strings"
	"time"
)

var PUPPET_INTERVAL = 31.0

type Master struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSL      bool   `yaml:"ssl"`
	Key      string `yaml:"key"`
	Ca       string `yaml:"ca"`
	Insecure bool   `yaml:"insecure"`
	Cert     string `yaml:"cert"`
}

type logEntry struct {
	NewValue string   `json:"new_value"`
	Property string   `json:"property"`
	File     string   `json:"file"`
	Line     string   `json:"line"`
	Tags     []string `json:"tags"`
	Time     []string `json:"time"`
	Level    string   `json:"level"`
	Source   string   `json:"source"`
	Message  string   `json:"message"`
}

type MessageTime struct {
	Message string
	Time    time.Time
}

func checkInterval(dates []time.Time, recurring int, interval int) bool {
	c := dates[0]
	for index, d := range dates {
		diff := c.Sub(d)
		dayBetween1 := math.Abs(math.Round(diff.Hours() / 24))
		if int(dayBetween1) == interval {
			if recurring == 1 {
				return true
			}
			newArr := dates[index:]
			return checkInterval(newArr, recurring+1, interval)
		}
	}
	return false
}

func checkIntervalHourly(dates []time.Time, recurring int, interval int) bool {
	c := dates[0]
	for index, d := range dates {
		diff := c.Sub(d)
		dayBetween1 := math.Abs(math.Round(diff.Hours()))
		if int(dayBetween1) == interval {
			if recurring == 1 {
				return true
			}
			newArr := dates[index:]
			return checkInterval(newArr, recurring+1, interval)
		}
	}
	return false

}

func GetMessageTimesForNode(certname string, master Master, sWarns bool, sErrors bool) []MessageTime {
	entries := GetLogEntryForNode(certname, master)
	messageEntries := []MessageTime{}
	if entries != nil {
		for _, e := range *entries {
			if !checkIsNoise(e, sWarns, sErrors) {
				for _, d := range e.Time {
					d2 := getTime(d)
					if d2 != nil {
						str := fmt.Sprintf("%s certname: %s message: %s %s %s", d2.String(), certname, e.Level, e.Source, e.Message)
						m := MessageTime{
							Message: str,
							Time:    *d2,
						}
						messageEntries = append(messageEntries, m)
					}
				}
			}
		}
	}

	sort.Slice(messageEntries, func(i, j int) bool {
		return messageEntries[i].Time.Before(messageEntries[j].Time)
	})
	return messageEntries
}

func GetHistoryForNode(certname string, master Master, sWarns bool, sErrors bool) {
	messageEntries := GetMessageTimesForNode(certname, master, sWarns, sErrors)

	for _, change := range messageEntries {
		fmt.Println(change.Message)
	}

}

func GetHistoryForAll(master Master, sWarns bool, sErrors bool) {
	certnames := GetCertNames(master)
	messageEntries := []MessageTime{}
	for _, c := range certnames {
		messageEntries = append(messageEntries, GetMessageTimesForNode(c, master, sWarns, sErrors)...)

	}
	sort.Slice(messageEntries, func(i, j int) bool {
		return messageEntries[i].Time.Before(messageEntries[j].Time)
	})
	for _, change := range messageEntries {
		fmt.Println(change.Message)
	}

}

func checkIsNoise(e logEntry, sWarns bool, sErrors bool) bool {
	if e.Message == "Applied catalog in x seconds" {
		return true
	}
	if e.Level == "info" && e.Source == "Puppet" {
		return true

	}
	if !sWarns && e.Level == "warning" {
		return true
	}
	if !sErrors && e.Level == "err" {
		return true
	}

	return false
}

func GetLogEntryForNode(certname string, master Master) *[]logEntry {
	entries := &[]logEntry{}
	reports := GetReportsForCertname(certname, master)
	for _, r := range reports {
		for _, l := range r.Logs.Data {
			e := pupDbEntryToLogEntry(l)
			AppendToLogEntries(entries, e)
		}
	}
	return entries
}

func GetContiniousChangesForNode(certname string, master Master, sWarns bool, sErrors bool) {
	entries := GetLogEntryForNode(certname, master)
	if entries != nil {
		check := false
		for _, e := range *entries {
			if !checkIsNoise(e, sWarns, sErrors) {
				pattern := seePatern(e.Time)
				if pattern != "" {
					check = true
					str := fmt.Sprintf("certname: %s | message: %s %s %s | pattern: %s",
						certname, e.Level, e.Source, e.Message, pattern)
					fmt.Println(str)
				}
			}
		}
		if !check {
			str := fmt.Sprintf("certname: %s | not recurring changes found", certname)
			fmt.Println(str)
		}
	}

}

func checkLogEntryAge(logTime time.Time) bool {
	today := time.Now()
	diff := today.Sub(logTime)
	if diff.Minutes() < PUPPET_INTERVAL {
		return true
	}
	return false
}

func checkDatesAreNextLoop(logTime time.Time, logTime2 time.Time) bool {
	diff := logTime2.Sub(logTime)
	if math.Abs(diff.Minutes()) < PUPPET_INTERVAL {
		return true
	}
	return false
}

func seePatern(times []string) string {
	dates := []time.Time{}
	for _, d := range times {
		d2 := getTime(d)
		if d2 != nil {
			dates = append(dates, *d2)
		}
	}
	if checkLogEntryAge(dates[0]) {
		if len(dates) >= 2 {
			if checkDatesAreNextLoop(dates[1], dates[0]) && checkDatesAreNextLoop(dates[2], dates[1]) {
				return "continious"
			}
		}

	}

	if len(dates) > 1 {
		if checkIntervalHourly(dates, 0, 1) {
			return "hourly"
		}

		if checkInterval(dates, 0, 1) {
			return "daily"
		}
		if checkInterval(dates, 0, 7) {
			return "weekly"
		}
	}

	return ""
}

func getTime(str string) *time.Time {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return &t
}

func AppendToLogEntries(entries *[]logEntry, entry logEntry) {
	check := false
	for index, e := range *entries {
		if compareEntries(e, entry) {
			check = true
			(*entries)[index].Time = append((*entries)[index].Time, entry.Time[0])
		}
	}
	if !check {
		(*entries) = append((*entries), entry)
	}
}

func EntryInSlice(entries []logEntry, entry logEntry) bool {
	for _, e := range entries {
		check := compareEntries(e, entry)
		if check {
			return true
		}
	}
	return false
}

func pupDbEntryToLogEntry(entry puppetdb.PuppetReportMetricsLogEntry) logEntry {
	e := logEntry{
		NewValue: entry.NewValue,
		Property: entry.Property,
		File:     entry.File,
		Line:     entry.Line,
		Tags:     entry.Tags,
		Time: []string{
			entry.Time,
		},
		Level:   entry.Level,
		Source:  entry.Source,
		Message: entry.Message,
	}
	if e.Level == "notice" && strings.HasPrefix(e.Message, "Applied catalog in ") {
		e.Message = "Applied catalog in x seconds"
	}
	return e
}

func compareEntries(e1 logEntry, e2 logEntry) bool {
	if e1.NewValue == e2.NewValue && e1.Property == e2.Property && e1.File == e1.File &&
		e1.Line == e2.Line && e1.Source == e2.Source &&
		e1.Message == e2.Message && e1.Level == e2.Level {
		return true
	}
	return false
}

func GetCertNames(master Master) []string {
	var cl *puppetdb.Client
	if !master.SSL {
		cl = puppetdb.NewClient(master.Host, master.Port, false)

	} else {
		if master.Insecure {
			cl = puppetdb.NewClientSSLInsecure(master.Host, master.Port, false)

		} else {
			cl = puppetdb.NewClientSSL(master.Host, master.Port, master.Key, master.Cert, master.Ca, false)

		}
	}
	nodes, err := cl.Nodes()
	if err != nil {
		log.Println(err.Error())
	}
	certnames := []string{}
	for _, n := range nodes {
		certnames = append(certnames, n.Certname)
	}
	return certnames
}

func GetReportsForCertname(certname string, master Master) []puppetdb.ReportJSON {
	var cl *puppetdb.Client
	if !master.SSL {
		cl = puppetdb.NewClient(master.Host, master.Port, false)

	} else {
		if master.Insecure {
			cl = puppetdb.NewClientSSLInsecure(master.Host, master.Port, false)

		} else {
			cl = puppetdb.NewClientSSL(master.Host, master.Port, master.Key, master.Cert, master.Ca, false)

		}
	}
	q := fmt.Sprintf("[\"=\",\"certname\",\"%s\"]", certname)
	reports, err := cl.Reports(q, nil)
	if err != nil {
		log.Println(err.Error())
	}
	sort.Slice(reports, func(i, j int) bool {
		return reports[i].ReceiveTime > reports[j].ReceiveTime
	})
	return reports
}
