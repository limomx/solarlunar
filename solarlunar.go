package solarlunar

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"
)

var MIN_YEAR = 1900
var MAX_YEAR = 2049

var DATELAYOUT = "2006-01-02"
var STARTDATESTR = "1900-01-30"

var TIANGAN = []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
var DIZHI = []string{"子", "丑", "寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥"}
var ANIMALS = []string{"鼠", "牛", "虎", "兔", "龙", "蛇", "马", "羊", "猴", "鸡", "狗", "猪"}
var DIZHI_MONTH = []string{"寅", "卯", "辰", "巳", "午", "未", "申", "酉", "戌", "亥", "子", "丑"}
var CHINESENUMBER = []string{"一", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "十二"}
var CHINESENUMBERSPECIAL = []string{"正", "二", "三", "四", "五", "六", "七", "八", "九", "十", "十一", "腊"}
var WEEBNUMBER = map[string]string{"Monday": "一", "Tuesday": "二", "Wednesday": "三", "Thursday": "四", "Friday": "五", "Saturday": "六", "Sunday": "日"}
var MONTHNUMBER = map[string]int{"January": 1, "February": 2, "March": 3, "April": 4, "May": 5, "June": 6, "July": 7, "August": 8, "September": 9, "October": 10, "November": 11, "December": 12}
var LUNAR_INFO = []int{
	0x04bd8, 0x04ae0, 0x0a570, 0x054d5, 0x0d260, 0x0d950, 0x16554, 0x056a0, 0x09ad0, 0x055d2,
	0x04ae0, 0x0a5b6, 0x0a4d0, 0x0d250, 0x1d255, 0x0b540, 0x0d6a0, 0x0ada2, 0x095b0, 0x14977,
	0x04970, 0x0a4b0, 0x0b4b5, 0x06a50, 0x06d40, 0x1ab54, 0x02b60, 0x09570, 0x052f2, 0x04970,
	0x06566, 0x0d4a0, 0x0ea50, 0x06e95, 0x05ad0, 0x02b60, 0x186e3, 0x092e0, 0x1c8d7, 0x0c950,
	0x0d4a0, 0x1d8a6, 0x0b550, 0x056a0, 0x1a5b4, 0x025d0, 0x092d0, 0x0d2b2, 0x0a950, 0x0b557,
	0x06ca0, 0x0b550, 0x15355, 0x04da0, 0x0a5d0, 0x14573, 0x052d0, 0x0a9a8, 0x0e950, 0x06aa0,
	0x0aea6, 0x0ab50, 0x04b60, 0x0aae4, 0x0a570, 0x05260, 0x0f263, 0x0d950, 0x05b57, 0x056a0,
	0x096d0, 0x04dd5, 0x04ad0, 0x0a4d0, 0x0d4d4, 0x0d250, 0x0d558, 0x0b540, 0x0b5a0, 0x195a6,
	0x095b0, 0x049b0, 0x0a974, 0x0a4b0, 0x0b27a, 0x06a50, 0x06d40, 0x0af46, 0x0ab60, 0x09570,
	0x04af5, 0x04970, 0x064b0, 0x074a3, 0x0ea50, 0x06b58, 0x055c0, 0x0ab60, 0x096d5, 0x092e0,
	0x0c960, 0x0d954, 0x0d4a0, 0x0da50, 0x07552, 0x056a0, 0x0abb7, 0x025d0, 0x092d0, 0x0cab5,
	0x0a950, 0x0b4a0, 0x0baa4, 0x0ad50, 0x055d9, 0x04ba0, 0x0a5b0, 0x15176, 0x052b0, 0x0a930,
	0x07954, 0x06aa0, 0x0ad50, 0x05b52, 0x04b60, 0x0a6e6, 0x0a4e0, 0x0d260, 0x0ea65, 0x0d530,
	0x05aa0, 0x076a3, 0x096d0, 0x04bd7, 0x04ad0, 0x0a4d0, 0x1d0b6, 0x0d250, 0x0d520, 0x0dd45,
	0x0b5a0, 0x056d0, 0x055b2, 0x049b0, 0x0a577, 0x0a4b0, 0x0aa50, 0x1b255, 0x06d20, 0x0ada0}

