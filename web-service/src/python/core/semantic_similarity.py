import gensim.downloader as api
from gensim.models import WordEmbeddingSimilarityIndex, TfidfModel
from gensim.similarities import SoftCosineSimilarity, SparseTermSimilarityMatrix
from gensim.corpora import TextDirectoryCorpus

def glove_similarity(input_path):
    # Load the model: downloading happens once, opening on each call
    glove = api.load("glove-wiki-gigaword-50")
    termsim_index = WordEmbeddingSimilarityIndex(glove)

    corpus = TextDirectoryCorpus(input_path, max_depth=0)
    dictionary = corpus.dictionary
    bow_corpus = list(corpus)

    tfidf = TfidfModel(dictionary=dictionary)
    similarity_matrix = SparseTermSimilarityMatrix(termsim_index, dictionary, tfidf)
    docsim_index = SoftCosineSimilarity(bow_corpus, similarity_matrix)

    batch_of_documents = [tfidf[x] for x in bow_corpus]
    return docsim_index[batch_of_documents]

def get_semantic_similarity_mat(input_path):
    try:
        return glove_similarity(input_path)
    except Exception as ex:
        raise Exception("Exception encountered when computing semantic similarity") from ex