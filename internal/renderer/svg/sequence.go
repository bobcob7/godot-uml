package svg

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/bobcob7/go-uml/internal/ast"
	"github.com/bobcob7/go-uml/internal/font"
	"github.com/bobcob7/go-uml/internal/theme"
)

// SequenceRenderer renders sequence diagrams to SVG.
type SequenceRenderer struct {
	resolver *theme.Resolver
}

// NewSequenceRenderer creates a new sequence diagram SVG renderer.
// If resolver is nil, the default Darcula theme is used.
func NewSequenceRenderer(resolver *theme.Resolver) *SequenceRenderer {
	if resolver == nil {
		resolver = theme.NewResolver(nil)
	}
	return &SequenceRenderer{resolver: resolver}
}

// participantBox holds layout info for a participant.
type participantBox struct {
	name   string
	alias  string
	kind   ast.ParticipantKind
	x      float64 // center x
	y      float64 // top of box
	width  float64
	height float64
}

// displayName returns the name to show for a participant.
func (p *participantBox) displayName() string {
	if p.alias != "" {
		return p.alias
	}
	return p.name
}

// centerX returns the x coordinate of the lifeline.
func (p *participantBox) centerX() float64 {
	return p.x + p.width/2
}

// bottomY returns the bottom of the participant box.
func (p *participantBox) bottomY() float64 {
	return p.y + p.height
}

// seqEvent represents something that occupies vertical space in the diagram.
type seqEvent struct {
	y      float64
	height float64
	stmt   ast.Statement
}

// activationRange tracks when a lifeline is active.
type activationRange struct {
	participant string
	startY      float64
	endY        float64
}

const (
	seqParticipantPadX = 20.0
	seqParticipantPadY = 8.0
	seqParticipantGap  = 40.0
	seqMessageSpacing  = 40.0
	seqArrowSize       = 8.0
	seqActivationWidth = 10.0
	seqFragmentPadding = 10.0
	seqNotePadding     = 8.0
	seqNoteMaxWidth    = 150.0
	seqTopMargin       = 20.0
	seqLeftMargin      = 20.0
	seqBottomMargin    = 20.0
	seqDividerHeight   = 30.0
	seqDelayHeight     = 30.0
	seqFragmentLabelH  = 20.0
	seqLifelineDash    = "5,5"
)

