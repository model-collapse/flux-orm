package florm

import (
	"time"
)

func From(bucket string, ss *FluxSession) FluxStream {
	return &FluxSinglePipe{
		op:      OpFrom,
		session: ss,
		params:  []string{fkvSq("bucket", bucket)},
	}
}

func FromC(bucket string, host string, org string, token string, ss *FluxSession) FluxStream {
	return &FluxSinglePipe{
		op:      OpFrom,
		session: ss,
		params:  []string{fkvSq("bucket", bucket), fkvSq("host", host), fkvSq("org", org), fkvSq("token", token)},
	}
}

func Buckets(ss *FluxSession) FluxStream {
	return &FluxSinglePipe{
		op:      OpBuckets,
		session: ss,
	}
}

// Timewindow
func Range(in FluxStream, start time.Time, end time.Time) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpRange,
		params:  []string{fkvS("start", start.Format(time.RFC3339)), fkvS("end", end.Format(time.RFC3339))},
	}
}

func RangeR(in FluxStream, start string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpRange,
		params:  []string{fkvS("start", start)},
	}
}

func Static(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpRange,
		params:  []string{fkvS("start", "0"), fkvS("stop", "1")},
	}
}

func Window(in FluxStream, every string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpWindow,
		params:  []string{fkvS("every", every)},
	}
}

func WindowF(in FluxStream, every string, period string, offset string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpWindow,
		params:  []string{fkvS("every", every), fkvS("period", period), fkvS("offset", offset)},
	}
}

func WindowA(in FluxStream, every string, period string, offset string, startColumn string, stopColumn string, createEmpty bool) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpWindow,
		params: []string{fkvS("every", every), fkvS("period", period), fkvS("offset", offset),
			fkvSq("startColumn", startColumn), fkvSq("stopColumn", stopColumn), fkvB("createEmpty", createEmpty)},
	}
}

// Mapping
func Filter(in FluxStream, fn string, onEmpty string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpFilter,
		params:  []string{fkvS("fn", fn), fkvSq("onEmpty", onEmpty)},
	}
}

func Fill(in FluxStream, col string, val string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpFill,
		params:  []string{fkvSq("col", col), fkvSq("val", val)},
	}
}

func FillWithPrev(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpFill,
		params:  []string{fkvSq("col", col), "usePrevious: true"},
	}
}

func HourSelection(in FluxStream, start string, stop string, timeCol string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpHourSelection,
		params:  []string{fkvS("start", start), fkvS("stop", stop), fkvSq("timeColumn", timeCol)},
	}
}

func Limit(in FluxStream, n int, offset int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpLimit,
		params:  []string{fkvI("n", n), fkvI("offset", offset)},
	}
}

func Tail(in FluxStream, n int, offset int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTail,
		params:  []string{fkvI("n", n), fkvI("offset", offset)},
	}
}

func Map(in FluxStream, fn string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMap,
		params:  []string{fkvS("fn", fn)},
	}
}

func TimeShift(in FluxStream, duration string, cols []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTimeShift,
		params:  []string{fkvS("duration", duration), fkvSJ("columns", cols)},
	}
}

func TruncateTimeColumn(in FluxStream, unit string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTruncateTimeColumn,
		params:  []string{fkvS("unit", unit)},
	}
}

// Edit
func Set(in FluxStream, key string, val string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSetString,
		params:  []string{fkvSq("key", key), fkvSq("value", val)},
	}
}

// Aggregation
func Mean(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMean,
		params:  nil,
	}
}

func MeanCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMean,
		params:  []string{fkvSq("column", col)},
	}
}

func Max(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMean,
		params:  nil,
	}
}

func MaxCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMean,
		params:  []string{fkvSq("column", col)},
	}
}

func Min(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMin,
		params:  nil,
	}
}

func MinCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMax,
		params:  []string{fkvSq("column", col)},
	}
}

func Sum(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSum,
		params:  nil,
	}
}

func SumCol(in FluxStream, col string) FluxStream {
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSum,
		params:  []string{fkvSq("column", col)},
	}
}

func Median(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMedian,
		params:  nil,
	}
}

func MedianCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMedian,
		params:  []string{fkvSq("column", col)},
	}
}

func Mode(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMode,
		params:  nil,
	}
}

func ModeCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMode,
		params:  []string{fkvSq("column", col)},
	}
}

func Quantile(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpQuantile,
		params:  nil,
	}
}

func QuantileCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpQuantile,
		params:  []string{fkvSq("column", col)},
	}
}

func Skew(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSkew,
		params:  nil,
	}
}

func SkewCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSkew,
		params:  []string{fkvSq("column", col)},
	}
}

func Spread(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSpread,
		params:  nil,
	}
}

func SpreadCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSpread,
		params:  []string{fkvSq("column", col)},
	}
}

func Stddev(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpStddev,
		params:  nil,
	}
}

func StddevCol(in FluxStream, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpStddev,
		params:  []string{fkvSq("column", column)},
	}
}

func TimeWeightedAvg(in FluxStream, unit string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTimeWeightedAvg,
		params:  []string{fkvSq("unit", unit)},
	}
}

func Integral(in FluxStream, unit string, col, timeCol string, interpolate string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTimeWeightedAvg,
		params:  []string{fkvSq("unit", unit)},
	}
}

func Reduce(in FluxStream, fn string, identity string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpReduce,
		params:  []string{fkvS("fn", fn), fkvS("identity", identity)},
	}
}

func HistogramQuantile(in FluxStream, quantile float64, countColumn string, upperBoundColumn string, valueColumn string, minValue float64) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpHistogramQuantile,
		params: []string{
			fkvF("quantile", quantile),
			fkvSq("countColumn", countColumn),
			fkvSq("upperBoundColumn", upperBoundColumn),
			fkvSq("valueColumn", valueColumn),
			fkvF("minValue", minValue),
		},
	}
}

func Count(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpCount,
		params:  nil,
	}
}

func CountCol(in FluxStream, col string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpCount,
		params:  nil,
	}
}

// Cross Tables
func Join(in FluxStream, with FluxStream, on []string, leftName, rightName string, method string) FluxStream {
	in.AddRef()
	return &FluxMultiplePipe{
		ins:     []FluxStream{in, with},
		session: in.Session(),
		op:      OpJoin,
		params:  []string{fkvSJ("on", on), fkvSq("leftName", leftName), fkvSq("rightName", rightName), fkvSq("method", method)},
	}
}

func Cov(in FluxStream, with FluxStream, on []string, pearsonR bool) FluxStream {
	in.AddRef()
	return &FluxMultiplePipe{
		ins:     []FluxStream{in, with},
		session: in.Session(),
		op:      OpCov,
		params:  []string{fkvSJ("on", on), fkvB("pearsonR", pearsonR)},
	}
}

func PearsonR(in FluxStream, with FluxStream, on []string) FluxStream {
	in.AddRef()
	return &FluxMultiplePipe{
		ins:     []FluxStream{in, with},
		session: in.Session(),
		op:      OpUnion,
		params:  []string{fkvSJ("on", on)},
	}
}

func Union(in FluxStream, with FluxStream) FluxStream {
	in.AddRef()
	return &FluxMultiplePipe{
		ins:     []FluxStream{in, with},
		session: in.Session(),
		op:      OpUnion,
		params:  nil,
	}
}

// Tranforms
func Group(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpGroup,
		params:  nil,
	}
}

func GroupBy(in FluxStream, by []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpGroup,
		params:  []string{fkvSJ("by", by)},
	}
}

func GroupExcept(in FluxStream, exc []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpGroup,
		params:  []string{fkvSJ("exec", exc)},
	}
}

func Keys(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKeys,
		params:  []string{},
	}
}

func KeysCol(in FluxStream, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKeys,
		params:  []string{fkvSq("column", column)},
	}
}

func KeyValues(in FluxStream, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKeys,
		params:  []string{fkvSJ("keyColumns", columns)},
	}
}

func DistinctCol(in FluxStream, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDistinct,
		params:  []string{fkvSq("column", column)},
	}
}

func Distinct(in FluxStream) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDistinct,
		params:  []string{},
	}
}

func Pivot(in FluxStream) FluxStream {
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpPivot,
		params:  []string{fkvSJ("rowKey", []string{"_time"}), fkvSJ("columnKey", []string{"_field"}), fkvSq("valueColumn", "_value")},
	}
}

func PivotC(in FluxStream, rowKey []string, colKey []string, valCol string) FluxStream {
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpPivot,
		params:  []string{fkvSJ("rowKey", rowKey), fkvSJ("columnKey", colKey), fkvSq("valueColumn", valCol)},
	}
}

// Schema
func Columns(in FluxStream, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpColumns,
		params:  []string{fkvSq("column", column)},
	}
}

