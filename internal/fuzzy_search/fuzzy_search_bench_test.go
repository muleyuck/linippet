package fuzzy_search

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/muleyuck/linippet/internal/linippet"
)

var benchCommandTemplates = []string{
	"docker run -d -p ${{port}}:${{container_port}} ${{image}}",
	"kubectl get pods -n ${{namespace}} -l app=${{label}}",
	"git log --oneline --graph --all -n ${{count}}",
	"curl -X POST -H 'Content-Type: application/json' -d '${{body}}' ${{url}}",
	"ssh -i ${{key}} ${{user}}@${{host}}",
	"tar -czf ${{archive}}.tar.gz ${{directory}}",
	"find ${{path}} -name '${{pattern}}' -type f",
	"aws s3 cp ${{source}} s3://${{bucket}}/${{key}} --recursive",
	"psql -h ${{host}} -U ${{user}} -d ${{database}} -c '${{query}}'",
	"python3 -m venv ${{venv_name}} && source ${{venv_name}}/bin/activate",
	"docker-compose -f ${{file}} up -d --build ${{service}}",
	"rsync -avz --progress ${{source}} ${{destination}}",
	"jq '${{filter}}' ${{file}} | head -n ${{lines}}",
	"grep -rn '${{pattern}}' ${{directory}} --include='${{glob}}'",
	"systemctl ${{action}} ${{service}}",
	"openssl req -x509 -newkey rsa:${{bits}} -keyout ${{key}} -out ${{cert}} -days ${{days}}",
	"npm run ${{script}} -- --env=${{environment}}",
	"go test -bench=${{pattern}} -benchmem -count=${{count}} ./${{package}}/...",
	"redis-cli -h ${{host}} -p ${{port}} ${{command}}",
	"ffmpeg -i ${{input}} -vf scale=${{width}}:${{height}} ${{output}}",
}

func generateBenchLinippets(n int) linippet.Linippets {
	linippets := make(linippet.Linippets, n)
	for i := range n {
		linippets[i] = linippet.Linippet{
			Id:      fmt.Sprintf("bench-id-%06d", i),
			Snippet: benchCommandTemplates[i%len(benchCommandTemplates)],
		}
	}
	return linippets
}

func BenchmarkFuzzySearch1000_MultiWord(b *testing.B) {
	linippets := generateBenchLinippets(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		FuzzySearch(context.Background(),"docker run port", linippets)
	}
}

func BenchmarkFuzzySearch1000_NoMatch(b *testing.B) {
	linippets := generateBenchLinippets(1000)

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		FuzzySearch(context.Background(),"zzzzxxx", linippets)
	}
}

func BenchmarkFuzzyMatch_SinglePair(b *testing.B) {
	snippet := "docker run -d -p ${{port}}:${{container_port}} ${{image}}"

	b.ReportAllocs()
	b.ResetTimer()
	for b.Loop() {
		fuzzyMatch("docker", snippet)
	}
}

func BenchmarkFuzzySearchScaling(b *testing.B) {
	for _, n := range []int{100, 500, 1000, 5000} {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			linippets := generateBenchLinippets(n)

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				FuzzySearch(context.Background(),"docker", linippets)
			}
		})
	}
}

func BenchmarkFuzzyMatchQueryLengthScaling(b *testing.B) {
	snippet := "docker-compose -f ${{file}} up -d --build ${{service}}"
	queries := []string{
		"d",
		"dock",
		"docker-co",
		"docker-compose up",
	}

	for _, q := range queries {
		b.Run(fmt.Sprintf("qlen=%d", len(q)), func(b *testing.B) {
			// For multi-word queries, test with first word only for fuzzyMatch
			query := strings.Fields(q)[0]

			b.ReportAllocs()
			b.ResetTimer()
			for b.Loop() {
				fuzzyMatch(query, snippet)
			}
		})
	}
}
