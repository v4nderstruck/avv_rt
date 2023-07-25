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

	hafasclient "github.com/zensayyy/avv_rt/hafas_client"
)

type Busstop struct {
	Name  string `json:"name"`
	pairs [][2]byte
}

type AppCtx struct {
	busstops map[string]Busstop
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

type Score struct {
	Id    string  `json:"id"`
	Name  string  `json:"name"`
	Score float64 `json:"score"`
}

func scoring(qpairs [][2]byte, busstop *Busstop, id string, c chan Score) {
	intersection := intersection(qpairs, busstop.pairs)
	score := (2 * float64(intersection)) / (float64(len(qpairs)) + float64(len(busstop.pairs)))
	s := Score{id, busstop.Name, score}
	c <- s
}

func autocomplete(ctx *gin.Context, app *AppCtx) {
	query := ctx.Query("q")
	if len(query) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"scores": []Score{}})

		return
	}
	query = normalizeBusstop(query)
	qpairs := pairedBusstops(query)

	log.Printf("Query %s", qpairs)
	ch := make(chan Score)
	var wg sync.WaitGroup
	for stopId, busstop := range app.busstops {
		id := stopId
		bus := busstop
		wg.Add(1)
		go func() {
			defer wg.Done()
			scoring(qpairs, &bus, id, ch)
		}()
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

func fetchDepartures(ctx *gin.Context, app *AppCtx, stopId string) {
	if stopId == "" {
		ctx.Status(http.StatusOK)
		return
	}

	if stop, ok := app.busstops[stopId]; ok {
		client := hafasclient.NewClient()
		departures, err := client.GetDepartures(stopId)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"stop": stop, "departures": departures})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"error": "Stop not found"})
}

func setupRoutes(app *AppCtx) *gin.Engine {
	r := gin.Default()
	r.GET("/autocomplete", func(ctx *gin.Context) {
		autocomplete(ctx, app)
	})

	r.GET("/departure/:stopId", func(ctx *gin.Context) {
		fetchDepartures(ctx, app, ctx.Param("stopId"))
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

	busstops := make(map[string]Busstop)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		scan := scanner.Text()
		split := strings.Split(scan, ";")
		id, name := split[0], split[1]

		busstops[id] = Busstop{
			Name:  name,
			pairs: pairedBusstops(normalizeBusstop(name)),
		}

	}

	log.Printf("Busstops loaded %d entries \n", len(busstops))
	app := AppCtx{busstops}
	r := setupRoutes(&app)
	r.Run(":8080")
}
