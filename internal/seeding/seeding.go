package seeding

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
)

//go:embed data/cutoffs.json
var cutoffsJSON []byte

//go:embed data/qualifiers.json
var qualifiersJSON []byte

// WaveBracket defines min and max PBE for a specific wave.
type WaveBracket struct {
	Wave   int     `json:"wave"`
	MinPBE float64 `json:"min_pbe"`
	MaxPBE float64 `json:"max_pbe"`
}

// Qualifier defines a recognized qualifier race and its percentage adjustment.
type Qualifier struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Adjustment float64  `json:"adjustment"`
	Distance   float64  `json:"distance,omitempty"`
	Technique  string   `json:"technique,omitempty"`
	IsOfficial bool     `json:"is_official"`
	Keywords   []string `json:"keywords,omitempty"`
}

// EventInfo represents basic metadata for Birkie/Korte events.
type EventInfo struct {
	Key        string  `json:"key"`
	Name       string  `json:"name"`
	DistanceKm float64 `json:"distance_km"`
	MaxWaves   int     `json:"max_waves"`
}

var SupportedEvents = []EventInfo{
	{Key: "birkie_skate", Name: "American Birkebeiner Skate 50K", DistanceKm: 50.0, MaxWaves: 8},
	{Key: "birkie_classic", Name: "American Birkebeiner Classic 53K", DistanceKm: 53.0, MaxWaves: 5},
	{Key: "korte_skate", Name: "Kortelopet Skate 29K", DistanceKm: 29.0, MaxWaves: 3},
	{Key: "korte_classic", Name: "Kortelopet Classic 29K", DistanceKm: 29.0, MaxWaves: 3},
}

// EventQualification holds the computed wave result for an event.
type EventQualification struct {
	EventKey           string  `json:"event_key"`
	EventName          string  `json:"event_name"`
	RawPercentBack     float64 `json:"raw_percent_back"`
	Adjustment         float64 `json:"adjustment"`
	FinalPBE           float64 `json:"final_pbe"`
	QualifiedWave      int     `json:"qualified_wave"`
	QualifiedWaveLabel string  `json:"qualified_wave_label"`
	CutoffPBE          float64 `json:"cutoff_pbe"`
	DistanceToNext     float64 `json:"distance_to_next"`
	HasNextWave        bool    `json:"has_next_wave"`
}

// RaceInputData is a decoupled representation of a race for seeding calculations.
type RaceInputData struct {
	ID           string  `json:"id"`
	RaceName     string  `json:"raceName"`
	RaceDate     string  `json:"raceDate"`
	RaceDistance float64 `json:"raceDistance"`
	PercentBack  float64 `json:"percentBack"`
}

// RaceSeeding holds a race with all its event qualifications.
type RaceSeeding struct {
	ID                  string                        `json:"id"`
	RaceName            string                        `json:"race_name"`
	RaceDate            string                        `json:"race_date"`
	RaceDistance        float64                       `json:"race_distance"`
	PercentBack         float64                       `json:"percent_back"`
	Season              string                        `json:"season"`
	IsOfficialQualifier bool                          `json:"is_official_qualifier"`
	MatchedQualifier    Qualifier                     `json:"matched_qualifier"`
	Qualifications      map[string]EventQualification `json:"qualifications"`
}

// SeasonSummary aggregates wave qualifications for a ski season.
type SeasonSummary struct {
	Season         string                 `json:"season"`
	BestWaves      map[string]int         `json:"best_waves"`
	BestWaveLabels map[string]string      `json:"best_wave_labels"`
	BestRaces      map[string]string      `json:"best_races"` // race name that earned the best wave
	Races          []RaceSeeding          `json:"races"`
}

// Store holds cached cutoffs and qualifiers.
type Store struct {
	Cutoffs    map[string][]WaveBracket
	Qualifiers []Qualifier
}

var (
	globalStore *Store
	storeOnce   sync.Once
)

