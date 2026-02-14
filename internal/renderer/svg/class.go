// Package svg implements SVG output rendering for PlantUML diagrams.
package svg

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/bobcob7/godot-uml/internal/ast"
	"github.com/bobcob7/godot-uml/internal/font"
	"github.com/bobcob7/godot-uml/internal/layout"
	"github.com/bobcob7/godot-uml/internal/theme"
)

const (
	cornerRadius     = 8
	compartmentGap   = 1
	stereotypeFontPx = 11
	visibilityWidth  = 14
	diagramPadding   = 20
)

// ClassRenderer renders class diagrams to SVG.
type ClassRenderer struct {
	resolver *theme.Resolver
}

// NewClassRenderer creates a renderer with the given theme resolver.
// If resolver is nil, Darcula defaults are used.
func NewClassRenderer(resolver *theme.Resolver) *ClassRenderer {
	if resolver == nil {
		resolver = theme.NewResolver(nil)
	}
	return &ClassRenderer{resolver: resolver}
}

// classBox holds measured dimensions and content for a class-like element.
type classBox struct {
	id         string
	name       string
	stereotype string
	abstract   bool
	kind       string // "class", "interface", "enum"
	fields     []memberLine
	methods    []memberLine
	width      float64
	height     float64
	nameH      float64
	fieldsH    float64
	methodsH   float64
}

type memberLine struct {
	visibility ast.Visibility
	modifier   ast.Modifier
	text       string
}

// noteBox holds a positioned note.
type noteBox struct {
	target string
	text   string
	left   bool
	width  float64
	height float64
}

// packageBox holds a positioned package.
type packageBox struct {
	name       string
	children   []string
	x, y, w, h float64
}

