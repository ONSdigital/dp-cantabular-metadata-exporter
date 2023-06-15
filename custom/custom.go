package custom

import "github.com/ONSdigital/dp-api-clients-go/v2/dataset"

// GenerateCustomTitle generates a title for custom datasets based on the requested dimensions
func GenerateCustomTitle(dims []dataset.VersionDimension) string {
	var title string
	lastIndex := len(dims) - 1
	for i := range dims {
		if i > 0 && i == lastIndex {
			title += " and "
		} else if i > 0 {
			title += ", "
		}
		title += dims[i].Label
	}
	return title
}