// Render writes the sequence diagram SVG to w.
func (r *SequenceRenderer) Render(w io.Writer, diagram *ast.Diagram) error {
	participants := r.collectParticipants(diagram)
	if len(participants) == 0 {
		return r.renderEmpty(w)
	}
	r.applySkinparams(diagram)
	pboxes := r.layoutParticipants(participants)
	pmap := make(map[string]*participantBox)
	for i := range pboxes {
		pmap[pboxes[i].name] = &pboxes[i]
		if pboxes[i].alias != "" {
			pmap[pboxes[i].alias] = &pboxes[i]
		}
	}
	events, activations := r.layoutEvents(diagram, pboxes, pmap)
	totalWidth, totalHeight := r.computeBounds(pboxes, events, activations)
	var sb strings.Builder
	fmt.Fprintf(&sb, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f"`, totalWidth, totalHeight)
	fmt.Fprintf(&sb, ` viewBox="0 0 %.0f %.0f">`, totalWidth, totalHeight)
	bgColor := r.resolver.ResolveColor("BackgroundColor")
	fmt.Fprintf(&sb, `<rect width="%.0f" height="%.0f" fill="%s"/>`, totalWidth, totalHeight, escSeq(bgColor))
	for i := range pboxes {
		r.renderParticipantBox(&sb, &pboxes[i])
	}
	lifelineEndY := r.lifelineEndY(events, pboxes)
	for i := range pboxes {
		r.renderLifeline(&sb, &pboxes[i], lifelineEndY)
	}
	for i := range activations {
		r.renderActivation(&sb, &activations[i], pmap)
	}
	msgNum := 0
	autonumber := false
	for _, ev := range events {
		switch s := ev.stmt.(type) {
		case *ast.Message:
			if autonumber {
				msgNum++
			}
			r.renderMessage(&sb, s, ev.y, pmap, autonumber, msgNum)
		case *ast.Note:
			r.renderSeqNote(&sb, s, ev.y, pmap)
		case *ast.Fragment:
			r.renderFragment(&sb, s, ev.y, ev.height, pmap, pboxes)
		case *ast.Divider:
			r.renderDivider(&sb, s, ev.y, totalWidth)
		case *ast.Delay:
			r.renderDelay(&sb, s, ev.y, totalWidth)
		case *ast.Autonumber:
			autonumber = true
			if s.Start != "" {
				msgNum = atoiSimple(s.Start) - 1
			}
		}
	}
	for i := range pboxes {
		r.renderParticipantBoxBottom(&sb, &pboxes[i], lifelineEndY)
	}
	sb.WriteString("</svg>")
	_, err := io.WriteString(w, sb.String())
	return err
}

func (r *SequenceRenderer) renderEmpty(w io.Writer) error {
	bgColor := r.resolver.ResolveColor("BackgroundColor")
	svg := fmt.Sprintf(`<svg xmlns="http://www.w3.org/2000/svg" width="100" height="100"><rect width="100" height="100" fill="%s"/></svg>`, escSeq(bgColor))
	_, err := io.WriteString(w, svg)
	return err
}

func (r *SequenceRenderer) applySkinparams(diagram *ast.Diagram) {
	for _, stmt := range diagram.Statements {
		if sp, ok := stmt.(*ast.Skinparam); ok {
			r.resolver.SetSkinparam(sp.Name, sp.Value)
		}
	}
}

// collectParticipants extracts ordered participants from the diagram.
// Explicit participant declarations come first, then implicit ones from messages.
func (r *SequenceRenderer) collectParticipants(diagram *ast.Diagram) []*ast.Participant {
	seen := make(map[string]bool)
	var result []*ast.Participant
	for _, stmt := range diagram.Statements {
		if p, ok := stmt.(*ast.Participant); ok {
			key := p.Name
			if !seen[key] {
				seen[key] = true
				result = append(result, p)
			}
		}
	}
	for _, stmt := range diagram.Statements {
		if m, ok := stmt.(*ast.Message); ok {
			for _, name := range []string{m.From, m.To} {
				if name != "" && !seen[name] {
					seen[name] = true
					result = append(result, &ast.Participant{Name: name, Kind: ast.ParticipantDefault})
				}
			}
		}
	}
	return result
}

// layoutParticipants computes the positions of participant boxes.
func (r *SequenceRenderer) layoutParticipants(participants []*ast.Participant) []participantBox {
	fontSize := float64(r.resolver.ResolveInt("FontSize", 13))
	boxes := make([]participantBox, len(participants))
	for i, p := range participants {
		displayName := p.Name
		if p.Alias != "" {
			displayName = p.Alias
		}
		size, _ := font.MeasureText(displayName, fontSize, font.FamilySans)
		boxes[i] = participantBox{
			name:   p.Name,
			alias:  p.Alias,
			kind:   p.Kind,
			width:  size.Width + seqParticipantPadX*2,
			height: size.Height + seqParticipantPadY*2,
		}
	}
	x := seqLeftMargin
	for i := range boxes {
		boxes[i].x = x
		boxes[i].y = seqTopMargin
		x += boxes[i].width + seqParticipantGap
	}
	return boxes
}

// layoutEvents assigns vertical Y positions to each diagram statement.
func (r *SequenceRenderer) layoutEvents(diagram *ast.Diagram, pboxes []participantBox, pmap map[string]*participantBox) ([]seqEvent, []activationRange) {
	var events []seqEvent
	var activations []activationRange
	activeStarts := make(map[string]float64)
	maxBottom := float64(0)
	for _, pb := range pboxes {
		if pb.bottomY() > maxBottom {
			maxBottom = pb.bottomY()
		}
	}
	curY := maxBottom + seqMessageSpacing
	for _, stmt := range diagram.Statements {
		switch s := stmt.(type) {
		case *ast.Message:
			events = append(events, seqEvent{y: curY, height: seqMessageSpacing, stmt: s})
			curY += seqMessageSpacing
		case *ast.Note:
			h := r.noteHeight(s)
			events = append(events, seqEvent{y: curY, height: h, stmt: s})
			curY += h
		case *ast.Fragment:
			h := r.fragmentHeight(s)
			events = append(events, seqEvent{y: curY, height: h, stmt: s})
			curY += h
		case *ast.Divider:
			events = append(events, seqEvent{y: curY, height: seqDividerHeight, stmt: s})
			curY += seqDividerHeight
		case *ast.Delay:
			events = append(events, seqEvent{y: curY, height: seqDelayHeight, stmt: s})
			curY += seqDelayHeight
		case *ast.Autonumber:
			events = append(events, seqEvent{y: curY, height: 0, stmt: s})
		case *ast.Activate:
			if s.Deactivate {
				if startY, ok := activeStarts[s.Target]; ok {
					activations = append(activations, activationRange{
						participant: s.Target,
						startY:     startY,
						endY:       curY,
					})
					delete(activeStarts, s.Target)
				}
			} else {
				activeStarts[s.Target] = curY
			}
		}
	}
	for name, startY := range activeStarts {
		activations = append(activations, activationRange{
			participant: name,
			startY:     startY,
			endY:       curY,
		})
	}
	return events, activations
}

func (r *SequenceRenderer) noteHeight(n *ast.Note) float64 {
	fontSize := float64(r.resolver.ResolveInt("FontSize", 13))
	size, _ := font.MeasureText(n.Text, fontSize, font.FamilySans)
	return size.Height + seqNotePadding*2 + 10
}

func (r *SequenceRenderer) fragmentHeight(f *ast.Fragment) float64 {
	h := seqFragmentLabelH + seqFragmentPadding
	count := len(f.Statements)
	for _, ep := range f.ElseParts {
		count += len(ep.Statements)
		h += seqFragmentLabelH // else divider
	}
	if count == 0 {
		count = 1
	}
	h += float64(count) * seqMessageSpacing
	h += seqFragmentPadding
	return h
}

func (r *SequenceRenderer) computeBounds(pboxes []participantBox, events []seqEvent, _ []activationRange) (float64, float64) {
	maxX := float64(0)
	for _, pb := range pboxes {
		right := pb.x + pb.width
		if right > maxX {
			maxX = right
		}
	}
	maxY := float64(0)
	for _, pb := range pboxes {
		if pb.bottomY() > maxY {
			maxY = pb.bottomY()
		}
	}
	for _, ev := range events {
		bottom := ev.y + ev.height
		if bottom > maxY {
			maxY = bottom
		}
	}
	maxY += seqMessageSpacing // space for bottom participant boxes
	for _, pb := range pboxes {
		maxY += pb.height
	}
	maxY += seqBottomMargin
	maxX += seqLeftMargin
	return maxX, maxY
}

func (r *SequenceRenderer) lifelineEndY(events []seqEvent, pboxes []participantBox) float64 {
	maxY := float64(0)
	for _, pb := range pboxes {
		if pb.bottomY() > maxY {
			maxY = pb.bottomY()
		}
	}
	for _, ev := range events {
		bottom := ev.y + ev.height
		if bottom > maxY {
			maxY = bottom
		}
	}
	return maxY + seqMessageSpacing/2
}

func (r *SequenceRenderer) renderParticipantBox(sb *strings.Builder, pb *participantBox) {
	bgColor := r.resolver.ResolveColor("ParticipantBackgroundColor")
	borderColor := r.resolver.ResolveColor("ParticipantBorderColor")
	fontColor := r.resolver.ResolveColor("ParticipantFontColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	borderWidth := r.resolver.ResolveInt("BorderWidth", 1)
	switch pb.kind {
	case ast.ParticipantActor:
		r.renderActorIcon(sb, pb, borderColor, fontColor, fontSize)
	default:
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="%d" rx="4"/>`,
			pb.x, pb.y, pb.width, pb.height, escSeq(bgColor), escSeq(borderColor), borderWidth)
		textX := pb.centerX()
		textY := pb.y + pb.height/2 + float64(fontSize)/3
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle">%s</text>`,
			textX, textY, fontSize, escSeq(fontColor), escSeq(pb.displayName()))
	}
}

func (r *SequenceRenderer) renderActorIcon(sb *strings.Builder, pb *participantBox, borderColor, fontColor string, fontSize int) {
	cx := pb.centerX()
	topY := pb.y + 4
	headR := 8.0
	fmt.Fprintf(sb, `<circle cx="%.1f" cy="%.1f" r="%.1f" fill="none" stroke="%s" stroke-width="1"/>`,
		cx, topY+headR, headR, escSeq(borderColor))
	bodyTop := topY + headR*2
	bodyBot := bodyTop + 12
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1"/>`,
		cx, bodyTop, cx, bodyBot, escSeq(borderColor))
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1"/>`,
		cx-10, bodyTop+4, cx+10, bodyTop+4, escSeq(borderColor))
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1"/>`,
		cx, bodyBot, cx-8, bodyBot+10, escSeq(borderColor))
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1"/>`,
		cx, bodyBot, cx+8, bodyBot+10, escSeq(borderColor))
	textY := pb.y + pb.height - 2
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle">%s</text>`,
		cx, textY, fontSize, escSeq(fontColor), escSeq(pb.displayName()))
}

