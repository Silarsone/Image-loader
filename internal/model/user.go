package model

import "io"

type User struct {
	ID          int
	Name        string
	Login       string
	Password    string
	Description string
	ImageUrls   []string
}

type Image struct {
	ID        int
	UserID    int
	Name      string
	Extension string
	Data      io.Reader
}
