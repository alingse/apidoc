// SPDX-License-Identifier: MIT

package lsp

import (
	"strings"

	"github.com/caixw/apidoc/v7/core"
	"github.com/caixw/apidoc/v7/internal/locale"
	"github.com/caixw/apidoc/v7/internal/lsp/protocol"
	"github.com/caixw/apidoc/v7/internal/lsp/search"
)

// textDocument/didChange
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_didChange
func (s *server) textDocumentDidChange(notify bool, in *protocol.DidChangeTextDocumentParams, out *interface{}) error {
	var f *folder
	for _, f = range s.folders {
		if strings.HasPrefix(string(in.TextDocument.URI), string(f.URI)) {
			break
		}
	}
	if f == nil {
		return newError(ErrInvalidRequest, locale.ErrFileNotFound, in.TextDocument.URI)
	}

	if !f.deleteURI(in.TextDocument.URI) {
		return nil
	}

	for _, blk := range in.Blocks() {
		f.parseBlock(blk)
	}
	f.srv.textDocumentPublishDiagnostics(f, in.TextDocument.URI)
	return nil
}

// textDocument/hover
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_hover
func (s *server) textDocumentHover(notify bool, in *protocol.HoverParams, out *protocol.Hover) error {
	for _, f := range s.folders {
		if search.Hover(f.doc, in.TextDocument.URI, in.TextDocumentPositionParams.Position, out) {
			return nil
		}
	}
	return nil
}

// textDocument/publishDiagnostics
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#textDocument_publishDiagnostics
func (s *server) textDocumentPublishDiagnostics(f *folder, uri core.URI) error {
	if s.clientCapabilities.TextDocument.PublishDiagnostics.RelatedInformation == false {
		return nil
	}

	p := &protocol.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: make([]protocol.Diagnostic, 0, len(f.errors)+len(f.warns)),
	}

	for _, err := range f.errors {
		if err.Location.URI == uri {
			p.Diagnostics = append(p.Diagnostics, protocol.Diagnostic{
				Range:    err.Location.Range,
				Message:  err.Error(),
				Severity: protocol.DiagnosticSeverityError,
			})
		}
	}

	for _, err := range f.warns {
		if err.Location.URI == uri {
			p.Diagnostics = append(p.Diagnostics, protocol.Diagnostic{
				Range:    err.Location.Range,
				Message:  err.Error(),
				Severity: protocol.DiagnosticSeverityWarning,
			})
		}
	}

	return s.Notify("textDocument/publishDiagnostics", p)
}
