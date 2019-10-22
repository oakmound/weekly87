package end

import (
	"fmt"
	"image/color"
	"math"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/alg/floatgeom"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/entities"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/mod"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/doodads"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/run"
	"github.com/oakmound/weekly87/internal/sfx"
)

func Init() {
	_, err := render.LoadSprite("", filepath.Join("raw", "wood_junk.png"))
	dlog.Error("Something went wrong loading in wood junk " + err.Error())
}

var (
	graveX   = 90.0
	graveY   = 300.0
	npcScale = 1.6 //TODO: remove need for this and rename
)

var stayInEndScene bool
var endSceneNextScene string

type intStringer struct {
	i *int
}

func (is intStringer) String() string {
	return strconv.Itoa(*is.i)
}

var (
	presentationX, presentationY, startY, pitX, pitY, hopDistance float64
)

// Scene is a scene in the same package as run to allow for easy variable access.
//If there is time at the e nd we can look at what vbariables this touches and get them exported or passed onwards so this can have its own package
var Scene = scene.Scene{
	Start: func(prevScene string, data interface{}) {
		r := records.Load()
		justVisiting := false

		outcome, ok := data.(run.Outcome)
		if !ok {
			justVisiting = true
			outcome = run.Outcome{R: r.LastRun}
		}
		dlog.Info("Just visting end game:", justVisiting)

		// STANDARD SETUP STUFF
		runInfo := outcome.R
		stayInEndScene = true
		endSceneNextScene = "inn"

		render.SetDrawStack(layer.Get()...)

		presentationX = float64(oak.ScreenWidth / 2)
		presentationY = float64(oak.ScreenHeight/2) + 20
		startY = float64(oak.ScreenHeight/2) + 20

		pitX = float64(oak.ScreenWidth) - 180.0
		pitY = float64(oak.ScreenHeight) - 75.0

		hopDistance = 100.0

		// Update info in the records with info from the most recent run
		// TODO: debate moving this to the run scene's end function

		sc := int64(runInfo.SectionsCleared)
		r.SectionsCleared += sc
		if sc > r.FarthestGoneInSections {
			r.FarthestGoneInSections = sc
		}

		// Display variables
		currentDeathToll := r.Deaths

		chestTotal := 0
		deadChests := 0
		if !justVisiting {
			for _, pl := range runInfo.Party.Players {
				playerChestValue := 0
				for _, j := range pl.ChestValues {
					playerChestValue += int(j)
				}
				if !pl.Alive {
					r.Deaths++
					deadChests += playerChestValue
				} else {
					chestTotal += playerChestValue
				}
			}
		}

		// For the next run TODO: move to run
		r.BaseSeed = int64(runInfo.SectionsCleared) + 1

		r.Wealth += chestTotal
		r.EnemiesDefeated += runInfo.EnemiesDefeated

		r.Store()

		fnt := render.DefFontGenerator.Copy()

		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()
		fnt.Color = render.FontColor("Black")
		fnt.Size = 28
		graves := fnt.Generate()

		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickColoredRTree(collision.DefTree, 4, labels.ColorMap)
		render.Draw(debugTree, layer.Play, 1000)

		// Make the graveyard backing
		endBackground, _ := render.LoadSprite("", filepath.Join("raw", "end_scene.png"))
		render.Draw(endBackground, layer.Ground)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := 16.0

		if !justVisiting {
			btn.New(menus.BtnCfgB,
				btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
				btn.Pos(menuX, menuY), btn.Text("Return To Inn"),
				btn.Binding(mouse.ClickOn, func(int, interface{}) int {
					stayInEndScene = false
					return 0
				}))
		}

		// totalChestValue := 20
		goldPit := floatgeom.NewRect2WH(670, float64(oak.ScreenHeight)-100, 330, 100)
		makeGoldParticles(r.Wealth, goldPit)

		//TODO: 2 layers: first current run
		// Second layer totals
		// Each layer has: sections_completed, enemies_defeated, chestvalue
		textY := 80.0
		textX := float64(oak.ScreenWidth) / 6

		titling := blueFnt.NewStrText("Last Run Info:", textX, textY)
		textX += 120
		render.Draw(titling, 2, 2)

		sectionText := blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(runInfo.SectionsCleared), textX, textY)
		textX += 200
		render.Draw(sectionText, 2, 2)

		enemy := blueFnt.NewStrText("Enemies Defeated: "+strconv.FormatInt(runInfo.EnemiesDefeated, 10), textX, textY)
		textX += 200
		render.Draw(enemy, 2, 2)

		chestValues := blueFnt.NewStrText("Chest Value: "+strconv.Itoa(chestTotal), textX, textY)
		textX += 200
		render.Draw(chestValues, 2, 2)

		// Current Run info
		textX = float64(oak.ScreenWidth) / 6
		textY += 30

		titling = blueFnt.NewStrText("Overall Stats:", textX, textY)
		textX += 120
		render.Draw(titling, 2, 2)

		sectionText = blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(int(r.SectionsCleared)), textX, textY)
		textX += 200
		render.Draw(sectionText, 2, 2)

		enemy = blueFnt.NewStrText("Enemies Defeated: "+strconv.FormatInt(r.EnemiesDefeated, 10), textX, textY)
		textX += 200
		render.Draw(enemy, 2, 2)

		chestValues = blueFnt.NewStrText("Chest Value: "+strconv.Itoa(r.Wealth), textX, textY)
		textX += 200
		render.Draw(chestValues, 2, 2)

		currentDeathTollp := &currentDeathToll
		render.Draw(graves.NewText(intStringer{currentDeathTollp}, 100, 350), 1, 10)

		debugElements := []render.Renderable{}
		// debug locations
		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{200, 200, 10, 255}))
		debugElements[0].SetPos(presentationX, presentationY)

		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{200, 200, 200, 255}))
		debugElements[1].SetPos(graveX, presentationY)

		debugElements = append(debugElements, render.NewColorBox(5, 5, color.RGBA{10, 200, 200, 255}))
		debugElements[2].SetPos(pitX, pitY)
		for _, r := range debugElements {
			render.Draw(r, layer.Debug, 18)
		}

		// A way to return to the inn scene
		dWidth := 150.0
		doodads.NewCustomInnDoor("inn", (float64(oak.ScreenWidth)-dWidth)/2, float64(oak.ScreenHeight)-20, dWidth, 20)
		// Block off the top of the inn from being walkable
		doodads.NewFurniture(0, 0, float64(oak.ScreenWidth), 187) // top of inn
		oak.ResetCommands()
		oak.AddCommand("debug", func(args []string) {
			dlog.Warn("Cheating to toggle debug mode")
			if debugTree.DrawDisabled {
				dlog.Warn("Debug turned off")
				debugTree.DrawDisabled = false
				for _, r := range debugElements {
					r.Undraw() // TODO: fix this
				}
				return
			}
			dlog.Warn("Debug turned on")
			debugTree.DrawDisabled = true
			for _, r := range debugElements {
				render.Draw(r, layer.Debug, 18)
			}

		})

		dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
		oak.AddCommand("help", func(args []string) {
			dlog.Info("Current Debug Commands are: ", strings.Join(oak.GetDebugKeys(), " , "))
		})

		if justVisiting == true {
			visitEnter(r.PartyComp)
			return
		}

		// Extra running info
		presentSpoils(runInfo.Party, currentDeathTollp, 0)

	},
	Loop: scene.BooleanLoop(&stayInEndScene),

	End: scene.GoToPtr(&endSceneNextScene),
}

