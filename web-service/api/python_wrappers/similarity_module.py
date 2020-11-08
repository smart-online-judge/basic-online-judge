import sys, string, nltk, codecs
from os import listdir, sep
from os.path import isdir, join, isfile, split
from datetime import datetime
from sklearn.feature_extraction.text import TfidfVectorizer
import matplotlib.pyplot as plt, numpy as np

def eprint(*args, **kwargs):
    print('ERR %s: ' % datetime.now().strftime("%d/%m/%Y %H:%M:%S"), end='', file=sys.stderr)
    print(*args, file=sys.stderr, **kwargs)
    sys.exit()

def cosine_similarity(file_objects, tokenizer_option):
    remove_punct_dict = dict((ord(punct), None) for punct in string.punctuation)

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

if __name__ == '__main__':
    INPUT_FOLDER = sys.argv[-1]
    if not INPUT_FOLDER.startswith(sep):
        INPUT_FOLDER = join(split(sys.argv[0])[0], INPUT_FOLDER)

    if not isdir(INPUT_FOLDER):
        eprint("Given path '%s' isn't a directory" % INPUT_FOLDER)

    FILE_NAMES = [f for f in listdir(INPUT_FOLDER) if isfile(join(INPUT_FOLDER, f))]
    if len(FILE_NAMES) > 20:
        eprint("Maximum amount of input files exceeded - allowed 20, received %d" % len(FILE_NAMES))

    OPTIONS_LIST = [x for x in sys.argv if x in ('-l, -s, --lemming, --stemming')]
    if len(OPTIONS_LIST) > 1:
        eprint("Incorrect working option passed, received %s" % OPTIONS_LIST)

    if len(OPTIONS_LIST) == 0:
        TOKENIZER_OPTION = 0 # None
    elif OPTIONS_LIST[0] in ('-l', '--lemming'):
        TOKENIZER_OPTION = 1 # lemming
    else:
        TOKENIZER_OPTION = 2 # stemming
        
    FILE_OBJECTS = (codecs.open(join(INPUT_FOLDER, f), 'rb', encoding = 'utf-8').read() for f in FILE_NAMES)
    RES_MAT = cosine_similarity(FILE_OBJECTS, TOKENIZER_OPTION)

    if '--external' in sys.argv:
        if len(RES_MAT) == 0:
            print('0')
            sys.exit()

        output_arr = [str(len(RES_MAT[0]))]
        for arr in RES_MAT:
            for val in arr:
                output_arr.append("%.2f" % val)
        print(','.join(output_arr))
    else:
        heatmap(FILE_NAMES, FILE_NAMES, RES_MAT)