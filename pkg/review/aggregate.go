package review

import (
	"math"
	"math/big"
	"sort"
)

type DimensionSummary struct {
	Count int     `json:"count"`
	Mean  float64 `json:"mean"`
}
type VariantSummary struct {
	Variant         string                      `json:"variant"`
	Available       int                         `json:"available"`
	Reviewed        int                         `json:"reviewed"`
	UniqueSubjects  int                         `json:"unique_subjects"`
	AnnotationCount int                         `json:"annotation_count"`
	Dimensions      map[string]DimensionSummary `json:"dimensions"`
}
type PairedComparison struct {
	VariantA        string   `json:"variant_a"`
	VariantB        string   `json:"variant_b"`
	Dimension       string   `json:"dimension"`
	ReviewedPairs   int      `json:"reviewed_pairs"`
	Wins            int      `json:"wins"`
	Ties            int      `json:"ties"`
	Losses          int      `json:"losses"`
	MeanDelta       float64  `json:"mean_delta"`
	MedianDelta     float64  `json:"median_delta"`
	Bootstrap95Low  *float64 `json:"bootstrap_95_low,omitempty"`
	Bootstrap95High *float64 `json:"bootstrap_95_high,omitempty"`
}
type ReviewerDimensionAgreement struct {
	Count                  int     `json:"count"`
	ExactAgreements        int     `json:"exact_agreements"`
	ExactAgreementRate     float64 `json:"exact_agreement_rate"`
	MeanAbsoluteDifference float64 `json:"mean_absolute_difference"`
	MaxAbsoluteDifference  int     `json:"max_absolute_difference"`
}
type ReviewerDisagreement struct {
	ReviewID   string `json:"review_id"`
	SubjectID  string `json:"subject_id"`
	Variant    string `json:"variant"`
	Dimension  string `json:"dimension"`
	ScoreA     int    `json:"score_a"`
	ScoreB     int    `json:"score_b"`
	Difference int    `json:"absolute_difference"`
}
type ReviewerOverlapSummary struct {
	ReviewerA     string                                `json:"reviewer_a"`
	ReviewerB     string                                `json:"reviewer_b"`
	OverlapItems  int                                   `json:"overlap_items"`
	Dimensions    map[string]ReviewerDimensionAgreement `json:"dimensions"`
	Disagreements []ReviewerDisagreement                `json:"disagreements,omitempty"`
}
type Report struct {
	QueueItems       int                      `json:"queue_items"`
	ReviewedItems    int                      `json:"reviewed_items"`
	AnnotationCount  int                      `json:"annotation_count"`
	OverlappingItems int                      `json:"overlapping_items"`
	Variants         []VariantSummary         `json:"variants"`
	Comparisons      []PairedComparison       `json:"comparisons"`
	ReviewerOverlaps []ReviewerOverlapSummary `json:"reviewer_overlaps"`
}

