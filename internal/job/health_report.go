package job

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tenz-io/gokit/logger"
)

type HealthReporter struct {
	c *cron.Cron
}

func NewHealthReporter(
	c *cron.Cron,
) *HealthReporter {
	return &HealthReporter{
		c: c,
	}
}

func (wg *HealthReporter) Init() error {
	log.Println("[HealthReporter] Init...")

	//generate weekly review every week on Saturday 23:40
	if entryID, err := wg.c.AddFunc("40 23 * * 6", wg.Generate); err != nil {
		return fmt.Errorf("add cron job(%d) error: %w", entryID, err)
	}

	wg.c.Start()

	return nil
}

func (wg *HealthReporter) Generate() {
	var (
		ctx = context.Background()
		now = time.Now()
		le  = logger.FromContext(ctx).WithFields(logger.Fields{
			"now": now,
		})
	)
	le.Info("start generate health report")

	return
}
