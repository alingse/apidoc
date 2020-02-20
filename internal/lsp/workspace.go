// SPDX-License-Identifier: MIT

package lsp

import (
	"github.com/caixw/apidoc/v6/doc"
	"github.com/caixw/apidoc/v6/internal/locale"
	"github.com/caixw/apidoc/v6/internal/lsp/protocol"
)

// The workspace/workspaceFolders request is sent from the server to the client to fetch the current open
// list of workspace folders. Returns null in the response if only a single file is open in the tool.
// Returns an empty array if a workspace is open but no folders are configured.
func (s *server) workspaceWorkspaceFolders() error {
	var folders []protocol.WorkspaceFolder
	if err := s.Send("workspace/workspaceFolders", nil, &folders); err != nil {
		return err
	}

	for _, f := range s.folders {
		if err := f.close(); err != nil {
			return err
		}
	}

	if len(folders) != 0 {
		for _, wf := range folders {
			s.folders = append(s.folders, &folder{
				WorkspaceFolder: wf,
				doc:             doc.New(),
			})
		}
	}
	return nil
}

// workspace/didChangeWorkspaceFolders
//
// https://microsoft.github.io/language-server-protocol/specifications/specification-current/#workspace_didChangeWorkspaceFolders
func (s *server) workspaceDidChangeWorkspaceFolders(notify bool, in *protocol.DidChangeWorkspaceFoldersParams, out *interface{}) error {
	if s.getState() != serverInitialized {
		return newError(ErrInvalidRequest, locale.ErrInvalidLSPState)
	}

	for _, folder := range in.Event.Removed {
		for index, f2 := range s.folders {
			if f2.Name == folder.Name && f2.URI == folder.URI {
				if err := f2.close(); err != nil {
					return err
				}
				s.folders = append(s.folders[:index], s.folders[index+1:]...)
			}
		}
	}

	for _, wf := range in.Event.Added {
		s.folders = append(s.folders, &folder{
			WorkspaceFolder: wf,
			doc:             doc.New(),
		})
	}

	return nil
}
