import sys, string, nltk
from os.path import isdir, join, isfile, split
from os import listdir, sep
from datetime import datetime
from sklearn.feature_extraction.text import TfidfVectorizer
import matplotlib.pyplot as plt, numpy as np

def eprint(*args, **kwargs):
    print('ERR %s: ' % datetime.now().strftime("%d/%m/%Y %H:%M:%S"), end ='')
    print(*args, file=sys.stderr, **kwargs)
    sys.exit()

def coisine_similarity(file_objects, tokenizer_option):
    remove_punct_dict = dict((ord(punct), None) for punct in string.punctuation)

    if tokenizer_option == 0: # None
        normalize_filter = lambda text: nltk.word_tokenize(text.lower().translate(remove_punct_dict))
    elif tokenizer_option == 1: # Lemming
        lemmer = nltk.stem.WordNetLemmatizer()
        LemTokens = lambda tokens: [lemmer.lemmatize(token) for token in tokens]
        normalize_filter = lambda text: LemTokens(nltk.word_tokenize(text.lower().translate(remove_punct_dict)))
    else: # Stemming
        stemmer = nltk.stem.porter.PorterStemmer()
        StemTokens = lambda tokens: [stemmer.stem(token) for token in tokens]
        normalize_filter = lambda text: StemTokens(nltk.word_tokenize(text.lower().translate(remove_punct_dict)))

    LemVectorizer = TfidfVectorizer(tokenizer=normalize_filter, stop_words='english')
    mat_sparse = LemVectorizer.fit_transform(file_objects)
    mat = (mat_sparse * mat_sparse.T).A
    return mat

def heatmap(x_labels, y_labels, values):
    fig, ax = plt.subplots()
    im = ax.imshow(values)

    ax.set_xticks(np.arange(len(x_labels)))
    ax.set_yticks(np.arange(len(y_labels)))

    ax.set_xticklabels(x_labels)
    ax.set_yticklabels(y_labels)

    # Rotate the tick labels and set their alignment.
    plt.setp(ax.get_xticklabels(), rotation=45, ha="right", fontsize=10,
         rotation_mode="anchor")

    for i in range(len(y_labels)):
        for j in range(len(x_labels)):
            text = ax.text(j, i, "%.2f"%values[i, j],
                           ha="center", va="center", color="w", fontsize=6)

    fig.tight_layout()
    plt.show()

if __name__ == '__main__':
    input_folder = sys.argv[-1]
    if not input_folder.startswith(sep):
        input_folder = join(split(sys.argv[0])[0], input_folder)

    if not isdir(input_folder):
        eprint("Given path '%s' isn't a directory", input_folder)

    file_names = [f for f in listdir(input_folder) if isfile(join(input_folder, f))]

    if len(file_names) > 20:
        eprint("Maximum amount of input files exceeded - allowed 20, received %d", len(file_names))

    option = [x for x in sys.argv if x in ('-l, -s, --lemming, --stemming')]
    
    if len(option) > 1:
        eprint("Icorrect option passed, received %s", option)

    if len(option) == 0:
        tokenizer_option = 0 # None
    elif option[0] in ('-l', '--lemming'):
        tokenizer_option = 1 # lemming
    else:
        tokenizer_option = 2 # stemming

    file_objects = (open(join(input_folder, f)).read() for f in file_names)
    mat = coisine_similarity(file_objects, tokenizer_option)

    heatmap(file_names, file_names, mat)






