package commandline

import (
	"fmt" // This example uses SQLite; adapt for your database
	"strings"
	"vertigo/pkg/dataitems"
)

func ParseShoeParams(args []string) (dataitems.Shoe, error) {
	shoe := dataitems.Shoe{}

	for _, arg := range args {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) != 2 {
			return shoe, fmt.Errorf("invalid parameter format, please use format name= brand= silhouette= (optional shoe_url= tags=): %s", arg)
		}

		key := parts[0]
		value := strings.Trim(parts[1], "\"")

		switch key {
		case "name":
			shoe.Name = value
		case "brand":
			shoe.Brand = value
		case "silhouette":
			shoe.Silhouette = value
		case "image_url":
			shoe.ImageUrl = value
		case "tags":
			shoe.Tags = value
		default:
			return shoe, fmt.Errorf("unknown parameter please use format name= brand= silhouette= (optional shoe_url= tags=): %s", key)
		}
	}

	return shoe, nil
}
