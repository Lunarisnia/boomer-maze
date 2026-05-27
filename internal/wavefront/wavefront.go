package wavefront

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lunarisnia/boomer-maze/internal/gmath"
)

type Object struct {
	// v 0.11526 0.700717 0.0677257
	Vertice []gmath.Vector3[float64]

	// NOTE: We will be only storing the first number in each axis for now
	// f 2491/2662/2491 2519/2674/2519 2492/2659/2492
	Faces []gmath.Vector3[int]
}

func parseRow(obj *Object, row string) error {
	detail := strings.Fields(row)
	if len(detail) == 0 {
		return nil
	}

	identifier := detail[0]
	switch identifier {
	case "v":
		if len(detail) < 4 {
			return fmt.Errorf("invalid vertex row: %q", row)
		}

		x, err := strconv.ParseFloat(detail[1], 64)
		if err != nil {
			return fmt.Errorf("invalid vertex x %q: %w", detail[1], err)
		}

		y, err := strconv.ParseFloat(detail[2], 64)
		if err != nil {
			return fmt.Errorf("invalid vertex y %q: %w", detail[2], err)
		}

		z, err := strconv.ParseFloat(detail[3], 64)
		if err != nil {
			return fmt.Errorf("invalid vertex z %q: %w", detail[3], err)
		}

		obj.Vertice = append(obj.Vertice, gmath.Vector3[float64]{X: x, Y: y, Z: z})
	case "f":
		if len(detail) < 4 {
			return fmt.Errorf("invalid face row: %q", row)
		}

		firstAxis := func(detail string) string {
			return strings.Split(detail, "/")[0]
		}

		x, err := strconv.Atoi(firstAxis(detail[1]))
		if err != nil {
			return fmt.Errorf("invalid face x %q: %w", detail[1], err)
		}

		y, err := strconv.Atoi(firstAxis(detail[2]))
		if err != nil {
			return fmt.Errorf("invalid face y %q: %w", detail[2], err)
		}

		z, err := strconv.Atoi(firstAxis(detail[3]))
		if err != nil {
			return fmt.Errorf("invalid face z %q: %w", detail[3], err)
		}

		obj.Faces = append(obj.Faces, gmath.Vector3[int]{X: x - 1, Y: y - 1, Z: z - 1})
	}

	return nil
}

func LoadModel(model string) (*Object, error) {
	fileByte, err := os.ReadFile(model)
	if err != nil {
		return nil, err
	}

	rows := strings.SplitSeq(string(fileByte), "\n")

	loadedObj := Object{}
	for row := range rows {
		err := parseRow(&loadedObj, row)
		if err != nil {
			return nil, err
		}
	}

	return &loadedObj, nil
}
