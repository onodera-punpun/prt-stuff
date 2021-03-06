package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/go2c/optparse"
	colorful "github.com/lucasb-eyer/go-colorful"
	"github.com/onodera-punpun/prt/ports"
)

// chartCommand generates a dependency grap.
// TODO: Replace current port name "." with the actual name.
func chartCommand(input []string) error {
	// Define valid arguments.
	o := optparse.New()
	//argd := o.Bool("duplicate", 'd', false)
	argn := o.Bool("no-alias", 'n', false)
	argt := o.String("type", 't', "svg")
	argh := o.Bool("help", 'h', false)

	// Parse arguments.
	_, err := o.Parse(input)
	if err != nil {
		return fmt.Errorf("invaild argument, use `-h` for a list of arguments")
	}

	// Print help.
	if *argh {
		fmt.Println("Usage: prt chart [arguments]")
		fmt.Println("")
		fmt.Println("arguments:")
		fmt.Println("  -d,   --duplicate       display duplicates as well")
		fmt.Println("  -n,   --no-alias        disable aliasing")
		fmt.Println("  -t,   --type            filetype to use")
		fmt.Println("  -h,   --help            print help and exit")

		return nil
	}

	p := ports.New(".")
	if err := p.Pkgfile.Parse(); err != nil {
		return err
	}

	// Get all ports.
	all, err := ports.All()
	if err != nil {
		return err
	}

	if err := p.ParseDepends(all, !*argn); err != nil {
		return err
	}

	// Set file to write to.
	f, err := os.OpenFile(p.Pkgfile.Name+".dot", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("could not create `%s`", p.Pkgfile.Name+".dot")
	}
	defer f.Close()

	// Prettify chart.
	fmt.Fprintf(f, "digraph G {\n")
	fmt.Fprintf(f, "\tgraph [\n")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "tcenter", "true")
	fmt.Fprintf(f, "\t\t%s=\"%f\"\n", "pad", 2.0)
	fmt.Fprintf(f, "\t]\n\n")

	// Prettify nodes.
	fmt.Fprint(f, "\tnode [\n")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "constraint", "false")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "fontcolor", "#111e38")
	fmt.Fprintf(f, "\t\t%s=\"%d\"\n", "penwidth", 3)
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "shape", "box")
	fmt.Fprintf(f, "\t]\n\n")

	// Prettify edges.
	fmt.Fprintf(f, "\tedge [\n")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "arrowhead", "dot")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "color", "#cee0e3")
	fmt.Fprintf(f, "\t\t%s=\"%s\"\n", "headport", "n")
	fmt.Fprintf(f, "\t\t%s=\"%d\"\n", "penwidth", 2)
	fmt.Fprintf(f, "\t]\n\n")

	pal, _ := colorful.SoftPalette(128)
	chartRecurse(&p, 0, f, pal)

	fmt.Fprintf(f, "}")

	f.Close()
	if *argt == "dot" {
		return nil
	}

	// Convert to chart.
	cmd := exec.Command("dot", p.Pkgfile.Name+".dot", "-T", *argt, "-o", p.
		Pkgfile.Name+"."+*argt)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("something went wrong with GrapViz")
	}

	// Remove dot file.
	if err := os.Remove(p.Pkgfile.Name + ".dot"); err != nil {
		return err
	}

	return nil
}

var chartCheck []*ports.Port

func chartRecurse(p *ports.Port, l int, f *os.File, pal []colorful.Color) {
outer:
	for _, d := range p.Depends {
		// Continue if already checked.
		for _, c := range chartCheck {
			if c.Pkgfile.Name == d.Pkgfile.Name {
				continue outer
			}
		}
		chartCheck = append(chartCheck, d)

		fmt.Fprintf(f, "\tnode [color=\"%s\"]\n", pal[l].Hex())
		fmt.Fprintf(f, "\t\"%s\"->\"%s\"\n", p.Location.Base(), d.Location.
			Base())

		chartRecurse(d, l+1, f, pal)
	}
}
