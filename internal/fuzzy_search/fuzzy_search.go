package fuzzy_search

import (
	"slices"
	"strings"

	"github.com/muleyuck/linippet/internal/linippet"
)

func backtraceMatch(dp map[int]map[int]int, query, target string, i, j int) []int {
	positions := []int{}

	for i > 0 && j > 0 {
		if query[i-1] == target[j-1] {
			positions = append([]int{j - 1}, positions...)
			i--
			j--
		} else if dp[i][j] == dp[i-1][j]+1 {
			i--
		} else {
			j--
		}
	}

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

	dp := make(map[int]map[int]int)
	for i := range queryLen + 1 {
		dp[i] = make(map[int]int, snippetLen+1)
		dp[i][0] = i
	}
	for j := 0; j <= snippetLen; j++ {
		dp[0][j] = 0
	}

	// Fill DP table
	for i := 1; i <= queryLen; i++ {
		for j := 1; j <= snippetLen; j++ {
			if queryLower[i-1] == snippetLower[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				// minimum delete or skip
				deletion := dp[i-1][j] + 1
				insertion := dp[i][j-1]
				dp[i][j] = min(deletion, insertion)
			}
		}
	}

	// minimum match index and cost
	minCost := dp[queryLen][0]
	bestPos := 0
	for j := 1; j <= snippetLen; j++ {
		if dp[queryLen][j] < minCost {
			minCost = dp[queryLen][j]
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
	matchPositions := backtraceMatch(dp, queryLower, snippetLower, queryLen, bestPos)
	score := calculateScore(query, snippet, matchPositions)
	return matchPositions, score
}

const BASE_MATCH_SCORE = 100

func isAlphaNumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

func calculateScore(query, snippet string, matches []int) int {
	if len(matches) == 0 {
		return 0
	}

	score := BASE_MATCH_SCORE
	// Exact match Bonus
	if snippet == query {
		score += 1000
	}
	// Prefix match Bonus
	prefixMatchLength := 0
	for i := 0; i < len(query) && i < len(snippet); i++ {
		if query[i] != snippet[i] {
			break
		}
		prefixMatchLength++
	}
	score += prefixMatchLength * 15

	consecutiveCount := 0
	for i, idx := range matches {
		// Continuous match Bonus
		if i > 0 {
			if idx == matches[i-1]+1 {
				consecutiveCount++
				score += consecutiveCount * consecutiveCount * 3
			} else {
				consecutiveCount = 0
			}
		}
		// Top of snippet word Bonus
		if idx == 0 || !isAlphaNumeric(rune(snippet[idx-1])) {
			score += 30
		}
		// Position Bonus (position that is matched is in front of snippet, higher score)
		score += int(float64(len(snippet)-idx) * (10.0 / float64(i+1)))
	}
	// Similarity length Bonus
	matchRatio := float64(len(query)) / float64(len(snippet))
	if matchRatio > 0.7 {
		score += int(matchRatio * 100)
	}
	// Top of char match Bonus
	if len(matches) > 0 && query[0] == snippet[matches[0]] {
		score += 25
	}
	// Bottom of char match Bonus
	if len(matches) > 0 && len(query) > 0 && query[len(query)-1] == snippet[matches[len(matches)-1]] {
		score += 15
	}

	return score
}

type SearchResult struct {
	Linippet linippet.Linippet
	Matches  []int
	score    int
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
			results = append(results, SearchResult{Linippet: linippet, Matches: allMatches, score: totalScore})
		}
	}

	// sort desc by score
	slices.SortFunc(results, func(a, b SearchResult) int {
		return b.score - a.score
	})
	return results
}
