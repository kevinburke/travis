package remoteci

import "time"

func round(f float64) int {
	if f < 0 {
		return int(f - 0.5)
	}
	return int(f + 0.5)
}

// GetEffectiveCost returns the cost in cents to pay an average San
// Francisco-based engineer to wait for the amount of time specified by d.
func GetEffectiveCost(d time.Duration) int {
	// https://www.glassdoor.com/Salaries/san-francisco-software-engineer-salary-SRCH_IL.0,13_IM759_KO14,31.htm
	yearlySalaryCents := float64(110554 * 100)
	// Estimate fully loaded costs add 40%.
	fullyLoadedSalary := yearlySalaryCents * 1.4

	workingDays := float64(49 * 5)
	hoursInWorkday := float64(8)
	salaryPerHour := fullyLoadedSalary * float64(d) / (workingDays * hoursInWorkday * float64(time.Hour))
	return round(salaryPerHour)
}
