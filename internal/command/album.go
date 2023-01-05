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
package command

import (
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/tickstep/aliyunpan-api/aliyunpan"
	"github.com/tickstep/aliyunpan-api/aliyunpan/apierror"
	"github.com/tickstep/aliyunpan/cmder"
	"github.com/tickstep/aliyunpan/cmder/cmdtable"
	"github.com/tickstep/aliyunpan/internal/config"
	"github.com/tickstep/library-go/logger"
	"github.com/urfave/cli"
	"os"
	"strconv"
	"time"
)

type (
	AlbumFileCategoryOption string
)

var (
	ImageOnlyOption      AlbumFileCategoryOption = "image"
	VideoOnlyOption      AlbumFileCategoryOption = "video"
	ImageVideoOnlyOption AlbumFileCategoryOption = "image_video"
	AllFileOption        AlbumFileCategoryOption = "none"
)

func CmdAlbum() cli.Command {
	return cli.Command{
		Name:      "album",
		Aliases:   []string{"abm"},
		Usage:     "相簿(Beta)",
		UsageText: cmder.App().Name + " album",
		Category:  "阿里云盘",
		Before:    ReloadConfigFunc,
		Action: func(c *cli.Context) error {
			cli.ShowCommandHelp(c, c.Command.Name)
			return nil
		},

		Subcommands: []cli.Command{
			{
				Name:      "list",
				Aliases:   []string{"ls"},
				Usage:     "展示相簿列表",
				UsageText: cmder.App().Name + " album list",
				Description: `
示例:

    展示相簿列表 
    aliyunpan album ls
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					RunAlbumList()
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "new",
				Aliases:   []string{""},
				Usage:     "创建相簿",
				UsageText: cmder.App().Name + " album new",
				Description: `
示例:

    新建相簿，名称为：我的相簿2022
    aliyunpan album new "我的相簿2022"

    新建相簿，名称为：我的相簿2022，描述为：存放2022所有文件
    aliyunpan album new "我的相簿2022" "存放2022所有文件"
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					RunAlbumCreate(c.Args().Get(0), c.Args().Get(1))
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "rm",
				Aliases:   []string{""},
				Usage:     "删除相簿",
				UsageText: cmder.App().Name + " album rm",
				Description: `
删除相簿，同名的相簿只会删除第一个符合条件的
示例:

    删除名称为"我的相簿2022"的相簿
    aliyunpan album rm "我的相簿2022"

    删除名称为"我的相簿2022-1" 和 "我的相簿2022-2"的相簿
    aliyunpan album rm "我的相簿2022-1" "我的相簿2022-2"
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					RunAlbumDelete(c.Args())
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "rename",
				Aliases:   []string{""},
				Usage:     "重命名相簿",
				UsageText: cmder.App().Name + " album rename",
				Description: `
重命名相簿，同名的相簿只会修改第一个符合条件的
示例:

    重命名相簿"我的相簿2022"为新的名称"我的相簿2022-new"
    aliyunpan album rename "我的相簿2022" "我的相簿2022-new"
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					RunAlbumRename(c.Args().Get(0), c.Args().Get(1))
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "list-file",
				Aliases:   []string{"lf"},
				Usage:     "展示相簿中的文件",
				UsageText: cmder.App().Name + " album list-file",
				Description: `
展示相簿中文件，同名的相簿只会展示第一个符合条件的
示例:

    展示相簿中文件"我的相簿2022"
    aliyunpan album list-file "我的相簿2022"
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					RunAlbumListFile(c.Args().Get(0))
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "rm-file",
				Aliases:   []string{"rf"},
				Usage:     "移除相簿中的文件",
				UsageText: cmder.App().Name + " album rm-file",
				Description: `
移除相簿中的文件，同名的相簿只会移除第一个符合条件的
示例:

    移除相簿 "我的相簿2022" 中的文件 1.png 2.png
    aliyunpan album rm-file 我的相簿2022 1.png 2.png
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					subArgs := c.Args()
					if len(subArgs) < 2 {
						fmt.Println("请指定移除的文件")
						return nil
					}
					RunAlbumRmFile(subArgs[0], subArgs[1:])
					return nil
				},
				Flags: []cli.Flag{},
			},
			{
				Name:      "add-file",
				Aliases:   []string{"af"},
				Usage:     "增加（文件/相册）网盘文件到相簿中",
				UsageText: cmder.App().Name + " album add-file",
				Description: `
增加文件到相簿中
示例:

    增加当前目录下的 1.png 2.png 文件到相簿 "我的相簿2022" 中
    aliyunpan album add-file 我的相簿2022 1.png 2.png

    增加当前目录下的 myFolder 文件夹下所有文件到相簿 "我的相簿2022" 中
    aliyunpan album add-file 我的相簿2022 myFolder
`,
				Action: func(c *cli.Context) error {
					if config.Config.ActiveUser() == nil {
						fmt.Println("未登录账号")
						return nil
					}
					subArgs := c.Args()
					if len(subArgs) < 2 {
						fmt.Println("请指定增加的文件")
						return nil
					}
					RunAlbumAddFile(subArgs[0], subArgs[1:], ImageVideoOnlyOption)
					return nil
				},
				Flags: []cli.Flag{},
			},
		},
	}
}

