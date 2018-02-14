package def

import "time"

// Configuration keeps all settings
type Configuration struct {
	PortNumber string
	GEFAddress string
	TimeOut    int64
	Apps       map[string]string
}

// SelectedJob is the format of JSON returned by GEF when inspecting a job
type SelectedJob struct {
	Job SingleJob
}

// JobVolume points to volumes bound to a particular job
type JobVolume struct {
	VolumeID string
	Name     string
}

// SingleJob is a single job object
type SingleJob struct {
	ID           string
	ConnectionID int
	ServiceID    string
	Created      time.Time
	Duration     int64
	State        *JobState
	InputVolume  []JobVolume
	OutputVolume []JobVolume
	Tasks        []Task
}

// JobState keeps information about a job state
type JobState struct {
	Status string
	Error  string
	Code   int
}

// Task contains tasks related to a specific job (used to serialize JSON)
type Task struct {
	ID             string
	Name           string
	ContainerID    string
	SwarmServiceID string
	Error          string
	ExitCode       int
	ConsoleOutput  string
}

// VolumeInspection is the format of JSON returned by GEF when inspecting a volume
type VolumeInspection struct {
	VolumeContent []VolumeItem
}

// VolumeItem is a JSON format used to keep information about the content of a volume
type VolumeItem struct {
	Name       string       `json:"name"`
	Size       int64        `json:"size"`
	Modified   time.Time    `json:"modified"`
	IsFolder   bool         `json:"isFolder"`
	Path       string       `json:"path"`
	FolderTree []VolumeItem `json:"folderTree"`
}

type Info struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}
