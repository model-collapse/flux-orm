package florm

import (
	"time"
)

// Timewindow
func (s *FluxSinglePipe) Static() FluxStream {
	return Static(s)
}

func (s *FluxSinglePipe) Range(start time.Time, stop time.Time) FluxStream {
	return Range(s, start, stop)
}

func (s *FluxSinglePipe) RangeR(start string) FluxStream {
	return RangeR(s, start)
}

func (s *FluxSinglePipe) Window(every string) FluxStream {
	return Window(s, every)
}

func (s *FluxSinglePipe) WindowF(every string, period string, offset string) FluxStream {
	return WindowF(s, every, period, offset)
}

func (s *FluxSinglePipe) WindowA(every string, period string, offset string, startColumn string, stopColumn string, createEmpty bool) FluxStream {
	return WindowA(s, every, period, offset, startColumn, stopColumn, createEmpty)
}

// Mapping
func (s *FluxSinglePipe) Filter(fn string, onEmpty string) FluxStream {
	return Filter(s, fn, onEmpty)
}

func (s *FluxSinglePipe) Fill(col string, val string) FluxStream {
	return Fill(s, col, val)
}

func (s *FluxSinglePipe) FillWithPrev(col string) FluxStream {
	return FillWithPrev(s, col)
}

func (s *FluxSinglePipe) HourSelection(start string, stop string, timeCol string) FluxStream {
	return HourSelection(s, start, stop, timeCol)
}

func (s *FluxSinglePipe) Limit(n int, offset int) FluxStream {
	return Limit(s, n, offset)
}

func (s *FluxSinglePipe) Tail(n int, offset int) FluxStream {
	return Tail(s, n, offset)
}

func (s *FluxSinglePipe) Map(fn string) FluxStream {
	return Map(s, fn)
}

func (s *FluxSinglePipe) TimeShift(duration string, cols []string) FluxStream {
	return TimeShift(s, duration, cols)
}

func (s *FluxSinglePipe) TruncateTimeColumn(unit string) FluxStream {
	return TruncateTimeColumn(s, unit)
}

// Edit
func (s *FluxSinglePipe) Set(key string, val string) FluxStream {
	return Set(s, key, val)
}

// Aggregation
func (s *FluxSinglePipe) Mean() FluxStream {
	return Mean(s)
}

func (s *FluxSinglePipe) MeanCol(col string) FluxStream {
	return MeanCol(s, col)
}

func (s *FluxSinglePipe) Max() FluxStream {
	return Max(s)
}

func (s *FluxSinglePipe) MaxCol(col string) FluxStream {
	return MaxCol(s, col)
}

func (s *FluxSinglePipe) Min() FluxStream {
	return Min(s)
}

func (s *FluxSinglePipe) MinCol(col string) FluxStream {
	return MinCol(s, col)
}

func (s *FluxSinglePipe) Sum() FluxStream {
	return Sum(s)
}

func (s *FluxSinglePipe) SumCol(col string) FluxStream {
	return SumCol(s, col)
}

func (s *FluxSinglePipe) Median() FluxStream {
	return Median(s)
}

func (s *FluxSinglePipe) MedianCol(col string) FluxStream {
	return MedianCol(s, col)
}

func (s *FluxSinglePipe) Mode() FluxStream {
	return Mode(s)
}

func (s *FluxSinglePipe) ModeCol(col string) FluxStream {
	return ModeCol(s, col)
}

func (s *FluxSinglePipe) Quantile() FluxStream {
	return Quantile(s)
}

func (s *FluxSinglePipe) QuantileCol(col string) FluxStream {
	return QuantileCol(s, col)
}

func (s *FluxSinglePipe) Skew() FluxStream {
	return Skew(s)
}

func (s *FluxSinglePipe) SkewCol(col string) FluxStream {
	return SkewCol(s, col)
}

func (s *FluxSinglePipe) Spread() FluxStream {
	return Spread(s)
}

func (s *FluxSinglePipe) SpreadCol(col string) FluxStream {
	return SpreadCol(s, col)
}

func (s *FluxSinglePipe) Stddev() FluxStream {
	return Stddev(s)
}

func (s *FluxSinglePipe) StddevCol(column string) FluxStream {
	return StddevCol(s, column)
}

func (s *FluxSinglePipe) TimeWeightedAvg(unit string) FluxStream {
	return TimeWeightedAvg(s, unit)
}

func (s *FluxSinglePipe) Integral(unit string, col, timeCol string, interpolate string) FluxStream {
	return Integral(s, unit, col, timeCol, interpolate)
}

func (s *FluxSinglePipe) Reduce(fn string, identity string) FluxStream {
	return Reduce(s, fn, identity)
}

func (s *FluxSinglePipe) HistogramQuantile(quantile float64, countColumn string, upperBoundColumn string, valueColumn string, minValue float64) FluxStream {
	return HistogramQuantile(s, quantile, countColumn, upperBoundColumn, valueColumn, minValue)
}

func (s *FluxSinglePipe) Count() FluxStream {
	return Count(s)
}

func (s *FluxSinglePipe) CountCol(col string) FluxStream {
	return CountCol(s, col)
}

// Cross Tables
func (s *FluxSinglePipe) Join(with FluxStream, on []string, leftName, rightName string, method string) FluxStream {
	return Join(s, with, on, leftName, rightName, method)
}