// investigate allows us to investigate and poke around the end game with our living characters
func investigate(party *players.Party) {
	fmt.Println("Ending actions can be taken here")

}

func visitEnter(pComp []players.PartyMember) {
	ptycon := players.PartyConstructor{
		Players:    players.ClassConstructor(pComp),
		MaxPlayers: len(pComp),
	}

	pty, err := ptycon.NewParty(true)
	if err != nil {
		dlog.Error(err)
		return
	}
	pc := newEndWalker(npcScale, pty.Players)

	stop1 := float64(oak.ScreenHeight) - 140

	stop2 := float64(oak.ScreenWidth/2) - 200

	pc.Front.Delta = physics.NewVector(0, -8)
	pc.Front.Bind(func(id int, _ interface{}) int {
		p, ok := event.GetEntity(id).(*entities.Interactive)
		if !ok {
			dlog.Error("Non-player sent to player binding")
		}
		_, y := p.GetPos()
		if y < stop1 {
			pc.State = overridable

			p.RSpace.Add(collision.Label(labels.Door), (func(s1, s2 *collision.Space) {
				_, ok := s2.CID.E().(*doodads.InnDoor)
				if !ok {
					dlog.Error("Non-door sent to inndoor binding")
					return
				}
				stayInEndScene = false
			}))

			pc.Front.Delta = physics.NewVector(-4, 0)
			pc.Front.Bind(func(id int, _ interface{}) int {
				if pc.State == playing {
					investigate(pty)
					return event.UnbindSingle
				}
				p, ok := event.GetEntity(id).(*entities.Interactive)
				if !ok {
					dlog.Error("Non-player sent to player binding")
				}
				x, _ := p.GetPos()

				if x < stop2 {

					pc.State = playing
					investigate(pty)
					return event.UnbindSingle

				}

				return 0
			}, "EnterFrame")
			return event.UnbindSingle
		}

		return 0
	}, "EnterFrame")

}