func (r *SequenceRenderer) renderParticipantBoxBottom(sb *strings.Builder, pb *participantBox, lifelineEndY float64) {
	bgColor := r.resolver.ResolveColor("ParticipantBackgroundColor")
	borderColor := r.resolver.ResolveColor("ParticipantBorderColor")
	fontColor := r.resolver.ResolveColor("ParticipantFontColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	borderWidth := r.resolver.ResolveInt("BorderWidth", 1)
	y := lifelineEndY
	switch pb.kind {
	case ast.ParticipantActor:
		botPb := *pb
		botPb.y = y
		r.renderActorIcon(sb, &botPb, borderColor, fontColor, fontSize)
	default:
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="%d" rx="4"/>`,
			pb.x, y, pb.width, pb.height, escSeq(bgColor), escSeq(borderColor), borderWidth)
		textX := pb.centerX()
		textY := y + pb.height/2 + float64(fontSize)/3
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle">%s</text>`,
			textX, textY, fontSize, escSeq(fontColor), escSeq(pb.displayName()))
	}
}

func (r *SequenceRenderer) renderLifeline(sb *strings.Builder, pb *participantBox, endY float64) {
	lineColor := r.resolver.ResolveColor("SequenceLifeLineBorderColor")
	cx := pb.centerX()
	startY := pb.bottomY()
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="%s"/>`,
		cx, startY, cx, endY, escSeq(lineColor), seqLifelineDash)
}

func (r *SequenceRenderer) renderActivation(sb *strings.Builder, a *activationRange, pmap map[string]*participantBox) {
	pb, ok := pmap[a.participant]
	if !ok {
		return
	}
	bgColor := r.resolver.ResolveColor("ParticipantBackgroundColor")
	borderColor := r.resolver.ResolveColor("ParticipantBorderColor")
	cx := pb.centerX()
	x := cx - seqActivationWidth/2
	h := a.endY - a.startY
	if h < 5 {
		h = 5
	}
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s" stroke="%s" stroke-width="1"/>`,
		x, a.startY, seqActivationWidth, h, escSeq(bgColor), escSeq(borderColor))
}

