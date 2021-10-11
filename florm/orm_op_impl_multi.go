package florm

import (
	"time"
)

// Timewindow
func (s *FluxMultiplePipe) Static() FluxStream {
	return Static(s)
}

func (s *FluxMultiplePipe) Range(start time.Time, end time.Time) FluxStream {
	return Range(s, start, end)
}

func (s *FluxMultiplePipe) RangeR(start string) FluxStream {
	return RangeR(s, start)
}

func (s *FluxMultiplePipe) Window(every string) FluxStream {
	return Window(s, every)
}

func (s *FluxMultiplePipe) WindowF(every string, period string, offset string) FluxStream {
	return WindowF(s, every, period, offset)
}

func (s *FluxMultiplePipe) WindowA(every string, period string, offset string, startColumn string, stopColumn string, createEmpty bool) FluxStream {
	return WindowA(s, every, period, offset, startColumn, stopColumn, createEmpty)
}

// Mapping
func (s *FluxMultiplePipe) Filter(fn string, onEmpty string) FluxStream {
	return Filter(s, fn, onEmpty)
}

func (s *FluxMultiplePipe) Fill(col string, val string) FluxStream {
	return Fill(s, col, val)
}

func (s *FluxMultiplePipe) FillWithPrev(col string) FluxStream {
	return FillWithPrev(s, col)
}

func (s *FluxMultiplePipe) HourSelection(start string, stop string, timeCol string) FluxStream {
	return HourSelection(s, start, stop, timeCol)
}

func (s *FluxMultiplePipe) Limit(n int, offset int) FluxStream {
	return Limit(s, n, offset)
}

func (s *FluxMultiplePipe) Tail(n int, offset int) FluxStream {
	return Tail(s, n, offset)
}

func (s *FluxMultiplePipe) Map(fn string) FluxStream {
	return Map(s, fn)
}

func (s *FluxMultiplePipe) TimeShift(duration string, cols []string) FluxStream {
	return TimeShift(s, duration, cols)
}

func (s *FluxMultiplePipe) TruncateTimeColumn(unit string) FluxStream {
	return TruncateTimeColumn(s, unit)
}

// Edit
func (s *FluxMultiplePipe) Set(key string, val string) FluxStream {
	return Set(s, key, val)
}

// Aggregation
func (s *FluxMultiplePipe) Mean() FluxStream {
	return Mean(s)
}

func (s *FluxMultiplePipe) MeanCol(col string) FluxStream {
	return MeanCol(s, col)
}

func (s *FluxMultiplePipe) Max() FluxStream {
	return Max(s)
}

func (s *FluxMultiplePipe) MaxCol(col string) FluxStream {
	return MaxCol(s, col)
}

func (s *FluxMultiplePipe) Min() FluxStream {
	return Min(s)
}

func (s *FluxMultiplePipe) MinCol(col string) FluxStream {
	return MinCol(s, col)
}

func (s *FluxMultiplePipe) Sum() FluxStream {
	return Sum(s)
}

func (s *FluxMultiplePipe) SumCol(col string) FluxStream {
	return SumCol(s, col)
}

func (s *FluxMultiplePipe) Median() FluxStream {
	return Median(s)
}

func (s *FluxMultiplePipe) MedianCol(col string) FluxStream {
	return MedianCol(s, col)
}

func (s *FluxMultiplePipe) Mode() FluxStream {
	return Mode(s)
}

func (s *FluxMultiplePipe) ModeCol(col string) FluxStream {
	return ModeCol(s, col)
}

func (s *FluxMultiplePipe) Quantile() FluxStream {
	return Quantile(s)
}

func (s *FluxMultiplePipe) QuantileCol(col string) FluxStream {
	return QuantileCol(s, col)
}

func (s *FluxMultiplePipe) Skew() FluxStream {
	return Skew(s)
}

func (s *FluxMultiplePipe) SkewCol(col string) FluxStream {
	return SkewCol(s, col)
}

func (s *FluxMultiplePipe) Spread() FluxStream {
	return Spread(s)
}

func (s *FluxMultiplePipe) SpreadCol(col string) FluxStream {
	return SpreadCol(s, col)
}

func (s *FluxMultiplePipe) Stddev() FluxStream {
	return Stddev(s)
}

func (s *FluxMultiplePipe) StddevCol(column string) FluxStream {
	return StddevCol(s, column)
}

func (s *FluxMultiplePipe) TimeWeightedAvg(unit string) FluxStream {
	return TimeWeightedAvg(s, unit)
}

func (s *FluxMultiplePipe) Integral(unit string, col, timeCol string, interpolate string) FluxStream {
	return Integral(s, unit, col, timeCol, interpolate)
}

func (s *FluxMultiplePipe) Reduce(fn string, identity string) FluxStream {
	return Reduce(s, fn, identity)
}

func (s *FluxMultiplePipe) HistogramQuantile(quantile float64, countColumn string, upperBoundColumn string, valueColumn string, minValue float64) FluxStream {
	return HistogramQuantile(s, quantile, countColumn, upperBoundColumn, valueColumn, minValue)
}

func (s *FluxMultiplePipe) Count() FluxStream {
	return Count(s)
}

func (s *FluxMultiplePipe) CountCol(col string) FluxStream {
	return CountCol(s, col)
}

