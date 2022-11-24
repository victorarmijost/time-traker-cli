package repl

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"
)

func (c *Handler) PrintList(list []string) {
	c.PrintMessage(SprintList(list))
}

func SprintList(list []string) string {
	buf := bytes.NewBufferString("")

	if len(list) == 0 {
		fmt.Fprint(buf, "Nothing to show!")
	}

	for i, l := range list {
		if i == len(list)-1 {
			fmt.Fprintf(buf, "%d. %s", i+1, l)
		} else {
			fmt.Fprintf(buf, "%d. %s\n", i+1, l)
		}
	}

	return buf.String()
}

func (c *Handler) PrintMap(m map[string]string) {
	c.PrintMessage(SprintMap(m))
}

func SprintMap(m map[string]string) string {
	buf := bytes.NewBufferString("")

	if len(m) == 0 {
		fmt.Fprint(buf, "Nothing to show!")
	}

	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for i, k := range keys {
		if i == len(m)-1 {
			fmt.Fprintf(buf, "- %s: %s", k, m[k])
		} else {
			fmt.Fprintf(buf, "- %s: %s\n", k, m[k])
		}
		i++
	}

	return buf.String()
}

func (c *Handler) GetPass(prompt string) string {
	fmt.Fprintf(c.terminal, "- %s: ", prompt)
	pass, err := term.ReadPassword(0)
	c.Br()
	c.Br()

	if err != nil {
		return ""
	}

	return string(pass)
}

func (c *Handler) PrintError(err error) {
	c.PrintErrorMsg(err.Error())
}

func (c *Handler) PrintErrorMsg(msg string) {
	c.PrintMessage(c.colorString(fmt.Sprintf("<< ERROR: %s >>", msg), RED))
}

func (c *Handler) PrintMessage(msg string) {
	fmt.Fprintln(c.terminal, msg)
	c.Br()
}

func (c *Handler) PrintInfoMessage(msg string) {
	lines := strings.Split(msg, "\n")
	if len(lines) == 1 {
		c.PrintMessage(c.colorString(fmt.Sprintf("[%s] **** %s ****", time.Now().Format("2006-01-02 15:04:05"), msg), BLUE))
	} else {
		c.PrintHighightedMessage("Result")
		for _, l := range lines {
			fmt.Fprintln(c.terminal, l)
		}
		c.Br()
	}
}

func (c *Handler) PrintTitle(msg string) {
	fmt.Fprintln(c.terminal, c.colorString(strings.Repeat("*", len(msg)+6), CYAN))
	fmt.Fprintf(c.terminal, c.colorString("** %s **\n", CYAN), msg)
	fmt.Fprintln(c.terminal, c.colorString(strings.Repeat("*", len(msg)+6), CYAN))
	c.Br()
}

func (c *Handler) GetInput(msg string) string {
	promptLock = true
	c.setTermPrompt(fmt.Sprintf("- %s: ", msg))
	r := c.getInput()
	c.Br()
	promptLock = false
	return r
}

type Selectable interface {
	GetElement(int) string
	Size() int
}

func (c *Handler) SelectFromList(l Selectable) int {
	c.PrintHighightedMessage("Selection list")

	if l.Size() == 0 {
		c.PrintMessage("No results!")
		return -1
	}

	selectMap := map[string]int{}
	for i := 0; i < l.Size(); i++ {
		fmt.Fprintf(c.terminal, "%d. %s\n", i+1, l.GetElement(i))
		selectMap[strconv.Itoa(i+1)] = i
	}

	c.Br()

	sid := c.GetInput(fmt.Sprintf("Select (1 - %d, q to quit)", l.Size()))

	if sid == "q" {
		c.PrintMessage("Canceled!")
		return -1
	}

	rid, ok := selectMap[sid]

	if !ok {
		c.PrintErrorMsg("Wrong id")
		return -1
	}

	return rid
}

type Searchable interface {
	Selectable
	Match(int, string) bool
}

func (c *Handler) PrintHighightedMessage(message string) {
	fmt.Fprintln(c.terminal, c.colorString(message, BLUE))
	fmt.Fprint(c.terminal, c.colorString(strings.Repeat("=", len(message)), BLUE))
	c.Br()
}

func (c *Handler) SearchFromList(prompt string, l Searchable) int {
	input := c.GetInput(prompt)

	c.PrintHighightedMessage("Search results")

	selectMap := map[string]int{}
	pi := 0
	for i := 0; i < l.Size(); i++ {
		if l.Match(i, input) {
			pi++
			fmt.Fprintf(c.terminal, "%d. %s\n", pi, l.GetElement(i))
			selectMap[strconv.Itoa(pi)] = i
		}
	}

	if pi == 0 {
		c.PrintMessage("No results!")
		return -1
	}

	c.Br()

	sid := c.GetInput(fmt.Sprintf("Select (1 - %d, q to quit)", pi))

	if sid == "q" {
		c.PrintMessage("Canceled!")
		return -1
	}

	rid, ok := selectMap[sid]

	if !ok {
		c.PrintErrorMsg("Wrong id")
		return -1
	}

	return rid
}

func (c *Handler) Br() {
	fmt.Fprintln(c.terminal)
}
