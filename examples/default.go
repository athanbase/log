package main
import "github.com/athanbase/log"

func main() {
	defer log.Sync()
	log.Info("demo1: ", log.String("app", "start ok"), log.Int("version", 2))
}