package fuzzy_search

import (
	"context"
	"slices"
	"testing"

	"github.com/muleyuck/linippet/internal/linippet"
)

func TestFuzzySearchRanking(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		snippets   []string
		wantHigher string
		wantLower  string
	}{
		{
			name:       "c aws: contiguous aws match beats scattered",
			query:      "c aws",
			snippets:   []string{"cat ~/.aws/config", "column -t -s ',' ${{csv_file}} | awk '{print $1, $2}' | sort"},
			wantHigher: "cat ~/.aws/config",
			wantLower:  "column -t -s ',' ${{csv_file}} | awk '{print $1, $2}' | sort",
		},
		{
			name:       "docker: prefix match beats non-prefix",
			query:      "docker",
			snippets:   []string{"docker run -d -p ${{port}} ${{image}}", "sudo docker-compose up -d --build ${{service}}"},
			wantHigher: "docker run -d -p ${{port}} ${{image}}",
			wantLower:  "sudo docker-compose up -d --build ${{service}}",
		},
		{
			name:       "git br: consecutive br match beats scattered",
			query:      "git br",
			snippets:   []string{"git branch -a | rg ${{pattern}}", "git rebase -i ${{branch}}"},
			wantHigher: "git branch -a | rg ${{pattern}}",
			wantLower:  "git rebase -i ${{branch}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linippets := make(linippet.Linippets, len(tt.snippets))
			for i, s := range tt.snippets {
				linippets[i] = linippet.Linippet{Id: s, Snippet: s}
			}

			results := FuzzySearch(context.Background(), tt.query, linippets)
			if len(results) < 2 {
				t.Fatalf("expected at least 2 results, got %d", len(results))
			}

			higherScore := -1
			lowerScore := -1
			for _, r := range results {
				if r.Linippet.Snippet == tt.wantHigher {
					higherScore = r.Score
				}
				if r.Linippet.Snippet == tt.wantLower {
					lowerScore = r.Score
				}
			}

			if higherScore == -1 {
				t.Fatalf("expected snippet %q in results", tt.wantHigher)
			}
			if lowerScore == -1 {
				t.Fatalf("expected snippet %q in results", tt.wantLower)
			}
			if higherScore <= lowerScore {
				t.Errorf("expected %q (score=%d) to rank higher than %q (score=%d)",
					tt.wantHigher, higherScore, tt.wantLower, lowerScore)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	t.Run("consecutive match scores higher than scattered", func(t *testing.T) {
		snippet := "abcdef"
		consecutive := []int{0, 1, 2} // "abc" consecutive
		scattered := []int{0, 2, 4}   // "a_c_e" scattered

		scoreConsecutive := calculateScore("abc", snippet, consecutive)
		scoreScattered := calculateScore("ace", snippet, scattered)

		if scoreConsecutive <= scoreScattered {
			t.Errorf("consecutive score (%d) should be higher than scattered (%d)",
				scoreConsecutive, scoreScattered)
		}
	})

	t.Run("boundary bonus by character class", func(t *testing.T) {
		// Match after whitespace should get bonusBoundaryWhite
		snippet := "foo bar"
		matches := []int{4} // 'b' after space
		score := calculateScore("b", snippet, matches)

		// Match after letter should get no boundary bonus
		snippet2 := "foobar"
		matches2 := []int{3} // 'b' after 'o'
		score2 := calculateScore("b", snippet2, matches2)

		if score <= score2 {
			t.Errorf("boundary match score (%d) should be higher than mid-word (%d)",
				score, score2)
		}
	})

	t.Run("delimiter boundary bonus", func(t *testing.T) {
		// Match after delimiter (/)
		snippet := "~/.aws/config"
		matches := []int{7} // 'c' after '/'
		score := calculateScore("c", snippet, matches)

		// Match in the middle of a word
		snippet2 := "abcdef"
		matches2 := []int{2} // 'c' after 'b'
		score2 := calculateScore("c", snippet2, matches2)

		if score <= score2 {
			t.Errorf("delimiter boundary score (%d) should be higher than mid-word (%d)",
				score, score2)
		}
	})

	t.Run("gap penalty reduces score", func(t *testing.T) {
		snippet := "abcdefghij"
		// Small gap: match a,b then skip to e
		smallGap := []int{0, 1, 4}
		// Large gap: match a,b then skip to i
		largeGap := []int{0, 1, 8}

		scoreSmallGap := calculateScore("abe", snippet, smallGap)
		scoreLargeGap := calculateScore("abi", snippet, largeGap)

		if scoreSmallGap <= scoreLargeGap {
			t.Errorf("small gap score (%d) should be higher than large gap (%d)",
				scoreSmallGap, scoreLargeGap)
		}
	})

	t.Run("exact match gets bonus", func(t *testing.T) {
		snippet := "docker"
		matches := []int{0, 1, 2, 3, 4, 5}
		score := calculateScore("docker", snippet, matches)

		// Partial match in longer snippet
		snippet2 := "docker run"
		matches2 := []int{0, 1, 2, 3, 4, 5}
		score2 := calculateScore("docker", snippet2, matches2)

		if score <= score2 {
			t.Errorf("exact match score (%d) should be higher than partial (%d)",
				score, score2)
		}
	})

	t.Run("empty matches returns zero", func(t *testing.T) {
		score := calculateScore("foo", "foobar", []int{})
		if score != 0 {
			t.Errorf("expected 0 for empty matches, got %d", score)
		}
	})

	t.Run("camelCase boundary bonus", func(t *testing.T) {
		// 'N' after lowercase 'm' → camelCase bonus
		snippet := "camelName"
		matches := []int{5} // 'N' after 'l'
		score := calculateScore("N", snippet, matches)

		// 'a' after lowercase 'c' → no bonus (lower→lower)
		snippet2 := "camelName"
		matches2 := []int{1} // 'a' after 'c'
		score2 := calculateScore("a", snippet2, matches2)

		if score <= score2 {
			t.Errorf("camelCase boundary score (%d) should be higher than mid-word (%d)",
				score, score2)
		}
	})

	t.Run("first character multiplier doubles boundary bonus", func(t *testing.T) {
		// First match at word boundary: bonus is multiplied by bonusFirstCharMul
		snippet := "foo bar"
		firstMatch := []int{4} // 'b' after space, as first matched char
		scoreFirst := calculateScore("b", snippet, firstMatch)

		// Mid-word first char: no boundary bonus at all
		snippetMid := "foobar"
		midMatch := []int{3} // 'b' after 'o', no boundary
		scoreMid := calculateScore("b", snippetMid, midMatch)

		if scoreFirst <= scoreMid {
			t.Errorf("first char at boundary (%d) should be much higher than first char mid-word (%d)",
				scoreFirst, scoreMid)
		}

		// Verify multiplier effect: boundary bonus should appear doubled.
		// scoreFirst = scoreMatch + bonusBoundaryWhite + bonusBoundaryWhite*(mul-1)
		//            = 16 + 10 + 10 = 36
		// scoreMid   = scoreMatch + 0 = 16
		expectedFirst := scoreMatch + bonusBoundaryWhite*bonusFirstCharMul
		if scoreFirst != expectedFirst {
			t.Errorf("first char score = %d, want %d (scoreMatch + boundary×mul)",
				scoreFirst, expectedFirst)
		}
	})

	t.Run("consecutive bonus inherits boundary bonus from run start", func(t *testing.T) {
		// "/aws" — 'a' gets delimiter boundary bonus, 'w' and 's' inherit it
		snippet := "~/.aws/config"
		matches := []int{3, 4, 5} // "aws" after '.'
		scoreWithBoundary := calculateScore("aws", snippet, matches)

		// "bcd" in "abcdef" — all mid-word, consecutive gets only bonusConsecutive
		snippet2 := "abcdef"
		matches2 := []int{1, 2, 3} // "bcd"
		scoreNoBoundary := calculateScore("bcd", snippet2, matches2)

		if scoreWithBoundary <= scoreNoBoundary {
			t.Errorf("consecutive after boundary (%d) should be higher than consecutive mid-word (%d)",
				scoreWithBoundary, scoreNoBoundary)
		}
	})

	t.Run("prefix match bonus", func(t *testing.T) {
		// Same matches at position 0, but one has a prefix match
		snippet := "docker run"
		matches := []int{0, 1, 2, 3} // "dock" — prefix match
		scorePrefix := calculateScore("dock", snippet, matches)

		// Same length match but not at prefix
		snippet2 := "sudo cker"
		matches2 := []int{5, 6, 7, 8} // "cker" — not a prefix
		scoreNoPrefix := calculateScore("cker", snippet2, matches2)

		if scorePrefix <= scoreNoPrefix {
			t.Errorf("prefix match score (%d) should be higher than non-prefix (%d)",
				scorePrefix, scoreNoPrefix)
		}
	})

	t.Run("similarity bonus applied when ratio exceeds threshold", func(t *testing.T) {
		// Short snippet where query covers most of it (ratio > 0.5)
		snippet := "git log"
		matches := []int{0, 1, 2, 3} // "git " — 4/7 = 0.57 > 0.5
		scoreHighRatio := calculateScore("git ", snippet, matches)

		// Long snippet where query covers little (ratio <= 0.5)
		snippet2 := "git log --oneline --graph"
		matches2 := []int{0, 1, 2, 3} // "git " — 4/25 = 0.16 < 0.5
		scoreLowRatio := calculateScore("git ", snippet2, matches2)

		if scoreHighRatio <= scoreLowRatio {
			t.Errorf("high similarity ratio score (%d) should be higher than low ratio (%d)",
				scoreHighRatio, scoreLowRatio)
		}
	})
}

func TestCharClassOf(t *testing.T) {
	tests := []struct {
		char byte
		want charClass
	}{
		{' ', charWhite},
		{'\t', charWhite},
		{'\n', charWhite},
		{'\r', charWhite},
		{'/', charDelimiter},
		{'.', charDelimiter},
		{',', charDelimiter},
		{':', charDelimiter},
		{';', charDelimiter},
		{'|', charDelimiter},
		{'-', charDelimiter},
		{'_', charDelimiter},
		{'a', charLower},
		{'z', charLower},
		{'A', charUpper},
		{'Z', charUpper},
		{'0', charNumber},
		{'9', charNumber},
		{'~', charNonWord},
		{'$', charNonWord},
		{'!', charNonWord},
		{'#', charNonWord},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			got := charClassOf(tt.char)
			if got != tt.want {
				t.Errorf("charClassOf(%q) = %d, want %d", tt.char, got, tt.want)
			}
		})
	}
}

