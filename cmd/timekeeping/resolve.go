// Copyright 2019 Gregory Doran <greg@gregorydoran.co.uk>. 
// All rights reserved.

package timekeeping

import (
	"github.com/aggronerd/housekeeper/pkg"
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
	summary string
}

type workTally struct {
	durationMap map[string]*tallyEntry
	total time.Duration
}

func newWorkTally() workTally {
	return workTally{
		durationMap: map[string]*tallyEntry{},
		total: time.Duration(0),
	}
}

func (w *workTally) Append(issueId string, summary string, duration time.Duration) {
	val, ok := w.durationMap[issueId]
	if !ok {
		w.durationMap[issueId] = &tallyEntry{
			duration: duration,
			summary: summary,
		}
 	} else {
 		val.duration += duration
	}
	w.total += duration
}

func (w *workTally) Print() {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Issue ID", "Summary", "Time Spent"})
	table.SetFooter([]string{"", "Total", w.total.String()})
	table.SetBorder(false)
	for _, key := range w.sortedKeys()  {
		entry := w.durationMap[key]
		table.Append([]string{
			key,
			pkg.TruncateString(entry.summary, 64),
			entry.duration.String(),
		})
	}
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

var resolveCmd = &cobra.Command{
	Use: "view",
	Short: "Shows how you have logged work in Jira",
	Run: func(cmd *cobra.Command, args []string) {
		year, month, day := time.Now().Date()
		location, err := time.LoadLocation("Europe/London")

		if err != nil {
			log.Fatal(err)
		}

		dayStart := time.Date(year, month, day, 0, 0, 0, 0, location)
		dayEnd := dayStart.Add(24 * time.Hour)

		client := pkg.GetClient()
		user, _, err := client.User.GetSelf()
		pkg.FatalIfErr(err)

		options := jira.SearchOptions{MaxResults: searchMaxResults}
		issues, result, err := client.Issue.Search(
			"worklogDate = now() and worklogAuthor = currentUser()", &options)
		pkg.FatalIfErr(err)

		if result.Total > searchMaxResults {
			// TODO: implement pagination
			log.Printf("WARNING: There are %d matching tickets but this only returns %d, " +
				              "results will be incomplete", result.Total, searchMaxResults)
		}

		tally := newWorkTally()
		for _, issue := range issues {
			workLogs, _, err := client.Issue.GetWorklogs(issue.ID)
			pkg.FatalIfErr(err)

			for _, workLog := range workLogs.Worklogs {
				workLogStarted := time.Time(*workLog.Started)
				if (workLogStarted.After(dayStart) || workLogStarted.Equal(dayStart)) &&
					workLogStarted.Before(dayEnd) &&
					workLog.Author.AccountID == user.AccountID {

					// This worklog is relevant
					tally.Append(issue.Key,
						issue.Fields.Summary,
						time.Duration(workLog.TimeSpentSeconds) * time.Second)
				}
			}
		}


		tally.Print()
	},
}
