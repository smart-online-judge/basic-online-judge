from sklearn.feature_extraction.text import TfidfVectorizer
from string import punctuation as string_punctuation
import nltk

remove_punct_dict = dict((ord(punct), None) for punct in string_punctuation)

def get_normalize_filter(type=0):
    normalize_filter = lambda text: nltk.word_tokenize(
        text.lower().translate(remove_punct_dict))

    if type == 0: # Default
        return normalize_filter

    elif type == 1: # Lemming
        lemmer = nltk.stem.WordNetLemmatizer()
        lem_tokens = lambda tokens: [lemmer.lemmatize(token) for token in tokens]
        return lambda text: lem_tokens(normalize_filter(text))

    else: # Stemming
        stemmer = nltk.stem.porter.PorterStemmer()
        stem_tokens = lambda tokens: [stemmer.stem(token) for token in tokens]
        return lambda text: stem_tokens(normalize_filter(text))

def cosine_similarity(file_objects, normalize_filter):
    lem_vectorizer = TfidfVectorizer(tokenizer=normalize_filter, stop_words='english')

    # TODO: handle empty files (ValueError: empty vocabulary;
    # perhaps the documents only contain stop words)

    mat_sparse = lem_vectorizer.fit_transform(file_objects)
    return (mat_sparse * mat_sparse.T).A

def get_syntactic_similarity_mat(file_paths):
    try:
        file_objects = (open(fp, 'rt', encoding='utf-8').read() for fp in file_paths)
        normalize_filter = get_normalize_filter()

        return cosine_similarity(file_objects, normalize_filter)
    except UnicodeDecodeError as ex:
        raise UnicodeDecodeError("Unable to decode given files") from ex
    except ValueError as ex:
        raise ValueError("Some files don't contain any valid words") from ex
    except Exception as ex:
        raise Exception("Exception encountered when computing syntactic similarity") from ex