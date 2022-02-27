package rollout

type Component interface {
	Name() string
	Rollout() error
	Healthcheck() error
	Rollback() error
}

type ErrComponent struct {
	Message       string
	ComponentName string
}

func (e ErrComponent) Error() string {
	return e.Message
}

func OrderedRollout(components []Component) error {
	completedComponents := make([]Component, 0)

	for _, c := range components {
		rolloutErr := c.Rollout()
		if rolloutErr != nil {
			if len(completedComponents) == 0 {
				rollbackErr := RollbackComponents([]Component{c})
				if rollbackErr != nil {
					return ErrComponent{
						Message:       rollbackErr.Error(),
						ComponentName: c.Name(),
					}
				}
			} else {
				rollbackComponents := append(completedComponents, c)
				rollbackErr := RollbackComponents(rollbackComponents)
				if rollbackErr != nil {
					return ErrComponent{
						Message:       rollbackErr.Error(),
						ComponentName: c.Name(),
					}
				}
			}
		}

		healthErr := c.Healthcheck()
		if healthErr != nil {
			if len(completedComponents) == 0 {
				rollbackErr := RollbackComponents([]Component{c})
				if rollbackErr != nil {
					return ErrComponent{
						Message:       rollbackErr.Error(),
						ComponentName: c.Name(),
					}
				}
			} else {
				rollbackComponents := append(completedComponents, c)
				rollbackErr := RollbackComponents(rollbackComponents)
				if rollbackErr != nil {
					return ErrComponent{
						Message:       rollbackErr.Error(),
						ComponentName: c.Name(),
					}
				}
			}
		}

		completedComponents = append(completedComponents, c)
	}

	return nil
}

// RollbackComponents attempts to roll back all components
// will stop and fail if a rollback errors to allow for
// human intervention.
func RollbackComponents(components []Component) error {
	for _, c := range components {
		err := c.Rollback()
		if err != nil {
			return err
		}
	}
	return nil
}