// Render produces SVG output for a class diagram.
func (r *ClassRenderer) Render(w io.Writer, diagram *ast.Diagram) error {
	fontSize := r.resolver.ResolveInt("ClassFontSize", 13)
	fontSizeF := float64(fontSize)
	padding := r.resolver.ResolveInt("ClassPadding", 10)
	paddingF := float64(padding)
	for _, stmt := range diagram.Statements {
		if sp, ok := stmt.(*ast.Skinparam); ok {
			r.resolver.SetSkinparam(sp.Name, sp.Value)
		}
	}
	var boxes []*classBox
	var rels []*ast.Relationship
	var notes []*noteBox
	var pkgs []*packageBox
	boxByName := map[string]*classBox{}
	for _, stmt := range diagram.Statements {
		switch s := stmt.(type) {
		case *ast.ClassDef:
			b := r.measureClass(s, fontSizeF, paddingF)
			boxes = append(boxes, b)
			boxByName[b.id] = b
		case *ast.InterfaceDef:
			b := r.measureInterface(s, fontSizeF, paddingF)
			boxes = append(boxes, b)
			boxByName[b.id] = b
		case *ast.EnumDef:
			b := r.measureEnum(s, fontSizeF, paddingF)
			boxes = append(boxes, b)
			boxByName[b.id] = b
		case *ast.Relationship:
			rels = append(rels, s)
			for _, name := range []string{s.Left, s.Right} {
				if _, exists := boxByName[name]; !exists && name != "" {
					b := r.measureImplicitClass(name, fontSizeF, paddingF)
					boxes = append(boxes, b)
					boxByName[b.id] = b
				}
			}
		case *ast.Note:
			nb := r.measureNote(s, fontSizeF, paddingF)
			notes = append(notes, nb)
		case *ast.Package:
			pb := &packageBox{name: s.Name}
			for _, child := range s.Statements {
				switch c := child.(type) {
				case *ast.ClassDef:
					b := r.measureClass(c, fontSizeF, paddingF)
					boxes = append(boxes, b)
					boxByName[b.id] = b
					pb.children = append(pb.children, b.id)
				case *ast.InterfaceDef:
					b := r.measureInterface(c, fontSizeF, paddingF)
					boxes = append(boxes, b)
					boxByName[b.id] = b
					pb.children = append(pb.children, b.id)
				case *ast.EnumDef:
					b := r.measureEnum(c, fontSizeF, paddingF)
					boxes = append(boxes, b)
					boxByName[b.id] = b
					pb.children = append(pb.children, b.id)
				}
			}
			pkgs = append(pkgs, pb)
		}
	}
	if len(boxes) == 0 {
		return r.writeEmptyDiagram(w)
	}
	g := &layout.Graph{}
	nodeByID := map[string]*layout.Node{}
	for _, b := range boxes {
		n := &layout.Node{ID: b.id, Width: b.width, Height: b.height}
		g.Nodes = append(g.Nodes, n)
		nodeByID[b.id] = n
	}
	for _, rel := range rels {
		if rel.Left != "" && rel.Right != "" {
			g.Edges = append(g.Edges, &layout.Edge{From: rel.Left, To: rel.Right, Label: rel.Label})
		}
	}
	layout.Layout(g, layout.DefaultOptions())
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64
	for _, n := range g.Nodes {
		if n.Virtual {
			continue
		}
		if n.X < minX {
			minX = n.X
		}
		if n.Y < minY {
			minY = n.Y
		}
		if n.X+n.Width > maxX {
			maxX = n.X + n.Width
		}
		if n.Y+n.Height > maxY {
			maxY = n.Y + n.Height
		}
	}
	noteOffset := 160.0
	for _, nb := range notes {
		if nb.target == "" {
			continue
		}
		if n, ok := nodeByID[nb.target]; ok {
			var noteX float64
			if nb.left {
				noteX = n.X - noteOffset
			} else {
				noteX = n.X + n.Width + 20
			}
			noteRight := noteX + nb.width
			if noteX < minX {
				minX = noteX
			}
			if noteRight > maxX {
				maxX = noteRight
			}
			noteY := n.Y
			noteBottom := noteY + nb.height
			if noteY < minY {
				minY = noteY
			}
			if noteBottom > maxY {
				maxY = noteBottom
			}
		}
	}
	for _, pb := range pkgs {
		r.computePackageBounds(pb, nodeByID, paddingF)
		if pb.x < minX {
			minX = pb.x
		}
		if pb.y < minY {
			minY = pb.y
		}
		if pb.x+pb.w > maxX {
			maxX = pb.x + pb.w
		}
		if pb.y+pb.h > maxY {
			maxY = pb.y + pb.h
		}
	}
	offsetX := -minX + diagramPadding
	offsetY := -minY + diagramPadding
	svgW := int(maxX - minX + 2*diagramPadding)
	svgH := int(maxY - minY + 2*diagramPadding)
	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">`, svgW, svgH, svgW, svgH)
	sb.WriteString("\n")
	bgColor := r.resolver.ResolveColor("BackgroundColor")
	fmt.Fprintf(&sb, `<rect width="%d" height="%d" fill="%s"/>`, svgW, svgH, bgColor)
	sb.WriteString("\n")
	for _, pb := range pkgs {
		r.renderPackage(&sb, pb, offsetX, offsetY, fontSizeF)
	}
	for _, rel := range rels {
		fromNode := nodeByID[rel.Left]
		toNode := nodeByID[rel.Right]
		if fromNode == nil || toNode == nil {
			continue
		}
		r.renderRelationship(&sb, rel, fromNode, toNode, offsetX, offsetY, fontSizeF)
	}
	for _, b := range boxes {
		n := nodeByID[b.id]
		if n == nil || n.Virtual {
			continue
		}
		r.renderClassBox(&sb, b, n.X+offsetX, n.Y+offsetY, fontSizeF, paddingF)
	}
	for _, nb := range notes {
		targetNode := nodeByID[nb.target]
		if targetNode == nil {
			continue
		}
		var noteX float64
		if nb.left {
			noteX = targetNode.X - noteOffset + offsetX
		} else {
			noteX = targetNode.X + targetNode.Width + 20 + offsetX
		}
		noteY := targetNode.Y + offsetY
		r.renderNote(&sb, nb, noteX, noteY, fontSizeF)
		var lineFromX, lineToX float64
		if nb.left {
			lineFromX = noteX + nb.width
			lineToX = targetNode.X + offsetX
		} else {
			lineFromX = noteX
			lineToX = targetNode.X + targetNode.Width + offsetX
		}
		lineY := noteY + nb.height/2
		arrowColor := r.resolver.ResolveColor("ArrowColor")
		fmt.Fprintf(&sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-dasharray="5,5"/>`,
			lineFromX, lineY, lineToX, lineY, arrowColor)
		sb.WriteString("\n")
	}
	sb.WriteString("</svg>\n")
	_, err := io.WriteString(w, sb.String())
	return err
}

func (r *ClassRenderer) writeEmptyDiagram(w io.Writer) error {
	bgColor := r.resolver.ResolveColor("BackgroundColor")
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100" viewBox="0 0 100 100">
<rect width="100" height="100" fill="%s"/>
</svg>
`, bgColor)
	_, err := io.WriteString(w, svg)
	return err
}

func (r *ClassRenderer) measureClass(cd *ast.ClassDef, fontSize, padding float64) *classBox {
	b := &classBox{
		id:         cd.Name,
		name:       cd.Name,
		stereotype: cd.Stereotype,
		abstract:   cd.Abstract,
		kind:       "class",
	}
	r.measureMembers(b, cd.Members, fontSize, padding)
	return b
}

func (r *ClassRenderer) measureInterface(id *ast.InterfaceDef, fontSize, padding float64) *classBox {
	b := &classBox{
		id:         id.Name,
		name:       id.Name,
		stereotype: id.Stereotype,
		kind:       "interface",
	}
	r.measureMembers(b, id.Members, fontSize, padding)
	return b
}

func (r *ClassRenderer) measureEnum(ed *ast.EnumDef, fontSize, padding float64) *classBox {
	b := &classBox{
		id:   ed.Name,
		name: ed.Name,
		kind: "enum",
	}
	r.measureMembers(b, ed.Members, fontSize, padding)
	return b
}

func (r *ClassRenderer) measureImplicitClass(name string, fontSize, padding float64) *classBox {
	b := &classBox{
		id:   name,
		name: name,
		kind: "class",
	}
	r.measureMembers(b, nil, fontSize, padding)
	return b
}

func (r *ClassRenderer) measureMembers(b *classBox, members []ast.Member, fontSize, padding float64) {
	lineH := fontSize + 4
	nameSize, _ := font.MeasureText(b.name, fontSize, font.FamilyBold)
	maxW := nameSize.Width + 2*padding
	b.nameH = lineH + 2*padding
	if b.stereotype != "" || b.kind == "interface" || b.kind == "enum" {
		b.nameH += float64(stereotypeFontPx) + 4
	}
	for _, m := range members {
		switch mem := m.(type) {
		case *ast.Field:
			text := formatField(mem)
			b.fields = append(b.fields, memberLine{
				visibility: mem.Visibility,
				modifier:   mem.Modifier,
				text:       text,
			})
		case *ast.Method:
			text := formatMethod(mem)
			b.methods = append(b.methods, memberLine{
				visibility: mem.Visibility,
				modifier:   mem.Modifier,
				text:       text,
			})
		}
	}
	if len(b.fields) > 0 {
		b.fieldsH = float64(len(b.fields))*lineH + padding
		for _, f := range b.fields {
			sz, _ := font.MeasureText(f.text, fontSize, font.FamilySans)
			w := sz.Width + visibilityWidth + 2*padding
			if w > maxW {
				maxW = w
			}
		}
	}
	if len(b.methods) > 0 {
		b.methodsH = float64(len(b.methods))*lineH + padding
		for _, m := range b.methods {
			sz, _ := font.MeasureText(m.text, fontSize, font.FamilySans)
			w := sz.Width + visibilityWidth + 2*padding
			if w > maxW {
				maxW = w
			}
		}
	}
	b.width = math.Max(maxW, 100)
	b.height = b.nameH
	if len(b.fields) > 0 {
		b.height += b.fieldsH + compartmentGap
	}
	if len(b.methods) > 0 {
		b.height += b.methodsH + compartmentGap
	}
}

func (r *ClassRenderer) measureNote(note *ast.Note, fontSize, padding float64) *noteBox {
	isLeft := note.Placement == ast.NoteLeft
	sz, _ := font.MeasureText(note.Text, fontSize, font.FamilySans)
	nb := &noteBox{
		target: note.Target,
		text:   note.Text,
		left:   isLeft,
		width:  sz.Width + 2*padding + 10,
		height: sz.Height + 2*padding,
	}
	if nb.width < 80 {
		nb.width = 80
	}
	if nb.height < 30 {
		nb.height = 30
	}
	return nb
}

func (r *ClassRenderer) computePackageBounds(pb *packageBox, nodeByID map[string]*layout.Node, padding float64) {
	if len(pb.children) == 0 {
		pb.w = 100
		pb.h = 60
		return
	}
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64
	for _, childID := range pb.children {
		n := nodeByID[childID]
		if n == nil {
			continue
		}
		if n.X < minX {
			minX = n.X
		}
		if n.Y < minY {
			minY = n.Y
		}
		if n.X+n.Width > maxX {
			maxX = n.X + n.Width
		}
		if n.Y+n.Height > maxY {
			maxY = n.Y + n.Height
		}
	}
	tabH := 25.0
	pb.x = minX - padding
	pb.y = minY - padding - tabH
	pb.w = (maxX - minX) + 2*padding
	pb.h = (maxY - minY) + 2*padding + tabH
}

func (r *ClassRenderer) renderClassBox(sb *strings.Builder, b *classBox, x, y, fontSize, padding float64) {
	bgColor := r.resolver.ResolveColor("ClassBackgroundColor")
	borderColor := r.resolver.ResolveColor("ClassBorderColor")
	fontColor := r.resolver.ResolveColor("ClassFontColor")
	borderW := r.resolver.ResolveInt("BorderWidth", 1)
	switch b.kind {
	case "interface":
		bgColor = r.resolver.ResolveColor("InterfaceBackgroundColor")
		borderColor = r.resolver.ResolveColor("InterfaceBorderColor")
		fontColor = r.resolver.ResolveColor("InterfaceFontColor")
	case "enum":
		bgColor = r.resolver.ResolveColor("EnumBackgroundColor")
		borderColor = r.resolver.ResolveColor("EnumBorderColor")
		fontColor = r.resolver.ResolveColor("EnumFontColor")
	}
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%d" ry="%d" fill="%s" stroke="%s" stroke-width="%d"/>`,
		x, y, b.width, b.height, cornerRadius, cornerRadius, bgColor, borderColor, borderW)
	sb.WriteString("\n")
	lineH := fontSize + 4
	nameY := y + padding
	stereotypeColor := r.resolver.ResolveColor("ClassStereotypeFontColor")
	switch {
	case b.kind == "interface":
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="%d" fill="%s" font-style="italic">&lt;&lt;interface&gt;&gt;</text>`,
			x+b.width/2, nameY+float64(stereotypeFontPx), stereotypeFontPx, stereotypeColor)
		sb.WriteString("\n")
		nameY += float64(stereotypeFontPx) + 4
	case b.kind == "enum":
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="%d" fill="%s" font-style="italic">&lt;&lt;enum&gt;&gt;</text>`,
			x+b.width/2, nameY+float64(stereotypeFontPx), stereotypeFontPx, stereotypeColor)
		sb.WriteString("\n")
		nameY += float64(stereotypeFontPx) + 4
	case b.stereotype != "":
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="%d" fill="%s" font-style="italic">&lt;&lt;%s&gt;&gt;</text>`,
			x+b.width/2, nameY+float64(stereotypeFontPx), stereotypeFontPx, stereotypeColor, escapeXML(b.stereotype))
		sb.WriteString("\n")
		nameY += float64(stereotypeFontPx) + 4
	}
	fontStyle := ""
	if b.abstract {
		fontStyle = ` font-style="italic"`
	}
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="%.0f" font-weight="bold" fill="%s"%s>%s</text>`,
		x+b.width/2, nameY+fontSize, fontSize, fontColor, fontStyle, escapeXML(b.name))
	sb.WriteString("\n")
	curY := y + b.nameH
	if len(b.fields) > 0 {
		fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%d"/>`,
			x, curY, x+b.width, curY, borderColor, borderW)
		sb.WriteString("\n")
		memberY := curY + padding/2
		for _, f := range b.fields {
			r.renderMemberLine(sb, f, x+padding, memberY+lineH-2, fontSize, fontColor)
			memberY += lineH
		}
		curY += b.fieldsH + compartmentGap
	}
	if len(b.methods) > 0 {
		fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%d"/>`,
			x, curY, x+b.width, curY, borderColor, borderW)
		sb.WriteString("\n")
		memberY := curY + padding/2
		for _, m := range b.methods {
			r.renderMemberLine(sb, m, x+padding, memberY+lineH-2, fontSize, fontColor)
			memberY += lineH
		}
	}
}

