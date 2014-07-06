package schedule

import (
	"fmt"
	"reflect"
	"time"
)

//获取距启动的时间（秒）
func getCountDown(cyc string, sm []int, ss []time.Duration) (countDown time.Duration, err error) { // {{{
	now := GetNow()
	var startTime time.Time
	var b bool //执行时间是否在当前时间之后的标志

	//按周期取整
	s := TruncDate(cyc, now)

	for i, st := range ss {
		if s.AddDate(0, sm[i], 0).Add(st).After(now) {
			//执行时间在当前时间之后，设置标志，跳出循环进行下一步
			b = true
			startTime = s.AddDate(0, sm[i], 0).Add(st)
			break
		}
	}

	if !b { //在当前时间周期内的执行时间全部小于当前时间，执行应该从下一周期开始
		//解析周期并取得距下一周期的时间
		switch {
		case cyc == "ss":
			startTime = s.AddDate(0, sm[0], 0).Add(time.Second).Add(ss[0])
		case cyc == "mi":
			//按分钟取整
			startTime = s.AddDate(0, sm[0], 0).Add(time.Minute).Add(ss[0])
		case cyc == "h":
			//按小时取整
			startTime = s.AddDate(0, sm[0], 0).Add(time.Hour).Add(ss[0])
		case cyc == "d":
			//按日取整
			startTime = s.AddDate(0, sm[0], 0).AddDate(0, 0, 1).Add(ss[0])
		case cyc == "m":
			//按月取整
			startTime = s.AddDate(0, sm[0], 0).AddDate(0, 1, 0).Add(ss[0])
		case cyc == "w":
			//按周取整
			startTime = s.AddDate(0, sm[0], 0).AddDate(0, 0, 7).Add(ss[0])
		case cyc == "q":
			//回头再处理
		case cyc == "y":
			//按年取整
			startTime = s.AddDate(0, sm[0], 0).AddDate(1, 0, 0).Add(ss[0])
		}

	}
	countDown = startTime.Sub(time.Now())

	return countDown, nil

} // }}}

//时间取整
func TruncDate(cyc string, now time.Time) time.Time { // {{{

	//解析周期并取得距下一周期的时间
	switch {
	case cyc == "ss":
		//按秒取整
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), 0, time.Local)
	case cyc == "mi":
		//按分钟取整
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, time.Local)

	case cyc == "h":
		//按小时取整
		return time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, time.Local)
	case cyc == "d":
		//按日取整
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	case cyc == "m":
		//按月取整
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	case cyc == "w":
		//按周取整
		return time.Date(now.Year(), now.Month(), now.Day()-int(now.Weekday()), 0, 0, 0, 0, time.Local)
	case cyc == "q":
		//回头再处理
	case cyc == "y":
		//按年取整
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}
	return time.Now()

} // }}}

//获取当前时间
func GetNow() time.Time { // {{{
	return time.Now().Local()
} // }}}

//CheckErr检查错误信息，若有错误则打印并抛出异常。
func CheckErr(info string, err error) { // {{{
	if err != nil {
		g.L.Errorln(info, err.Error())
		panic(err)
	}
} // }}}

//PrintErr打印错误信息
func PrintErr(info string, err error) {
	if err != nil {
		g.L.Errorln(info, err.Error())
	}
}

//Copy复制对象
//来自github.com/jinzhu/copier
func Copy(copy_to interface{}, copy_from interface{}) (err error) {
	var (
		is_slice    bool
		from_typ    reflect.Type
		to_typ      reflect.Type
		elem_amount int
	)

	from := reflect.ValueOf(copy_from)
	to := reflect.ValueOf(copy_to)
	from_elem := reflect.Indirect(from)
	to_elem := reflect.Indirect(to)

	if to_elem.Kind() == reflect.Slice {
		is_slice = true
		if from_elem.Kind() == reflect.Slice {
			from_typ = from_elem.Type().Elem()
			elem_amount = from_elem.Len()
		} else {
			from_typ = from_elem.Type()
			elem_amount = 1
		}

		to_typ = to_elem.Type().Elem()
	} else {
		from_typ = from_elem.Type()
		to_typ = to_elem.Type()
		elem_amount = 1
	}

	for e := 0; e < elem_amount; e++ {
		var dest, source reflect.Value
		if is_slice {
			if from_elem.Kind() == reflect.Slice {
				source = from_elem.Index(e)
			} else {
				source = from_elem
			}
		} else {
			source = from_elem
		}

		if is_slice {
			dest = reflect.New(to_typ).Elem()
		} else {
			dest = to_elem
		}

		for i := 0; i < from_typ.NumField(); i++ {
			field := from_typ.Field(i)
			if !field.Anonymous {
				name := field.Name
				from_field := source.FieldByName(name)
				to_field := dest.FieldByName(name)
				to_method := dest.Addr().MethodByName(name)
				if from_field.IsValid() && to_field.IsValid() {
					to_field.Set(from_field)
				}

				if from_field.IsValid() && to_method.IsValid() {
					to_method.Call([]reflect.Value{from_field})
				}
			}
		}

		for i := 0; i < dest.NumField(); i++ {
			field := to_typ.Field(i)
			if !field.Anonymous {
				name := field.Name
				from_method := source.Addr().MethodByName(name)
				to_field := dest.FieldByName(name)

				if from_method.IsValid() && to_field.IsValid() {
					values := from_method.Call([]reflect.Value{})
					if len(values) >= 1 {
						to_field.Set(values[0])
					}
				}
			}
		}

		if is_slice {
			to_elem.Set(reflect.Append(to_elem, dest))
		}
	}
	return
}

//打印调度信息
func printSchedule(scds map[int64]*Schedule) { // {{{
	for _, scd := range scds {
		fmt.Println(scd.Name, "\tjobs=", scd.JobCnt, " tasks=", scd.TaskCnt)
		//打印调度中的作业信息
		for j := scd.Job; j != nil; {
			fmt.Println("\t--------------------------------------")
			fmt.Println("\t", j.Name)
			//打印作业中的任务信息
			for _, t := range j.Tasks {
				fmt.Println("\t\t", t.Name)

				fmt.Print("\t\t\t[")
				//打印任务依赖链
				for _, r := range t.RelTasks {
					fmt.Print(r.Name, ",")

				}
				fmt.Print("]\n")
			}
			fmt.Print("\n")
			j = j.NextJob

		}

	}
} // }}}
