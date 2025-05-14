package protocol

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/umahmood/haversine"
)

func randAccuracy() float64 {
	return 10 * rand.Float64()
}
func randPos(from, to [2]float64) [2]float64 {
	x := rand.Float64()
	return [2]float64{from[0] + (to[0]-from[0])*x, from[1] + (to[1]-from[1])*x}
}
func randRange(from, to int) int {
	return from + rand.Intn(to-from+1)
}
func geoDistance(from, to [2]float64) float64 {
	dist, _ := haversine.Distance(haversine.Coord{Lat: from[1], Lon: from[0]}, haversine.Coord{Lat: to[1], Lon: to[0]})
	return dist * 1000
}
func genTrackAlgorithm(loc []location, dis int64) []string {
	ret := make([]string, 0)
	cur := rand.Intn(len(loc))
	pos := loc[cur]
	l := int64(0)
	t := time.Now().Add(-30 * time.Minute).UnixMilli()
	last := int32(-1)
	ret = append(ret, fmt.Sprintf("%s-%d-%.1f", strings.ReplaceAll(pos.Location, ",", "-"), t, randAccuracy()))
	for l < dis {
		edge := pos.Edge
		x := edge[rand.Intn(len(edge))]
		for x == last {
			x = edge[rand.Intn(len(edge))]
		}
		last = pos.ID
		nxt := loc[x]
		var from [2]float64
		var to [2]float64
		{
			str := strings.Split(pos.Location, ",")
			from[0], _ = strconv.ParseFloat(str[0], 64)
			from[1], _ = strconv.ParseFloat(str[1], 64)
		}
		{
			str := strings.Split(nxt.Location, ",")
			to[0], _ = strconv.ParseFloat(str[0], 64)
			to[1], _ = strconv.ParseFloat(str[1], 64)
		}
		l += int64(geoDistance(from, to))
		randpos := from
		for range 10 {
			nrandpos := randPos(randpos, to)
			randpos = nrandpos
			t += int64(geoDistance(randpos, nrandpos)) / int64(randRange(1, 5)) * 1000
			ret = append(ret, fmt.Sprintf("%.7f-%.7f-%d-%.1f", randpos[0], randpos[1], t, randAccuracy()))
		}
		t += int64(geoDistance(randpos, to)) / int64(randRange(1, 5)) * 1000
		ret = append(ret, fmt.Sprintf("%.7f-%.7f-%d-%.1f", to[0], to[1], t, randAccuracy()))
		pos = nxt
		//pos = nxt
	}
	t += int64(randRange(1, 5)) * 1000
	ret = append(ret, fmt.Sprintf("%s-%d-%.1f", strings.ReplaceAll(pos.Location, ",", "-"), t, randAccuracy()))
	return ret
}

// resource/map.json does not exist or corrupted will lead to crash!!!
// returning a specific json
func readLocation() []location {
	file, err := os.Open("resource/map.json")
	if err != nil {
		panic(err)
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		panic(err)
	}
	var arr []location
	if err := json.Unmarshal(bytes, &arr); err != nil {
		panic(err)
	}
	return arr
}

func genTrack(distance int64) string {
	ret, _ := json.Marshal(genTrackAlgorithm(readLocation(), distance))
	return string(ret)
}
