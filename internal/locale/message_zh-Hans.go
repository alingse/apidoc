// SPDX-License-Identifier: MIT

package locale

import "golang.org/x/text/language"

var zhHans = map[string]string{
	// 与 flag 包相关的处理
	CmdUsage: `%s 是一个 RESTful API 文档生成工具

用法：
apidoc cmd [args]

cmd 为子命令，args 为传递给子命令的参数，目前支持以下子命令。
%s

源代码采用 MIT 开源许可证，发布于 %s
详细信息可访问官网 %s
`,
	CmdHelpUsage:    "显示帮助信息",
	CmdVersionUsage: "显示版本信息",
	CmdLangUsage:    "显示所有支持的语言",
	CmdLocaleUsage:  "显示所有支持的本地化内容",
	CmdDetectUsage:  "根据目录下的内容生成配置文件",
	CmdTestUsage:    "测试语法的正确性",
	CmdMockUsage:    "启用 Mock 服务",
	CmdBuildUsage:   "生成文档内容",
	Version:         "版本：%s\n文档：%s\n提交：%s\nGo：%s",
	CmdNotFound:     "子命令 %s 未找到\n",

	VersionInCompatible: "当前程序与配置文件中指定的版本号不兼容",
	Complete:            "完成！文档保存在：%s，总用时：%v",
	ConfigWriteSuccess:  "配置内容成功写入 %s",
	TestSuccess:         "语法没有问题！",
	LangID:              "ID",
	LangName:            "名称",
	LangExts:            "扩展名",

	// 错误信息，可能在地方用到
	ErrRequired:              "不能为空",
	ErrInvalidFormat:         "格式不正确",
	ErrDirNotExists:          "目录不存在",
	ErrUnsupportedInputLang:  "不支持的输入语言：%s",
	ErrNotFoundEndFlag:       "找不到结束符号",
	ErrNotFoundSupportedLang: "该目录下没有支持的语言文件",
	ErrDirIsEmpty:            "目录下没有需要解析的文件",
	ErrInvalidValue:          "无效的值",
	ErrPathNotMatchParams:    "地址参数不匹配",
	ErrDuplicateValue:        "重复的值",
	ErrMessage:               "%s 位于 %s",
	ErrNotFound:              "未找到该值",

	// logs
	InfoPrefix:    "[信息] ",
	WarnPrefix:    "[警告] ",
	ErrorPrefix:   "[错误] ",
	SuccessPrefix: "[成功] ",
}

func init() {
	addLocale(language.MustParse("zh-Hans"), zhHans)

	// 大部分的系统都采用 zh-cn 作为语言标记，
	// 但是 golang.org/x/text 现在不能将 zh-cn 自动转换成 zh-hans
	addLocale(language.MustParse("zh-cn"), zhHans)
}