/**
* 日柱依照高氏公式计算, 需要世纪常数, 月基常数
* https://gss0.bdstatic.com/94o3dSag_xI4khGkpoWK1HF6hhy/baike/s%3D247/sign=4c4b339e92eef01f49141fc1d7fc99e0/a686c9177f3e670931f9cf053dc79f3df9dc554b.jpg
**/

// 月基数
var M = map[int]int{3: 0, 5: 1, 7: 2, 9: 4, 11: 5, 13: 6, 4: 31, 6: 326, 8: 33, 10: 34, 12: 35, 14: 37}

//世纪常数, 预算好
var X = map[int]int{17: -3, 18: 41, 19: 25, 20: 9, 21: -6, 22: 38, 23: 22, 24: 6, 25: -9, 26: 35}

//干支表, 60年一轮回

func LunarToSolar(date string, leapMonthFlag bool) string {
	loc, _ := time.LoadLocation("Local")
	lunarTime, err := time.ParseInLocation(DATELAYOUT, date, loc)
	if err != nil {
		fmt.Println(err.Error())
	}
	lunarYear := lunarTime.Year()
	lunarMonth := MONTHNUMBER[lunarTime.Month().String()]
	lunarDay := lunarTime.Day()
	err = checkLunarDate(lunarYear, lunarMonth, lunarDay, leapMonthFlag)

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	offset := 0

	for i := MIN_YEAR; i < lunarYear; i++ {
		yearDaysCount := getYearDays(i) // 求阴历某年天数
		offset += yearDaysCount
	}
	//计算该年闰几月
	leapMonth := getLeapMonth(lunarYear)
	if leapMonthFlag && leapMonth != lunarMonth {
		panic("您输入的闰月标志有误！")
	}
	if leapMonth == 0 || (lunarMonth < leapMonth) || (lunarMonth == leapMonth && !leapMonthFlag) {
		for i := 1; i < lunarMonth; i++ {
			tempMonthDaysCount := getMonthDays(lunarYear, uint(i))
			offset += tempMonthDaysCount
		}

		// 检查日期是否大于最大天
		if lunarDay > getMonthDays(lunarYear, uint(lunarMonth)) {
			panic("不合法的农历日期！")
		}
		offset += lunarDay // 加上当月的天数
	} else { //当年有闰月，且月份晚于或等于闰月
		for i := 1; i < lunarMonth; i++ {
			tempMonthDaysCount := getMonthDays(lunarYear, uint(i))
			offset += tempMonthDaysCount
		}
		if lunarMonth > leapMonth {
			temp := getLeapMonthDays(lunarYear) // 计算闰月天数
			offset += temp                      // 加上闰月天数

			if lunarDay > getMonthDays(lunarYear, uint(lunarMonth)) {
				panic("不合法的农历日期！")
			}
			offset += lunarDay
		} else { // 如果需要计算的是闰月，则应首先加上与闰月对应的普通月的天数
			// 计算月为闰月
			temp := getMonthDays(lunarYear, uint(lunarMonth)) // 计算非闰月天数
			offset += temp

			if lunarDay > getLeapMonthDays(lunarYear) {
				panic("不合法的农历日期！")
			}
			offset += lunarDay
		}
	}

	myDate, err := time.ParseInLocation(DATELAYOUT, STARTDATESTR, loc)
	if err != nil {
		fmt.Println(err.Error())
	}

	dayDuaration, _ := time.ParseDuration("24h")
	myDate = myDate.Add(dayDuaration * time.Duration(offset))
	return myDate.Format(DATELAYOUT)
}

// 星期
func WeedToChinese() string {
	return "星期" + WEEBNUMBER[time.Now().Weekday().String()]
}

