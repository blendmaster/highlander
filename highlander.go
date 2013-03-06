/* Highlander: prominent topology for prominent people
 * Author: Steven Ruppert
 * For CSCI 447 Scientific Visualization, Spring 2013
 * Colorado School Of Mines
 */
package main

import (
	"fmt"
	"log"
	//"image"
	"image/color"
	"image/png"
	//"io"
	"highlander/prominence"
	"os"
)

func main() {
	img, err := png.Decode(os.Stdin)
	if err != nil {
		log.Fatalln("Couldn't read png image from Stdin!", err)
	}


	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			fmt.Println(img.At(x, y).(color.Gray16).Y)
		}
	}

	// TODO pass threshold and image as command line args
	fmt.Println(prominence.ProminentFeatures(img), 10)
}
