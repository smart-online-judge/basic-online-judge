import warnings
warnings.filterwarnings("ignore")

import sys, numpy as np
from os import listdir as os_listdir, sep as os_sep, getcwd as os_getcwd
from os.path import isdir, join, isfile, getsize
from json import dumps as json_dumps

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, end='', **kwargs)
    sys.exit(1)

MAX_ALLOWED_FILES_SIZE = (1 << 20) * 10 # 10Mb

def collect_verify_files(input_folder):
    if not isdir(input_folder):
        eprint("Given path '%s' isn't a directory" % INPUT_FOLDER)

    dir_size = 0
    file_names, file_paths = [], []
    for f in os_listdir(input_folder):
        fp = join(input_folder, f)
        if isfile(fp):
            file_names.append(f)
            file_paths.append(fp)
            dir_size += getsize(fp)

    if dir_size > MAX_ALLOWED_FILES_SIZE:
        eprint("Total files size({}) exceeds allowed limit({})".format(
            dir_size, MAX_ALLOWED_FILES_SIZE))

    if len(file_names) < 2:
        eprint("Files amount is lower than 2 - nothing to compare")

    return file_names, file_paths

from core import get_syntactic_similarity_mat
from core import get_semantic_similarity_mat
from core import get_averaged_similarity_mat

if __name__ == '__main__':
    INPUT_FOLDER = sys.argv[-1]
    if not INPUT_FOLDER.startswith(os_sep):
        INPUT_FOLDER = join(os_getcwd(), INPUT_FOLDER)

    FILE_NAMES, FILE_PATHS = collect_verify_files(INPUT_FOLDER)

    try:
        syn_mat = get_syntactic_similarity_mat(FILE_PATHS)
        sem_mat = get_semantic_similarity_mat(INPUT_FOLDER)
        res_mat = get_averaged_similarity_mat(syn_mat, sem_mat)
        print(json_dumps(res_mat), end='')
    except Exception as ex:
        eprint(ex)