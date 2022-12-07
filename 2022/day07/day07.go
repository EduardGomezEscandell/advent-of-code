// Package day07 solves day 7 of AoC2022.
package day07

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/array"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/channel"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/charray"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/fun"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/input"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/stack"
	"github.com/EduardGomezEscandell/AdventOfCode/2022/utils/tree"
)

const (
	today    = 7
	fileName = "input.txt"
)

type size int

// Part1 solves the first half of the problem.
func Part1(input <-chan input.Line) (int, error) {
	p := NewParser()

	for ln := range input {
		if err := ln.Err(); err != nil {
			return 0, fmt.Errorf("failed to read line %d: %v", p.LineNo+1, err)
		}
		p.NextLine(ln.Str())
	}
	channel.Exhaust(input) // Sync point

	dfsComputeSizes(&p.Fs.Root)

	fmt.Println(String(p.Fs))

	return int(dfsSumSmallDirs(&p.Fs.Root)), nil
}

// Part2 solves the second half of the problem.
func Part2(input <-chan input.Line) (int, error) {
	channel.Exhaust(input)
	return 0, nil
}

/// ---------- Tree exploration implementation ------------------
func dfsComputeSizes(root *FsNode) size {
	if !root.Value.IsDir {
		return root.Value.Size
	}

	ch := make(chan size)
	var wg sync.WaitGroup

	array.Foreach(root.Children, func(child **FsNode) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- dfsComputeSizes(*child)
		}()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()

	root.Value.Size = charray.Reduce(ch, fun.Add[size])
	return root.Value.Size
}

func dfsSumSmallDirs(root *FsNode) size {
	if !root.Value.IsDir {
		return 0
	}

	ch := make(chan size)
	var wg sync.WaitGroup

	array.Foreach(root.Children, func(child **FsNode) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch <- dfsSumSmallDirs(*child)
		}()
	})

	go func() {
		wg.Wait()
		close(ch)
	}()
	children := charray.Reduce(ch, fun.Add[size])
	if root.Value.Size > 100000 {
		return children
	}
	return children + root.Value.Size
}

/// ---------- Parser implementation ------------------

type Cmd int
type LineType int

type FsNode = tree.Node[FileDescriptor]
type Filesystem = tree.Tree[FileDescriptor]
type Path = stack.Stack[*FsNode]

type FileDescriptor struct {
	Name  string
	Size  size
	IsDir bool
}

// Parser is a finite state machine that transitions between
// states defined as Cmd. Each state can have some variables,
// contained in Args.
//
// It constructs a filesystem (stored as a tree), where each
// file or directory is a FileDescriptor (a node in the tree)
//
// To transition, a line must be fed into it.
type Parser struct {
	// Current command and its args
	Cmd  Cmd
	Args []string

	// Filesystem
	Fs      Filesystem
	cwdPath Path

	// For better error reporting
	LineNo int
}

func NewParser() Parser {
	var p Parser
	p.Fs.Root.Value.Name = "/"
	p.Fs.Root.Value.IsDir = true
	return p
}

const (
	BLANK Cmd = iota
	LS
	CD
	ERROR
)

func (c Cmd) String() string {
	switch c {
	case BLANK:
		return "blank state"
	case LS:
		return "command ls"
	case CD:
		return "command cd"
	case ERROR:
		return "error state"
	}
	panic("Unknown state")
}

const (
	CMD LineType = iota
	RESPONSE
	EMPTY
)

func (t LineType) String() string {
	switch t {
	case CMD:
		return "command"
	case RESPONSE:
		return "response"
	case EMPTY:
		return "empty"
	}
	panic("Unknown state")
}

func (p Parser) Cwd() *FsNode {
	return p.cwdPath.Peek()
}

func (p *Parser) NextLine(line string) (err error) {
	defer func() {
		if err == nil {
			return
		}
		path := array.Map(p.cwdPath.Data(), func(f *FsNode) string { return f.Value.Name })
		err = fmt.Errorf("error in line %d\nline: %s\ncwd:  %s\ncmd:  %s\nargs: %s\nerr: %v", p.LineNo, line, strings.Join(path, "/"), p.Cmd, p.Args, err)
	}()
	p.LineNo++
	lineType, args := lineType(line)

	switch lineType {
	case EMPTY:
		p.Cmd = BLANK
		p.Args = nil
		return nil
	case CMD:
		return p.evalCommand(args[0], args[1:]...)
	case RESPONSE:
		return p.evalResponse(args...)
	}

	panic("Unreachable")
}

func (p *Parser) evalResponse(args ...string) (err error) {
	defer func() {
		if err == nil {
			return
		}
		err = fmt.Errorf("parsing response: %v", err)
	}()
	return evalResponseMap[p.Cmd](p, args...)
}

func (p *Parser) evalCommand(verb string, args ...string) (err error) {
	defer func() {
		if err == nil {
			return
		}
		err = fmt.Errorf("parsing command: %v", err)
	}()

	eval, found := evalCmdMap[verb]
	if !found {
		verbs := []string{}
		for v := range evalCmdMap {
			verbs = append(verbs, v)
		}
		return fmt.Errorf("unknown verb %q.\nExpected any of: %s", verb, strings.Join(verbs, ", "))
	}

	return eval(p, args...)
}

func lineType(line string) (ld LineType, args []string) {
	data := strings.Split(line, " ")
	if len(data) == 0 {
		return EMPTY, nil
	}
	if data[0] != "$" {
		return RESPONSE, data
	}
	return CMD, data[1:]
}

