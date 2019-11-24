package datafaker

import (
	"errors"
	"strconv"
)

type FakedData struct {
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Number int    `json:"number"`
}

const (
	descBase  string = "Hello this is desc"
	titleBase string = "Hello this is title"
)

func fakeDesc(i int) string {
	return descBase + " " + strconv.Itoa(i)
}

func fakeTitle(i int) string {
	return titleBase + " " + strconv.Itoa(i)
}

func FakeDatas(i int) ([]FakedData, error) {
	if i <= 0 {
		return []FakedData{}, errors.New("Paste in integer larger than 0 please")
	}
	result := []FakedData{}
	for j := 0; j < i; j++ {
		result = append(result, FakedData{
			Title:  fakeTitle(j),
			Desc:   fakeDesc(j),
			Number: j,
		})
	}
	return result, nil
}