func Aggregate(keys []KeyEntry, annotations []Annotation, dimensions []Dimension) Report {
	keyByID := map[string]KeyEntry{}
	available := map[string]int{}
	for _, key := range keys {
		keyByID[key.ReviewID] = key
		available[key.Variant]++
	}
	type score struct {
		count int
		sum   map[string]*big.Int
	}
	itemScores := map[string]*score{}
	for _, annotation := range annotations {
		current := itemScores[annotation.ReviewID]
		if current == nil {
			current = &score{sum: map[string]*big.Int{}}
			itemScores[annotation.ReviewID] = current
		}
		current.count++
		for name, value := range annotation.Scores {
			if current.sum[name] == nil {
				current.sum[name] = new(big.Int)
			}
			current.sum[name].Add(current.sum[name], big.NewInt(int64(value)))
		}
	}
	variants := make([]string, 0, len(available))
	for variant := range available {
		variants = append(variants, variant)
	}
	sort.Strings(variants)
	variantItems := map[string]map[string]*score{}
	// A subject may legitimately have multiple independently identified items
	// for one variant. Keep every item average instead of allowing map
	// iteration order to choose one survivor.
	subjectScores := map[string]map[string]map[string][]*big.Rat{}
	reviewIDs := make([]string, 0, len(itemScores))
	for id := range itemScores {
		reviewIDs = append(reviewIDs, id)
	}
	sort.Strings(reviewIDs)
	for _, id := range reviewIDs {
		item := itemScores[id]
		key := keyByID[id]
		if variantItems[key.Variant] == nil {
			variantItems[key.Variant] = map[string]*score{}
		}
		variantItems[key.Variant][id] = item
		if subjectScores[key.SubjectID] == nil {
			subjectScores[key.SubjectID] = map[string]map[string][]*big.Rat{}
		}
		if subjectScores[key.SubjectID][key.Variant] == nil {
			subjectScores[key.SubjectID][key.Variant] = map[string][]*big.Rat{}
		}
		for dimension, total := range item.sum {
			average := new(big.Rat).SetFrac(total, big.NewInt(int64(item.count)))
			subjectScores[key.SubjectID][key.Variant][dimension] = append(
				subjectScores[key.SubjectID][key.Variant][dimension], average,
			)
		}
	}
	report := Report{QueueItems: len(keys), ReviewedItems: len(itemScores), AnnotationCount: len(annotations)}
	for _, variant := range variants {
		summary := VariantSummary{Variant: variant, Available: available[variant], Reviewed: len(variantItems[variant]), Dimensions: map[string]DimensionSummary{}}
		subjects := map[string]struct{}{}
		ids := make([]string, 0, len(variantItems[variant]))
		for id := range variantItems[variant] {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			item := variantItems[variant][id]
			subjects[keyByID[id].SubjectID] = struct{}{}
			for dimension, total := range item.sum {
				current := summary.Dimensions[dimension]
				current.Count++
				current.Mean += bigIntFloat64(total) / float64(item.count)
				summary.Dimensions[dimension] = current
			}
			summary.AnnotationCount += item.count
		}
		summary.UniqueSubjects = len(subjects)
		for dimension, current := range summary.Dimensions {
			current.Mean /= float64(current.Count)
			summary.Dimensions[dimension] = current
		}
		report.Variants = append(report.Variants, summary)
	}
	for left := 0; left < len(variants); left++ {
		for right := left + 1; right < len(variants); right++ {
			for _, dimension := range dimensions {
				var deltas []float64
				subjects := make([]string, 0, len(subjectScores))
				for subject := range subjectScores {
					subjects = append(subjects, subject)
				}
				sort.Strings(subjects)
				for _, subject := range subjects {
					byVariant := subjectScores[subject]
					a, aOK := byVariant[variants[left]]
					b, bOK := byVariant[variants[right]]
					if aOK && bOK && len(a[dimension.Name]) > 0 && len(b[dimension.Name]) > 0 {
						delta := new(big.Rat).Sub(meanRationals(a[dimension.Name]), meanRationals(b[dimension.Name]))
						value, _ := new(big.Float).SetRat(delta).Float64()
						deltas = append(deltas, value)
					}
				}
				report.Comparisons = append(report.Comparisons, summarizeDeltas(variants[left], variants[right], dimension.Name, deltas))
			}
		}
	}
	report.OverlappingItems, report.ReviewerOverlaps = summarizeReviewerOverlap(keyByID, annotations, dimensions)
	return report
}

func bigIntFloat64(value *big.Int) float64 {
	result, _ := new(big.Float).SetInt(value).Float64()
	return result
}

func meanRationals(values []*big.Rat) *big.Rat {
	total := new(big.Rat)
	for _, value := range values {
		total.Add(total, value)
	}
	return total.Quo(total, big.NewRat(int64(len(values)), 1))
}