func Keep(in FluxStream, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKeep,
		params:  []string{fkvSJ("columns", columns)},
	}
}

func Drop(in FluxStream, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDrop,
		params:  []string{fkvSJ("columns", columns)},
	}
}

func Duplicate(in FluxStream, column string, as string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDrop,
		params:  []string{fkvSq("column", column), fkvSq("as", as)},
	}
}

func Rename(in FluxStream, column string, as string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDrop,
		params:  []string{fkvSq("column", column), fkvSq("as", as)},
	}
}

// Series Op
func ChandeMomentumOscillator(in FluxStream, n int, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpChandeMomentumOscillator,
		params:  []string{fkvI("n", n), fkvSJ("columns", columns)},
	}
}

func Covariance(in FluxStream, columns []string, pearsonR bool, dstCol string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpCovariance,
		params:  []string{fkvSJ("columns", columns), fkvB("pearsonR", pearsonR), fkvSq("dstCol", dstCol)},
	}
}

func CumulativeSum(in FluxStream, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpCumulativeSum,
		params:  []string{fkvSJ("columns", columns)},
	}
}

func Derivative(in FluxStream, unit string, nonNegative bool, columns []string, timeCol string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDerivative,
		params:  []string{fkvS("unit", unit), fkvB("nonNegative", nonNegative), fkvSJ("columns", columns), fkvSq("timeCol", timeCol)},
	}
}

func Difference(in FluxStream, nonNegative bool, columns []string, keepFirst bool) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDifference,
		params:  []string{fkvB("nonNegative", nonNegative), fkvSJ("columns", columns), fkvB("keepFirst", keepFirst)},
	}
}

func DoubleEMA(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpDoubleEMA,
		params:  []string{fkvI("n", n)},
	}
}

func Elapsed(in FluxStream, unit string, timeCol string, dstCol string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpElapsed,
		params:  []string{fkvS("unit", unit), fkvSq("timeCol", timeCol), fkvSq("dstCol", dstCol)},
	}
}

func ExponentialMovingAverage(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpExponentialMovingAverage,
		params:  []string{fkvI("n", n)},
	}
}

func Histogram(in FluxStream, column string, upperBoundDstCol string, countCol string, bins []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpHistogram,
		params:  []string{fkvSq("column", column), fkvSq("upperBoundDstCol", upperBoundDstCol), fkvSq("countCol", countCol), fkvSJ("bins", bins)},
	}
}

func HoltWinters(in FluxStream, n int, seasonality int, interval string, withFit bool, timeColumn string, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpHoltWinters,
		params:  []string{fkvI("n", n), fkvI("seasonality", seasonality), fkvS("interval", interval), fkvB("withFit", withFit), fkvSq("timeColumn", timeColumn), fkvSq("column", column)},
	}
}

func Increase(in FluxStream, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpIncrease,
		params:  []string{fkvSJ("columns", columns)},
	}
}

func KaufmansAMA(in FluxStream, n int, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKaufmansAMA,
		params:  []string{fkvI("n", n), fkvS("column", column)},
	}
}

func KaufmansER(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpKaufmansER,
		params:  []string{fkvI("n", n)},
	}
}

func MovingAverage(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpMovingAverage,
		params:  []string{fkvI("n", n)},
	}
}

func RelativeStrengthIndex(in FluxStream, n int, columns []string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpRelativeStrengthIndex,
		params:  []string{fkvI("n", n), fkvSJ("columns", columns)},
	}
}

func StateCount(in FluxStream, fn string, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpStateCount,
		params:  []string{fkvS("fn", fn), fkvS("column", column)},
	}
}

func StateDuration(in FluxStream, fn string, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpStateDuration,
		params:  []string{fkvS("fn", fn), fkvS("column", column)},
	}
}

func TimedMovingAverage(in FluxStream, every string, period string, column string) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTimedMovingAverage,
		params:  []string{fkvS("every", every), fkvS("period", period), fkvSq("column", column)},
	}
}

func TripleEMA(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTripleEMA,
		params:  []string{fkvI("n", n)},
	}
}

func TripleExponentialDerivative(in FluxStream, n int) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpTripleExponentialDerivative,
		params:  []string{fkvI("n", n)},
	}
}

// Sort
func Sort(in FluxStream, columns []string, desc bool) FluxStream {
	in.AddRef()
	return &FluxSinglePipe{
		in:      in,
		session: in.Session(),
		op:      OpSort,
		params:  []string{fkvSJ("columns", columns), fkvB("desc", desc)},
	}
}
