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
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	Positive bayesian.Class = "Positive"
	Negative bayesian.Class = "Negative"
	Neutral  bayesian.Class = "Neutral"
)

const DATA_FILE = "./sentiment-data/sentiment-classifier.dmp"

type Analyzer struct {
	classifier *bayesian.Classifier
}

// Classifies a string
func (a *Analyzer) Classify(s string) int {
	if len(s) <= 2 {
		return 0
	}
	tokens := tokenize(s)

	_, likely, _ := a.classifier.LogScores(tokens)

	sentiment := 0
	// Positive, Negative, Neutral was the order in which classes were defined
	switch likely {
	case 0:
		sentiment = 1
	case 1:
		sentiment = -1
	case 2:
		sentiment = 0
	}

	return sentiment
}

// Sets up and trains a new analyzer to classify sentiment
func NewAnalyzer() Analyzer {
	a := Analyzer{}

	// Get the training data if not present
	_, err := os.Stat(DATA_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			a.downloadDataSet()
		}
	}

	c, err := bayesian.NewClassifierFromFile(DATA_FILE)
	if err == nil {
		a.classifier = c
	} else {
		// Note: Nothing will be trained at this point, but we'll still have a classifier that can be trained
		a.classifier = bayesian.NewClassifier(Positive, Negative, Neutral)
	}

	return a
}

// Retrieves training data (which is much too large to keep in GitHub)
func (a *Analyzer) downloadDataSet() {
	os.Mkdir("./sentiment-data", 0777)
	out, oErr := os.Create(DATA_FILE)
	defer out.Close()
	if oErr == nil {
		r, rErr := http.Get("https://s3.amazonaws.com/socialharvest/public-data/sentiment/sentiment-classifier.dmp")
		defer r.Body.Close()
		if rErr == nil {
			_, nErr := io.Copy(out, r.Body)
			if nErr != nil {
				err := os.Remove(DATA_FILE)
				if err != nil {
					log.Println(err)
				}
			}
			r.Body.Close()
		} else {
			log.Println(rErr)
		}
		out.Close()
	} else {
		log.Println(oErr)
	}
}

// Splits apart a string to train/classify it by word and word n-grams.
func tokenize(s string) []string {
	tokens := []string{}

	tokenSlice := strings.Split(s, " ")
	for k, v := range tokenSlice {
		tokens = append(tokens, v)
		if len(tokenSlice)-1 > k && len(v) > 1 {
			ngram := tokenSlice[k] + " " + tokenSlice[k+1] //+ " " + tokenSlice[k+2]
			tokens = append(tokens, ngram)
		}
	}

	return tokens
}