// GetStore returns the singleton Store initialized once.
func GetStore() *Store {
	storeOnce.Do(func() {
		var cutoffs map[string][]WaveBracket
		if err := json.Unmarshal(cutoffsJSON, &cutoffs); err != nil {
			panic(fmt.Sprintf("failed to parse cutoffs.json: %v", err))
		}

		var qualifiers []Qualifier
		if err := json.Unmarshal(qualifiersJSON, &qualifiers); err != nil {
			panic(fmt.Sprintf("failed to parse qualifiers.json: %v", err))
		}

		globalStore = &Store{
			Cutoffs:    cutoffs,
			Qualifiers: qualifiers,
		}
	})
	return globalStore
}

// GetSeason calculates the Nordic ski season (July 1 - June 30) for a given date.
// e.g. "2024-02-15" -> "2023-2024", "2024-08-10" -> "2024-2025".
func GetSeason(dateStr string) string {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return "Unknown"
	}

	formats := []string{
		"2006-01-02",
		"2006/01/02",
		time.RFC3339,
		"2006-01-02T15:04:05",
		"01/02/2006",
	}

	var parsed time.Time
	var err error
	for _, layout := range formats {
		parsed, err = time.Parse(layout, dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		// Fallback: try extracting 4-digit year
		if len(dateStr) >= 4 {
			return dateStr[:4]
		}
		return "Unknown"
	}

	year := parsed.Year()
	if parsed.Month() >= time.July {
		return fmt.Sprintf("%d-%d", year, year+1)
	}
	return fmt.Sprintf("%d-%d", year-1, year)
}

// MatchQualifier inspects a race name and distance, returning the best matching Qualifier.
func MatchQualifier(raceName string, raceDistance float64) Qualifier {
	store := GetStore()
	cleanName := strings.ToLower(strings.TrimSpace(raceName))

	if cleanName == "" {
		return Qualifier{ID: "none", Name: "Non-Qualifier", Adjustment: 0.0, IsOfficial: false}
	}

	// 1. Explicit rollerski detection: Rollerski races are never official on-snow Birkie qualifiers
	if strings.Contains(cleanName, "rollerski") || strings.Contains(cleanName, "roller ski") || strings.Contains(cleanName, "roller-ski") {
		return Qualifier{
			ID:         "none",
			Name:       "Non-Qualifier (Rollerski)",
			Adjustment: 0.0,
			IsOfficial: false,
		}
	}

	// 2. Minimum distance threshold: Official Birkie qualifiers are marathons/half-marathons (shortest is Sisu 15K).
	// If the race is explicitly entered with distance < 14 km, it cannot be an official Birkie qualifier.
	if raceDistance > 0 && raceDistance < 14.0 {
		return Qualifier{
			ID:         "none",
			Name:       "Non-Qualifier (<15km)",
			Adjustment: 0.0,
			IsOfficial: false,
		}
	}

	var candidates []Qualifier
	bestLen := 0

	for i := range store.Qualifiers {
		q := store.Qualifiers[i]
		if q.ID == "none" || !q.IsOfficial {
			continue
		}
		for _, kw := range q.Keywords {
			kwLower := strings.ToLower(kw)
			if strings.Contains(cleanName, kwLower) {
				// Distance sanity check: If both race and qualifier have distances, don't match if diff > 15 km
				if raceDistance > 0 && q.Distance > 0 && math.Abs(raceDistance-q.Distance) > 15.0 {
					continue
				}
				if len(kwLower) > bestLen {
					bestLen = len(kwLower)
					candidates = []Qualifier{q}
				} else if len(kwLower) == bestLen {
					candidates = append(candidates, q)
				}
			}
		}
	}

	if len(candidates) == 0 {
		return Qualifier{
			ID:         "none",
			Name:       "Non-Qualifier",
			Adjustment: 0.0,
			IsOfficial: false,
		}
	}

	if len(candidates) == 1 {
		return candidates[0]
	}

	// If multiple candidates share the longest keyword match, use raceDistance to pick closest
	if raceDistance > 0 {
		bestCandidate := candidates[0]
		minDiff := math.Abs(raceDistance - bestCandidate.Distance)
		for _, c := range candidates[1:] {
			if c.Distance > 0 {
				diff := math.Abs(raceDistance - c.Distance)
				if diff < minDiff {
					minDiff = diff
					bestCandidate = c
				}
			}
		}
		return bestCandidate
	}

	return candidates[0]
}

