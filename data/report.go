package data

import "time"

type Report struct {
	Server string
	ImageName string
	Registry string
	Summary string
	ImageCreationDate string
	Os string
	OsVersion string
	Created time.Time

	ImageAllowed bool
}