func SolarToChineseLuanr(date string) string {

	// fmt.Println(date)

	lunarYear, lunarMonth, lunarDay, leapMonth, leapMonthFlag := calculateLunar(date)
	result := cyclical(lunarYear+lunarDay-lunarDay) + "年"
	if leapMonthFlag && (lunarMonth == leapMonth) {
		result += "闰"
	}

	// result += CHINESENUMBERSPECIAL[lunarMonth-1] + "月"
	// result += DIZHI_MONTH[lunarMonth-1] + "月"

	result += chineseMonthTD(lunarMonth-1, TIANGAN[(lunarYear-1900+36)%10])
	//日柱
	result += chineseDayTD(date) + "日"
	// fmt.Println(lunarDay)
	// result += chineseDayString(lunarDay) + "日"
	return result
}

func SolarToSimpleLuanr(date string) string {
	lunarYear, lunarMonth, lunarDay, leapMonth, leapMonthFlag := calculateLunar(date)
	result := strconv.Itoa(lunarYear) + "年"
	if leapMonthFlag && (lunarMonth == leapMonth) {
		result += "闰"
	}
	if lunarMonth < 10 {
		result += "0" + strconv.Itoa(lunarMonth) + "月"
	} else {
		result += strconv.Itoa(lunarMonth) + "月"
	}
	if lunarDay < 10 {
		result += "0" + strconv.Itoa(lunarDay) + "日"
	} else {
		result += strconv.Itoa(lunarDay) + "日"
	}
	return result
}

func SolarToLuanr(date string) (string, bool) {
	lunarYear, lunarMonth, lunarDay, leapMonth, leapMonthFlag := calculateLunar(date)
	result := strconv.Itoa(lunarYear) + "-"
	if lunarMonth < 10 {
		result += "0" + strconv.Itoa(lunarMonth) + "-"
	} else {
		result += strconv.Itoa(lunarMonth) + "-"
	}
	if lunarDay < 10 {
		result += "0" + strconv.Itoa(lunarDay)
	} else {
		result += strconv.Itoa(lunarDay)
	}

	if leapMonthFlag && (lunarMonth == leapMonth) {
		return result, true
	} else {
		return result, false
	}
}

func calculateLunar(date string) (lunarYear, lunarMonth, lunarDay, leapMonth int, leapMonthFlag bool) {
	loc, _ := time.LoadLocation("Local")
	i := 0
	temp := 0
	leapMonthFlag = false
	isLeapYear := false

	myDate, err := time.ParseInLocation(DATELAYOUT, date, loc)
	if err != nil {
		fmt.Println(err.Error())
	}
	startDate, err := time.ParseInLocation(DATELAYOUT, STARTDATESTR, loc)
	if err != nil {
		fmt.Println(err.Error())
	}

	offset := daysBwteen(myDate, startDate)
	for i = MIN_YEAR; i < MAX_YEAR; i++ {
		temp = getYearDays(i) //求当年农历年天数
		if offset-temp < 1 {
			break
		} else {
			offset -= temp
		}
	}
	lunarYear = i

	leapMonth = getLeapMonth(lunarYear) //计算该年闰哪个月

	//设定当年是否有闰月
	if leapMonth > 0 {
		isLeapYear = true
	} else {
		isLeapYear = false
	}

	for i = 1; i <= 12; i++ {
		if i == leapMonth+1 && isLeapYear {
			temp = getLeapMonthDays(lunarYear)
			isLeapYear = false
			leapMonthFlag = true
			i--
		} else {
			temp = getMonthDays(lunarYear, uint(i))
		}
		offset -= temp
		if offset <= 0 {
			break
		}
	}
	offset += temp
	lunarMonth = i
	lunarDay = offset
	return
}

func checkLunarDate(lunarYear, lunarMonth, lunarDay int, leapMonthFlag bool) error {
	if (lunarYear < MIN_YEAR) || (lunarYear > MAX_YEAR) {
		return errors.New("非法农历年份！")
	}
	if (lunarMonth < 1) || (lunarMonth > 12) {
		return errors.New("非法农历月份！")
	}
	if (lunarDay < 1) || (lunarDay > 30) { // 中国的月最多30天
		return errors.New("非法农历天数！")
	}

	leap := getLeapMonth(lunarYear) // 计算该年应该闰哪个月
	if (leapMonthFlag == true) && (lunarMonth != leap) {
		return errors.New("非法闰月！")
	}
	return nil
}