func (r *ClassRenderer) renderMemberLine(sb *strings.Builder, ml memberLine, x, y, fontSize float64, fontColor string) {
	annotationColor := r.resolver.ResolveColor("AnnotationColor")
	visIcon := visibilityIcon(ml.visibility)
	visColor := visibilityColor(ml.visibility, annotationColor)
	if visIcon != "" {
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%.0f" fill="%s">%s</text>`,
			x, y, fontSize, visColor, visIcon)
	}
	textX := x + visibilityWidth
	decoration := ""
	if ml.modifier == ast.ModifierStatic {
		decoration = ` text-decoration="underline"`
	}
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%.0f" fill="%s"%s>%s</text>`,
		textX, y, fontSize, fontColor, decoration, escapeXML(ml.text))
	sb.WriteString("\n")
}

func (r *ClassRenderer) renderRelationship(sb *strings.Builder, rel *ast.Relationship, from, to *layout.Node, offsetX, offsetY, fontSize float64) {
	arrowColor := r.resolver.ResolveColor("ArrowColor")
	thickness := r.resolver.ResolveInt("ArrowThickness", 1)
	fromCX := from.X + from.Width/2 + offsetX
	fromCY := from.Y + from.Height/2 + offsetY
	toCX := to.X + to.Width/2 + offsetX
	toCY := to.Y + to.Height/2 + offsetY
	fromPt := edgePoint(from.X+offsetX, from.Y+offsetY, from.Width, from.Height, toCX, toCY)
	toPt := edgePoint(to.X+offsetX, to.Y+offsetY, to.Width, to.Height, fromCX, fromCY)
	dashAttr := ""
	if rel.Type == ast.RelDependency || rel.Type == ast.RelRealization {
		dashAttr = ` stroke-dasharray="7,4"`
	}
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="%d"%s/>`,
		fromPt.x, fromPt.y, toPt.x, toPt.y, arrowColor, thickness, dashAttr)
	sb.WriteString("\n")
	r.renderArrowHead(sb, rel, fromPt, toPt, arrowColor)
	if rel.Label != "" {
		arrowFontSize := r.resolver.ResolveInt("ArrowFontSize", 11)
		labelX := (fromPt.x + toPt.x) / 2
		labelY := (fromPt.y+toPt.y)/2 - 5
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="%d" fill="%s">%s</text>`,
			labelX, labelY, arrowFontSize, arrowColor, escapeXML(rel.Label))
		sb.WriteString("\n")
	}
	if rel.LeftCard != "" {
		r.renderCardinality(sb, rel.LeftCard, fromPt, toPt, true, arrowColor)
	}
	if rel.RightCard != "" {
		r.renderCardinality(sb, rel.RightCard, fromPt, toPt, false, arrowColor)
	}
}

