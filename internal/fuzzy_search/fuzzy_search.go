package fuzzy_search

import (
	"slices"
	"strings"

	"github.com/muleyuck/linippet/internal/linippet"
)

// fzf-aligned scoring constants
const (
	scoreMatch        = 16 // base score per matched character
	scoreGapStart     = -3 // penalty for starting a gap
	scoreGapExtension = -1 // penalty per additional gap character

	bonusBoundaryWhite = 10 // match after whitespace
	bonusBoundaryDelim = 9  // match after delimiter (/ . , : ; |)
	bonusBoundary      = 8  // match after other non-alphanumeric
	bonusCamelCase     = 7  // lower→upper transition
	bonusConsecutive   = 4  // consecutive match
	bonusFirstCharMul  = 2  // multiplier for first matched character bonus

	bonusExactMatch    = 100 // exact match bonus
	bonusPrefixChar    = 4   // per-character prefix match bonus
	bonusSimilarityMax = 30  // max similarity length bonus
)

// charClass categorizes characters for boundary bonus calculation.
type charClass int

const (
	charWhite     charClass = iota // whitespace
	charDelimiter                  // / . , : ; | - _
	charNonWord                    // other non-alphanumeric
	charLower                      // lowercase letter
	charUpper                      // uppercase letter
	charNumber                     // digit
)

func charClassOf(c byte) charClass {
	switch {
	case c == ' ' || c == '\t' || c == '\n' || c == '\r':
		return charWhite
	case c == '/' || c == '.' || c == ',' || c == ':' || c == ';' || c == '|' || c == '-' || c == '_':
		return charDelimiter
	case c >= 'a' && c <= 'z':
		return charLower
	case c >= 'A' && c <= 'Z':
		return charUpper
	case c >= '0' && c <= '9':
		return charNumber
	default:
		return charNonWord
	}
}

func isAlphaNumClass(c charClass) bool {
	return c == charLower || c == charUpper || c == charNumber
}

func bonusFor(prevClass, currClass charClass) int {
	if !isAlphaNumClass(currClass) {
		return 0
	}
	switch {
	case prevClass == charWhite:
		return bonusBoundaryWhite
	case prevClass == charDelimiter:
		return bonusBoundaryDelim
	case !isAlphaNumClass(prevClass):
		return bonusBoundary
	case prevClass == charLower && currClass == charUpper:
		return bonusCamelCase
	default:
		return 0
	}
}

func backtraceMatch(dp []int, cols int, query, target string, i, j int) []int {
	positions := []int{}

	for i > 0 && j > 0 {
		if query[i-1] == target[j-1] {
			positions = append(positions, j-1)
			i--
			j--
		} else if dp[i*cols+j] == dp[(i-1)*cols+j]+1 {
			i--
		} else {
			j--
		}
	}

	slices.Reverse(positions)
	return positions
}

func fuzzyMatch(query, snippet string) ([]int, int) {
	queryLen, snippetLen := len(query), len(snippet)
	if queryLen == 0 {
		return []int{}, 0
	}
	if snippetLen == 0 {
		return nil, 0
	}

	queryLower := strings.ToLower(query)
	snippetLower := strings.ToLower(snippet)

	cols := snippetLen + 1
	dp := make([]int, (queryLen+1)*cols)
	for i := range queryLen + 1 {
		dp[i*cols] = i
	}
	// dp[0*cols+j] is already 0 from make

	// Fill DP table
	for i := 1; i <= queryLen; i++ {
		for j := 1; j <= snippetLen; j++ {
			if queryLower[i-1] == snippetLower[j-1] {
				dp[i*cols+j] = dp[(i-1)*cols+(j-1)]
			} else {
				// minimum delete or skip
				deletion := dp[(i-1)*cols+j] + 1
				insertion := dp[i*cols+(j-1)]
				dp[i*cols+j] = min(deletion, insertion)
			}
		}
	}

	// minimum match index and cost
	minCost := dp[queryLen*cols]
	bestPos := 0
	for j := 1; j <= snippetLen; j++ {
		if dp[queryLen*cols+j] < minCost {
			minCost = dp[queryLen*cols+j]
			bestPos = j
		}
		if minCost == 0 {
			break
		}
	}

	// Exclude cost more than 0
	if minCost != 0 {
		return nil, 0
	}
	matchPositions := backtraceMatch(dp, cols, queryLower, snippetLower, queryLen, bestPos)
	score := calculateScore(query, snippet, matchPositions)
	return matchPositions, score
}

