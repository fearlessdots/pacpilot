package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"fmt"
	"strings"

	// External modules
	color "github.com/gookit/color"
	wordwrap "github.com/mitchellh/go-wordwrap"
)

//
//// DISPLAY VARIABLES
//

var (
	// Gray
	grayHex = "#808080"
	gray    = color.HEX(grayHex)

	// Light Gray
	lightGrayHex = "#c8c4a9"
	lightGray    = color.HEX(lightGrayHex)

	// Orange
	orangeHex = "#ffa860"
	orange    = color.HEX(orangeHex)

	// Blue
	blueHex = "#55aaff"
	blue    = color.HEX(blueHex)

	// Green
	greenHex = "#55ff7f"
	green    = color.HEX(greenHex)

	// Red
	redHex = "#ff5050"
	red    = color.HEX(redHex)

	// Light Copper
	lightCopperHex = "#ffaa7f"
	lightCopper    = color.HEX(lightCopperHex)

	// Salmon Pink
	salmonPinkHex = "#ff8a82"
	salmonPink    = color.HEX(salmonPinkHex)

	// Coral
	coralHex = "#ff7458"
	coral    = color.HEX(coralHex)

	// Light Brown
	paleLimeHex = "#cee89c"
	paleLime    = color.HEX(paleLimeHex)
)

//
//// DISPLAY FUNCTIONS
//

func returnText(msg string, indentLevel int) string {
	indent := strings.Repeat("    ", indentLevel)

	return fmt.Sprintf(indent + msg)
}

func showText(msg string, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)

	fmt.Println(indent + msg)
}

func returnTextWrapped(msg string, factor float64, program Program) string {
	terminalDimensions, _ := getTerminalDimensions(program)

	wrapLimit := uint(float64(terminalDimensions.width) * factor)
	wrappedText := wordwrap.WrapString(msg, wrapLimit)

	return wrappedText
}

func showTextWrapped(msg string, factor float64, program Program) {
	terminalDimensions, _ := getTerminalDimensions(program)

	wrapLimit := uint(float64(terminalDimensions.width) * factor)
	wrappedText := wordwrap.WrapString(msg, wrapLimit)

	fmt.Println(wrappedText)
}

func printColoredMessage(c color.RGBColor, msg string, indentLevel int) {
	indent := strings.Repeat("    ", indentLevel)
	c.Println(indent + msg)
}

func showInfo(msg string, indentLevel int) {
	printColoredMessage(gray, msg, indentLevel)
}

func showAttention(msg string, indentLevel int) {
	printColoredMessage(orange, msg, indentLevel)
}

func showInfoSectionTitle(msg string, indentLevel int) {
	printColoredMessage(lightGray, msg, indentLevel)
}

func showSuccess(msg string, indentLevel int) {
	printColoredMessage(blue, msg, indentLevel)
}

func showError(msg string, indentLevel int) {
	printColoredMessage(red, msg, indentLevel)
}

func space() {
	fmt.Println("")
}

func hr(char string, factor float64, program Program) {
	terminalDimensions, _ := getTerminalDimensions(program)

	horizontalLine := strings.Repeat(string(char), int(float64(terminalDimensions.width)*factor))
	showText(horizontalLine, program.indentLevel)
}