func (r *ClassRenderer) renderCardinality(sb *strings.Builder, card string, from, to point, nearFrom bool, color string) {
	t := 0.1
	if !nearFrom {
		t = 0.9
	}
	cx := from.x + t*(to.x-from.x)
	cy := from.y + t*(to.y-from.y) - 8
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" text-anchor="middle" font-family="sans-serif" font-size="11" fill="%s">%s</text>`,
		cx, cy, color, escapeXML(card))
	sb.WriteString("\n")
}

func (r *ClassRenderer) renderArrowHead(sb *strings.Builder, rel *ast.Relationship, from, to point, color string) {
	dir := rel.Direction
	switch rel.Type {
	case ast.RelInheritance:
		if dir == ast.ArrowLeft {
			drawTriangle(sb, from, to, color, true)
		} else {
			drawTriangle(sb, to, from, color, true)
		}
	case ast.RelRealization:
		if dir == ast.ArrowLeft {
			drawTriangle(sb, from, to, color, true)
		} else {
			drawTriangle(sb, to, from, color, true)
		}
	case ast.RelComposition:
		if dir == ast.ArrowLeft {
			drawDiamond(sb, from, to, color, true)
		} else {
			drawDiamond(sb, to, from, color, true)
		}
	case ast.RelAggregation:
		if dir == ast.ArrowLeft {
			drawDiamond(sb, from, to, color, false)
		} else {
			drawDiamond(sb, to, from, color, false)
		}
	case ast.RelDependency, ast.RelAssociation:
		if dir == ast.ArrowLeft || dir == ast.ArrowBoth {
			drawOpenArrow(sb, from, to, color)
		}
		if dir == ast.ArrowRight || dir == ast.ArrowBoth {
			drawOpenArrow(sb, to, from, color)
		}
	}
}

