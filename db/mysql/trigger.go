package mysql

import (
	"fmt"

	"github.com/feitian124/goapi/database/schema"
	"github.com/pkg/errors"
)

const triggerSQL = `
	SELECT
	  trigger_name,
	  action_timing,
	  event_manipulation,
	  event_object_table,
	  action_orientation,
	  action_statement
	FROM information_schema.triggers
	WHERE event_object_schema = ?
	AND event_object_table = ?
`

func (d *DB) Triggers(tableName string) ([]*schema.Trigger, error) {
	triggerRows, err := d.db.Query(triggerSQL, d.Schema.Name, tableName)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer triggerRows.Close()
	var triggers []*schema.Trigger
	for triggerRows.Next() {
		var (
			triggerName              string
			triggerActionTiming      string
			triggerEventManipulation string
			triggerEventObjectTable  string
			triggerActionOrientation string
			triggerActionStatement   string
			triggerDef               string
		)
		err = triggerRows.Scan(&triggerName, &triggerActionTiming, &triggerEventManipulation, &triggerEventObjectTable,
			&triggerActionOrientation, &triggerActionStatement)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		triggerDef = fmt.Sprintf("CREATE TRIGGER %s %s %s ON %s\nFOR EACH %s\n%s", triggerName, triggerActionTiming,
			triggerEventManipulation, triggerEventObjectTable, triggerActionOrientation, triggerActionStatement)
		trigger := &schema.Trigger{
			Name: triggerName,
			Def:  triggerDef,
		}
		triggers = append(triggers, trigger)
	}
	return triggers, nil
}
