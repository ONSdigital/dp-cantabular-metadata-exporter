package custom

import "github.com/ONSdigital/dp-api-clients-go/v2/dataset"

func GenerateCustomTitle(dims []dataset.VersionDimension) string {
	var title string
	l := len(dims)
	for i, d := range dims {
		if i == (l - 1) {
			title += " and "
		} else if i > 0 {
			title += ", "
		}
		title += d.Label
	}
	return title
}
