package prisoner

import (
	"fmt"
	"math/rand"

	"github.com/lucidfy/lucid/pkg/facade/logger"
	cli "github.com/urfave/cli/v2"
)

type PrisonerLoop struct {
	Command *cli.Command
}

type Result struct {
	PrisonersWhoFoundTheirNumber []int
}

func Prisoner() *PrisonerLoop {
	var cc PrisonerLoop
	cc.Command = &cli.Command{
		Name:   "prisoner",
		Usage:  "Spawn Prisoners and calculate",
		Action: cc.Handle,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "count",
				Value: 100,
				Usage: `Number of prisoners (i.e: 100)`,
			},
		},
	}
	return &cc
}

func (cc PrisonerLoop) Handle(c *cli.Context) error {
	var boxes []int
	prisoners_count := c.Int("count")
	// max_visit := prisoners_count / 2

	// now, append those prisoners' number and place them on each boxes
	logger.Info("Preparing boxes!")
	for i := 0; i < prisoners_count; i++ {
		boxes = append(boxes, i)
	}

	// randomize the content of the boxes
	logger.Info("Randomizing the box contents!")
	rand.Shuffle(len(boxes), func(i int, j int) {
		boxes[i], boxes[j] = boxes[j], boxes[i]
	})

	// logger.Info("Extract boxes", boxes, len(boxes))

	result := &Result{}

	// each prisoner, should only be allowed to
	// find their number in the loop at least 50%
	// of their population
	logger.Info("Prisoner's starting to find their boxes!")
	for p := 0; p < prisoners_count; p++ {
		// each prisoner, they can choose which box to start with
		// source := rand.NewSource(time.Now().UnixNano())
		box_idx := rand.Intn(len(boxes) - 1)
		result = iAmRunningLikeIdiotFindingMyNumber(result, p, boxes, box_idx, 0)
	}

	res := (float32(len(result.PrisonersWhoFoundTheirNumber)) / float32(prisoners_count)) * 100
	logger.Info(fmt.Sprintf("Probability: %.2f%%", res))

	return nil
}

func iAmRunningLikeIdiotFindingMyNumber(result *Result, prisoner int, boxes []int, box_idx int, loop int) *Result {
	if loop == 0 {
		if boxes[box_idx] == prisoner {
			// logger.Info(fmt.Sprintf("Prisoner [%d], found it immediately under box [%d]", prisoner, box_idx))
			result.PrisonersWhoFoundTheirNumber = append(result.PrisonersWhoFoundTheirNumber, prisoner)
			return result
		}
		// else {
		// 	logger.Info(fmt.Sprintf("Prisoner [%d]:", prisoner))
		// }
	}

	// logger.Info(fmt.Sprintf(" -> [box %d] contains [%d]", box_idx, boxes[box_idx]))

	if boxes[box_idx] == prisoner {
		// logger.Info(" -> Found it!")
		result.PrisonersWhoFoundTheirNumber = append(result.PrisonersWhoFoundTheirNumber, prisoner)
		return result
	}

	if loop == (len(boxes) / 2) {
		// logger.Info(" -> Oh no, prisoner is not allowed to open a box anymore!")
		return result
	}

	loop += 1
	next_box_id := boxes[box_idx]
	iAmRunningLikeIdiotFindingMyNumber(result, prisoner, boxes, next_box_id, loop)
	return result
}
