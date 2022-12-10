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
	"strings"

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

// completions are the completion things we can support.
//
// Since we only support completion of the functions within our
// standard-library they will not change, and we can calculate
// the complete list once and reuse it.
var completions []protocol.CompletionItem

// initialize is called to setup a new buffer.
func initialize(context *glsp.Context, params *protocol.InitializeParams) (any, error) {
	capabilities := handler.CreateServerCapabilities()

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

// lspStart launches our LSP server in the foreground, and doesn't return.
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

// textDocumentCompletion should return available completions.
//
// Since we only offer completion of the functions defined within
// our standard library we only calculate the (sorted) list once,
// and reuse it thereafter.
func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {

	// If we've already discovered our completions then return them.
	if len(completions) > 0 {
		return completions, nil
	}

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
	completions = make([]protocol.CompletionItem, len(keys))

	// We're only going to provide completion of things
	// which are functions.
	kind := protocol.CompletionItemKindFunction

	// Now we have a list of sorted things.
	for i, key := range keys {

		// Save the details in our global completions-array
		completions[i] = protocol.CompletionItem{
			Label:  key,
			Kind:   &kind,
			Detail: &key,
		}
	}

	// And return them.
	return completions, nil
}

// textDocumentHover is called when the client hovers over a token.
//
// We need to find out what text is being hovered over, and return
// something "useful" to the client.
func textDocumentHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {

	// Get the file the user is visiting.
	_uri, err := uri.Parse(params.TextDocument.URI)
	if err != nil {
		return nil, err
	}

	// open the file, and read the content
	var content []byte
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

	// Build up the arguments to the procedure.
	args := ""

	if len(prc.Args) > 0 {

		for _, arg := range prc.Args {
			args += " " + arg.ToString()
		}
		args = strings.TrimSpace(args)
		args = " (" + args + ")"
	}

	// The text we'll show - name args, and help.
	help := fmt.Sprintf("**%s%s**\n%s", token, args, prc.Help)

	return &protocol.Hover{
		Contents: protocol.MarkupContent{
			Kind:  protocol.MarkupKindMarkdown,
			Value: help,
		},
	}, nil
}
