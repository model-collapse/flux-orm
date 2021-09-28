package florm

import (
	"time"
)

type FluxOp int

const (
	// Source
	OpFrom FluxOp = iota
	OpBuckets

	// Timewindow
	OpRange
	OpWindow

	// Mapping
	OpFilter
	OpFill
	OpHourSelection
	OpLimit
	OpTail
	OpMap
	OpTimeShift
	OpTruncateTimeColumn

	// Edit
	OpSet

	// Aggregation
	OpMean
	OpMax
	OpMin
	OpSum
	OpMedian
	OpMode
	OpQuantile
	OpSkew
	OpSpread
	OpStddev
	OpTimeWeightedAvg
	OpIntegral
	OpReduce
	OpHistogramQuantile
	OpCount

	// Cross Tables
	OpJoin
	OpCov
	OpPearsonR
	OpUnion

	// Tranforms
	OpGroup
	OpKeys
	OpKeyValues
	OpPivot

	// Schema
	OpColumns
	OpKeep
	OpDrop
	OpDuplicate
	OpRename

	// Series Op
	OpChandeMomentumOscillator
	OpCovariance
	OpCumulativeSum
	OpDerivative
	OpDifference
	OpDoubleEMA
	OpElapsed
	OpExponentialMovingAverage
	OpHistogram
	OpHoltWinters
	OpIncrease
	OpKaufmansAMA
	OpKaufmansER
	OpMovingAverage
	OpRelativeStrengthIndex
	OpStateCount
	OpStateDuration
	OpTimedMovingAverage
	OpTripleEMA
	OpTripleExponentialDerivative

	// Sort
	OpSort

	//Yield
	OpYield
)

var opToStr = map[FluxOp]string{
	// Source
	OpFrom: "from",
	// Timewindow
	OpRange:  "range",
	OpWindow: "window",

	// Mapping
	OpFilter:             "filter",
	OpFill:               "fill",
	OpHourSelection:      "hourSelection",
	OpLimit:              "limit",
	OpTail:               "tail",
	OpMap:                "map",
	OpTimeShift:          "timeShift",
	OpTruncateTimeColumn: "truncateTimeColumn",

	// Edit
	OpSet: "set",

	// Aggregation
	OpMean:              "mean",
	OpMax:               "max",
	OpMin:               "min",
	OpSum:               "sum",
	OpMedian:            "median",
	OpMode:              "mode",
	OpQuantile:          "quantile",
	OpSkew:              "skew",
	OpSpread:            "spread",
	OpStddev:            "stddev",
	OpTimeWeightedAvg:   "timeWeightedAvg",
	OpIntegral:          "integral",
	OpReduce:            "reduce",
	OpHistogramQuantile: "histogramQuantile",
	OpCount:             "count",

	// Cross Tables
	OpJoin:     "join",
	OpCov:      "cov",
	OpPearsonR: "pearsonR",
	OpUnion:    "union",

	// Tranforms
	OpGroup:     "group",
	OpKeys:      "keys",
	OpKeyValues: "keyValues",
	OpPivot:     "pivot",

	// Schema
	OpColumns:   "columns",
	OpKeep:      "keep",
	OpDrop:      "drop",
	OpDuplicate: "duplicate",
	OpRename:    "rename",

	// Series Op
	OpChandeMomentumOscillator:    "chandeMomentumOscillator",
	OpCovariance:                  "covariance",
	OpCumulativeSum:               "cumulativeSum",
	OpDerivative:                  "derivative",
	OpDifference:                  "difference",
	OpDoubleEMA:                   "doubleEMA",
	OpElapsed:                     "elapsed",
	OpExponentialMovingAverage:    "exponentialMovingAverage",
	OpHistogram:                   "histogram",
	OpHoltWinters:                 "holtWinters",
	OpIncrease:                    "increase",
	OpKaufmansAMA:                 "kaufmansAMA",
	OpKaufmansER:                  "kaufmansER",
	OpMovingAverage:               "movingAverage",
	OpRelativeStrengthIndex:       "relativeStrengthIndex",
	OpStateCount:                  "stateCount",
	OpStateDuration:               "stateDuration",
	OpTimedMovingAverage:          "timedMovingAverage",
	OpTripleEMA:                   "tripleEMA",
	OpTripleExponentialDerivative: "tripleExponentialDerivative",

	// Sort
	OpSort: "sort",

	// Yeild
	OpYield: "yield",
}

