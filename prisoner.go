package prisoner

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

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
	source := rand.NewSource(time.Now().UnixNano())
	rand.New(source).Shuffle(len(boxes), func(i int, j int) {
		boxes[i], boxes[j] = boxes[j], boxes[i]
	})

	// logger.Info("Extract boxes", boxes, len(boxes))

	var wg sync.WaitGroup
	pch := make(chan int)

	// each prisoner, should only be allowed to
	// find their number in the loop at least 50%
	// of their population
	logger.Info("Prisoner's starting to find their boxes!")
	for p := 0; p < prisoners_count; p++ {
		// each prisoner, they can choose which box to start with
		source := rand.NewSource(time.Now().UnixNano())
		box_idx := rand.New(source).Intn(len(boxes) - 1)
		wg.Add(1)
		go iAmRunningLikeIdiotFindingMyNumber(pch, &wg, p, boxes, box_idx, 0)
	}

	go func() {
		wg.Wait()
		close(pch)
	}()

	result := &Result{}

	for prisoner_number := range pch {
		// it means, this prisoner cant find their number
		if prisoner_number == -1 {
			continue
		}

		result.PrisonersWhoFoundTheirNumber = append(result.PrisonersWhoFoundTheirNumber, prisoner_number)
	}

	res := (float32(len(result.PrisonersWhoFoundTheirNumber)) / float32(prisoners_count)) * 100
	logger.Info(fmt.Sprintf("Probability: %.2f%%", res))

	return nil
}

func iAmRunningLikeIdiotFindingMyNumber(pch chan int, wg *sync.WaitGroup, prisoner int, boxes []int, box_idx int, loop int) {
	/*if loop == 0 {
		if boxes[box_idx] == prisoner {
			logger.Info(fmt.Sprintf("Prisoner [%d], found it immediately under box [%d]", prisoner, box_idx))
			return
		} else {
		    logger.Info(fmt.Sprintf("Prisoner [%d]:", prisoner))
		}
	}*/
	// logger.Info(fmt.Sprintf(" -> [box %d] contains [%d]", box_idx, boxes[box_idx]))

	if boxes[box_idx] == prisoner {
		// logger.Info(" -> Found it!")
		defer wg.Done()
		pch <- prisoner
		logger.Info(fmt.Sprintf("Prisoner [%d] found their number at box [%d]", prisoner, box_idx))
		return
	}

	if loop == (len(boxes) / 2) {
		// logger.Info(" -> Oh no, prisoner is not allowed to open a box anymore!")
		defer wg.Done()
		pch <- -1
		logger.Info(fmt.Sprintf("Prisoner [%d] cant find the number", prisoner))
		return
	}

	loop += 1
	next_box_id := boxes[box_idx]
	iAmRunningLikeIdiotFindingMyNumber(pch, wg, prisoner, boxes, next_box_id, loop)
}
