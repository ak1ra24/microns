package utils

import (
	"os"
	"strconv"
	"strings"
)

func ParsePid(path string) (int, error) {
	fileinfo, err := os.Lstat(path)

	if fileinfo.Mode()&os.ModeSymlink != 0 {
		originfile, err := os.Readlink(path)
		if err != nil {
			return 0, err
		}
		str_pid := strings.Replace(originfile, "/proc/", "", 1)
		str_pid = strings.Replace(str_pid, "/ns/net", "", 1)

		pid, err := strconv.Atoi(str_pid)
		if err != nil {
			return 0, err
		}

		return pid, nil
	}

	return 0, err
}
