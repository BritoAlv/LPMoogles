package main

import "math"

type ModelTfIdf struct {
	// create a map that will store the idf values for each word.
	idf            map[string]int
	keys           []string
	tf             map[string]map[string]int
	totalDocuments int
}

func ConstructormodelTfIdf(arr []ResultFromDto) *ModelTfIdf {
	idf := make(map[string]int)
	td := make(map[string]map[string]int, len(arr))
	numberOfDocuments := len(arr)
	for _, doc := range arr {
		tfmap := make(map[string]int)
		wordsInDoc := splitInWords(doc.Text())
		for _, word := range wordsInDoc {
			// if the word is not in the idf map, add it.
			if _, ok := idf[word]; !ok {
				idf[word] = 1
			}
			if _, ok := tfmap[word]; !ok {
				tfmap[word] = 0
			}
			tfmap[word]++
		}
		td[doc.Name()] = tfmap
	}
	keys := Keys[string](idf)
	return &ModelTfIdf{
		idf, keys, td, numberOfDocuments}
}

func tfIdfQuery(query string, model ModelTfIdf) []float64 {
	tfmap := make(map[string]int)
	wordsInDoc := splitInWords(query)
	for _, word := range wordsInDoc {
		if _, ok := tfmap[word]; !ok {
			tfmap[word] = 0
		}
		tfmap[word]++
	}
	result := make([]float64, len(model.keys))
	for index, word := range model.keys {
		result[index] = 0
		if v, ok := tfmap[word]; ok {
			result[index] += float64(v) * math.Log(float64(model.totalDocuments)/float64(model.idf[word]+1))
		}
	}
	return result
}

func cos_sim(a []float64, b []float64) float64 {
	// calculate the cosine similarity between a and b.
	var dotProduct = float64(0)
	var magA = float64(0)
	var magB = float64(0)
	for i := 0; i < len(a); i++ {
		dotProduct += a[i] * b[i]
		magA += a[i] * a[i]
		magB += b[i] * b[i]
	}
	if magA == 0 || magB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(magA) * math.Sqrt(magB))
}

func tf_idf_doc(index string,  model  ModelTfIdf) []float64 {
	result := make([]float64, len(model.keys))
	for ind, word := range model.keys {
		result[ind] = 0
		if v, ok := model.tf[index][word]; ok {
			result[ind] += float64(v) * math.Log(float64(model.totalDocuments)/float64(model.idf[word]+1))
		}
	}
	return result
}
