package utils

import (
	"os"
	"strconv"
	"strings"
)

// ParsePid func is Parse Docker pid
func ParsePid(path string) (int, error) {
	fileinfo, err := os.Lstat(path)

	if fileinfo.Mode()&os.ModeSymlink != 0 {
		originfile, err := os.Readlink(path)
		if err != nil {
			return 0, err
		}
		strPid := strings.Replace(originfile, "/proc/", "", 1)
		strPid = strings.Replace(strPid, "/ns/net", "", 1)

		pid, err := strconv.Atoi(strPid)
		if err != nil {
			return 0, err
		}

		return pid, nil
	}

	return 0, err
}
