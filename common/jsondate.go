package common

import (
	"fmt"
	"time"
)

const timeFormatYMDhms = "2006-01-02 15:04:05" // Time format used by json
type JsonDate struct {
	time.Time // Time types used by json
}

// JSONDATE deserialization
func (self *JsonDate) UnmarshalJSON(data []byte) (err error) {
	newTime, _ := time.ParseInLocation("\""+timeFormatYMDhms+"\"", string(data), time.Local)
	*&self.Time = newTime
	return
}

// JSONDATE serialization
func (self JsonDate) MarshalJSON() ([]byte, error) {
	timeStr := fmt.Sprintf("\"%s\"", self.Format(timeFormatYMDhms))
	return []byte(timeStr), nil
}

//Output string
func (self JsonDate) String() string {
	return self.Time.Format(timeFormatYMDhms)
}

// ConvertJsonDate Time转为JsonDate
func ConvertJsonDate(date time.Time) JsonDate {
	return JsonDate{Time: date}
}

// ConvertJsonDateArray []Time转为[]JsonDate
func ConvertJsonDateArray(date []time.Time) []JsonDate {
	var dayArr []JsonDate
	for _, day := range date {
		dayArr = append(dayArr, JsonDate{Time: day})
	}
	return dayArr
}