func summarizeReviewerOverlap(keys map[string]KeyEntry, annotations []Annotation, dimensions []Dimension) (int, []ReviewerOverlapSummary) {
	byReview := map[string]map[string]Annotation{}
	for _, annotation := range annotations {
		if byReview[annotation.ReviewID] == nil {
			byReview[annotation.ReviewID] = map[string]Annotation{}
		}
		byReview[annotation.ReviewID][annotation.Reviewer] = annotation
	}
	ids := make([]string, 0, len(byReview))
	for id := range byReview {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	type accumulator struct {
		count, exact, maxDelta int
		totalDelta             float64
	}
	type pair struct {
		a, b          string
		items         int
		dimensions    map[string]*accumulator
		disagreements []ReviewerDisagreement
	}
	pairs := map[string]*pair{}
	overlapping := 0
	for _, id := range ids {
		reviewersByID := byReview[id]
		if len(reviewersByID) < 2 {
			continue
		}
		overlapping++
		reviewers := make([]string, 0, len(reviewersByID))
		for reviewer := range reviewersByID {
			reviewers = append(reviewers, reviewer)
		}
		sort.Strings(reviewers)
		for left := 0; left < len(reviewers); left++ {
			for right := left + 1; right < len(reviewers); right++ {
				a, b := reviewers[left], reviewers[right]
				pairKey := a + "\x00" + b
				current := pairs[pairKey]
				if current == nil {
					current = &pair{a: a, b: b, dimensions: map[string]*accumulator{}}
					pairs[pairKey] = current
				}
				current.items++
				key := keys[id]
				for _, dimension := range dimensions {
					delta := absInt(reviewersByID[a].Scores[dimension.Name] - reviewersByID[b].Scores[dimension.Name])
					metric := current.dimensions[dimension.Name]
					if metric == nil {
						metric = &accumulator{}
						current.dimensions[dimension.Name] = metric
					}
					metric.count++
					metric.totalDelta += float64(delta)
					metric.maxDelta = max(metric.maxDelta, delta)
					if delta == 0 {
						metric.exact++
						continue
					}
					current.disagreements = append(current.disagreements, ReviewerDisagreement{ReviewID: id, SubjectID: key.SubjectID, Variant: key.Variant, Dimension: dimension.Name, ScoreA: reviewersByID[a].Scores[dimension.Name], ScoreB: reviewersByID[b].Scores[dimension.Name], Difference: delta})
				}
			}
		}
	}
	pairKeys := make([]string, 0, len(pairs))
	for key := range pairs {
		pairKeys = append(pairKeys, key)
	}
	sort.Strings(pairKeys)
	summaries := make([]ReviewerOverlapSummary, 0, len(pairKeys))
	for _, pairKey := range pairKeys {
		current := pairs[pairKey]
		summary := ReviewerOverlapSummary{ReviewerA: current.a, ReviewerB: current.b, OverlapItems: current.items, Dimensions: map[string]ReviewerDimensionAgreement{}, Disagreements: current.disagreements}
		for dimension, metric := range current.dimensions {
			summary.Dimensions[dimension] = ReviewerDimensionAgreement{Count: metric.count, ExactAgreements: metric.exact, ExactAgreementRate: float64(metric.exact) / float64(metric.count), MeanAbsoluteDifference: float64(metric.totalDelta) / float64(metric.count), MaxAbsoluteDifference: metric.maxDelta}
		}
		summaries = append(summaries, summary)
	}
	return overlapping, summaries
}

func summarizeDeltas(a, b, dimension string, deltas []float64) PairedComparison {
	result := PairedComparison{VariantA: a, VariantB: b, Dimension: dimension, ReviewedPairs: len(deltas)}
	if len(deltas) == 0 {
		return result
	}
	sorted := append([]float64(nil), deltas...)
	sort.Float64s(sorted)
	for _, delta := range deltas {
		result.MeanDelta += delta
		if delta > 0 {
			result.Wins++
		} else if delta < 0 {
			result.Losses++
		} else {
			result.Ties++
		}
	}
	result.MeanDelta /= float64(len(deltas))
	middle := len(sorted) / 2
	result.MedianDelta = sorted[middle]
	if len(sorted)%2 == 0 {
		result.MedianDelta = (sorted[middle-1] + sorted[middle]) / 2
	}
	if len(deltas) >= 30 {
		low, high := deterministicBootstrap95(deltas, 2000)
		result.Bootstrap95Low, result.Bootstrap95High = &low, &high
	}
	return result
}
func deterministicBootstrap95(values []float64, samples int) (float64, float64) {
	means := make([]float64, samples)
	var state uint64 = 0x9e3779b97f4a7c15
	for sample := range samples {
		var total float64
		for range values {
			state ^= state << 13
			state ^= state >> 7
			state ^= state << 17
			total += values[state%uint64(len(values))]
		}
		means[sample] = total / float64(len(values))
	}
	sort.Float64s(means)
	return means[int(math.Floor(.025*float64(samples)))], means[min(samples-1, int(math.Ceil(.975*float64(samples)))-1)]
}
func absInt(value int) int {
	if value < 0 {
		return -value
	}
	return value
}
