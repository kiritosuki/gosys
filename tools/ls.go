package tools

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/spf13/cobra"
)

// lsOptions 是选项集合结构体
type lsOptions struct {
	lsAll  bool
	lsLong bool
}

// lenList 用来维护字段最长长度 用于 -l 的格式化输出
// userName groupName size name
var lenList = make(map[string]int)

// opts 全局选项集合变量
var opts = &lsOptions{}

// lsCmd 是 ls 命令
var lsCmd = &cobra.Command{
	Use:   "ls [paths...]",
	Short: "list files in directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runLs(args)
	},
}

// init 用来初始化 lsCmd 配置
func init() {
	lsCmd.Flags().BoolVarP(&opts.lsAll, "all", "a", false, "list all files in directory including hidden files")
	lsCmd.Flags().BoolVarP(&opts.lsLong, "long", "l", false, "list files int directory with long model\n"+
		"permission | username | group | size(/bytes) | updateTime | name")
}

// runLs 是 ls 命令实际核心执行逻辑
func runLs(args []string) error {
	if len(args) == 0 {
		return singleLs("")
	} else if len(args) == 1 {
		return singleLs(args[0])
	} else {
		for _, arg := range args {
			absPath, err := filepath.Abs(arg)
			if err != nil {
				return err
			}
			fmt.Println(absPath + ": ")
			err = singleLs(arg)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// singleLs 是 ls 后面只跟 一个或零个参数 的实现
func singleLs(arg string) error {
	path := "."
	if arg != "" {
		path = arg
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	// 维护 lenList 便于 -l 的格式化输出
	for _, file := range entries {
		infoMap, err := getAllInfo(file)
		if err != nil {
			return err
		}
		// lenList 包括 userName groupName size name 的长度
		if lenList["userName"] < len(infoMap["userName"]) {
			lenList["userName"] = len(infoMap["userName"])
		}
		if lenList["groupName"] < len(infoMap["groupName"]) {
			lenList["groupName"] = len(infoMap["groupName"])
		}
		if lenList["size"] < len(infoMap["size"]) {
			lenList["size"] = len(infoMap["size"])
		}
		if lenList["name"] < len(infoMap["name"]) {
			lenList["name"] = len(infoMap["name"])
		}
	}

	for _, entry := range entries {
		name := entry.Name()
		if !opts.lsAll && name[0] == '.' {
			continue
		}
		if entry.IsDir() {
			name += "/"
		}
		if opts.lsLong {
			// 获取 entry 的字段
			entryMap, err := getAllInfo(entry)
			if err != nil {
				return err
			}
			pmt := entryMap["pmt"]
			userName := entryMap["userName"]
			groupName := entryMap["groupName"]
			size := entryMap["size"]
			updateTime := entryMap["updateTime"]
			// name 已经有了

			// 打印字符串 格式化输出
			resStr := fmt.Sprintf("%s  %-*s  %-*s  %*s  %s  %s",
				pmt, lenList["userName"], userName, lenList["groupName"], groupName, lenList["size"], size, updateTime, name)
			fmt.Println(resStr)
		} else {
			fmt.Print(name + "  ")
		}
	}
	fmt.Println()
	return nil
}

// getAllInfo 获取 entry 的所有信息
// -l 选项
// 分别是：
// 权限 用户 组 文件大小(/bytes) 修改时间 文件名
// pmt userName groupName size updateTime name
func getAllInfo(entry os.DirEntry) (map[string]string, error) {
	argsMap := make(map[string]string)
	info, err := entry.Info()
	if err != nil {
		return nil, err
	}

	// 获取 类型和权限
	// -rwxrwxrwx
	// 写 读 执行
	pmtStr := "----------"
	pmtRunes := []rune(pmtStr)
	// mode 是 类型 + 权限 码
	mode := info.Mode()
	if mode.IsDir() {
		pmtRunes[0] = 'd'
	}
	// 修改权限
	if mode&0400 != 0 {
		pmtRunes[1] = 'r'
	}
	if mode&0200 != 0 {
		pmtRunes[2] = 'w'
	}
	if mode&0100 != 0 {
		pmtRunes[3] = 'x'
	}
	if mode&0040 != 0 {
		pmtRunes[4] = 'r'
	}
	if mode&0020 != 0 {
		pmtRunes[5] = 'w'
	}
	if mode&0010 != 0 {
		pmtRunes[6] = 'x'
	}
	if mode&0004 != 0 {
		pmtRunes[7] = 'r'
	}
	if mode&0002 != 0 {
		pmtRunes[8] = 'w'
	}
	if mode&0001 != 0 {
		pmtRunes[9] = 'x'
	}

	argsMap["pmt"] = string(pmtRunes)

	// 类型断言
	t, ok := info.Sys().(*syscall.Stat_t)
	// 表明不是 *syscall.Stat_t 说明不是类 Unix 系统
	if !ok {
		return nil, errors.New("illegal operation system")
	}

	uid := t.Uid
	gid := t.Gid

	// 获取文件所属用户名
	me, err := user.LookupId(strconv.Itoa(int(uid)))
	if err != nil {
		return nil, err
	}
	userName := me.Username
	argsMap["userName"] = userName

	// 获取文件所属组名
	myGroup, err := user.LookupGroupId(strconv.Itoa(int(gid)))
	if err != nil {
		return nil, err
	}
	groupName := myGroup.Name
	argsMap["groupName"] = groupName

	// 获取文件大小
	size := info.Size()
	argsMap["size"] = strconv.Itoa(int(size))

	// 获取更新时间
	updateTime := info.ModTime()
	formattedTime := updateTime.Format("2006-01-02 15:04:05")
	argsMap["updateTime"] = formattedTime

	// 获取文件名
	argsMap["name"] = info.Name()

	return argsMap, nil
}
