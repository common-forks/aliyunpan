// Copyright (c) 2020 tickstep.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package config

import (
	"encoding/hex"
	"github.com/olekukonko/tablewriter"
	"github.com/tickstep/aliyunpan/cmder/cmdtable"
	"github.com/tickstep/library-go/converter"
	"github.com/tickstep/library-go/crypto"
	"github.com/tickstep/library-go/ids"
	"github.com/tickstep/library-go/logger"
	"os"
	"strconv"
	"strings"
)

func (pl *PanUserList) String() string {
	builder := &strings.Builder{}

	tb := cmdtable.NewTable(builder)
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_RIGHT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_CENTER})
	tb.SetHeader([]string{"#", "uid", "用户名", "昵称"})

	for k, userInfo := range *pl {
		tb.Append([]string{strconv.Itoa(k + 1), userInfo.UserId, userInfo.AccountName, userInfo.Nickname})
	}

	tb.Render()

	return builder.String()
}

// AverageParallel 返回平均的下载最大并发量
func AverageParallel(parallel, downloadLoad int) int {
	if downloadLoad < 1 {
		return 1
	}

	p := parallel / downloadLoad
	if p < 1 {
		return 1
	}
	return p
}

func stripPerSecond(sizeStr string) string {
	i := strings.LastIndex(sizeStr, "/")
	if i < 0 {
		return sizeStr
	}
	return sizeStr[:i]
}

func showMaxRate(size int64) string {
	if size <= 0 {
		return "不限制"
	}
	return converter.ConvertFileSize(size, 2) + "/s"
}

// EncryptString 加密
func EncryptString(text string) string {
	if text == "" {
		return ""
	}
	d := []byte(text)
	key := []byte(ids.GetUniqueId("aliyunpan", 16))
	r, e := crypto.EncryptAES(d, key)
	if e != nil {
		return text
	}
	return hex.EncodeToString(r)
}

// DecryptString 解密
func DecryptString(text string) string {
	defer func() {
		if err := recover(); err != nil {
			logger.Verboseln("decrypt string failed, maybe the key has been changed")
		}
	}()

	if text == "" {
		return ""
	}
	d, _ := hex.DecodeString(text)

	// use the machine unique id as the key
	// but in some OS, this key will be changed if you reinstall the OS
	key := []byte(ids.GetUniqueId("aliyunpan", 16))
	r, e := crypto.DecryptAES(d, key)
	if e != nil {
		return text
	}
	return string(r)
}

// isFolderExist 判断文件夹是否存在
func IsFolderExist(pathStr string) bool {
	fi, err := os.Stat(pathStr)
	if err != nil {
		if os.IsExist(err) {
			return fi.IsDir()
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return fi.IsDir()
}
