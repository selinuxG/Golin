package crack

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/sijms/go-ora/v2"
)

func oracle(ctx context.Context, cancel context.CancelFunc, ip, user, passwd string, port int) {
	defer func() {
		wg.Done()
		<-ch
	}()
	select {
	case <-ctx.Done():
		return
	default:
	}

	dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s", user, passwd, ip, port, "orcl")
	db, err := sql.Open("oracle", dsn)
	if err == nil {
		a, err := db.Query("select * from v$version")
		if err == nil {
			fmt.Println(dsn)
			fmt.Println(a)
			end(ip, user, passwd, port)
			//cancel()
		}
	}

}