func TestBonusFor(t *testing.T) {
	tests := []struct {
		name      string
		prevClass charClass
		currClass charClass
		want      int
	}{
		{"non-alnum current returns 0", charWhite, charNonWord, 0},
		{"white to lower", charWhite, charLower, bonusBoundaryWhite},
		{"white to upper", charWhite, charUpper, bonusBoundaryWhite},
		{"white to number", charWhite, charNumber, bonusBoundaryWhite},
		{"delimiter to lower", charDelimiter, charLower, bonusBoundaryDelim},
		{"delimiter to upper", charDelimiter, charUpper, bonusBoundaryDelim},
		{"nonword to lower", charNonWord, charLower, bonusBoundary},
		{"lower to upper (camelCase)", charLower, charUpper, bonusCamelCase},
		{"lower to lower (same type)", charLower, charLower, 0},
		{"upper to upper (same type)", charUpper, charUpper, 0},
		{"upper to lower", charUpper, charLower, 0},
		{"number to lower", charNumber, charLower, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bonusFor(tt.prevClass, tt.currClass)
			if got != tt.want {
				t.Errorf("bonusFor(%d, %d) = %d, want %d",
					tt.prevClass, tt.currClass, got, tt.want)
			}
		})
	}
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		snippet   string
		wantMatch bool
	}{
		{
			name:      "exact match",
			query:     "docker",
			snippet:   "docker",
			wantMatch: true,
		},
		{
			name:      "partial match",
			query:     "dock",
			snippet:   "docker run",
			wantMatch: true,
		},
		{
			name:      "case insensitive",
			query:     "Docker",
			snippet:   "docker run",
			wantMatch: true,
		},
		{
			name:      "no match",
			query:     "xyz",
			snippet:   "docker run",
			wantMatch: false,
		},
		{
			name:      "empty query",
			query:     "",
			snippet:   "docker",
			wantMatch: true,
		},
		{
			name:      "empty snippet",
			query:     "docker",
			snippet:   "",
			wantMatch: false,
		},
		{
			name:      "scattered characters match",
			query:     "dkr",
			snippet:   "docker",
			wantMatch: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, score := fuzzyMatch(tt.query, tt.snippet)
			if tt.wantMatch {
				if matches == nil {
					t.Errorf("expected match for query=%q snippet=%q, got nil", tt.query, tt.snippet)
				}
				if tt.query != "" && score <= 0 {
					t.Errorf("expected positive score for match, got %d", score)
				}
			} else {
				if matches != nil {
					t.Errorf("expected no match for query=%q snippet=%q, got matches=%v score=%d",
						tt.query, tt.snippet, matches, score)
				}
			}
		})
	}
}

