package cache

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func inCIDR(addr, network string) (bool, error) {
	_, subnet, err := net.ParseCIDR(network)
	if err == nil {
		ip := net.ParseIP(addr)
		if ip != nil {
			if subnet.Contains(ip) {
				return true, nil
			}
		}
		return false, fmt.Errorf("invalid IP address")
	}
	return false, err
}

func equal(val1, val2 string) bool {
	return val1 == val2
}

func lowerEqual(val1, val2 string) bool {
	return equal(strings.ToLower(val1), strings.ToLower(val2))
}

func contain(str, substr string) bool {
	return strings.Contains(str, substr)
}

func in(str, list string) bool {
	l := strings.Split(list, ",")
	for _, value := range l {
		if str == value {
			return true
		}
	}
	return false
}

func startWith(str, pref string) bool {
	return strings.HasPrefix(str, pref)
}

func endWith(str, suff string) bool {
	return strings.HasSuffix(str, suff)
}

func expresion(exp, str string) (bool, error) {
	re, err := regexp.Compile(exp)
	if err == nil {
		if re.MatchString(str) {
			return true, nil
		}
	}
	return false, err
}

func parseFloats(val1, val2 string) (float64, float64, error) {
	f1, err1 := strconv.ParseFloat(val1, 64)
	if err1 != nil {
		return 0, 0, err1
	}
	f2, err2 := strconv.ParseFloat(val2, 64)
	if err2 != nil {
		return 0, 0, err2
	}
	return f1, f2, nil
}

func compare(operator, val1, val2 string) bool {
	switch operator {
	case "==":
		return equal(val1, val2)
	case "!=":
		return !equal(val1, val2)
	case "<>":
		return !equal(val1, val2)
	case "::":
		return lowerEqual(val1, val2)
	case "!!":
		return !lowerEqual(val1, val2)
	case "contains":
		return contain(val1, val2)
	case "not contain":
		return !contain(val1, val2)
	case "in":
		return in(val1, val2)
	case "not in":
		return !in(val1, val2)
	case "start with":
		return startWith(val1, val2)
	case "not start with":
		return !startWith(val1, val2)
	case "end with":
		return endWith(val1, val2)
	case "not end with":
		return !endWith(val1, val2)
	case "regexp":
		matched, err := expresion(val2, val1)
		if err != nil {
			return false
		}
		return matched
	case "not regexp":
		matched, err := expresion(val2, val1)
		if err != nil {
			return false
		}
		return matched
	case "<":
		f1, f2, err := parseFloats(val1, val2)
		if err != nil {
			return false
		}
		return f1 < f2
	case ">":
		f1, f2, err := parseFloats(val1, val2)
		if err != nil {
			return false
		}
		return f1 > f2
	case "<=":
		f1, f2, err := parseFloats(val1, val2)
		if err != nil {
			return false
		}
		return f1 <= f2
	case ">=":
		f1, f2, err := parseFloats(val1, val2)
		if err != nil {
			return false
		}
		return f1 >= f2
	case "exist":
		return true
	case "in cidr":
		matched, err := inCIDR(val1, val2)
		if err == nil {
			return matched
		}
		return false
	case "not in cidr":
		matched, err := inCIDR(val1, val2)
		if err == nil {
			return !matched
		}
		return false
	default:
		return false
	}
}

func evalElement(elem, field, operator, value string) bool {
	if gjson.Get(elem, field).Exists() {
		return compare(operator, gjson.Get(elem, field).String(), value)
	} else if operator == "not exist" {
		return true
	}
	return false
}