// Cross Tables
func (s *FluxMultiplePipe) Join(with FluxStream, on []string, leftName, rightName string, method string) FluxStream {
	return Join(s, with, on, leftName, rightName, method)
}

func (s *FluxMultiplePipe) Cov(with FluxStream, on []string, pearsonR bool) FluxStream {
	return Cov(s, with, on, pearsonR)
}

func (s *FluxMultiplePipe) PearsonR(with FluxStream, on []string) FluxStream {
	return PearsonR(s, with, on)
}

func (s *FluxMultiplePipe) Union(with FluxStream) FluxStream {
	return Union(s, with)
}

// Tranforms
func (s *FluxMultiplePipe) Group() FluxStream {
	return Group(s)
}

func (s *FluxMultiplePipe) GroupBy(by []string) FluxStream {
	return GroupBy(s, by)
}

func (s *FluxMultiplePipe) GroupExcept(exc []string) FluxStream {
	return GroupExcept(s, exc)
}

func (s *FluxMultiplePipe) Keys(column string) FluxStream {
	return Keys(s, column)
}

func (s *FluxMultiplePipe) KeyValues(columns []string) FluxStream {
	return KeyValues(s, columns)
}

func (s *FluxMultiplePipe) Pivot() FluxStream {
	return Pivot(s)
}

func (s *FluxMultiplePipe) PivotC(rowKey []string, colKey []string, valCol string) FluxStream {
	return PivotC(s, rowKey, colKey, valCol)
}

// Schema
func (s *FluxMultiplePipe) Columns(column string) FluxStream {
	return Columns(s, column)
}

func (s *FluxMultiplePipe) Keep(columns []string) FluxStream {
	return Keep(s, columns)
}

func (s *FluxMultiplePipe) Drop(columns []string) FluxStream {
	return Drop(s, columns)
}

func (s *FluxMultiplePipe) Duplicate(column string, as string) FluxStream {
	return Duplicate(s, column, as)
}

func (s *FluxMultiplePipe) Rename(column string, as string) FluxStream {
	return Rename(s, column, as)
}

// Series Op
func (s *FluxMultiplePipe) ChandeMomentumOscillator(n int, columns []string) FluxStream {
	return ChandeMomentumOscillator(s, n, columns)
}

func (s *FluxMultiplePipe) Covariance(columns []string, pearsonR bool, dstCol string) FluxStream {
	return Covariance(s, columns, pearsonR, dstCol)
}

func (s *FluxMultiplePipe) CumulativeSum(columns []string) FluxStream {
	return CumulativeSum(s, columns)
}

func (s *FluxMultiplePipe) Derivative(unit string, nonNegative bool, columns []string, timeCol string) FluxStream {
	return Derivative(s, unit, nonNegative, columns, timeCol)
}

func (s *FluxMultiplePipe) Difference(nonNegative bool, columns []string, keepFirst bool) FluxStream {
	return Difference(s, nonNegative, columns, keepFirst)
}

func (s *FluxMultiplePipe) DoubleEMA(n int) FluxStream {
	return DoubleEMA(s, n)
}

func (s *FluxMultiplePipe) Elapsed(unit string, timeCol string, dstCol string) FluxStream {
	return Elapsed(s, unit, timeCol, dstCol)
}

func (s *FluxMultiplePipe) ExponentialMovingAverage(n int) FluxStream {
	return ExponentialMovingAverage(s, n)
}

func (s *FluxMultiplePipe) Histogram(column string, upperBoundDstCol string, countCol string, bins []string) FluxStream {
	return Histogram(s, column, upperBoundDstCol, countCol, bins)
}

func (s *FluxMultiplePipe) HoltWinters(n int, seasonality int, interval string, withFit bool, timeColumn string, column string) FluxStream {
	return HoltWinters(s, n, seasonality, interval, withFit, timeColumn, column)
}

func (s *FluxMultiplePipe) Increase(columns []string) FluxStream {
	return Increase(s, columns)
}

func (s *FluxMultiplePipe) KaufmansAMA(n int, column string) FluxStream {
	return KaufmansAMA(s, n, column)
}

func (s *FluxMultiplePipe) KaufmansER(n int) FluxStream {
	return KaufmansER(s, n)
}

func (s *FluxMultiplePipe) MovingAverage(n int) FluxStream {
	return MovingAverage(s, n)
}

func (s *FluxMultiplePipe) RelativeStrengthIndex(n int, columns []string) FluxStream {
	return RelativeStrengthIndex(s, n, columns)
}

func (s *FluxMultiplePipe) StateCount(fn string, column string) FluxStream {
	return StateCount(s, fn, column)
}

func (s *FluxMultiplePipe) StateDuration(fn string, column string) FluxStream {
	return StateDuration(s, fn, column)
}

func (s *FluxMultiplePipe) TimedMovingAverage(every string, period string, column string) FluxStream {
	return TimedMovingAverage(s, every, period, column)
}

func (s *FluxMultiplePipe) TripleEMA(n int) FluxStream {
	return TripleEMA(s, n)
}

func (s *FluxMultiplePipe) TripleExponentialDerivative(n int) FluxStream {
	return TripleExponentialDerivative(s, n)
}

// Sort
func (s *FluxMultiplePipe) Sort(columns []string, desc bool) FluxStream {
	return Sort(s, columns, desc)
}
