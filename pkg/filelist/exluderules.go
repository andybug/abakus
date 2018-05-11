// Copyright © 2018 Andrew Fields <andy@andybug.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package filelist

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	sll "github.com/emirpasic/gods/lists/singlylinkedlist"
	"github.com/emirpasic/gods/stacks/arraystack"
	"gopkg.in/yaml.v2"
)

// IGNORE_FILE is the name of the file that can be in each directory
// where the user can specify what files to exclude
const IGNORE_FILE = ".abakusignore"

// ignoreFile defines the abakus ignorefile format
// version must be 1
// excludes is a list of rules (like .gititgnore)
type ignoreFile struct {
	Version  uint32
	Excludes []string
}

// readRules returns the exclude rules for a directory
// if there is an abakus ignore file, it is read and the rules added
// if not, an empty rule object is returned
func readRules(dir string) (*excludeRules, error) {
	rules := newExcludeRules(dir)
	ignoreFilePath := filepath.Join(dir, IGNORE_FILE)

	file, err := os.Open(ignoreFilePath)
	defer file.Close()
	if err != nil {
		// if no abakus ignore file, return empty rules
		if os.IsNotExist(err) {
			return rules, nil
		}
		// otherwise, there's something wrong
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	f := ignoreFile{}
	err = yaml.Unmarshal(bytes, &f)
	if err != nil {
		return nil, err
	}

	if f.Version != 1 {
		errMsg := fmt.Sprintf("Ignore file version %u not supported: %s",
			f.Version, ignoreFilePath)
		return nil, errors.New(errMsg)
	}

	for _, rule := range f.Excludes {
		rules.add(rule)
	}

	return rules, nil
}

// excludeRules defines the list of exclusion rules added at a path
// in the tree (in the IGNORE_FILE for that dir)
type excludeRules struct {
	path  string
	rules *sll.List
}

// newExcludeRules returns an empty excludeRules object
func newExcludeRules(path string) *excludeRules {
	return &excludeRules{
		path:  path,
		rules: sll.New(),
	}
}

// add converts the given rule to a regex then adds it to the list
// for this dir
func (er *excludeRules) add(rule string) {
	if rule[0] == '/' {
		// only match file in this dir
		rule = fmt.Sprintf("^%s$", filepath.Join(er.path, rule[1:]))
	} else {
		// match files in any dir under this point
		rule = fmt.Sprintf("^.*/%s$", rule)
	}
	er.rules.Add(rule)
}

// exclude returns true if the given absolute path to a file matches
// one of the rules (so should be dropped)
func (er *excludeRules) exclude(fileName string) bool {
	it := er.rules.Iterator()
	for it.Next() {
		rule := it.Value().(string)
		matched, _ := regexp.MatchString(rule, fileName)
		if matched {
			return true
		}
	}

	return false
}

// excludeRulesStack tracks the exclude rules for each dir
// as it is visited. each dir should have rules pushed on
// to the stack when entered and popped when left
type excludeRulesStack struct {
	stack *arraystack.Stack
}

// newExcludeRulesStack returns empty stack
func newExcludeRulesStack() *excludeRulesStack {
	return &excludeRulesStack{
		stack: arraystack.New(),
	}
}

// push adds the rules to the top of the stack
func (ers *excludeRulesStack) push(rules *excludeRules) {
	ers.stack.Push(rules)
}

// pop removes rules from the top of the stack
func (ers *excludeRulesStack) pop() {
	ers.stack.Pop()
}

// exclude checks the absolute path to the file against all of the
// rules in the stack. returns true if the file should be excluded
func (ers *excludeRulesStack) exclude(fileName string) bool {
	it := ers.stack.Iterator()
	for it.Next() {
		rule := it.Value().(*excludeRules)
		if rule.exclude(fileName) {
			return true
		}
	}

	return false
}