// EvaluateWaveForEvent calculates wave qualification for a specific event.
func EvaluateWaveForEvent(eventKey string, rawPercentBack float64, adjustment float64) EventQualification {
	store := GetStore()
	brackets, ok := store.Cutoffs[eventKey]

	eventName := eventKey
	for _, e := range SupportedEvents {
		if e.Key == eventKey {
			eventName = e.Name
			break
		}
	}

	finalPBE := math.Round((rawPercentBack+adjustment)*100) / 100
	rawRounded := math.Round(rawPercentBack*100) / 100

	if !ok || len(brackets) == 0 {
		return EventQualification{
			EventKey:           eventKey,
			EventName:          eventName,
			RawPercentBack:     rawRounded,
			Adjustment:         adjustment,
			FinalPBE:           finalPBE,
			QualifiedWave:      0,
			QualifiedWaveLabel: "Unknown",
			CutoffPBE:          0,
			DistanceToNext:     0,
			HasNextWave:        false,
		}
	}

	qualifiedWave := -1
	var matchedBracket WaveBracket
	var prevBracket *WaveBracket

	for i := range brackets {
		b := brackets[i]
		if (b.Wave == 1 && finalPBE <= b.MaxPBE) || (finalPBE >= b.MinPBE && finalPBE <= b.MaxPBE) || (finalPBE <= b.MaxPBE && qualifiedWave == -1) {
			qualifiedWave = b.Wave
			matchedBracket = b
			if i > 0 {
				prevBracket = &brackets[i-1]
			}
			break
		}
	}

	if qualifiedWave == -1 {
		lastBracket := brackets[len(brackets)-1]
		qualifiedWave = lastBracket.Wave + 1
		matchedBracket = WaveBracket{
			Wave:   qualifiedWave,
			MinPBE: lastBracket.MaxPBE + 0.01,
			MaxPBE: math.MaxFloat64,
		}
		prevBracket = &lastBracket
	}

	label := fmt.Sprintf("Wave %d", qualifiedWave)
	if qualifiedWave > brackets[len(brackets)-1].Wave {
		label = fmt.Sprintf("Wave %d+ (Unseeded)", brackets[len(brackets)-1].Wave)
	}

	var distanceToNext float64
	hasNext := false
	if prevBracket != nil {
		hasNext = true
		distanceToNext = math.Round((finalPBE-prevBracket.MaxPBE)*100) / 100
		if distanceToNext < 0 {
			distanceToNext = 0
		}
	}

	return EventQualification{
		EventKey:           eventKey,
		EventName:          eventName,
		RawPercentBack:     rawRounded,
		Adjustment:         adjustment,
		FinalPBE:           finalPBE,
		QualifiedWave:      qualifiedWave,
		QualifiedWaveLabel: label,
		CutoffPBE:          matchedBracket.MaxPBE,
		DistanceToNext:     distanceToNext,
		HasNextWave:        hasNext,
	}
}

