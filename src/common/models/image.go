package models

import (
	"fmt"
	"strings"
)

// Image holds each part (project, repo, tag) of an image name
type Image struct {
	Project string
	Repo    string
	Tag     string
}

// ParseImage parses an image name such as 'library/app:v1.0' to a structure with
// project, repo, and tag fields
func ParseImage(image string) (*Image, error) {
	repo := strings.SplitN(image, "/", 2)
	if len(repo) < 2 {
		return nil, fmt.Errorf("unable to parse image from string: %s", image)
	}
	i := strings.SplitN(repo[1], ":", 2)
	res := &Image{
		Project: repo[0],
		Repo:    i[0],
	}
	if len(i) == 2 {
		res.Tag = i[1]
	}
	return res, nil
}
