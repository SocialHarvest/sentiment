// Copyright 2015 Tom Maiaroto, Shift8Creative

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sentiment

import (
	"github.com/jbrukh/bayesian"
	. "gopkg.in/check.v1"
	"testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type SentimentSuite struct {
	testPhrases []map[string]interface{}
	a           Analyzer
}

var _ = Suite(&SentimentSuite{})

func (s *SentimentSuite) SetUpSuite(c *C) {
	s.a = NewAnalyzer()

	// 0 = Neutral
	// 1 = Positive
	// -1 = Negative
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "this is pretty cool", "expected": 1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "♥", "expected": 1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "that's so horrible!", "expected": -1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "Off to Belize today for vacation!  WooHoo!", "expected": 1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "just ate a peanut that tasted like dish soap. that can't be a good thing.", "expected": -1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "Just at work bored", "expected": -1})
	s.testPhrases = append(s.testPhrases, map[string]interface{}{"phrase": "the weather outside is meh", "expected": -1})
}

func (s *SentimentSuite) TestNewAnalyzer(c *C) {
	analyzer := NewAnalyzer()
	learned := analyzer.classifier.Learned()

	c.Assert(analyzer.classifier, FitsTypeOf, bayesian.NewClassifier(Positive, Negative, Neutral))
	c.Assert(learned, Not(Equals), 0)
}

func (s *SentimentSuite) TestClassify(c *C) {
	for _, v := range s.testPhrases {
		r := s.a.Classify(v["phrase"].(string))
		c.Assert(r, Equals, v["expected"].(int), Commentf("for test phrase: %s\n(0 = Neutral, 1 = Positive, -1 = Negative)", v["phrase"]))
	}

	r := s.a.Classify("")
	c.Assert(r, Equals, 0)
}

func (s *SentimentSuite) BenchmarkNewAnalyzer(c *C) {
	for n := 0; n < c.N; n++ {
		_ = NewAnalyzer()
	}
}

func (s *SentimentSuite) BenchmarkClassify(c *C) {
	for n := 0; n < c.N; n++ {
		s.a.Classify("I deserve good things. I am entitled to my share of happiness. I refuse to beat myself up. I am attractive person. I am fun to be with.")
	}
}