// GetEventAdjustment computes the event-specific adjustment when a race is applied to a target event.
func GetEventAdjustment(targetEvent string, qualifier Qualifier, raceName string) (float64, string) {
	// If the qualifier is not official or is "none"
	if qualifier.ID == "none" || !qualifier.IsOfficial {
		return 0.0, "Non-Qualifier (0.0%)"
	}

	// If the user explicitly chose one of the directional transfer options:
	switch qualifier.ID {
	case "birkie_classic_to_skate":
		if targetEvent == "birkie_skate" {
			return 8.0, "+8.0% (Classic to Skate)"
		}
		return 0.0, "Direct (0.0%)"
	case "birkie_skate_to_classic":
		if targetEvent == "birkie_classic" {
			return -8.0, "-8.0% (Skate to Classic)"
		}
		return 0.0, "Direct (0.0%)"
	case "birkie_to_korte":
		if targetEvent == "korte_skate" || targetEvent == "korte_classic" {
			return -15.0, "-15.0% (Birkie to Korte)"
		}
		return 0.0, "Direct (0.0%)"
	}

	qID := qualifier.ID
	rLower := strings.ToLower(raceName)

	isBirkieSkate := qID == "birkie_skate" || (strings.Contains(rLower, "birkie") && (strings.Contains(rLower, "skate") || strings.Contains(rLower, "freestyle")))
	isBirkieClassic := qID == "birkie_classic" || (strings.Contains(rLower, "birkie") && strings.Contains(rLower, "classic"))
	isKorteSkate := qID == "korte_skate" || (strings.Contains(rLower, "korte") && (strings.Contains(rLower, "skate") || strings.Contains(rLower, "freestyle")))
	isKorteClassic := qID == "korte_classic" || (strings.Contains(rLower, "korte") && strings.Contains(rLower, "classic"))

	if isBirkieSkate {
		switch targetEvent {
		case "birkie_skate":
			return 0.0, "Direct (0.0%)"
		case "birkie_classic":
			return -8.0, "-8.0% (Skate to Classic)"
		case "korte_skate", "korte_classic":
			return -15.0, "-15.0% (Birkie to Korte)"
		default:
			return 0.0, "Direct (0.0%)"
		}
	}

	if isBirkieClassic {
		switch targetEvent {
		case "birkie_classic":
			return 0.0, "Direct (0.0%)"
		case "birkie_skate":
			return 8.0, "+8.0% (Classic to Skate)"
		case "korte_skate", "korte_classic":
			return -15.0, "-15.0% (Birkie to Korte)"
		default:
			return 0.0, "Direct (0.0%)"
		}
	}

	if isKorteSkate {
		switch targetEvent {
		case "korte_skate":
			return 0.0, "Direct (0.0%)"
		case "korte_classic":
			return -8.0, "-8.0% (Skate to Classic)"
		default:
			return 0.0, "Direct (0.0%)"
		}
	}

	if isKorteClassic {
		switch targetEvent {
		case "korte_classic":
			return 0.0, "Direct (0.0%)"
		case "korte_skate":
			return 8.0, "+8.0% (Classic to Skate)"
		default:
			return 0.0, "Direct (0.0%)"
		}
	}

	// Regional qualifiers (Mora, Seeley, Noque, etc.)
	baseAdj := qualifier.Adjustment
	isClassic := strings.Contains(qID, "classic") || strings.Contains(rLower, "classic") || qualifier.Technique == "classic"

	if isClassic {
		switch targetEvent {
		case "birkie_classic", "korte_classic":
			return baseAdj, fmt.Sprintf("%+.1f%%", baseAdj)
		case "birkie_skate", "korte_skate":
			return baseAdj + 8.0, fmt.Sprintf("%+.1f%% (+8%% technique)", baseAdj+8.0)
		}
	} else {
		// Skate qualifier
		switch targetEvent {
		case "birkie_skate", "korte_skate":
			return baseAdj, fmt.Sprintf("%+.1f%%", baseAdj)
		case "birkie_classic", "korte_classic":
			return baseAdj - 8.0, fmt.Sprintf("%+.1f%% (-8%% technique)", baseAdj-8.0)
		}
	}

	return baseAdj, fmt.Sprintf("%+.1f%%", baseAdj)
}