func (r *SequenceRenderer) renderMessage(sb *strings.Builder, m *ast.Message, y float64, pmap map[string]*participantBox, autonumber bool, msgNum int) {
	fromPb := pmap[m.From]
	toPb := pmap[m.To]
	if fromPb == nil || toPb == nil {
		return
	}
	arrowColor := r.resolver.ResolveColor("ArrowColor")
	fontColor := r.resolver.ResolveColor("FontColor")
	fontSize := r.resolver.ResolveInt("ArrowFontSize", 11)
	x1 := fromPb.centerX()
	x2 := toPb.centerX()
	dashAttr := ""
	if m.Dashed {
		dashAttr = ` stroke-dasharray="6,4"`
	}
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1"%s/>`,
		x1, y, x2, y, escSeq(arrowColor), dashAttr)
	r.drawSeqArrowHead(sb, x1, x2, y, arrowColor)
	label := m.Label
	if autonumber && msgNum > 0 {
		label = fmt.Sprintf("%d. %s", msgNum, label)
	}
	if label != "" {
		midX := (x1 + x2) / 2
		labelY := y - 5
		fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle">%s</text>`,
			midX, labelY, fontSize, escSeq(fontColor), escSeq(label))
	}
}

func (r *SequenceRenderer) drawSeqArrowHead(sb *strings.Builder, x1, x2, y float64, color string) {
	if x2 > x1 {
		fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s"/>`,
			x2, y, x2-seqArrowSize, y-seqArrowSize/2, x2-seqArrowSize, y+seqArrowSize/2, escSeq(color))
	} else {
		fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s"/>`,
			x2, y, x2+seqArrowSize, y-seqArrowSize/2, x2+seqArrowSize, y+seqArrowSize/2, escSeq(color))
	}
}

