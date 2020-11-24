import sys, nltk
from os import listdir as os_listdir, sep as os_sep
from string import punctuation as string_punctuation
from os.path import isdir, join, isfile, split, getsize
from json import dumps as json_dumps
from sklearn.feature_extraction.text import TfidfVectorizer
import matplotlib.pyplot as plt, numpy as np

MAX_ALLOWED_FILES_SIZE = (1 << 20) * 10 # 10Mb

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)
    sys.exit(1)

def cosine_similarity(file_objects, tokenizer_option):
    remove_punct_dict = dict((ord(punct), None) for punct in string_punctuation)

    if tokenizer_option == 0: # None
        normalize_filter = lambda text: nltk.word_tokenize(
            text.lower().translate(remove_punct_dict))

    elif tokenizer_option == 1: # Lemming
        lemmer = nltk.stem.WordNetLemmatizer()
        lem_tokens = lambda tokens: [lemmer.lemmatize(token) for token in tokens]
        normalize_filter = lambda text: lem_tokens(
            nltk.word_tokenize(text.lower().translate(remove_punct_dict)))

    else: # Stemming
        stemmer = nltk.stem.porter.PorterStemmer()
        stem_tokens = lambda tokens: [stemmer.stem(token) for token in tokens]
        normalize_filter = lambda text: stem_tokens(
            nltk.word_tokenize(text.lower().translate(remove_punct_dict)))

    lem_vectorizer = TfidfVectorizer(tokenizer=normalize_filter, stop_words='english')

    # TODO: handle empty files (ValueError: empty vocabulary; perhaps the documents only contain stop words)
    mat_sparse = lem_vectorizer.fit_transform(file_objects)
    return (mat_sparse * mat_sparse.T).A

def heatmap(x_labels, y_labels, values):
    fig, axes = plt.subplots()
    _ = axes.imshow(values)

    axes.set_xticks(np.arange(len(x_labels)))
    axes.set_yticks(np.arange(len(y_labels)))

    axes.set_xticklabels(x_labels)
    axes.set_yticklabels(y_labels)

    # Rotate the tick labels and set their alignment.
    plt.setp(axes.get_xticklabels(), rotation=45, ha="right", fontsize=10,
             rotation_mode="anchor")

    for i in range(len(y_labels)):
        for j in range(len(x_labels)):
            axes.text(j, i, "%.2f"%values[i, j],
                      ha="center", va="center", color="w", fontsize=6)

    fig.tight_layout()
    plt.show()

def collect_files(input_folder):
    if not isdir(input_folder):
        eprint("Given path '%s' isn't a directory" % INPUT_FOLDER)

    dir_size = 0
    file_names = []
    file_paths = []
    for f in os_listdir(input_folder):
        fp = join(input_folder, f)
        if isfile(fp):
            file_names.append(f)
            file_paths.append(fp)
            dir_size += getsize(fp)

    if dir_size > MAX_ALLOWED_FILES_SIZE:
        eprint("Total files size({}) exceeds allowed limit({})".format(dir_size, MAX_ALLOWED_FILES_SIZE))
    if len(file_names) < 2:
        eprint("Files amount is lower than 2 - nothing to compare")

    return file_names, file_paths

def produce_similarity_mat(FILE_OBJECTS, TOKENIZER_OPTION):
    try:
        res_mat = cosine_similarity(FILE_OBJECTS, TOKENIZER_OPTION).tolist()
        res_mat = [[round(x, 2) for x in arr] for arr in res_mat]
        print(json_dumps(res_mat))
    except UnicodeDecodeError as ex:
        eprint("Unable to decode given files: %s" % ex)
    except ValueError as ex:
        eprint("Some files don't contain any valid words: %s" % ex)
    except Exception as ex:
        eprint("Exception encountered: %s" % ex)

def visualize_similarity_mat(file_objects, file_names, tokenizer_option):
    try:
        res_mat = cosine_similarity(file_objects, tokenizer_option)
        heatmap(file_names, file_names, res_mat)
    except Exception as ex:
        eprint("Exception encountered: %s" % ex)

if __name__ == '__main__':
    INPUT_FOLDER = sys.argv[-1]
    if not INPUT_FOLDER.startswith(os_sep):
        INPUT_FOLDER = join(split(sys.argv[0])[0], INPUT_FOLDER)

    FILE_NAMES, FILE_PATHS = collect_files(INPUT_FOLDER)

    OPTIONS_LIST = [x for x in sys.argv if x in ('-l, -s, --lemming, --stemming')]
    if len(OPTIONS_LIST) > 1:
        eprint("Incorrect working option passed, received %s" % OPTIONS_LIST)

    if not OPTIONS_LIST:
        TOKENIZER_OPTION = 0 # None
    elif OPTIONS_LIST[0] in ('-l', '--lemming'):
        TOKENIZER_OPTION = 1 # lemming
    else:
        TOKENIZER_OPTION = 2 # stemming

    FILE_OBJECTS = (open(fp, encoding='utf-8').read() for fp in FILE_PATHS)
    if '--external' in sys.argv:
        produce_similarity_mat(FILE_OBJECTS, TOKENIZER_OPTION)
    else:
        visualize_similarity_mat(FILE_OBJECTS, FILE_NAMES, TOKENIZER_OPTION)