// EvaluateRaceSeeding evaluates a single race across all 4 supported events.
func EvaluateRaceSeeding(race RaceInputData, customQualifierID string) RaceSeeding {
	var qualifier Qualifier
	if customQualifierID != "" && customQualifierID != "none" {
		store := GetStore()
		for _, q := range store.Qualifiers {
			if q.ID == customQualifierID {
				qualifier = q
				break
			}
		}
	}
	if qualifier.ID == "" {
		qualifier = MatchQualifier(race.RaceName, race.RaceDistance)
	}

	isOfficial := qualifier.IsOfficial && qualifier.ID != "none"

	season := GetSeason(race.RaceDate)
	quals := make(map[string]EventQualification, len(SupportedEvents))

	for _, event := range SupportedEvents {
		adj, _ := GetEventAdjustment(event.Key, qualifier, race.RaceName)
		evQual := EvaluateWaveForEvent(event.Key, race.PercentBack, adj)
		if !isOfficial {
			evQual.QualifiedWaveLabel = fmt.Sprintf("Est. %s*", evQual.QualifiedWaveLabel)
		}
		quals[event.Key] = evQual
	}

	return RaceSeeding{
		ID:                  race.ID,
		RaceName:            race.RaceName,
		RaceDate:            race.RaceDate,
		RaceDistance:        race.RaceDistance,
		PercentBack:         race.PercentBack,
		Season:              season,
		IsOfficialQualifier: isOfficial,
		MatchedQualifier:    qualifier,
		Qualifications:      quals,
	}
}

// GroupAndEvaluateSeasons groups races by ski season and determines best wave achievements.
func GroupAndEvaluateSeasons(races []RaceInputData) []SeasonSummary {
	seasonsMap := make(map[string][]RaceSeeding)

	for _, r := range races {
		seeding := EvaluateRaceSeeding(r, "")
		seasonsMap[seeding.Season] = append(seasonsMap[seeding.Season], seeding)
	}

	// Sort season names descending (e.g. 2024-2025, 2023-2024...)
	seasons := make([]string, 0, len(seasonsMap))
	for s := range seasonsMap {
		seasons = append(seasons, s)
	}
	sort.Slice(seasons, func(i, j int) bool {
		return seasons[i] > seasons[j]
	})

	summaries := make([]SeasonSummary, 0, len(seasons))
	for _, season := range seasons {
		raceList := seasonsMap[season]

		// Sort races by date descending
		sort.Slice(raceList, func(i, j int) bool {
			return raceList[i].RaceDate > raceList[j].RaceDate
		})

		bestWaves := make(map[string]int)
		bestWaveLabels := make(map[string]string)
		bestRaces := make(map[string]string)

		// Initialize best waves
		for _, ev := range SupportedEvents {
			bestWaves[ev.Key] = 999
			bestWaveLabels[ev.Key] = "Not Qualified"
			bestRaces[ev.Key] = ""
		}

		for _, race := range raceList {
			if !race.IsOfficialQualifier {
				continue
			}
			for _, ev := range SupportedEvents {
				q, ok := race.Qualifications[ev.Key]
				if ok && q.QualifiedWave > 0 && q.QualifiedWave <= ev.MaxWaves {
					if q.QualifiedWave < bestWaves[ev.Key] {
						bestWaves[ev.Key] = q.QualifiedWave
						bestWaveLabels[ev.Key] = q.QualifiedWaveLabel
						bestRaces[ev.Key] = race.RaceName
					}
				}
			}
		}

		// Clean up any event where no valid wave was achieved
		for _, ev := range SupportedEvents {
			if bestWaves[ev.Key] == 999 {
				bestWaves[ev.Key] = 0
			}
		}

		summaries = append(summaries, SeasonSummary{
			Season:         season,
			BestWaves:      bestWaves,
			BestWaveLabels: bestWaveLabels,
			BestRaces:      bestRaces,
			Races:          raceList,
		})
	}

	return summaries
}