func (r *SequenceRenderer) renderSeqNote(sb *strings.Builder, n *ast.Note, y float64, pmap map[string]*participantBox) {
	bgColor := r.resolver.ResolveColor("NoteBackgroundColor")
	borderColor := r.resolver.ResolveColor("NoteBorderColor")
	fontColor := r.resolver.ResolveColor("NoteFontColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	size, _ := font.MeasureText(n.Text, float64(fontSize), font.FamilySans)
	noteW := size.Width + seqNotePadding*2
	if noteW > seqNoteMaxWidth {
		noteW = seqNoteMaxWidth
	}
	noteH := size.Height + seqNotePadding*2
	pb := pmap[n.Target]
	if pb == nil {
		return
	}
	cx := pb.centerX()
	var noteX float64
	switch n.Placement {
	case ast.NoteLeft:
		noteX = cx - noteW - 15
	case ast.NoteRight:
		noteX = cx + 15
	case ast.NoteOver:
		noteX = cx - noteW/2
	}
	fold := 8.0
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="%s" stroke="%s" stroke-width="1"/>`,
		noteX, y,
		noteX+noteW-fold, y,
		noteX+noteW, y+fold,
		noteX+noteW, y+noteH,
		noteX, y+noteH,
		escSeq(bgColor), escSeq(borderColor))
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="none" stroke="%s" stroke-width="1"/>`,
		noteX+noteW-fold, y,
		noteX+noteW-fold, y+fold,
		noteX+noteW, y+fold,
		escSeq(borderColor))
	fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="5,5"/>`,
		cx, y+noteH/2, noteX+noteW, y+noteH/2, escSeq(borderColor))
	textX := noteX + seqNotePadding
	textY := y + seqNotePadding + float64(fontSize)
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s">%s</text>`,
		textX, textY, fontSize, escSeq(fontColor), escSeq(n.Text))
}