func presentSpoils(party *players.Party, graveCount *int, index int) {
	if index > len(party.Players)-1 {
		investigate(party)
		return
	}
	p := party.Players[index]
	p.CID = p.Init()

	s := p.Swtch.Copy()
	p.Swtch = s.Modify(mod.Scale(npcScale, npcScale)).(*render.Switch)
	p.Interactive.R = p.Swtch
	//Player enters stage right
	p.SetPos(float64(oak.ScreenWidth-64), startY)
	render.Draw(p.R, layer.Play, 20)
	p.ChestsHeight = 0
	for _, r := range p.Chests {
		_, h := r.GetDims()
		p.ChestsHeight += float64(h)
		chestHeight := p.ChestsHeight
		r.(*render.Sprite).Vector = r.Attach(p.Vector, -3, -chestHeight)
		render.Draw(r, layer.Play, 21)
	}

	fmt.Printf("\nCharacter %d walking through and is Alive:%t ", index, p.Alive)

	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.ShiftPos(-2, 0)

		p.Swtch.Set("walkLT")
		if len(ply.ChestValues) > 0 {
			p.Swtch.Set("walkHold")
		}
		if !ply.Alive {
			p.Swtch.Set("deadLT")

		}

		// TODO: the following
		// Player comes to middle ish

		// Do something with spoils or look sad if dead
		// If alive: throw chests with hop and cheer
		//		Card explaining the amount in their chests?
		//		Chest explodes into money and goes into pit
		// If dead: whomp whomp
		// 		eulogoy? name, class, run?

		// Next person in party starts process
		// You walk to end point (graves or bottom)
		// When reach your end point destroy self

		if ply.R.X() < presentationX {

			if p.Alive {
				p.Swtch.Set("standRT")
				if len(p.ChestValues) > 0 {
					hop(p)

				} else {
					sfx.Play("ohWell")
				}

				t := time.Now().Add(time.Second)
				p.CheckedBind(func(ply *players.Player, _ interface{}) int {

					if time.Now().After(t) {
						livingExit(p)
						return event.UnbindSingle
					}
					return 0
				}, "EnterFrame")
			} else {
				*graveCount++
				deadMovement(p)
			}

			presentSpoils(party, graveCount, index+1)
			return event.UnbindSingle

		}
		return 0
	}, "EnterFrame")

}

func hop(p *players.Player) {
	sfx.Play("chestHop1")
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {
		if ply.Y() < presentationY-hopDistance {
			tossChests(ply)
			ply.CheckedBind(func(plyz *players.Player, _ interface{}) int {

				if plyz.Y() > presentationY {

					return event.UnbindSingle
				}
				plyz.ShiftPos(0, 4)
				return 0
			}, "EnterFrame")
			return event.UnbindSingle
		}

		ply.ShiftPos(0, -4)
		return 0
	}, "EnterFrame")
}
func tossChests(p *players.Player) {

	for _, c := range p.Chests {
		c.(*render.Sprite).Vector = c.(*render.Sprite).Vector.Detach()
	}
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		done := false
		for _, cr := range p.Chests {

			cr.ShiftX((pitX - presentationX) / 100.0)
			cr.ShiftY((pitY - presentationY + hopDistance) / 100.0)

			if cr.X() > pitX {
				explodeChest(cr.X(), cr.Y())
				cr.Undraw()
				done = true
			}
		}
		if done {
			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")
}
func livingExit(p *players.Player) {
	p.Swtch.Set("walkLT")
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.ShiftPos(0, 2)
		if ply.R.Y() > float64(oak.ScreenHeight) {

			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")
}

func deadMovement(p *players.Player) {
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		// ply.R.Undraw()

		ply.ShiftPos(-2, 0)
		if ply.R.X() < graveX {
			deathSprites(ply.R.X(), ply.R.Y())
			sfx.Play("dissappear1")
			ply.R.Undraw()
			return event.UnbindSingle
		}

		return 0
	}, "EnterFrame")

}

