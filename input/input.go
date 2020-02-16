// SPDX-License-Identifier: MIT

// Package input 用于处理输入的文件，从代码中提取基本的注释内容。
//
// 多行注释和单行注释在处理上会有一定区别：
//  - 单行注释，风格相同且相邻的注释会被合并成一个注释块；
//  - 单行注释，风格不相同且相邻的注释会被按注释风格多个注释块；
//  - 多行注释，即使两个注释释块相邻也会被分成两个注释块来处理。
package input

import (
	"bytes"
	"io/ioutil"
	"sync"

	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"

	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/internal/lang"
	"github.com/caixw/apidoc/v6/message"
)

// Parse 分析从 input 中获取的代码块
//
// 所有与解析有关的错误均通过 h 输出。
// 如果是配置文件的错误，则通过 error 返回
func Parse(h *message.Handler, opt ...*Options) (*doc.Doc, error) {
	for _, item := range opt {
		if err := item.sanitize(); err != nil {
			return nil, err
		}
	}

	blocks := buildBlock(h, opt...)
	d := doc.New()
	wg := sync.WaitGroup{}

	for blk := range blocks {
		wg.Add(1)
		go func(b doc.Block) {
			parseBlock(d, &b, h)
			wg.Done()
		}(blk)
	}

	wg.Wait()

	if err := d.Sanitize(); err != nil {
		h.Error(message.Erro, err)
	}

	return d, nil
}

func parseBlock(d *doc.Doc, block *doc.Block, h *message.Handler) {
	if err := d.ParseBlock(block); err != nil {
		h.Error(message.Erro, err)
	}
}

// 分析源代码，获取注释块。
//
// 当所有的代码块已经放入 Block 之后，Block 会被关闭。
func buildBlock(h *message.Handler, opt ...*Options) chan doc.Block {
	data := make(chan doc.Block, 500)

	go func() {
		wg := &sync.WaitGroup{}
		for _, o := range opt {
			parseOptions(data, h, wg, o)
		}
		wg.Wait()

		close(data)
	}()

	return data
}

// 分析每个配置项对应的内容
func parseOptions(data chan doc.Block, h *message.Handler, wg *sync.WaitGroup, o *Options) {
	for _, path := range o.paths {
		wg.Add(1)
		go func(path string) {
			parseFile(data, h, path, o)
			wg.Done()
		}(path)
	}
}

// 分析 path 指向的文件。
//
// NOTE: parseFile 内部不能有协程处理代码。
func parseFile(channel chan doc.Block, h *message.Handler, path string, o *Options) {
	data, err := readFile(path, o.encoding)
	if err != nil {
		h.Error(message.Erro, message.WithError(path, "", 0, err))
		return
	}

	ret := lang.Parse(path, data, o.blocks, h)
	for line, data := range ret {
		channel <- doc.Block{
			File: path,
			Line: line,
			Data: data,
		}
	}
}

// 以指定的编码方式读取内容。
func readFile(path string, enc encoding.Encoding) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if enc == nil || enc == encoding.Nop {
		return data, nil
	}

	reader := transform.NewReader(bytes.NewReader(data), enc.NewDecoder())
	return ioutil.ReadAll(reader)
}
