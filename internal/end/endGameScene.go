package end

import (
	"fmt"
	"image/color"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/200sc/go-dist/floatrange"
	"github.com/200sc/go-dist/intrange"
	"github.com/oakmound/oak/collision"
	"github.com/oakmound/oak/dlog"
	"github.com/oakmound/oak/event"
	"github.com/oakmound/oak/mouse"
	"github.com/oakmound/oak/physics"
	"github.com/oakmound/oak/render/particle"
	"github.com/oakmound/oak/shape"

	"github.com/oakmound/oak"
	"github.com/oakmound/oak/entities/x/btn"
	"github.com/oakmound/oak/render"
	"github.com/oakmound/oak/scene"
	"github.com/oakmound/weekly87/internal/characters/labels"
	"github.com/oakmound/weekly87/internal/characters/players"
	"github.com/oakmound/weekly87/internal/dtools"
	"github.com/oakmound/weekly87/internal/layer"
	"github.com/oakmound/weekly87/internal/menus"
	"github.com/oakmound/weekly87/internal/records"
	"github.com/oakmound/weekly87/internal/run"
)

func Init() {
	_, err := render.LoadSprite("", filepath.Join("raw", "wood_junk.png"))
	dlog.Error("OH MY " + err.Error())
}

var (
	graveX = 90.0
	graveY = 300.0
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

		outcome := data.(run.Outcome)
		fmt.Println(outcome.R)
		runInfo := outcome.R
		stayInEndScene = true
		endSceneNextScene = "inn"

		render.SetDrawStack(layer.Get()...)

		presentationX = float64(oak.ScreenWidth / 2)
		presentationY = float64(oak.ScreenHeight/2) + 20
		startY = float64(oak.ScreenHeight/2) + 20

		pitX = float64(oak.ScreenWidth) - 180.0
		pitY = float64(oak.ScreenHeight) - 100.0

		hopDistance = 100.0
		// partyJson, _ := json.Marshal(runInfo.Party)

		// Update info in the records with info from the most recent run
		// TODO: debate moving this to the run scene's end function
		r := records.Load()
		sc := int64(runInfo.SectionsCleared)
		r.SectionsCleared += sc
		if sc > r.FarthestGoneInSections {
			r.FarthestGoneInSections = sc
		}

		// Display variables
		currentDeathToll := r.Deaths

		chestTotal := 0
		deadChests := 0
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

		// For the next run
		r.BaseSeed = int64(runInfo.SectionsCleared) + 1
		r.Store()

		fnt := render.DefFontGenerator.Copy()

		fnt.Color = render.FontColor("Blue")
		fnt.Size = 14
		blueFnt := fnt.Generate()
		fnt.Color = render.FontColor("Black")
		fnt.Size = 28
		graves := fnt.Generate()

		// nextscene = "inn"
		render.SetDrawStack(layer.Get()...)
		debugTree := dtools.NewThickColoredRTree(collision.DefTree, 4, labels.ColorMap)
		render.Draw(debugTree, layer.Play, 1000)

		// Make the Inn backing
		innBackground, _ := render.LoadSprite("", filepath.Join("raw", "end_scene.png"))
		render.Draw(innBackground, layer.Ground)

		// textBacking := render.NewColorBox(textBackingX, oak.ScreenHeight*2/3, color.RGBA{120, 120, 120, 190})
		// textBacking.SetPos(float64(oak.ScreenWidth)*0.33, 40)
		// render.Draw(textBacking, 1)

		menuX := (float64(oak.ScreenWidth) - menus.BtnWidthA) / 2
		menuY := 16.0

		btn.New(menus.BtnCfgB,
			btn.TxtOff(menus.BtnWidthA/8, menus.BtnHeightA/3),
			btn.Pos(menuX, menuY), btn.Text("Return To Inn"),
			btn.Binding(mouse.ClickOn, func(int, interface{}) int {
				stayInEndScene = false
				return 0
			}))

		textY := 120.0

		//TODO: 2 layers: first current run
		// Second layer totals
		// Each layer has: sections_completed, enemies_defeated, chestvalue

		textX := float64(oak.ScreenWidth) / 8

		titling := blueFnt.NewStrText("Last Run Info:", textX, textY)
		textX += 120
		render.Draw(titling, 2, 2)

		sectionText := blueFnt.NewStrText("Sections Cleared: "+strconv.Itoa(runInfo.SectionsCleared), textX, textY)
		textX += 200
		render.Draw(sectionText, 2, 2)

		enemy := blueFnt.NewStrText("Enemies Defeated: "+strconv.Itoa(runInfo.EnemiesDefeated), textX, textY)
		textX += 200
		render.Draw(enemy, 2, 2)

		chestValues := blueFnt.NewStrText("Chest Value: "+strconv.Itoa(chestTotal), textX, textY)
		textX += 200
		render.Draw(chestValues, 2, 2)

		currentDeathTollp := &currentDeathToll
		render.Draw(graves.NewText(intStringer{currentDeathTollp}, 180, 300), 1, 10)

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

		// Extra running info
		presentSpoils(runInfo.Party, currentDeathTollp, 0)

	},
	Loop: scene.BooleanLoop(&stayInEndScene),

	End: scene.GoToPtr(&endSceneNextScene),
}

func presentSpoils(party *players.Party, graveCount *int, index int) {
	if index > len(party.Players)-1 {
		return
	}
	p := party.Players[index]
	p.CID = p.Init()

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

		ply.Swtch.Set("walkLT")
		if len(ply.ChestValues) > 0 {
			ply.Swtch.Set("walkHold")
		}
		if !ply.Alive {
			ply.Swtch.Set("deadLT")

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

				if len(p.ChestValues) > 0 {
					hop(p)

				}

				t := time.Now().Add(time.Second)
				p.CheckedBind(func(ply *players.Player, _ interface{}) int {

					if time.Now().After(t) {
						liviningExit(p)
						return event.UnbindSingle
					}
					return 0
				}, "EnterFrame")
			} else {
				*graveCount++
				fmt.Printf("Graves is now %d\n ", *graveCount)
				deadMovement(p)
			}

			presentSpoils(party, graveCount, index+1)
			return event.UnbindSingle

		}
		return 0
	}, "EnterFrame")

}

func hop(p *players.Player) {
	p.Swtch.Set("standRT")
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {
		if ply.Y() < presentationY-hopDistance {
			tossChests(ply)
			ply.CheckedBind(func(plyz *players.Player, _ interface{}) int {

				if plyz.Y() > presentationY {
					plyz.Swtch.Set("walkLT")
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
func liviningExit(p *players.Player) {
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		ply.ShiftPos(0, 2)
		if ply.R.Y() > float64(oak.ScreenHeight) {

			return event.UnbindSingle
		}
		return 0
	}, "EnterFrame")
}

// var deathParticles particle.Generator

func deadMovement(p *players.Player) {
	p.CheckedBind(func(ply *players.Player, _ interface{}) int {

		// ply.R.Undraw()

		ply.ShiftPos(-2, 0)
		if ply.R.X() < graveX {
			deathSprites(ply.R.X(), ply.R.Y())
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
