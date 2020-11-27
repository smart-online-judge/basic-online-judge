import difflib, sys, codecs
from datetime import datetime
from os import sep
from os.path import join, isfile, split

def eprint(*args, **kwargs):
    print('ERR %s: ' % datetime.now().strftime("%d/%m/%Y %H:%M:%S"), end='', file=sys.stderr)
    print(*args, file=sys.stderr, **kwargs)
    sys.exit()

def difference_comparison(input_files, difference_option):
    #  a unified diff only includes modified lines and a bit of context.
    compare_func = difflib.unified_diff
    if difference_option == 1:
        compare_func = difflib.ndiff # avoids noise
    elif difference_option == 2:
        compare_func = difflib.Differ().compare

    list1, list2 = [], []
    idx = None
    for idx, lines in enumerate(zip(input_files[0], input_files[1])):
        if not idx+1 % 100 == 0:
            list1.append(lines[0])
            list2.append(lines[1])
        else:
            diff = compare_func(list1, list2)
            list1, list2 = [], []
            print(''.join(list(diff)))

    if idx and idx+1 % 100 != 0:
        diff = compare_func(list1, list2)
        print(''.join(list(diff)))

if __name__ == '__main__':
    if len(sys.argv) < 3 or len(sys.argv) > 4:
        eprint('Wrong count of input arguments(%d), 3-4 expected' % len(sys.argv))

    INPUT_FILES = sys.argv[-2:]
    RUNNING_DIR = split(sys.argv[0])[0]
    for idx, f in enumerate(INPUT_FILES):
        if not f.startswith(sep):
            INPUT_FILES[idx] = join(RUNNING_DIR, f)
        if not isfile(INPUT_FILES[idx]):
            eprint("Given path '%s' isn't a file" % INPUT_FILES[idx])

    ARGV = sys.argv[:-2]
    DIFFERENCE_OPTION = 3
    if '--ndiff' in ARGV:
        DIFFERENCE_OPTION = 1
    elif '--compare' in ARGV:
        DIFFERENCE_OPTION = 2
    elif '--unified' in ARGV:
        DIFFERENCE_OPTION = 3

    INPUT_FILES = [codecs.open(f, encoding='utf-8') for f in INPUT_FILES]
    difference_comparison(INPUT_FILES, DIFFERENCE_OPTION)
