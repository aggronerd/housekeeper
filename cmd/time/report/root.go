// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>.
// All rights reserved.

package report

import (
	"fmt"
	"github.com/andygrunwald/go-jira"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"log"
	"os"
	"sort"
	"time"
)

const searchMaxResults = 100

type tallyEntry struct {
	duration time.Duration
	summary  string
}

type workTally struct {
	durationMap map[string]*tallyEntry
	total       time.Duration
}

func newWorkTally() workTally {
	return workTally{
		durationMap: map[string]*tallyEntry{},
		total:       time.Duration(0),
	}
}

func (w *workTally) Append(issueID string, summary string, duration time.Duration) {
	val, ok := w.durationMap[issueID]
	if !ok {
		w.durationMap[issueID] = &tallyEntry{
			duration: duration,
			summary:  summary,
		}
	} else {
		val.duration += duration
	}
	w.total += duration
}

// Print outputs the details of the workTally as a table
func (w *workTally) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Issue ID", "Summary", "Time Spent"})
	table.SetFooter([]string{"", "Total", w.total.String()})
	table.SetBorder(false)
	for _, key := range w.sortedKeys() {
		entry := w.durationMap[key]
		table.Append([]string{
			key,
			truncateString(entry.summary, 64),
			entry.duration.String(),
		})
	}
	fmt.Println("")
	table.Render()
}

func (w *workTally) sortedKeys() []string {
	var keys []string
	for key := range w.durationMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func runReport(_ *cobra.Command, _ []string) {
	client := getJiraClient()
	user, _, err := client.User.GetSelf()
	fatalIfErr(err)

	options := jira.SearchOptions{MaxResults: searchMaxResults}
	jiraDate := dayStart.Format("2006-01-02")
	jqlQuery := fmt.Sprintf("worklogDate = '%s' and worklogAuthor = currentUser()", jiraDate)
	issues, result, err := client.Issue.Search(jqlQuery, &options)
	fatalIfErr(err)

	if result.Total > searchMaxResults {
		// TODO: implement pagination
		log.Printf("WARNING: There are %d matching tickets but this only returns %d, "+
			"results will be incomplete", result.Total, searchMaxResults)
	}

	tally := newWorkTally()
	for _, issue := range issues {
		workLogs, _, err := client.Issue.GetWorklogs(issue.ID)
		fatalIfErr(err)

		for _, workLog := range workLogs.Worklogs {
			workLogStarted := time.Time(*workLog.Started)
			if (workLogStarted.After(dayStart) || workLogStarted.Equal(dayStart)) &&
				workLogStarted.Before(dayEnd) &&
				workLog.Author.AccountID == user.AccountID {

				// This work log is relevant
				tally.Append(issue.Key,
					issue.Fields.Summary,
					time.Duration(workLog.TimeSpentSeconds)*time.Second)
			}
		}
	}

	tally.Print()
}

// Cmd is a command which gives nicely detailed summary of Jira tickets for the day
var Cmd = &cobra.Command{
	Use:   "report",
	Short: "Shows how you have logged time against tickets in Jira",
	Run:   runReport,
}

var dateParameter string
var dayStart time.Time
var dayEnd time.Time

func initParameters() {
	var date time.Time
	// TODO: remove hard-coded time-zone
	location, err := time.LoadLocation("Europe/London")
	fatalIfErr(err)
	if dateParameter == "" {
		date = time.Now()
	} else {
		date, err = time.ParseInLocation("20060102", dateParameter, location)

		if err != nil {
			log.Fatalf(
				"Could not parse date, must be in the form yyyyMMdd (eg. \"20121225\")")
		}
	}
	year, month, day := date.Date()
	dayStart = time.Date(year, month, day, 0, 0, 0, 0, location)
	dayEnd = dayStart.Add(24 * time.Hour)
}

func init() {
	cobra.OnInitialize(initParameters)

	Cmd.PersistentFlags().StringVar(
		&dateParameter,
		"date",
		"",
		"Date for the report in the format dd/mm/yyyy (default is today)")
}