func (r *SequenceRenderer) renderFragment(sb *strings.Builder, f *ast.Fragment, y, height float64, pmap map[string]*participantBox, pboxes []participantBox) {
	borderColor := r.resolver.ResolveColor("ParticipantBorderColor")
	fontColor := r.resolver.ResolveColor("FontColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	minX, maxX := r.fragmentSpan(f, pmap, pboxes)
	fragX := minX - seqFragmentPadding
	fragW := (maxX - minX) + seqFragmentPadding*2
	if fragW < 100 {
		fragW = 100
	}
	fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="none" stroke="%s" stroke-width="1"/>`,
		fragX, y, fragW, height, escSeq(borderColor))
	label := fragmentLabel(f.Kind)
	if f.Condition != "" {
		label += " [" + f.Condition + "]"
	}
	labelW, _ := font.MeasureText(label, float64(fontSize), font.FamilySans)
	tagW := labelW.Width + 16
	tagH := seqFragmentLabelH
	fmt.Fprintf(sb, `<polygon points="%.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f %.1f,%.1f" fill="none" stroke="%s" stroke-width="1"/>`,
		fragX, y,
		fragX+tagW, y,
		fragX+tagW, y+tagH-5,
		fragX+tagW-5, y+tagH,
		fragX, y+tagH,
		escSeq(borderColor))
	fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" font-weight="bold">%s</text>`,
		fragX+8, y+tagH-5, fontSize, escSeq(fontColor), escSeq(label))
	if len(f.ElseParts) > 0 {
		stmtCount := len(f.Statements)
		if stmtCount == 0 {
			stmtCount = 1
		}
		elseY := y + seqFragmentLabelH + seqFragmentPadding + float64(stmtCount)*seqMessageSpacing
		for _, ep := range f.ElseParts {
			fmt.Fprintf(sb, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="5,5"/>`,
				fragX, elseY, fragX+fragW, elseY, escSeq(borderColor))
			elseLabel := "else"
			if ep.Condition != "" {
				elseLabel += " [" + ep.Condition + "]"
			}
			fmt.Fprintf(sb, `<text x="%.1f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s">%s</text>`,
				fragX+8, elseY+float64(fontSize)+2, fontSize, escSeq(fontColor), escSeq(elseLabel))
			epCount := len(ep.Statements)
			if epCount == 0 {
				epCount = 1
			}
			elseY += seqFragmentLabelH + float64(epCount)*seqMessageSpacing
		}
	}
}

func (r *SequenceRenderer) fragmentSpan(f *ast.Fragment, pmap map[string]*participantBox, pboxes []participantBox) (float64, float64) {
	names := r.collectFragmentParticipants(f)
	if len(names) == 0 && len(pboxes) >= 2 {
		return pboxes[0].centerX(), pboxes[len(pboxes)-1].centerX()
	}
	if len(names) == 0 && len(pboxes) == 1 {
		return pboxes[0].x, pboxes[0].x + pboxes[0].width
	}
	minX := math.MaxFloat64
	maxX := -math.MaxFloat64
	for name := range names {
		if pb, ok := pmap[name]; ok {
			cx := pb.centerX()
			if cx-pb.width/2 < minX {
				minX = cx - pb.width/2
			}
			if cx+pb.width/2 > maxX {
				maxX = cx + pb.width/2
			}
		}
	}
	if minX == math.MaxFloat64 {
		return pboxes[0].centerX(), pboxes[len(pboxes)-1].centerX()
	}
	return minX, maxX
}

func (r *SequenceRenderer) collectFragmentParticipants(f *ast.Fragment) map[string]bool {
	names := make(map[string]bool)
	for _, stmt := range f.Statements {
		if m, ok := stmt.(*ast.Message); ok {
			names[m.From] = true
			names[m.To] = true
		}
	}
	for _, ep := range f.ElseParts {
		for _, stmt := range ep.Statements {
			if m, ok := stmt.(*ast.Message); ok {
				names[m.From] = true
				names[m.To] = true
			}
		}
	}
	return names
}

func (r *SequenceRenderer) renderDivider(sb *strings.Builder, d *ast.Divider, y, totalWidth float64) {
	fontColor := r.resolver.ResolveColor("FontColor")
	borderColor := r.resolver.ResolveColor("SequenceLifeLineBorderColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	midY := y + seqDividerHeight/2
	fmt.Fprintf(sb, `<line x1="0" y1="%.1f" x2="%.0f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="5,5"/>`,
		midY, totalWidth, midY, escSeq(borderColor))
	if d.Text != "" {
		size, _ := font.MeasureText(d.Text, float64(fontSize), font.FamilySans)
		rectW := size.Width + 20
		rectH := size.Height + 8
		rectX := totalWidth/2 - rectW/2
		rectY := midY - rectH/2
		bgColor := r.resolver.ResolveColor("BackgroundColor")
		fmt.Fprintf(sb, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" fill="%s"/>`,
			rectX, rectY, rectW, rectH, escSeq(bgColor))
		fmt.Fprintf(sb, `<text x="%.0f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle" font-weight="bold">%s</text>`,
			totalWidth/2, midY+float64(fontSize)/3, fontSize, escSeq(fontColor), escSeq(d.Text))
	}
}

func (r *SequenceRenderer) renderDelay(sb *strings.Builder, d *ast.Delay, y, totalWidth float64) {
	fontColor := r.resolver.ResolveColor("FontColor")
	fontSize := r.resolver.ResolveInt("FontSize", 13)
	midY := y + seqDelayHeight/2
	if d.Text != "" {
		fmt.Fprintf(sb, `<text x="%.0f" y="%.1f" font-family="sans-serif" font-size="%d" fill="%s" text-anchor="middle" font-style="italic">%s</text>`,
			totalWidth/2, midY+float64(fontSize)/3, fontSize, escSeq(fontColor), escSeq(d.Text))
	}
	fmt.Fprintf(sb, `<line x1="0" y1="%.1f" x2="%.0f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="2,4"/>`,
		y, totalWidth, y, escSeq(fontColor))
	fmt.Fprintf(sb, `<line x1="0" y1="%.1f" x2="%.0f" y2="%.1f" stroke="%s" stroke-width="1" stroke-dasharray="2,4"/>`,
		y+seqDelayHeight, totalWidth, y+seqDelayHeight, escSeq(fontColor))
}

func fragmentLabel(kind ast.FragmentKind) string {
	switch kind {
	case ast.FragmentAlt:
		return "alt"
	case ast.FragmentLoop:
		return "loop"
	case ast.FragmentPar:
		return "par"
	case ast.FragmentBreak:
		return "break"
	case ast.FragmentRef:
		return "ref"
	case ast.FragmentGroup:
		return "group"
	default:
		return "fragment"
	}
}

// escSeq escapes strings for safe SVG embedding.
func escSeq(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, `"`, "&quot;")
	return s
}

func atoiSimple(s string) int {
	n := 0
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