var evalCmdMap = map[string]func(*Parser, ...string) error{
	"ls": func(p *Parser, args ...string) error {
		p.Args = nil
		if len(args) != 0 {
			p.Cmd = ERROR
			return fmt.Errorf("unexpected token afer ls: ls accepts no arguments.")
		}
		p.Cmd = LS
		return nil
	},
	"cd": func(p *Parser, args ...string) error {
		if len(args) == 0 {
			p.Cmd = ERROR
			return fmt.Errorf("unexpected endline after 'cd': cd needs exactly one argument")
		}
		if len(args) > 1 {
			p.Cmd = ERROR
			return fmt.Errorf("unexpected token after %q: cd needs exactly one argument", args[0])
		}
		return p.cd(args[0])
	},
}

var evalResponseMap = map[Cmd]func(*Parser, ...string) error{
	BLANK: func(*Parser, ...string) error { return errors.New("unexpected response blank state") },
	CD:    func(*Parser, ...string) error { return errors.New("unexpected response to command cd") },
	ERROR: func(*Parser, ...string) error { return errors.New("unexpected response after error state") },
	LS: func(p *Parser, fields ...string) error {
		if len(fields) != 2 {
			return fmt.Errorf("expected two tokens per line, found %d", len(fields))
		}
		// Is it a directory?
		if fields[0] == "dir" {
			return p.mkdir(fields[1])
		}
		// Is it a file?
		sz, err := strconv.Atoi(fields[0])
		if err != nil {
			return fmt.Errorf("Unrecognized first token. Should be either 'dir' or an integer")
		}
		return p.touch(fields[1], size(sz))
	},
}

// Filesystem implementation

func String(fs Filesystem) string {
	return str(fs.Root, "")
}

func str(fn FsNode, indent string) string {
	if !fn.Value.IsDir {
		return fmt.Sprintf("%s- %s (file, size=%d)", indent, fn.Value.Name, fn.Value.Size)
	}
	s := fmt.Sprintf("%s- %s (dir)", indent, fn.Value.Name)
	children := strings.Join(array.Map(fn.Children, func(child *FsNode) string { return str(*child, indent+"  ") }), "\n")
	if len(children) == 0 {
		return s
	}
	return s + "\n" + children
}

func getChild(dir FsNode, filename string) int {
	return array.FindIf(dir.Children, func(f *FsNode) bool { return f.Value.Name == filename })
}

func (p *Parser) mkdir(name string) error {
	idx := getChild(*p.Cwd(), name)
	if idx != -1 {
		original := p.Cwd().Children[idx]
		if !original.Value.IsDir {
			return fmt.Errorf("Attempting to overwrite file %q as a directory", name)
		}
	}

	p.Cwd().Children = append(p.Cwd().Children, &FsNode{
		Value: FileDescriptor{
			Name:  name,
			IsDir: true,
		},
	})

	return nil
}

func (p *Parser) touch(name string, sz size) error {
	idx := getChild(*p.Cwd(), name)
	if idx != -1 {
		original := p.Cwd().Children[idx]
		if original.Value.IsDir {
			return fmt.Errorf("Attempting to overwrite directory %q as a file", name)
		}
		if sz != original.Value.Size {
			return fmt.Errorf("Attempting to overwrite file %q with diferent size:\n\toriginal:%d\n\tnew:     %d", name, original.Value.Size, sz)
		}
		return nil
	}
	p.Cwd().Children = append(p.Cwd().Children, &FsNode{
		Value: FileDescriptor{
			Name: name,
			Size: sz,
		},
	})
	return nil
}

func (p *Parser) cd(target string) error {
	// cd /
	if target == "/" {
		p.cwdPath = Path{}
		p.cwdPath.Push(&p.Fs.Root)
		return nil
	}

	// cd ..
	if target == ".." {
		if p.cwdPath.Size() <= 1 {
			return fmt.Errorf("attempting to cd .. from root")
		}
		p.cwdPath.Pop()
		return nil
	}

	// cd directory
	idx := getChild(*p.Cwd(), target)
	if idx == -1 {
		p.Cmd = ERROR
		return fmt.Errorf("cd: directory %q does not exist.", target)
	}
	child := p.Cwd().Children[idx]
	if !child.Value.IsDir {
		p.Cmd = ERROR
		return fmt.Errorf("cd: file %q is not a directory.", target)
	}
	p.cwdPath.Push(child)
	return nil
}

/// ---------- Here be boilerplate ------------------

type problemResult struct {
	id  int
	res string
	err error
}

// Main is the entry point to today's problem solution.
func Main(stdout io.Writer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	reader, err := DataReader()
	if err != nil {
		return err
	}

	ch, err := input.ReadDataAsync(ctx, reader, 3) // Reading in groups of three to optimize 2nd part
	if err != nil {
		return err
	}

	channels := channel.Split(ctx, ch, 2)

	resultCh := make(chan problemResult)
	go func() {
		result, err := Part1(channels[0])
		channel.Exhaust(channels[0])
		if err != nil {
			resultCh <- problemResult{0, "", err}
			cancel()
		}
		resultCh <- problemResult{0, fmt.Sprintf("Result of part 1: %v", result), nil}
	}()

	go func() {
		result, err := Part2(channels[1])
		channel.Exhaust(channels[1])
		if err != nil {
			resultCh <- problemResult{1, "", err}
			cancel()
		}
		resultCh <- problemResult{1, fmt.Sprintf("Result of part 2: %v", result), nil}
	}()

	var results [2]problemResult
	for i := 0; i < 2; i++ {
		r := <-resultCh
		results[r.id] = r
	}

	for _, v := range results {
		if v.err != nil {
			return v.err
		}
		fmt.Fprintln(stdout, v.res)
	}

	return nil
}

// DataReader is a wrapper around input.DataFile made to be
// easily mocked.
var DataReader = func() (r io.ReadCloser, e error) {
	f, e := input.DataFile(today, fileName)
	if e != nil {
		return nil, e
	}
	return f, nil
}