func TestFuzzyMatchPositions(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		snippet       string
		wantPositions []int
	}{
		{
			name:          "exact match positions",
			query:         "abc",
			snippet:       "abc",
			wantPositions: []int{0, 1, 2},
		},
		{
			name:          "prefix match positions",
			query:         "doc",
			snippet:       "docker",
			wantPositions: []int{0, 1, 2},
		},
		{
			name:          "empty query returns empty positions",
			query:         "",
			snippet:       "docker",
			wantPositions: []int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches, _ := fuzzyMatch(tt.query, tt.snippet)
			if !slices.Equal(matches, tt.wantPositions) {
				t.Errorf("fuzzyMatch(%q, %q) positions = %v, want %v",
					tt.query, tt.snippet, matches, tt.wantPositions)
			}
		})
	}
}

func TestFuzzySearchEdgeCases(t *testing.T) {
	makeLinippets := func(snippets ...string) linippet.Linippets {
		l := make(linippet.Linippets, len(snippets))
		for i, s := range snippets {
			l[i] = linippet.Linippet{Id: s, Snippet: s}
		}
		return l
	}

	t.Run("empty query returns empty results", func(t *testing.T) {
		results := FuzzySearch(context.Background(), "", makeLinippets("docker run", "git log"))
		if len(results) != 0 {
			t.Errorf("expected 0 results for empty query, got %d", len(results))
		}
	})

	t.Run("whitespace-only query returns empty results", func(t *testing.T) {
		results := FuzzySearch(context.Background(), "   ", makeLinippets("docker run"))
		if len(results) != 0 {
			t.Errorf("expected 0 results for whitespace query, got %d", len(results))
		}
	})

	t.Run("cancelled context returns nil", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		results := FuzzySearch(ctx, "docker", makeLinippets("docker run", "docker ps"))
		if results != nil {
			t.Errorf("expected nil for cancelled context, got %v", results)
		}
	})

	t.Run("multi-word partial failure excludes snippet", func(t *testing.T) {
		// "docker xyz" — "docker" matches but "xyz" doesn't
		results := FuzzySearch(context.Background(), "docker xyz", makeLinippets("docker run"))
		if len(results) != 0 {
			t.Errorf("expected 0 results when one word doesn't match, got %d", len(results))
		}
	})

	t.Run("equal score uses snippet length as tiebreaker", func(t *testing.T) {
		// Both start with "git" at the same position, should tie on score.
		// Shorter snippet should rank first.
		short := "git log"
		long := "git log --oneline --graph"
		results := FuzzySearch(context.Background(), "git", makeLinippets(long, short))
		if len(results) < 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].Linippet.Snippet != short {
			t.Errorf("expected shorter snippet %q first, got %q (scores: %d vs %d)",
				short, results[0].Linippet.Snippet,
				results[0].Score, results[1].Score)
		}
	})
}
