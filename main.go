package main

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"
	"math/rand"
	"time"
	"math"
)

func main() {
	sizeArray := make([]int, 6)
	sizeArray[0] = 5
	sizeArray[1] = 7
	sizeArray[2] = 10
	sizeArray[3] = 11
	sizeArray[4] = 12
	sizeArray[5] = 14

	collisionAverages := performCollisionAttacks(60, sizeArray)
	preImageAverages := performPreImageAttacks(60, sizeArray)

	fmt.Println("Collision averages:")
	for i := 0; i < 6; i++ {
		fmt.Printf("\t%d: %f\n", sizeArray[i], collisionAverages[i])
	}
	fmt.Println("Pre Image averages:")
	for i := 0; i < 6; i++ {
		fmt.Printf("\t%d: %f\n", sizeArray[i], preImageAverages[i])
	}

	logp, err1 := plot.New()
	logp.Title.Text = "Average Iterations (log2)"
	logp.X.Label.Text = "Hash length (in bits)"
	logp.Y.Label.Text = "# of iterations (log2)"
	err1 = plotutil.AddLinePoints(logp,
		"PreImage", convertToPlotterXYLog(sizeArray, preImageAverages),
		"Collision", convertToPlotterXYLog(sizeArray, collisionAverages))
	if err1 != nil {
		panic(err1)
	}
	logp.Legend.Top = true
	logp.Legend.Left = true

	p, err := plot.New()
	p.Title.Text = "Average Iterations"
	p.X.Label.Text = "Hash length (in bits)"
	p.Y.Label.Text = "# of iterations"
	err = plotutil.AddLinePoints(p,
		"PreImage", convertToPlotterXY(sizeArray, preImageAverages),
		"Collision", convertToPlotterXY(sizeArray, collisionAverages))
	if err != nil {
		panic(err)
	}
	p.Legend.Top = true
	p.Legend.Left = true
	if err := p.Save(4*vg.Inch, 4*vg.Inch, "points.pdf"); err != nil {
		panic(err)
	}
	if err1 := logp.Save(4*vg.Inch, 4*vg.Inch, "pointsLog.pdf"); err != nil {
		panic(err1)
	}
}

func convertToPlotterXY(xs []int, ys []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i := 0; i < len(xs); i++ {
		pts[i].X = float64(xs[i])
		pts[i].Y = ys[i]
	}
	return pts
}

func convertToPlotterXYLog(xs []int, ys []float64) plotter.XYs {
	pts := make(plotter.XYs, len(xs))
	for i := 0; i < len(xs); i++ {
		pts[i].X = float64(xs[i])
		pts[i].Y = math.Log2(ys[i])
	}
	return pts
}

func performCollisionAttacks(n int, sizeArray []int) []float64 {
	averages := make([]float64, len(sizeArray))
	for i := 0; i < len(sizeArray); i++ {
		fmt.Printf("Performing collision attack with size %d\n", sizeArray[i])
		counts := make([]int, n)
		for j := 0; j < n; j++ {
			_, _, _, count := collisionAttack(sizeArray[i])
			counts[j] = count
		}
		averages[i] = getAverage(counts)
	}
	return averages
}

func performPreImageAttacks(n int, sizeArray []int) []float64 {
	averages := make([]float64, len(sizeArray))
	for i := 0; i < len(sizeArray); i++ {
		fmt.Printf("Performing pre image attack with size %d\n", sizeArray[i])
		counts := make([]int, n)
		for j := 0; j < n; j++ {
			_, _, _, count := preImageAttack(sizeArray[i])
			counts[j] = count
		}
		averages[i] = getAverage(counts)
	}
	return averages
}

func getAverage(counts []int) float64 {
	total := 0
	for i := 0; i < len(counts); i++ {
		total += counts[i]
	}
	return float64(total) / float64(len(counts))
}

func shaone(input string, n int) []byte {
	bytes := sha1.Sum([]byte(input))
	return truncate(bytes, n)
}

func truncate(bytes [20]byte, n int) []byte {
	leftOverBits := n % 8
	if leftOverBits == 0 {
		return bytes[:(n / 8)]
	}
	resultBits := bytes[:(n/8)+1]
	//mask bits to truncate to n
	switch leftOverBits {
	case 1:
		//0000 0001
		resultBits[len(resultBits)-1] &= 0x01
		break
	case 2:
		//0000 0011
		resultBits[len(resultBits)-1] &= 0x03
		break
	case 3:
		//0000 0111
		resultBits[len(resultBits)-1] &= 0x07
		break
	case 4:
		//0000 1111
		resultBits[len(resultBits)-1] &= 0x0f
		break
	case 5:
		//0001 1111
		resultBits[len(resultBits)-1] &= 0x1f
		break
	case 6:
		//0011 1111
		resultBits[len(resultBits)-1] &= 0x3f
		break
	case 7:
		//0111 1111
		resultBits[len(resultBits)-1] &= 0x7f
		break
	}
	return resultBits
}

func preImageAttack(n int) (string, string, []byte, int) {
	originalString := generateRandomString()
	originalHash := shaone(originalString, n)
	count := 0
	for true {
		count++
		randomString := generateRandomString()
		hash := shaone(randomString, n)
		if bytes.Compare(hash, originalHash) == 0 {
			// fmt.Printf("Hash (%x) == (%x)\n", hash, originalHash)
			return originalString, randomString, hash, count
		}
	}
	return "", "", nil, 0
}

func collisionAttack(n int) (string, string, []byte, int) {
	usedHashes := make(map[string]string)
	count := 0
	for true {
		count++
		randomString := generateRandomString()
		hash := shaone(randomString, n)
		hashString := fmt.Sprint(hash)
		if val, ok := usedHashes[hashString]; ok {
			return val, randomString, hash, count
		}
		usedHashes[hashString] = randomString
	}
	return "", "", nil, count
}

func generateRandomString() string {
	rand.Seed(time.Now().UTC().UnixNano())
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, 20)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