type FluxOperatable interface {
	// Timewindow
	Range(start time.Time, end time.Time) FluxStream
	RangeR(start string) FluxStream
	Static() FluxStream

	Window(every string) FluxStream
	WindowF(every string, period string, offset string) FluxStream
	WindowA(every string, period string, offset string, startColumn string, stopColumn string, createEmpty bool) FluxStream

	// Mapping
	Filter(fn string, onEmpty string) FluxStream

	Fill(col string, val string) FluxStream
	FillWithPrev(col string) FluxStream

	HourSelection(start string, stop string, timeCol string) FluxStream

	Limit(n int, offset int) FluxStream

	Tail(n int, offset int) FluxStream

	Map(fn string) FluxStream

	TimeShift(duration string, cols []string) FluxStream

	TruncateTimeColumn(unit string) FluxStream

	// Edit
	Set(key string, val string) FluxStream

	// Aggregation
	Mean() FluxStream
	MeanCol(col string) FluxStream

	Max() FluxStream
	MaxCol(col string) FluxStream

	Min() FluxStream
	MinCol(col string) FluxStream

	Sum() FluxStream
	SumCol(col string) FluxStream

	Median() FluxStream
	MedianCol(col string) FluxStream

	Mode() FluxStream
	ModeCol(col string) FluxStream

	Quantile() FluxStream
	QuantileCol(col string) FluxStream

	Skew() FluxStream
	SkewCol(col string) FluxStream

	Spread() FluxStream
	SpreadCol(col string) FluxStream

	Stddev() FluxStream
	StddevCol(column string) FluxStream

	TimeWeightedAvg(uint string) FluxStream

	Integral(unit string, col, timeCol string, interpolate string) FluxStream

	Reduce(fn string, indentity string) FluxStream

	HistogramQuantile(quantile float64, countColumn string, upperBoundColumn string, valueColumn string, minValue float64) FluxStream

	Count() FluxStream
	CountCol(col string) FluxStream

	// Cross Tables
	Join(with FluxStream, on []string, leftName, rightName string, method string) FluxStream
	Cov(with FluxStream, on []string, pearsonR bool) FluxStream
	PearsonR(with FluxStream, on []string) FluxStream
	Union(with FluxStream) FluxStream

	// Tranforms
	Group() FluxStream
	GroupBy(by []string) FluxStream
	GroupExcept(exc []string) FluxStream

	Keys(column string) FluxStream
	KeyValues(columns []string) FluxStream
	Pivot(rowKey []string, colKey []string, valCol string) FluxStream

	// Schema
	Columns(column string) FluxStream
	Keep(columns []string) FluxStream
	Drop(columns []string) FluxStream
	Duplicate(column string, as string) FluxStream
	Rename(column string, as string) FluxStream

	// Series Op
	ChandeMomentumOscillator(n int, columns []string) FluxStream
	Covariance(columns []string, pearsonR bool, dstCol string) FluxStream
	CumulativeSum(columns []string) FluxStream
	Derivative(unit string, nonNegative bool, columns []string, timeCol string) FluxStream
	Difference(nonNegative bool, columns []string, keepFirst bool) FluxStream
	DoubleEMA(n int) FluxStream
	Elapsed(unit string, timeCol string, dstCol string) FluxStream
	ExponentialMovingAverage(n int) FluxStream
	Histogram(column string, upperBoundDstCol string, countCol string, bins []string) FluxStream
	HoltWinters(n int, seasonality int, interval string, withFit bool, timeColumn string, column string) FluxStream
	Increase(columns []string) FluxStream
	KaufmansAMA(n int, column string) FluxStream
	KaufmansER(n int) FluxStream
	MovingAverage(n int) FluxStream
	RelativeStrengthIndex(n int, columns []string) FluxStream
	StateCount(fn string, column string) FluxStream
	StateDuration(fn string, column string) FluxStream
	TimedMovingAverage(every string, period string, column string) FluxStream
	TripleEMA(n int) FluxStream
	TripleExponentialDerivative(n int) FluxStream

	// Sort
	Sort(columns []string, desc bool) FluxStream
}