func (r *ClassRenderer) renderPackage(sb *strings.Builder, pb *packageBox, offsetX, offsetY, fontSize float64) {
	x := pb.x + offsetX
	y := pb.y + offsetY
	bgColor := r.resolver.ResolveColor("PackageBackgroundColor")
	borderColor := r.resolver.ResolveColor("PackageBorderColor")
	fontColor := r.resolver.ResolveColor("PackageFontColor")
	tabW := 80.0
	tabH := 20.0
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s"/>`,
		x, y, tabW, tabH, bgColor, borderColor)
	sb.WriteString("\n")
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" fill-opacity="0.3"/>`,
		x, y+tabH, pb.w, pb.h-tabH, bgColor, borderColor)
	sb.WriteString("\n")
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%.0f" fill="%s">%s</text>`,
		x+5, y+tabH-5, fontSize, fontColor, escapeXML(pb.name))
	sb.WriteString("\n")
}

func (r *ClassRenderer) renderNote(sb *strings.Builder, nb *noteBox, x, y, fontSize float64) {
	bgColor := r.resolver.ResolveColor("NoteBackgroundColor")
	borderColor := r.resolver.ResolveColor("NoteBorderColor")
	fontColor := r.resolver.ResolveColor("NoteFontColor")
	fold := 10.0
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s"/>`,
		x, y,
		x+nb.width-fold, y,
		x+nb.width, y+fold,
		x+nb.width, y+nb.height,
		x, y+nb.height,
		bgColor, borderColor)
	sb.WriteString("\n")
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s"/>`,
		x+nb.width-fold, y,
		x+nb.width-fold, y+fold,
		x+nb.width, y+fold,
		bgColor, borderColor)
	sb.WriteString("\n")
	lines := strings.Split(nb.text, "\n")
	lineH := fontSize + 4
	textY := y + fontSize + 5
	for _, line := range lines {
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%.0f" fill="%s">%s</text>`,
			x+5, textY, fontSize, fontColor, escapeXML(line))
		sb.WriteString("\n")
		textY += lineH
	}
}

