package common

import (
	"fmt"

	"github.com/jinzhu/now"

	//"goenv/tool"

	"strconv"
	"time"
)

// WorkDayOfMonth 查询一个月当中每周工作日具体日期
func WorkDayOfMonth(currentTime string) ([]interface{}, int) {
	currentDate, _ := time.Parse("2006-01", currentTime)
	year := currentDate.Year()
	month, _ := strconv.Atoi(currentDate.Month().String())
	count := DaysOfMonth(year, month)
	bmonth := now.New(currentDate).BeginningOfMonth()
	var monthArr []string
	var tempArr []string
	var weekArr []interface{}
	for i := 0; i < count; i++ {
		dd := bmonth.AddDate(0, 0, i)
		//除去星期天
		if dd.Weekday().String() != "Saturday" && dd.Weekday().String() != "Sunday" {
			monthArr = append(monthArr, dd.Format("2006-01-02"))
			if len(tempArr) > 0 {
				aa, _ := time.Parse("2006-01-02", tempArr[len(tempArr)-1])
				dd, _ = time.Parse("2006-01-02", dd.Format("2006-01-02"))
				if dd.Sub(aa).Hours()/24 > 1 {
					weekArr = append(weekArr, tempArr)
					tempArr = nil
				}
			}
			tempArr = append(tempArr, dd.Format("2006-01-02"))
		}
	}
	return weekArr, count
}

// DaysOfMonth 获取月天数
func DaysOfMonth(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30

		} else {
			days = 31
		}
	} else {
		if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
			days = 29
		} else {
			days = 28
		}
	}
	return days
}

// func main() {
// 	start, _ := time.Parse("2006-01-02", "2021-10-01")
// 	end, _ := time.Parse("2006-01-02", "2021-10-31")
// 	total, days := CalcWorkHour(start, end)
// 	fmt.Println(fmt.Sprintf("总计：%.2f个小时，%d天", total, days))
// }

//https://www.cnblogs.com/Nihility/p/14695646.html

func CalcWorkHour(begin, end time.Time) []string {
	var dayArr []string
	var currentTime = begin

	for {
		//https://cloud.tencent.com/developer/article/1467743
		if currentTime.After(end) {
			break
		}
		// 周六周日
		if currentTime.Weekday() == time.Sunday || currentTime.Weekday() == time.Saturday {
			// nothing
			dayArr = append(dayArr, "")
		} else {
			dayArr = append(dayArr, currentTime.Format("2006-01-02"))
		}
		//currentTime = currentTime.Add(24 * time.Hour)
		currentTime = currentTime.AddDate(0, 0, 1)
	}

	return dayArr
}

// GetBetweenDates 获取开始日期和结束日期的所有日期列表
func GetBetweenDates(begin, end time.Time) []time.Time {
	var dayArr []time.Time
	var currentTime = begin

	for {
		//https://cloud.tencent.com/developer/article/1467743
		if currentTime.After(end) {
			break
		}

		dayArr = append(dayArr, currentTime)
		//currentTime = currentTime.Add(24 * time.Hour)
		currentTime = currentTime.AddDate(0, 0, 1)
	}

	return dayArr
}

// GetWeekDay 获得传入日期这周的开始和结束日期
func GetWeekDay(begin time.Time) (string, string) {
	//now := time.Now()

	//fmt.Printf("now.Year()-------%v\n", now.Year())
	//fmt.Printf("now.Month()-------%v\n", now.Month())
	//fmt.Printf("now.Day()-------%v\n", now.Day())

	offset := int(time.Monday - begin.Weekday()) //周一减去 当天周几   当前时间-到周一天数=周一日期
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastoffset := int(time.Saturday - begin.Weekday()) //周六减去 当天周几    结果+1为周日
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}

	firstOfWeek := time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	lastOfWeeK := time.Date(begin.Year(), begin.Month(), begin.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastoffset+1)
	f := firstOfWeek.Unix()
	l := lastOfWeeK.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

// WeekStart 传入n年，m周，返回m周开始日期(周一)
func WeekStart(year, week int) time.Time {
	//  https://www.mianshigee.com/question/127307wea
	// 从今年年中开始：
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	//获取当年7.1号，属于周几，周日为0-6=周一，
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		//非周日处理    例：7.1为周三 -3+1天=周一日期
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// 几周内的差异：
	_, w := t.ISOWeek() //7.1日为当年那一周
	t = t.AddDate(0, 0, (week-w)*7)
	return ConvertLocalTime(fmt.Sprintf("%s 00:00:00", t.Format("2006-01-02")))
}

// WeekRange 传入n年，m周，返回m周开始(周一)和结束日期(周日)
func WeekRange(year, week int) (start, end time.Time) {
	start = WeekStart(year, week)
	end = ConvertLocalTime(fmt.Sprintf("%s 23:59:59", start.AddDate(0, 0, 6).Format("2006-01-02")))
	return
}

// GetWeek 获取传入时间属于n年 m周 datetime:2006-01-02
// 例 2022-01-01=> 2021_52   2022-02-21=> 2022_8
func GetWeek(datetime string) (y, w int) {
	return ConvertLocalTime(fmt.Sprintf("%s 00:00:00", datetime)).ISOWeek()
}

// GetMonthDay 获得当前月的初始和结束日期
func GetMonthDay(datetime time.Time) (string, string) {

	currentYear, currentMonth, _ := datetime.Date()
	currentLocation := datetime.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	f := firstOfMonth.Unix()
	l := lastOfMonth.Unix()
	return time.Unix(f, 0).Format("2006-01-02") + " 00:00:00", time.Unix(l, 0).Format("2006-01-02") + " 23:59:59"
}

// GetMonthWeekCount 返回当月日期有n周(包括年初年末、月初月末跨周的按最多计算) datetime:2006-01
// 例：2022-01，6，[2021_52 2022_1 2022_2 2022_3 2022_4 2022_5]
func GetMonthWeekCount(datetime string) (string, int, []string) {

	dt := ConvertLocalTime(fmt.Sprintf("%s-01 00:00:00", datetime))
	_, monthlast := GetMonthDay(dt)
	monthlastday := ConvertLocalTime(monthlast)
	year, weekidx := GetWeek(fmt.Sprintf("%s-01", datetime))
	// 2021 52

	//循环，直到取到最后一周周日大于或等于当月最后一天才结束
	var weekArry []string
	weeecount := 0
	for {
		week := weekidx + weeecount
		var weeky, weekd int
		if week > 52 { //跨年处理
			weeky = year + 1
			weekd = week - 52
		} else {
			weeky = year
			weekd = week
		}
		_, dayl := WeekRange(weeky, weekd)
		weekArry = append(weekArry, fmt.Sprintf("%v_%v", weeky, weekd))
		if dayl.Unix() >= monthlastday.Unix() {
			break
		}
		weeecount++
	}
	return datetime, len(weekArry), weekArry
}

// ConvertLocalTime 转为北京时间，datetime:2006-01-02 15:04:05
func ConvertLocalTime(datetime string) time.Time {
	timeLayout := "2006-01-02 15:04:05"
	//loc, _ := time.LoadLocation("Local")

	//因为目标电脑没有安装go的话会导致这个问题出现
	//有些电脑上面使用time.LoadLocation会失败，因为缺少一个文件，可以使用下面的方法替代
	// golang 时间missing Location in call to Date - dz45693 - 博客园
    // https://www.cnblogs.com/majiang/p/15576618.html
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("CST", 8*3600) //替换上海时区
	}
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	return tmp
}