func deathSprites(x, y float64) {
	ptGenLife := intrange.NewLinear(40, 60)
	ptColor := color.RGBA{200, 200, 200, 255}
	ptColorRand := color.RGBA{0, 0, 0, 0}
	newPf := floatrange.NewLinear(5, 10)
	ptLife := floatrange.NewLinear(100, 200)
	angle := floatrange.NewLinear(0, 360)
	speed := floatrange.NewLinear(1, 4)
	size := intrange.Constant(3)
	layerFn := func(v physics.Vector) int {
		return layer.Effect
	}
	particle.NewColorGenerator(
		particle.Pos(x, y),
		particle.Duration(ptGenLife),
		particle.LifeSpan(ptLife),
		particle.Angle(angle),
		particle.Speed(speed),
		particle.Layer(layerFn),
		particle.Shape(shape.Square),
		particle.Size(size),
		particle.Color(ptColor, ptColorRand, ptColor, ptColorRand),
		particle.NewPerFrame(newPf)).Generate(0)
}

func explodeChest(x, y float64) {
	sp, err := render.GetSprite(filepath.Join("raw", "wood_junk.png"))
	if err != nil {
		dlog.Error(err)
		return
	}
	sfx.Play("chestExplode")
	explodeSprite(x, y, sp)

}
func explodeSprite(x, y float64, sprite *render.Sprite) {
	layerFn := func(v physics.Vector) int {
		return layer.Effect
	}
	ptGenLife := intrange.NewLinear(40, 60)
	sg := particle.NewSpriteGenerator(
		particle.NewPerFrame(floatrange.NewSpread(8, 0)),
		particle.Pos(x, y),
		particle.LifeSpan(floatrange.NewSpread(20, 5)),
		particle.Angle(floatrange.NewSpread(0, 360)),
		particle.Speed(floatrange.NewSpread(2, .5)),
		particle.Spread(3, 2),
		particle.Duration(ptGenLife),
		particle.Layer(layerFn),
		particle.Sprite(sprite),
		particle.SpriteRotation(floatrange.Constant(10)),
	)
	sg.Generate(layer.Effect)
}

// makeGoldParticles creates the appropriate amount of collision particles within the given location
func makeGoldParticles(goldCount int, location floatgeom.Rect2) {
	debug := collision.NewRect2Space(location, 0)
	debug.UpdateLabel(collision.Label(labels.Ornament))
	collision.Add(debug)

	center := location.Center()

	//TODO: make this an actual fxn probably making it a log of goldCount
	particleCount := int(math.Log(float64(10.0 * goldCount)))

	colorOpts := particle.And(
		particle.NewPerFrame(floatrange.NewConstant(float64(particleCount))),
		particle.Limit(particleCount),
		particle.InfiniteLifeSpan(),
		particle.Spread(location.W()/2+8, location.H()/2),
		particle.Shape(shape.Diamond),
		particle.Size(intrange.NewConstant(4)),
		particle.Speed(floatrange.NewConstant(0)),
		particle.Pos(center.X(), center.Y()),
		particle.Color(color.RGBA{200, 200, 0, 255}, color.RGBA{0, 0, 0, 0},
			color.RGBA{200, 200, 0, 255}, color.RGBA{0, 0, 0, 0}),
	)
	shiftFactor := floatrange.NewLinear(0, 6)
	pg := particle.NewCollisionGenerator(
		particle.NewColorGenerator(colorOpts),
		particle.Fragile(false),
		particle.HitMap(map[collision.Label]collision.OnHit{
			labels.PC: func(a, b *collision.Space) {
				// b.CID.Trigger("Attacked", hitEffects)
				p, ok := event.GetEntity(int(b.CID)).(*entities.Interactive)
				if !ok {
					dlog.Error("A non player is colliding with gold?")
					return
				}
				goldPiece := particle.Lookup(int(a.CID))
				d := p.Delta.Copy().Scale(shiftFactor.Poll())
				goldPiece.ShiftX(d.X())
				goldPiece.ShiftY(d.Y())
			},
		}),
	)
	pg.Generate(layer.Play)

}