type point struct {
	x, y float64
}

func edgePoint(rx, ry, rw, rh, targetX, targetY float64) point {
	cx := rx + rw/2
	cy := ry + rh/2
	dx := targetX - cx
	dy := targetY - cy
	if dx == 0 && dy == 0 {
		return point{cx, cy}
	}
	absDx := math.Abs(dx)
	absDy := math.Abs(dy)
	scaleX := (rw / 2) / absDx
	scaleY := (rh / 2) / absDy
	if absDx < 0.001 {
		if dy > 0 {
			return point{cx, ry + rh}
		}
		return point{cx, ry}
	}
	if absDy < 0.001 {
		if dx > 0 {
			return point{rx + rw, cy}
		}
		return point{rx, cy}
	}
	scale := math.Min(scaleX, scaleY)
	return point{cx + dx*scale, cy + dy*scale}
}

func drawTriangle(sb *strings.Builder, tip, from point, color string, filled bool) {
	size := 12.0
	angle := math.Atan2(from.y-tip.y, from.x-tip.x)
	spread := math.Pi / 7
	x1 := tip.x + size*math.Cos(angle+spread)
	y1 := tip.y + size*math.Sin(angle+spread)
	x2 := tip.x + size*math.Cos(angle-spread)
	y2 := tip.y + size*math.Sin(angle-spread)
	fill := color
	if filled {
		fill = "white"
	}
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="1"/>`,
		tip.x, tip.y, x1, y1, x2, y2, fill, color)
	sb.WriteString("\n")
}

func drawDiamond(sb *strings.Builder, tip, from point, color string, filled bool) {
	size := 10.0
	angle := math.Atan2(from.y-tip.y, from.x-tip.x)
	spread := math.Pi / 5
	mx := tip.x + size*math.Cos(angle)
	my := tip.y + size*math.Sin(angle)
	x1 := tip.x + size*0.6*math.Cos(angle+spread)
	y1 := tip.y + size*0.6*math.Sin(angle+spread)
	x2 := tip.x + size*0.6*math.Cos(angle-spread)
	y2 := tip.y + size*0.6*math.Sin(angle-spread)
	fill := color
	if !filled {
		fill = "white"
	}
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="1"/>`,
		tip.x, tip.y, x1, y1, mx, my, x2, y2, fill, color)
	sb.WriteString("\n")
}

