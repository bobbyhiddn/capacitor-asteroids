/*
 * I know this file looks awful, but I'm not sure how else to generate graphics with lines.
 * If someone knows a better way to do this, please let me know, or submit a pull request.
 */
package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

func DrawSpaceship(screen *ebiten.Image, centerX, centerY, angle float64) {
	// All the "window_height/x" is used to scale the ship to the window size
	// TODO: make an easier way to scale items to the window size
	ship_point_x := centerX + math.Cos(angle)*float64(window_height/36)
	ship_point_y := centerY + math.Sin(angle)*float64(window_height/36)

	rear_line_center_x := centerX - math.Cos(angle)*float64(window_height/114)
	rear_line_center_y := centerY - math.Sin(angle)*float64(window_height/114)

	line_one_x := math.Cos(angle-math.Pi/8) * float64(window_height/20)
	line_one_y := math.Sin(angle-math.Pi/8) * float64(window_height/20)
	line_two_x := math.Cos(angle+math.Pi/8) * float64(window_height/20)
	line_two_y := math.Sin(angle+math.Pi/8) * float64(window_height/20)
	line_three_x := math.Cos(angle-math.Pi/2) * float64(window_height/62)
	line_three_y := math.Sin(angle-math.Pi/2) * float64(window_height/62)
	line_four_x := math.Cos(angle+math.Pi/2) * float64(window_height/62)
	line_four_y := math.Sin(angle+math.Pi/2) * float64(window_height/62)

	ebitenutil.DrawLine(screen, ship_point_x, ship_point_y, ship_point_x-line_one_x, ship_point_y-line_one_y, color.White)
	ebitenutil.DrawLine(screen, ship_point_x, ship_point_y, ship_point_x-line_two_x, ship_point_y-line_two_y, color.White)
	ebitenutil.DrawLine(screen, rear_line_center_x, rear_line_center_y, rear_line_center_x-line_three_x, rear_line_center_y-line_three_y, color.White)
	ebitenutil.DrawLine(screen, rear_line_center_x, rear_line_center_y, rear_line_center_x-line_four_x, rear_line_center_y-line_four_y, color.White)
}

func DrawThrusters(screen *ebiten.Image, centerX, centerY, angle float64) {
	thruster_point_x := centerX - math.Cos(angle)*float64(window_height/40)
	thruster_point_y := centerY - math.Sin(angle)*float64(window_height/40)

	line_one_x := -math.Cos(angle-math.Pi/8) * float64(window_height/60)
	line_one_y := -math.Sin(angle-math.Pi/8) * float64(window_height/60)
	line_two_x := -math.Cos(angle+math.Pi/8) * float64(window_height/60)
	line_two_y := -math.Sin(angle+math.Pi/8) * float64(window_height/60)

	ebitenutil.DrawLine(screen, thruster_point_x, thruster_point_y, thruster_point_x-line_one_x, thruster_point_y-line_one_y, color.White)
	ebitenutil.DrawLine(screen, thruster_point_x, thruster_point_y, thruster_point_x-line_two_x, thruster_point_y-line_two_y, color.White)
}

// TODO: Replace these functions with ones like above, generating the lines
func DrawAsteroidSmall(screen *ebiten.Image, centerX, centerY, angle float64) {
	ebitenutil.DrawRect(screen, centerX, centerY, 10, 10, color.White)
}

func DrawAsteroidMedium(screen *ebiten.Image, centerX, centerY, angle float64) {
	ebitenutil.DrawRect(screen, centerX, centerY, 20, 20, color.White)
}

func DrawAsteroidLarge(screen *ebiten.Image, centerX, centerY, angle float64) {
	ebitenutil.DrawRect(screen, centerX, centerY, 30, 30, color.White)
}
