// This file provides our LSP support.
//
// When yal is invoked with the "-lsp" flag we call lspStart(),
// which provides simple completion and hover support.
//

package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/skx/yal/primitive"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"
	"go.lsp.dev/uri"

	// Must include a backend implementation. See kutil's logging/ for other options.
	_ "github.com/tliron/kutil/logging/simple"
)

// lsName contains the name of our LSP handler
const lsName = "yal"

// handler contains the pointer to our handler
var handler protocol.Handler

// initialize is called to setup a new buffer.
func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    "yal",
			Version: &version,
		},
	}, nil
}

func lspStart() {
	logging.Configure(1, nil)

	handler = protocol.Handler{
		// generic
		Initialize: initialize,

		// Yal Specific
		TextDocumentCompletion: textDocumentCompletion,
		TextDocumentHover:      textDocumentHover,
	}

	server := server.NewServer(&handler, lsName, false)

	err := server.RunStdio()
	if err != nil {
		fmt.Printf("Internal Error Running LSP Process\n%s", err)
		os.Exit(1)
	}
}

// textDocumentCompletion should return appropriate completions.
//
// However we just always return a list of all known functions, the
// client can sort it out.
func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {

	// Build up a list of all things known in the environment
	keys := []string{}

	// Save the known "things", because we want show them in sorted-order.
	items := ENV.Items()
	for k := range items {
		keys = append(keys, k)
	}

	// sort the known-things (i.e. environment keys)
	sort.Strings(keys)

	// Create the return value
	out := make([]protocol.CompletionItem, len(keys))

	// The kind of completion we have
	kind := protocol.CompletionItemKindFunction

	// Now we have a list of sorted things.
	for i, key := range keys {

		out[i] = protocol.CompletionItem{
			Label:  key,
			Kind:   &kind,
			Detail: &key,
		}
	}
	return out, nil
}

// textDocumentHover is called when the client hovers over a token
//
// We need to find out what text is being hovered over, and return
// something "useful" to the client.
func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {

	// Get the file
	_uri, err := uri.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// open the file
	var content []byte

	// read the content
	content, err = os.ReadFile(_uri.Filename())
	if err != nil {
		return nil, err
	}

	// We'll build up the current line, being hovered on here
	var curLine uint32
	line := ""

	// count the newlines to get the current line.
	for _, chr := range content {
		if chr == '\n' {
			curLine++
			continue
		}
		if curLine == params.Position.Line {
			line += string(chr)
		}
	}

	// current line is empty?  Then abort
	if line == "" {
		return nil, nil
	}

	// Right now we have the line we want the token
	//
	// Assume we have a line like "(this is (cake))"
	// and position points to the "c" we want to have the
	// whole token
	//
	token := ""

	for i, chr := range line {
		if chr == rune(' ') || chr == rune('(') || chr == rune(')') || chr == rune('\t') {
			if uint32(i) > params.Position.Character {
				break
			}
			token = ""
			continue
		}
		token += string(chr)
	}

	// Find the details of the function, if we can
	info, ok := ENV.Get(token)
	if !ok {
		return nil, nil
	}

	// Is it a procedure?
	prc, ok2 := info.(*primitive.Procedure)
	if !ok2 {
		return nil, nil
	}

	// The text we'll show.
	help := fmt.Sprintf("**%s**\n%s", token, prc.Help)

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: help,
		},
	}, nil
}