func (s *FluxSinglePipe) Cov(with FluxStream, on []string, pearsonR bool) FluxStream {
	return Cov(s, with, on, pearsonR)
}

func (s *FluxSinglePipe) PearsonR(with FluxStream, on []string) FluxStream {
	return PearsonR(s, with, on)
}

func (s *FluxSinglePipe) Union(with FluxStream) FluxStream {
	return Union(s, with)
}

// Tranforms
func (s *FluxSinglePipe) Group() FluxStream {
	return Group(s)
}

func (s *FluxSinglePipe) GroupBy(by []string) FluxStream {
	return GroupBy(s, by)
}

func (s *FluxSinglePipe) GroupExcept(exc []string) FluxStream {
	return GroupExcept(s, exc)
}

func (s *FluxSinglePipe) Keys() FluxStream {
	return Keys(s)
}

func (s *FluxSinglePipe) KeysCol(column string) FluxStream {
	return KeysCol(s, column)
}

func (s *FluxSinglePipe) KeyValues(columns []string) FluxStream {
	return KeyValues(s, columns)
}

func (s *FluxSinglePipe) Distinct() FluxStream {
	return Distinct(s)
}

func (s *FluxSinglePipe) DistinctCol(col string) FluxStream {
	return DistinctCol(s, col)
}

func (s *FluxSinglePipe) Pivot() FluxStream {
	return Pivot(s)
}

func (s *FluxSinglePipe) PivotC(rowKey []string, colKey []string, valCol string) FluxStream {
	return PivotC(s, rowKey, colKey, valCol)
}

// Schema
func (s *FluxSinglePipe) Columns(column string) FluxStream {
	return Columns(s, column)
}

func (s *FluxSinglePipe) Keep(columns []string) FluxStream {
	return Keep(s, columns)
}

func (s *FluxSinglePipe) Drop(columns []string) FluxStream {
	return Drop(s, columns)
}

func (s *FluxSinglePipe) Duplicate(column string, as string) FluxStream {
	return Duplicate(s, column, as)
}

func (s *FluxSinglePipe) Rename(column string, as string) FluxStream {
	return Rename(s, column, as)
}

// Series Op
func (s *FluxSinglePipe) ChandeMomentumOscillator(n int, columns []string) FluxStream {
	return ChandeMomentumOscillator(s, n, columns)
}

func (s *FluxSinglePipe) Covariance(columns []string, pearsonR bool, dstCol string) FluxStream {
	return Covariance(s, columns, pearsonR, dstCol)
}

func (s *FluxSinglePipe) CumulativeSum(columns []string) FluxStream {
	return CumulativeSum(s, columns)
}

func (s *FluxSinglePipe) Derivative(unit string, nonNegative bool, columns []string, timeCol string) FluxStream {
	return Derivative(s, unit, nonNegative, columns, timeCol)
}

func (s *FluxSinglePipe) Difference(nonNegative bool, columns []string, keepFirst bool) FluxStream {
	return Difference(s, nonNegative, columns, keepFirst)
}

func (s *FluxSinglePipe) DoubleEMA(n int) FluxStream {
	return DoubleEMA(s, n)
}

func (s *FluxSinglePipe) Elapsed(unit string, timeCol string, dstCol string) FluxStream {
	return Elapsed(s, unit, timeCol, dstCol)
}

func (s *FluxSinglePipe) ExponentialMovingAverage(n int) FluxStream {
	return ExponentialMovingAverage(s, n)
}

func (s *FluxSinglePipe) Histogram(column string, upperBoundDstCol string, countCol string, bins []string) FluxStream {
	return Histogram(s, column, upperBoundDstCol, countCol, bins)
}

func (s *FluxSinglePipe) HoltWinters(n int, seasonality int, interval string, withFit bool, timeColumn string, column string) FluxStream {
	return HoltWinters(s, n, seasonality, interval, withFit, timeColumn, column)
}

func (s *FluxSinglePipe) Increase(columns []string) FluxStream {
	return Increase(s, columns)
}

func (s *FluxSinglePipe) KaufmansAMA(n int, column string) FluxStream {
	return KaufmansAMA(s, n, column)
}

func (s *FluxSinglePipe) KaufmansER(n int) FluxStream {
	return KaufmansER(s, n)
}

func (s *FluxSinglePipe) MovingAverage(n int) FluxStream {
	return MovingAverage(s, n)
}

func (s *FluxSinglePipe) RelativeStrengthIndex(n int, columns []string) FluxStream {
	return RelativeStrengthIndex(s, n, columns)
}

func (s *FluxSinglePipe) StateCount(fn string, column string) FluxStream {
	return StateCount(s, fn, column)
}

func (s *FluxSinglePipe) StateDuration(fn string, column string) FluxStream {
	return StateDuration(s, fn, column)
}

func (s *FluxSinglePipe) TimedMovingAverage(every string, period string, column string) FluxStream {
	return TimedMovingAverage(s, every, period, column)
}

func (s *FluxSinglePipe) TripleEMA(n int) FluxStream {
	return TripleEMA(s, n)
}

func (s *FluxSinglePipe) TripleExponentialDerivative(n int) FluxStream {
	return TripleExponentialDerivative(s, n)
}

// Sort
func (s *FluxSinglePipe) Sort(columns []string, desc bool) FluxStream {
	return Sort(s, columns, desc)
}