func calculateScore(query, snippet string, matches []int) int {
	if len(matches) == 0 {
		return 0
	}

	score := 0
	consecutiveStartBonus := 0 // bonus of the character that started the consecutive run
	prevMatchIdx := -1

	for i, pos := range matches {
		// Base match score
		charScore := scoreMatch

		// Context bonus from character class transition
		var prevClass charClass
		if pos == 0 {
			prevClass = charWhite // treat start of string as whitespace boundary
		} else {
			prevClass = charClassOf(snippet[pos-1])
		}
		currClass := charClassOf(snippet[pos])
		ctxBonus := bonusFor(prevClass, currClass)

		if i > 0 && pos == prevMatchIdx+1 {
			// Consecutive match: inherit the max of context bonus,
			// the consecutive-run start bonus, and base consecutive bonus.
			// This is a key fzf property where consecutive matches inherit
			// the boundary bonus of the character that started the run.
			bonus := max(ctxBonus, consecutiveStartBonus, bonusConsecutive)
			charScore += bonus
		} else {
			// Non-consecutive: apply gap penalty
			if i > 0 {
				gapLen := pos - prevMatchIdx - 1
				charScore += scoreGapStart + scoreGapExtension*(gapLen-1)
			}
			charScore += ctxBonus
			consecutiveStartBonus = ctxBonus
		}

		// First matched character gets bonus multiplied
		if i == 0 {
			charScore += ctxBonus * (bonusFirstCharMul - 1)
		}

		score += charScore
		prevMatchIdx = pos
	}

	// Global bonuses

	// Exact match
	if strings.EqualFold(query, snippet) {
		score += bonusExactMatch
	}

	// Prefix match
	queryLower := strings.ToLower(query)
	snippetLower := strings.ToLower(snippet)
	prefixLen := 0
	for i := 0; i < len(queryLower) && i < len(snippetLower); i++ {
		if queryLower[i] != snippetLower[i] {
			break
		}
		prefixLen++
	}
	score += prefixLen * bonusPrefixChar

	// Similarity length bonus
	matchRatio := float64(len(query)) / float64(len(snippet))
	if matchRatio > 0.5 {
		score += int(matchRatio * bonusSimilarityMax)
	}

	return score
}

type SearchResult struct {
	Linippet linippet.Linippet
	Matches  []int
	Score    int
}

func FuzzySearch(query string, linippets linippet.Linippets) []SearchResult {
	// split query by whitespace
	queries := strings.Fields(query)
	if len(queries) == 0 {
		return []SearchResult{}
	}
	results := make([]SearchResult, 0)

	for _, linippet := range linippets {
		allMatched := true
		allMatches := make([]int, 0)
		totalScore := 0
		for _, q := range queries {
			// TODO: use cache
			matches, score := fuzzyMatch(q, linippet.Snippet)
			if matches == nil {
				allMatched = false
				break
			}
			allMatches = append(allMatches, matches...)
			totalScore += score
		}
		if allMatched {
			results = append(results, SearchResult{Linippet: linippet, Matches: allMatches, Score: totalScore})
		}
	}

	// sort desc by score, then asc by snippet length as tiebreaker
	slices.SortFunc(results, func(a, b SearchResult) int {
		if a.Score != b.Score {
			return b.Score - a.Score
		}
		return len(a.Linippet.Snippet) - len(b.Linippet.Snippet)
	})
	return results
}