// 计算该月总天数
func getMonthDays(lunarYeay int, month uint) int {
	if (month > 31) || (month < 0) {
		fmt.Println("error month")
	}
	// 0X0FFFF[0000 {1111 1111 1111} 1111]中间12位代表12个月，1为大月，0为小月
	bit := 1 << (16 - month)
	if ((LUNAR_INFO[lunarYeay-1900] & 0x0FFFF) & bit) == 0 {
		return 29
	} else {
		return 30
	}
}

// 计算阴历年的总天数
func getYearDays(year int) int {
	sum := 29 * 12
	for i := 0x8000; i >= 0x8; i >>= 1 {
		if (LUNAR_INFO[year-1900] & 0xfff0 & i) != 0 {
			sum++
		}
	}
	return sum + getLeapMonthDays(year)
}

//	计算阴历年闰月多少天
func getLeapMonthDays(year int) int {
	if getLeapMonth(year) != 0 {
		if (LUNAR_INFO[year-1900] & 0xf0000) == 0 {
			return 29
		} else {
			return 30
		}
	} else {
		return 0
	}
}

//	计算阴历年闰哪个月 1-12 , 没闰传回 0
func getLeapMonth(year int) int {
	return (int)(LUNAR_INFO[year-1900] & 0xf)
}

// 计算差的天数
func daysBwteen(myDate time.Time, startDate time.Time) int {
	subValue := float64(myDate.Unix()-startDate.Unix())/86400.0 + 0.5
	return int(subValue)
}

// 月份转为天干地支
func chineseMonthTD(lunar int, year string) string {

	JY := []string{"丙", "丁", "戊", "己", "庚", "辛", "壬", "癸", "甲", "乙", "丙", "丁"}
	YG := []string{"戊", "己", "庚", "辛", "壬", "癸", "甲", "乙", "丙", "丁", "戊", "己"}
	BX := []string{"庚", "辛", "壬", "癸", "甲", "乙", "丙", "丁", "戊", "己", "庚", "辛"}
	DR := []string{"壬", "癸", "甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸"}
	WK := []string{"甲", "乙", "丙", "丁", "戊", "己", "庚", "辛", "壬", "癸", "甲", "乙"}
	Pool := map[string][]string{"甲": JY, "己": JY, "乙": YG, "庚": YG, "丙": BX, "辛": BX, "丁": DR, "壬": DR, "戊": WK, "癸": WK}

	m := lunar
	endM := DIZHI_MONTH[m]

	lunarStr := Pool[year][lunar] + endM + "月"

	return lunarStr
}

func getDate(date string) (year int, month int, day int) {
	loc, _ := time.LoadLocation("Local")
	myDate, err := time.ParseInLocation(DATELAYOUT, date, loc)
	if err != nil {
		fmt.Println(err.Error())
	}
	year = myDate.Year()
	month = int(myDate.Month())
	day = myDate.Day()
	return
}

//日柱计算
func chineseDayTD(date string) string {
	year, month, day := getDate(date)
	//公元年数后两位数
	s := year % 10e1
	u := s % 4
	m := M[month]
	d := day
	// _x := year / 100
	x := X[int(math.Floor(float64(year/100+1)))]
	// 高氏公式
	r := int(math.Floor(float64(s/4))*6) + 5*(u+int(math.Floor(float64(s/4))*3)) + m + d + x

	return TIANGAN[r%10-1] + DIZHI[r%12-1]
}

func cyclicalm(num int) string {
	return TIANGAN[num%10] + DIZHI[num%12] + ANIMALS[num%12]
}

func cyclical(year int) string {
	num := year - 1900 + 36
	return cyclicalm(num)
}

func chineseDayString(day int) string {
	chineseTen := []string{"初", "十", "廿", "三"}
	n := 0
	if day%10 == 0 {
		n = 9
	} else {
		n = day%10 - 1
	}
	if day > 30 {
		return ""
	}
	if day == 20 {
		return "二十"
	} else if day == 10 {
		return "初十"
	} else {
		return chineseTen[day/10] + CHINESENUMBER[n]
	}
}
