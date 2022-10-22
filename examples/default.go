package main
import "github.com/athanbase/log"

func main() {
	defer log.Sync()
	log.Info("demo: 1", log.String("app", "start ok"), log.Int("major version", 2))
}