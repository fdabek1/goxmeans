/*
 Package goxmeans implements a simple library for the xmeans algorithm.

 See Dan Pelleg and Andrew Moore: X-means: Extending K-means with Efficient Estimation of the Number of Clusters. 
*/
package goxmeans

import (
	"bufio"
	"code.google.com/p/gomatrix/matrix"
	"errors"
	"fmt"
	"goxmeans/matutil"
	"io"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

// Atof64 is shorthand for ParseFloat(s, 64)
func Atof64(s string) (f float64, err error) {
	f64, err := strconv.ParseFloat(s, 64)
	return float64(f64), err
}

// Load loads a tab delimited text file of floats into a slice.
// Assume last column is the target.
// For now, we limit ourselves to two columns
func Load(fname string) (*matrix.DenseMatrix, error) {
	datamatrix := matrix.Zeros(1, 1)
	data := make([]float64, 2048)
	idx := 0

	fp, err := os.Open(fname)
	if err != nil {
		return datamatrix, err
	}
	defer fp.Close()

	r := bufio.NewReader(fp)
	linenum := 1
	eof := false
	for !eof {
		var line string
		line, err := r.ReadString('\n')
		if err == io.EOF {
			err = nil
			eof = true
			break
		} else if err != nil {
			return datamatrix, errors.New(fmt.Sprintf("means: reading linenum %d: %v", linenum, err))
		}

		linenum++
		l1 := strings.TrimRight(line, "\n")
		l := strings.Split(l1, "\t")
		if len(l) < 2 {
			return datamatrix, errors.New(fmt.Sprintf("means: linenum %d has only %d elements", linenum, len(line)))
		}

		// for now assume 2 dimensions only
		f0, err := Atof64(string(l[0]))
		if err != nil {
			return datamatrix, errors.New(fmt.Sprintf("means: cannot convert %s to float64.", l[0]))
		}
		f1, err := Atof64(string(l[1]))
		if err != nil {
			return datamatrix, errors.New(fmt.Sprintf("means: cannot convert %s to float64.", l[linenum][1]))
		}

		if linenum >= len(data) {
			data = append(data, f0, f1)
		} else {
			data[idx] = f0
			idx++
			data[idx] = f1
			idx++
		}
	}
	numcols := 2
	datamatrix = matrix.MakeDenseMatrix(data, len(data)/numcols, numcols)
	return datamatrix, nil
}

// RandCentroids picks random centroids based on the  min and max values in the matrix
// and return a k by cols matrix of the centroids.
func RandCentroids(mat *matrix.DenseMatrix, k int) *matrix.DenseMatrix {
	_, cols := mat.GetSize()
	centroids := matrix.Zeros(k, cols)

	for colnum := 0; colnum < cols; colnum++ {
		r := matutil.ColSlice(mat, colnum)

		minj := float64(0)
		// min value from column
		for _, val := range r {
			minj = math.Min(minj, val)
		}

		// max value from column
		maxj := float64(0)
		for _, val := range r {
			maxj = math.Max(maxj, val)
		}

		// create a slice of random centroids 
		// based on maxj + minJ * random num to stay in range
		// TODO: Better randomization or choose centroids 
		// from datapoints.
		rands := make([]float64, k)
		for i := 0; i < k; i++ {
			randint := float64(rand.Int())
			rf := (maxj - minj) * randint
			for rf > maxj {
				if rf > maxj*3 {
					rf = rf / maxj
				} else {
					rf = rf / 2
				}
			}
			rands[i] = rf
		}
		for h := 0; h < k; h++ {
			centroids.Set(h, colnum, rands[h])
		}
	}
	return centroids
}

/* TODO: An interface for all distances 
   should be in a separate distance package
type Distance interface {
	Distance()
}

type CentroidMaker interface {
	MakeCentroids()
	k int // number of centroids
	dataSet *matrix.DenseMatrix  // set of data points
}*/


// TODO: Create Distance interface so that any distance metric, Euclidean, Jacard, etc. can be passed
// kmeans takes a matrix as input data and attempts to find the best convergence on a set of k centroids.
//func kmeans(data *matrix.DenseMatrix, k int, dist Distance, maker CentroidMaker) (centroids  *matrix.DenseMatrix, clusterAssignment *matrix.DenseMatrix) {
// Get something working with Euclidean and RandCentroids
func kmeans(dataSet *matrix.DenseMatrix, k int) {
	numRows, numCols = dataSet.GetSize()
    //Pseudo Code
	//clusterAssignment - create mat to assign data points to a centroid, also holds SE of each point
	//clusterChanged = true
	//centroids = RandCentroids(dataSet, k)
	/* for ; clusterChanged ; {
	     clusterChanged = false
         for i := 0; i < numRows; {  // assign each data point to a centroid
        	 minDist := float64(0)
             minIndex := -1
             for j := 0; j < k; j++ {  // check distance against each centroid
     	         distJ := matutil.EuclidDist(centroids.getRowVector(j), dataSet.GetRowVector(i))
                 if distJ < minDist {
                     minDist = distJ
                     minIndex = j
	             } 
            	 if clusterAssignment.Get(i, 0) != minIndex {
	                 clusterChanged = true
	             }
                 clusterAssignment.Set(i,0) = minIndex
	             clusterAssignment.Set(i,1) = math.Pow(minDist, 2)
            	 //TODO: Write SetRowVector(row int, value float64[])
	         }
         }
         for c := 0; c < k; k++ {
	         pointsInCluster := all non-zero data points in the current cluster c into a matrix
         	 centroids.SetRowVector(c,  mean(pointsInCluster, axis=0)) #assign centroid to mean 
	     }
	    return centroids, clusterAssignment
     }
	*/
}

func ComputeCentroid(mat *matrix.DenseMatrix) (*matrix.DenseMatrix, error) {
	rows, _ := mat.GetSize()
	vectorSum := matutil.SumCols(mat)
	if rows == 0 {
		return vectorSum, errors.New("No points inputted")
	}
	vectorSum.Scale(1.0 / float64(rows))
	return vectorSum, nil
}

/*func kmeans(data *matrix.DenseMatrix, k int, dist Distance, centroids func(mat *matrix.DenseMatrix, howmany int)) {

}*/

