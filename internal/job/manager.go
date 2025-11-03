package job

import "log"

type Manager struct {
	weeklyReviewGenerator *HealthReporter
}

func NewManager(
	weeklyReviewGenerator *HealthReporter,
) *Manager {
	return &Manager{
		weeklyReviewGenerator: weeklyReviewGenerator,
	}
}

func (m *Manager) Init() error {
	log.Println("[Job Manager] Init...")
	if err := m.weeklyReviewGenerator.Init(); err != nil {
		return err
	}

	return nil
}
