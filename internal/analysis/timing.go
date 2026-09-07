package analysis

import "time"

type CommitData struct {
	Timestamp time.Time
}

type TimeAnalysis struct {
	Morning   int
	Afternoon int
	Evening   int
	Night     int
	Total     int
}

func AnalyzeTiming(commits []CommitData) TimeAnalysis {
	var result TimeAnalysis
	for _, commit := range commits {
		hour := commit.Timestamp.Hour()
		switch {
		case hour >= 5 && hour < 12:
			result.Morning++
		case hour >= 12 && hour < 17:
			result.Afternoon++
		case hour >= 17 && hour < 22:
			result.Evening++
		default:
			result.Night++
		}
	}
	result.Total = len(commits)
	return result
}

func Percentage(count int, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(count) / float64(total) * 100
}

func DominantPeriod(result TimeAnalysis) string {
	if result.Total == 0 {
		return "unknown"
	}
	period := "morning"
	max := result.Morning
	if result.Afternoon > max {
		max = result.Afternoon
		period = "afternoon"
	}
	if result.Evening > max {
		max = result.Evening
		period = "evening"
	}
	if result.Night > max {
		period = "night"
	}
	return period
}
