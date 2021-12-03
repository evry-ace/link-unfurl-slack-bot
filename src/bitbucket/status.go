package bitbucket

const (
	// StatusSuccess represents the success status
	StatusSuccess = "SUCCESSFUL"
	// StatusFailed represents the failed status
	StatusFailed = "FAILED"
	// StatusRunning represents the running status
	StatusInProgress = "INPROGRESS"
	// StatusUnknown represents the unknown status
	StatusUnknown = "UNKNOWN"
)

// Status represents a Bitbucket commit status
type Status struct {
	State       string `json:"state"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	URL         string `json:"url"`
	Description string `json:"description"`
	DateAdded   int    `json:"dateAdded"`
}

// StatusList is a list of Status
type StatusList struct {
	Size       int      `json:"size"`
	Limit      int      `json:"limit"`
	IsLastPage bool     `json:"isLastPage"`
	Start      int      `json:"start"`
	Values     []Status `json:"values"`
}

// State returns the aggregated state of all statuses in the list
func (s StatusList) State() string {
	isRunning := false
	isSuccess := false
	isFailed := false

	for _, status := range s.Values {
		switch status.State {
		case StatusInProgress:
			isRunning = true
		case StatusSuccess:
			isSuccess = true
		case StatusFailed:
			isFailed = true
		}
	}

	// If there is a failed or running status, return that. Otherwise return the
	// success status if it exists
	if isFailed {
		return StatusFailed
	} else if isRunning {
		return StatusInProgress
	} else if isSuccess {
		return StatusSuccess
	}

	return StatusUnknown
}
