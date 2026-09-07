package analysis

import (
	"sort"
	"time"
)

type ActivityAnalysis struct {
	TotalCommits       int
	AveragePerDay      float64
	MostActiveDay      string
	MostActiveWeekday  string
	LongestInactiveGap time.Duration
	MaxCommitsInDay    int
}

func AnalyzeActivity(commits []CommitData) ActivityAnalysis {
	result := ActivityAnalysis{TotalCommits: len(commits)}
	if len(commits) == 0 {
		return result
	}
	days := make(map[string]int)
	weekdays := make(map[time.Weekday]int)
	for _, c := range commits {
		day := c.Timestamp.Format("2006-01-02")
		days[day]++
		weekdays[c.Timestamp.Weekday()]++
	}
	var dates []time.Time
	for day := range days {
		date, err := time.Parse("2006-01-02", day)
		if err == nil {
			dates = append(dates, date)
		}
	}
	sort.Slice(dates, func(i, j int) bool { return dates[i].Before(dates[j]) })
	if len(dates) >= 2 {
		first := dates[0]
		last := dates[len(dates)-1]
		totalDays := last.Sub(first).Hours()/24 + 1
		result.AveragePerDay = float64(len(commits)) / totalDays
	}
	mostActiveDay := ""
	mostActiveCount := 0
	for day, count := range days {
		if count > mostActiveCount {
			mostActiveDay = day
			mostActiveCount = count
		}
	}
	result.MostActiveDay = mostActiveDay
	var mostActiveWeekday time.Weekday
	mostActiveWeekdayCount := 0
	for weekday, count := range weekdays {
		if count > mostActiveWeekdayCount {
			mostActiveWeekday = weekday
			mostActiveWeekdayCount = count
		}
	}
	result.MostActiveWeekday = mostActiveWeekday.String()
	maxCommits := 0
	for _, count := range days {
		if count > maxCommits {
			maxCommits = count
		}
	}
	result.MaxCommitsInDay = maxCommits
	for i := 1; i < len(dates); i++ {
		gap := dates[i].Sub(dates[i-1])
		if gap > result.LongestInactiveGap {
			result.LongestInactiveGap = gap
		}
	}
	return result
}
