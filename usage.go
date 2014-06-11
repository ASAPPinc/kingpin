package kingpin

import (
	"bytes"
	"fmt"
	"go/doc"
	"io"
	"strings"
)

func (c *Application) Usage(w io.Writer) {
	c.writeHelp(guessWidth(w), w)
}

func (c *Application) CommandUsage(w io.Writer, command string) {
	cmd, ok := c.commands[command]
	if !ok {
		Fatalf("unknown command '%s'", command)
	}
	s := []string{formatArgsAndFlags(c.Name, c.argGroup, c.flagGroup)}
	s = append(s, formatArgsAndFlags(cmd.name, cmd.argGroup, cmd.flagGroup))
	fmt.Fprintf(w, "usage: %s\n", strings.Join(s, " "))
	if cmd.help != "" {
		fmt.Fprintf(w, "\n%s\n", cmd.help)
	}
	cmd.writeHelp(guessWidth(w), w)
}

func (c *Application) writeHelp(width int, w io.Writer) {
	s := []string{formatArgsAndFlags(c.Name, c.argGroup, c.flagGroup)}
	if len(c.commands) > 0 {
		s = append(s, "<command>", "[<flags>]", "[<args> ...]")
	}

	helpSummary := ""
	if c.Help != "" {
		helpSummary = "\n\n" + c.Help
	}
	fmt.Fprintf(w, "usage: %s%s\n", strings.Join(s, " "), helpSummary)

	c.flagGroup.writeHelp(2, width, w)
	c.argGroup.writeHelp(2, width, w)

	if len(c.commands) > 0 {
		fmt.Fprintf(w, "\nCommands:\n")
		c.helpCommands(width, w)
	}
}

func (c *Application) helpCommands(width int, w io.Writer) {
	for _, cmd := range c.commandOrder {
		fmt.Fprintf(w, "  %s\n", formatArgsAndFlags(cmd.name, cmd.argGroup, cmd.flagGroup))
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, cmd.help, "", "", width-4)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		for _, line := range lines {
			fmt.Fprintf(w, "    %s\n", line)
		}
		fmt.Fprintf(w, "\n")
	}
}

func (f *flagGroup) writeHelp(indent, width int, w io.Writer) {
	if len(f.long) == 0 {
		return
	}

	fmt.Fprintf(w, "\nFlags:\n")
	l := 0
	for _, flag := range f.long {
		if fl := len(formatFlag(flag)); fl > l {
			l = fl
		}
	}

	l += 3 + indent

	indentStr := strings.Repeat(" ", l)

	for _, flag := range f.flagOrder {
		prefix := fmt.Sprintf("  %-*s", l-2, formatFlag(flag))
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, flag.help, "", "", width-l)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		fmt.Fprintf(w, "%s%s\n", prefix, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "%s%s\n", indentStr, line)
		}
	}
}

func (f *flagGroup) gatherFlagSummary() (out []string) {
	for _, flag := range f.flagOrder {
		if flag.required {
			fb, ok := flag.value.(boolFlag)
			if ok && fb.IsBoolFlag() {
				out = append(out, fmt.Sprintf("--%s", flag.name))
			} else {
				out = append(out, fmt.Sprintf("--%s=%s", flag.name, flag.formatPlaceHolder()))
			}
		}
	}
	if len(f.long) != len(out) {
		out = append(out, "[<flags>]")
	}
	return
}

func (a *argGroup) writeHelp(indent, width int, w io.Writer) {
	if len(a.args) == 0 {
		return
	}

	fmt.Fprintf(w, "\nArgs:\n")
	l := 0
	for _, arg := range a.args {
		if al := len(arg.name) + 2; al > l {
			l = al
			if !arg.required {
				l += 2
			}
		}
	}

	l += 3 + indent

	indentStr := strings.Repeat(" ", l)

	for _, arg := range a.args {
		argString := "<" + arg.name + ">"
		if !arg.required {
			argString = "[" + argString + "]"
		}
		prefix := fmt.Sprintf("  %-*s", l-2, argString)
		buf := bytes.NewBuffer(nil)
		doc.ToText(buf, arg.help, "", "", width-l)
		lines := strings.Split(strings.TrimRight(buf.String(), "\n"), "\n")
		fmt.Fprintf(w, "%s%s\n", prefix, lines[0])
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "%s%s\n", indentStr, line)
		}
	}

}

func (c *CmdClause) writeHelp(width int, w io.Writer) {
	c.flagGroup.writeHelp(2, width, w)
	c.argGroup.writeHelp(2, width, w)
}

func formatArgsAndFlags(name string, args *argGroup, flags *flagGroup) string {
	s := []string{name}
	s = append(s, flags.gatherFlagSummary()...)
	depth := 0
	for _, arg := range args.args {
		h := "<" + arg.name + ">"
		if !arg.required {
			h = "[" + h
			depth++
		}
		s = append(s, h)
	}
	s[len(s)-1] = s[len(s)-1] + strings.Repeat("]", depth)
	return strings.Join(s, " ")
}

func formatFlag(flag *FlagClause) string {
	flagString := ""
	if flag.shorthand != 0 {
		flagString += fmt.Sprintf("-%c, ", flag.shorthand)
	}
	flagString += fmt.Sprintf("--%s", flag.name)
	fb, ok := flag.value.(boolFlag)
	if !ok || !fb.IsBoolFlag() {
		flagString += fmt.Sprintf("=%s", flag.formatPlaceHolder())
	}
	return flagString
}
