package initresources

import (
	"OrderManagementSystem/internal/db"
	"sync"
)

var initOnce sync.Once

func InitAll() {
	initOnce.Do(func() {
		_ = db.SqlConnection()
	})

}