func RunAlbumList() {
	activeUser := GetActiveUser()
	records, err := activeUser.PanClient().AlbumListGetAll(&aliyunpan.AlbumListParam{})
	if err != nil {
		fmt.Printf("获取相簿列表失败: %s\n", err)
		return
	}

	tb := cmdtable.NewTable(os.Stdout)
	tb.SetHeader([]string{"#", "ALBUM_ID", "名称", "文件数量", "创建日期", "修改日期"})
	tb.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_CENTER, tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT})
	for k, record := range records {
		tb.Append([]string{strconv.Itoa(k + 1), record.AlbumId, record.Name, strconv.Itoa(record.FileCount),
			record.CreatedAtStr(), record.UpdatedAtStr()})
	}
	tb.Render()
}

func RunAlbumCreate(name, description string) {
	if name == "" {
		fmt.Printf("相簿名称不能为空\n")
		return
	}

	activeUser := GetActiveUser()
	_, err := activeUser.PanClient().AlbumCreate(&aliyunpan.AlbumCreateParam{
		Name:        name,
		Description: description,
	})
	if err != nil {
		fmt.Printf("创建相簿失败: %s\n", err)
		return
	}
	fmt.Printf("创建相簿成功: %s\n", name)
}

func RunAlbumDelete(nameList []string) {
	if len(nameList) == 0 {
		fmt.Printf("相簿名称不能为空\n")
		return
	}

	activeUser := GetActiveUser()
	records, err := activeUser.PanClient().AlbumListGetAll(&aliyunpan.AlbumListParam{})
	if err != nil {
		fmt.Printf("获取相簿列表失败: %s\n", err)
		return
	}

	for _, record := range records {
		for i, name := range nameList {
			if name == record.Name {
				nameList = append(nameList[:i], nameList[i+1:]...)
				_, err := activeUser.PanClient().AlbumDelete(&aliyunpan.AlbumDeleteParam{
					AlbumId: record.AlbumId,
				})
				if err != nil {
					fmt.Printf("删除相簿失败: %s\n", name)
					return
				} else {
					fmt.Printf("删除相簿成功: %s\n", name)
				}
				break
			}
		}
	}
}

func getAlbumFromName(activeUser *config.PanUser, name string) *aliyunpan.AlbumEntity {
	records, err := activeUser.PanClient().AlbumListGetAll(&aliyunpan.AlbumListParam{})
	if err != nil {
		fmt.Printf("获取相簿列表失败: %s\n", err)
		return nil
	}

	for _, record := range records {
		if name == record.Name {
			return record
		}
	}
	return nil
}

func RunAlbumRename(name, newName string) {
	if len(name) == 0 {
		fmt.Printf("相簿名称不能为空\n")
		return
	}
	if len(newName) == 0 {
		fmt.Printf("相簿名称不能为空\n")
		return
	}

	activeUser := GetActiveUser()
	record := getAlbumFromName(activeUser, name)
	if record == nil {
		return
	}
	_, err := activeUser.PanClient().AlbumEdit(&aliyunpan.AlbumEditParam{
		AlbumId:     record.AlbumId,
		Description: record.Description,
		Name:        newName,
	})
	if err != nil {
		fmt.Printf("重命名相簿失败: %s\n", name)
		return
	} else {
		fmt.Printf("重命名相簿成功: %s -> %s\n", name, newName)
	}
}

func RunAlbumListFile(name string) {
	if len(name) == 0 {
		fmt.Printf("相簿名称不能为空\n")
		return
	}

	activeUser := GetActiveUser()
	record := getAlbumFromName(activeUser, name)
	if record == nil {
		return
	}

	fileList, er := activeUser.PanClient().AlbumListFileGetAll(&aliyunpan.AlbumListFileParam{
		AlbumId: record.AlbumId,
	})
	if er != nil {
		fmt.Printf("获取相簿文件列表失败：%s\n", er)
		return
	}
	renderTable(opLs, false, "", fileList)
}

