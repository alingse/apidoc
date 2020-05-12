// SPDX-License-Identifier: MIT

package token

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/caixw/apidoc/v7/core"
)

type (
	anonymous struct {
		Attr1 intAttr `apidoc:"attr1,attr,usage"`
		Elem1 intTag  `apidoc:"elem1,elem,usage"`
	}

	intTag struct {
		BaseTag
		Value int `apidoc:"-"`
	}

	stringTag struct {
		BaseTag
		Value string `apidoc:"-"`
	}

	errTag struct {
		BaseTag
		Value int `apidoc:"-"`
	}

	// NOTE: objectTag 作为普通对象嵌套了 Decoder 等实例，本身不能实现这些接口。
	objectTag struct {
		BaseTag
		RootName struct{}  `apidoc:"apidoc,meta,usage-root"`
		ID       intAttr   `apidoc:"id,attr,usage"`
		Name     stringTag `apidoc:"name,elem,usage"`
	}

	intAttr struct {
		BaseAttribute
		Value int `apidoc:"-"`
	}

	stringAttr struct {
		BaseAttribute
		Value string `apidoc:"-"`
	}

	errAttr struct {
		BaseAttribute
		Value int `apidoc:"-"`
	}

	// NOTE: objectTag 作为普通对象嵌套了 Decoder 等实例，本身不能实现这些接口。
	objectAttr struct {
		BaseAttribute
		RootName struct{}  `apidoc:"apidoc,meta,usage-root"`
		ID       intTag    `apidoc:"id,attr,usage"`
		Name     stringTag `apidoc:"name,elem,usage"`
	}
)

var (
	_ AttrEncoder = &intAttr{}
	_ AttrDecoder = &intAttr{}
	_ Encoder     = &intTag{}
	_ Decoder     = &intTag{}

	_ AttrEncoder = &stringAttr{}
	_ AttrDecoder = &stringAttr{}
	_ Encoder     = &stringTag{}
	_ Decoder     = &stringTag{}

	_ Sanitizer = &objectTag{}
)

func (i *intTag) DecodeXML(p *Parser, start *StartElement) (*EndElement, error) {
	for {
		t, _, err := p.Token()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem, nil
			}
		case *String:
			v, err := strconv.Atoi(strings.TrimSpace(elem.Value))
			if err != nil {
				return nil, err
			}
			i.Value = v
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *intAttr) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	v, err := strconv.Atoi(strings.TrimSpace(attr.Value.Value))
	if err != nil {
		return p.WithError(attr.Value.Start, attr.Value.End, attr.Name.Value, err)
	}
	i.Value = v
	return nil
}

func (i *intTag) EncodeXML() (string, error) {
	return strconv.Itoa(i.Value), nil
}

func (i *intAttr) EncodeXMLAttr() (string, error) {
	return strconv.Itoa(i.Value), nil
}

func (i *stringTag) DecodeXML(p *Parser, start *StartElement) (*EndElement, error) {
	for {
		t, _, err := p.Token()
		if err == io.EOF {
			return nil, nil
		} else if err != nil {
			return nil, err
		}

		switch elem := t.(type) {
		case *EndElement:
			if elem.Name.Value == start.Name.Value {
				return elem, nil
			}
		case *String:
			i.Value = elem.Value
		default:
			panic(fmt.Sprintf("无效的类型 %s", reflect.TypeOf(t)))
		}
	}
}

func (i *stringAttr) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	i.Value = attr.Value.Value
	return nil
}

func (i *stringTag) EncodeXML() (string, error) {
	return i.Value, nil
}

func (i *stringAttr) EncodeXMLAttr() (string, error) {
	return i.Value, nil
}

func (o *objectTag) Sanitize(p *Parser) error {
	o.ID.Value++
	return nil
}

func (o *objectAttr) Sanitize(p *Parser) error {
	o.ID.Value++
	return nil
}

func (t *errTag) DecodeXML(p *Parser, start *StartElement) (core.Position, error) {
	return core.Position{}, errors.New("Decoder.DecodeXML")
}

func (t *errAttr) DecodeXMLAttr(p *Parser, attr *Attribute) error {
	return errors.New("AttrDecoder.DecodeXMLAttr")
}

func (t *errTag) EncodeXML() (string, error) {
	return "", errors.New("Encoder.EncodeXML")
}

func (t *errAttr) EncodeXMLAttr() (string, error) {
	return "", errors.New("AttrEncoder.EncodeXMLAttr")
}
