# Run this script in the service-related container before the first start.
# Uploads nltk data, gensim model for Python algorithms to the ~/home directory.

import nltk
nltk.download('punkt')

import gensim.downloader as api
glove = api.load("glove-wiki-gigaword-50")