func RunAlbumRmFile(name string, nameList []string) {
	if len(name) == 0 {
		fmt.Printf("相簿名称不能为空\n")
		return
	}
	if len(nameList) == 0 {
		fmt.Printf("指定文件不能为空\n")
		return
	}

	activeUser := GetActiveUser()
	album := getAlbumFromName(activeUser, name)
	if album == nil {
		return
	}

	fileList, er := activeUser.PanClient().AlbumListFileGetAll(&aliyunpan.AlbumListFileParam{
		AlbumId: album.AlbumId,
	})
	if er != nil {
		fmt.Printf("获取相簿文件列表失败：%s\n", er)
		return
	}
	param := &aliyunpan.AlbumDeleteFileParam{
		AlbumId:       album.AlbumId,
		DriveFileList: []aliyunpan.FileBatchActionParam{},
	}
	for _, file := range fileList {
		if len(nameList) == 0 {
			break
		}
		for i, name := range nameList {
			if name == file.FileName {
				nameList = append(nameList[:i], nameList[i+1:]...)
				param.AddFileItem(file.DriveId, file.FileId)
				break
			}
		}
	}

	// 1-500 范围
	if len(param.DriveFileList) == 0 {
		fmt.Printf("没有符合的文件\n")
		return
	}
	// delete file
	_, e := activeUser.PanClient().AlbumDeleteFile(param)
	if e != nil {
		fmt.Printf("删除相簿文件失败：%s\n", e)
		return
	}
	fmt.Printf("删除相簿文件成功：%s\n", name)
}

// RunAlbumAddFile 增加网盘文件到相簿
func RunAlbumAddFile(albumName string, filePathList []string, filterOption AlbumFileCategoryOption) {
	activeUser := GetActiveUser()

	if albumName == "" {
		fmt.Printf("必须指定相簿名称\n")
		return
	}
	album := getAlbumFromName(activeUser, albumName)
	if album == nil {
		fmt.Printf("相簿不存在\n")
		return
	}

	paths, err := makePathAbsolute(activeUser.ActiveDriveId, filePathList...)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(paths) == 0 {
		fmt.Printf("没有有效的文件\n")
		return
	}

	fmt.Printf("正在获取增加的文件信息，该操作可能会非常耗费时间，请耐心等待...\n")
	param := &aliyunpan.AlbumAddFileParam{
		AlbumId:       album.AlbumId,
		DriveFileList: []aliyunpan.FileBatchActionParam{},
	}
	for k := range paths {
		filePath := paths[k]
		fileInfo, apierr := activeUser.PanClient().FileInfoByPath(activeUser.ActiveDriveId, filePath)
		if apierr != nil {
			fmt.Printf("获取文件信息失败: %s\n", filePath)
			continue
		}
		if fileInfo.IsFile() {
			// file
			if isFileMatchCondition(fileInfo, filterOption) {
				param.AddFileItem(fileInfo.DriveId, fileInfo.FileId)
			}
		} else {
			// folder
			activeUser.PanClient().FilesDirectoriesRecurseList(activeUser.ActiveDriveId, fileInfo.Path, func(depth int, _ string, fd *aliyunpan.FileEntity, apiError *apierror.ApiError) bool {
				if apiError != nil {
					logger.Verbosef("%s\n", apiError)
					return true
				}
				if !fd.IsFolder() {
					if isFileMatchCondition(fd, filterOption) {
						param.AddFileItem(fd.DriveId, fd.FileId)
					}
				}
				time.Sleep(2 * time.Second)
				return true
			})
		}
		time.Sleep(2 * time.Second)
	}

	if len(param.DriveFileList) == 0 {
		fmt.Printf("没有符合的文件\n")
		return
	}
	// add file
	_, e := activeUser.PanClient().AlbumAddFile(param)
	if e != nil {
		fmt.Printf("增加相簿文件失败：%s\n", e)
		return
	}
	fmt.Printf("增加相簿文件成功：%s\n", albumName)
}

func isFileMatchCondition(fileInfo *aliyunpan.FileEntity, filterOption AlbumFileCategoryOption) bool {
	if fileInfo == nil {
		return false
	}
	if filterOption == ImageOnlyOption {
		return fileInfo.Category == "image"
	} else if filterOption == VideoOnlyOption {
		return fileInfo.Category == "video"
	} else if filterOption == ImageVideoOnlyOption {
		return fileInfo.Category == "image" || fileInfo.Category == "video"
	} else if filterOption == AllFileOption {
		return true
	}
	return false
}
