// SPDX-License-Identifier: MIT

package build

import (
	"bytes"
	"encoding/xml"
	"strings"
	"time"

	"github.com/issue9/errwrap"
	"github.com/issue9/version"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/ast"
	"github.com/caixw/apidoc/v7/internal/docs"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/openapi"
	"github.com/caixw/apidoc/v7/internal/xmlenc"
)

// 几种输出的类型
const (
	APIDocXML   = "apidoc+xml"
	OpenapiYAML = "openapi+yaml"
	OpenapiJSON = "openapi+json"
)

type marshaler func(*ast.APIDoc) ([]byte, error)

// Output 指定了渲染输出的相关设置项。
type Output struct {
	// 文档的版本号
	//
	// 该值会覆盖文档中 apidoc.version 的值，方便用户通过代码层面进行版本号同步，
	// 该值无法通过配置文件设置，只能由代码进行设置。
	Version string `yaml:"-"`

	// 导出的文件类型格式，默认为 apidoc 的 XML 文件。
	Type string `yaml:"type,omitempty"`

	// 文档的保存路径
	//
	// 仅适用本地路径
	Path core.URI `yaml:"path"`

	// 只输出该标签的文档，若为空，则表示所有。
	Tags []string `yaml:"tags,omitempty"`

	// xslt 文件地址
	//
	// 默认值为 https://apidoc.tools/docs/ 下当前版本的 apidoc.xsl，比如：
	//  https://apidoc.tools/docs/v7/apidoc.xsl
	//
	// NOTE: 仅针对 xml 类型的输出文件
	Style string `yaml:"style,omitempty"`

	// 命名空间的相关设置
	//
	// 当 namespace 为 true 时会在文档中输出以 core.XMLNamespace 作为命名空间的值，
	// 如果还指定了 NamespacePrefix 则会以此值作为前缀值。
	// NamespacePrefix 仅在 Namespace 为 true 时才启作用。
	//
	// NOTE: 仅针对 Type = APIDocXML
	Namespace       bool   `yaml:"namespace,omitempty"`
	NamespacePrefix string `yaml:"namespace-prefix,omitempty"`

	procInst []string  // 保存所有 xml 的指令内容，包括编码信息
	marshal  marshaler // Type 对应的转换函数
	xml      bool      // 是否为 xml 内容
}

func (o *Output) contains(tags ...string) bool {
	if len(o.Tags) == 0 {
		return true
	}

	for _, t := range o.Tags {
		for _, tag := range tags {
			if tag == t {
				return true
			}
		}
	}
	return false
}

func (o *Output) sanitize() error {
	if o.Type == "" {
		o.Type = APIDocXML
	}

	if o.Version != "" {
		if !version.SemVerValid(o.Version) {
			return core.NewError(locale.ErrInvalidFormat).WithField("version")
		}
	}

	switch o.Type {
	case APIDocXML:
		o.marshal = o.apidocMarshaler
	case OpenapiJSON:
		o.marshal = openapi.JSON
	case OpenapiYAML:
		o.marshal = openapi.YAML
	default:
		return core.NewError(locale.ErrInvalidValue).WithField("type")
	}

	o.xml = strings.HasSuffix(o.Type, "+xml")
	if o.xml {
		if o.Style == "" {
			o.Style = docs.StylesheetURL(core.OfficialURL)
		}

		o.procInst = []string{
			xml.Header,
			`<?xml-stylesheet type="text/xsl" href="` + o.Style + `"?>`,
		}
	}

	if len(o.Path) > 0 {
		scheme, _ := o.Path.Parse()
		if scheme != core.SchemeFile && scheme != "" {
			return core.NewError(locale.ErrInvalidURIScheme, scheme).WithField("path")
		}
	}

	return nil
}

func (o *Output) apidocMarshaler(d *ast.APIDoc) ([]byte, error) {
	if !o.Namespace {
		return xmlenc.Encode("\t", d, "", "")
	}
	return xmlenc.Encode("\t", d, core.XMLNamespace, o.NamespacePrefix)
}

func (o *Output) buffer(d *ast.APIDoc) (*bytes.Buffer, error) {
	filterDoc(d, o)

	if o.Version != "" {
		d.Version = &ast.VersionAttribute{Value: xmlenc.String{Value: o.Version}}
	}

	d.Created = &ast.DateAttribute{Value: ast.Date{Value: time.Now()}}
	d.APIDoc = &ast.APIDocVersionAttribute{Value: xmlenc.String{Value: ast.Version}}

	data, err := o.marshal(d)
	if err != nil {
		return nil, err
	}

	var buf errwrap.Buffer
	if o.xml {
		for _, v := range o.procInst {
			buf.WString(v).WByte('\n')
		}
	}
	buf.WBytes(data)
	if buf.Err != nil {
		return nil, buf.Err
	}

	return &buf.Buffer, nil
}

func filterDoc(d *ast.APIDoc, o *Output) {
	if len(o.Tags) == 0 {
		return
	}

	tags := make([]*ast.Tag, 0, len(o.Tags))
	for _, tag := range d.Tags {
		if o.contains(tag.Name.V()) {
			tags = append(tags, tag)
		}
	}
	d.Tags = tags

	apis := make([]*ast.API, 0, len(d.APIs))
LOOP:
	for _, api := range d.APIs {
		for _, tag := range api.Tags {
			if o.contains(tag.V()) {
				apis = append(apis, api)
				continue LOOP
			}
		}
	}
	d.APIs = apis
}
