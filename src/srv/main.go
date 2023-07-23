package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

type Busstop struct {
	name  string
	pairs [][2]byte
}

type AppCtx struct {
	busstops []Busstop
	n        int
}

func normalizeBusstop(busstop string) string {
	stop := strings.ToLower(busstop)
	stop = strings.ReplaceAll(stop, " ", "")
	stop = strings.ReplaceAll(stop, ",", "")
	return stop
}

func pairedBusstops(busstop string) [][2]byte {
	pairs := make([][2]byte, len(busstop)-1)

	for j := 0; j < len(busstop)-1; j++ {
		pairs[j] = [2]byte{busstop[j], busstop[j+1]}
	}
	return pairs
}

func prepareBusstops(busstops []string, n int) []Busstop {
	pairs := make([]Busstop, len(busstops))

	for i := 0; i < n; i++ {
		stop := normalizeBusstop(busstops[i])
		pairs[i] = Busstop{busstops[i], pairedBusstops(stop)}
	}
	return pairs
}

func intersection(a, b [][2]byte) int {
	c := 0
	for _, x := range a {
		for _, y := range b {
			if x == y {
				c++
			}
		}
	}
	return c
}

func scoring(qpairs [][2]byte, busstop *Busstop, c chan Score) {
	intersection := intersection(qpairs, busstop.pairs)
	score := (2 * float64(intersection)) / (float64(len(qpairs)) + float64(len(busstop.pairs)))
	s := Score{busstop.name, score}
	c <- s
}

type Score struct {
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func autocomplete(ctx *gin.Context, app *AppCtx) {
	query := ctx.Query("q")
	if len(query) == 0 {
		ctx.Status(http.StatusTeapot)
		return
	}
	query = normalizeBusstop(query)
	qpairs := pairedBusstops(query)

	ch := make(chan Score)
	var wg sync.WaitGroup
	for i := 0; i < app.n; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			scoring(qpairs, &app.busstops[index], ch)
		}(i)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	results := make([]Score, 0)
	for val := range ch {
		results = append(results, val)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	t := gin.H{"scores": results[0:10]}
	ctx.JSON(http.StatusOK, t)
}

func setupRoutes(app *AppCtx) *gin.Engine {
	r := gin.Default()
	r.GET("/autocomplete", func(ctx *gin.Context) {
		autocomplete(ctx, app)
	})
	return r
}

func main() {
	file, err := os.Open("data/Busstops.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	busstops := make([]string, 1000)

	scanner := bufio.NewScanner(file)
	i := 0
	for ; scanner.Scan(); i++ {
		if i == cap(busstops) {
			nb := make([]string, len(busstops)*2, cap(busstops)*2)
			copy(nb, busstops)
			busstops = nb
		}
		busstops[i] = scanner.Text()
	}

	log.Printf("Busstops loaded %d entries \n", i)
	log.Printf("%s \n", busstops[0])
	pairs := prepareBusstops(busstops, i)

	log.Printf("%s %d", pairs[0].name, len(pairs[0].pairs))
	app := AppCtx{pairs, i}

	r := setupRoutes(&app)
	r.Run(":8080")
}