func drawOpenArrow(sb *strings.Builder, tip, from point, color string) {
	size := 10.0
	angle := math.Atan2(from.y-tip.y, from.x-tip.x)
	spread := math.Pi / 7
	x1 := tip.x + size*math.Cos(angle+spread)
	y1 := tip.y + size*math.Sin(angle+spread)
	x2 := tip.x + size*math.Cos(angle-spread)
	y2 := tip.y + size*math.Sin(angle-spread)
	fmt.Fprintf(sb, `<polyline points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="none" stroke="%s" stroke-width="1"/>`,
		x1, y1, tip.x, tip.y, x2, y2, color)
	sb.WriteString("\n")
}

func formatField(f *ast.Field) string {
	s := f.Name
	if f.Type != "" {
		s += " : " + f.Type
	}
	return s
}

func formatMethod(m *ast.Method) string {
	s := m.Name + "(" + m.Params + ")"
	if m.ReturnType != "" {
		s += " : " + m.ReturnType
	}
	return s
}

func visibilityIcon(v ast.Visibility) string {
	switch v {
	case ast.VisibilityPublic:
		return "+"
	case ast.VisibilityPrivate:
		return "-"
	case ast.VisibilityProtected:
		return "#"
	case ast.VisibilityPackage:
		return "~"
	default:
		return ""
	}
}

func visibilityColor(v ast.Visibility, defaultColor string) string {
	switch v {
	case ast.VisibilityPublic:
		return "#6A8759"
	case ast.VisibilityPrivate:
		return "#CC7832"
	case ast.VisibilityProtected:
		return "#FFC66D"
	case ast.VisibilityPackage:
		return "#6897BB"
	default:
		return defaultColor
	}
}

func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
