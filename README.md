# sentiment

This package is part of Social Harvest, but has been made available as a stand alone package under a separate license. 

This package performs simple sentiment analysis using a Naive Bayes classifier. Training data is held in memory so no database is required. It heavily 
relies upon the ```github.com/jbrukh/bayesian``` package and is more or less providing you with a specific implementation of that package.    
Various freely available word lists and corpus data have been used to train the classifier, including:

* [An opinion lexicon by Minqing Hu and Bing Liu](http://www.cs.uic.edu/~liub/FBS/sentiment-analysis.html)    
* [AFINN: An affective lexicon by Finn Årup Nielsen](http://neuro.imm.dtu.dk/wiki/AFINN)    
* [A word list from SentiWordNet](http://sentiwordnet.isti.cnr.it/)    
* [Russian training data from Eismont Polina, Efanova Iuliia, Konovalova Svetlana, Losev Viktor and Velichko Alena of Saint Petersburg State University of Aerospace Instrumentation, Department of Applied Linguistics originally for the SentiStrength project](http://sentistrength.wlv.ac.uk)    
* [A sentiment data set based on Twitter from Ibrahim Naji](http://thinknook.com/twitter-sentiment-analysis-training-corpus-dataset-2012-09-22)    
* [Movie reviews polarity (v2) by Pang/Lee ACL 2004](http://www.cs.cornell.edu/People/pabo/movie-review-data)    

Training data from previous Social Harvest sentiment analysis (discussed here http://www.slideshare.net/shift8/mongodb-machine-learning) and custom word lists for emoticons have also been used to train this classifier. 

The corpus data used for training has not been included with this package for multiple reasons. First, it would unnecessarily add to the size of the repository. Second, some of the data 
can not be distributed as per Twitter's ToS. Additional training can be performed so it is possible to re-train and even add to the training of this classifier. However, the training 
shipped with this package is privately maintained for these reasons. The training shipped with this package was targeted at detecting sentiment in short messages (social media) and 
should work well for many cases without the need for additional training.

## Implementation Notes

Keep in mind that this classifier is in memory and there has been a good deal of training, so the memory allocation is going to be considerable. Data gets loaded from a file on disk 
upon ```NewAnalyzer()``` being called. It is recommended to re-use this analyzer so that unecessary amounts of memory are not allocated and application performance isn't hindered.

There is an up-front cost to initilize an analyzer (memory allocation). So it is best to setup an analyzer early on in your application. This way the delay is at startup rather 
than at some random point in your application which could create delays elsewhere and impact other things.

The dump file is too large to store in GitHub so it is downloaded from the internet. The training data is about 125MB so keep in mind that you will need enough RAM to cover this and 
anything else your application and system needs.

## To Train Additional Data

To add training data on top of the existing set, you can access the underlying ```Classifier``` from the ```github.com/jbrukh/bayesian``` package. For example:    
```
a := NewAnalyzer()    
a.classifier.Learn([]string{"some","tokenized","string"}, sentiment.Positive)    
a.classifier.Learn([]string{"some","tokenized","string"}, sentiment.Negative)    
a.classifier.Learn([]string{"some","tokenized","string"}, sentiment.Neutral)
```

You have complete access to all methods from this classifier package. Again, this sentiment package heavily relies upon this other bayesian package and mostly provides training and 
convenience (by tokenizing input). There aren't a lot of additional things going on here.