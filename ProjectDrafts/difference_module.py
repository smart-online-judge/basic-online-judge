import difflib, sys, codecs
from os import sep
from os.path import isdir, join, isfile, split

def eprint(*args, **kwargs):
    print('ERR %s: ' % datetime.now().strftime("%d/%m/%Y %H:%M:%S"), end ='')
    print(*args, file=sys.stderr, **kwargs)
    sys.exit()

def difference_comparison(input_files, difference_option):
    # https://pymotw.com/2/difflib/

    compare_func = difflib.unified_diff #  a unified diff only includes modified lines and a bit of context.
    if difference_option == 1:
        compare_func = difflib.ndiff # avoids noise
    elif difference_option == 2:
        compare_func = difflib.Differ().compare

    l1, l2 = [], []
    for n, lines in enumerate(zip(input_files[0],input_files[1])):
        if not (n+1 % 100 == 0):
            l1.append(lines[0])
            l2.append(lines[1])
        else:
            diff = compare_func(l1, l2)
            l1, l2 = '', ''
            print(''.join(list(diff)))

    if n+1 % 100 != 0:
        diff = compare_func(l1, l2)
        print(''.join(list(diff)))

if __name__ == '__main__':
    if len(sys.argv) < 3 or len(sys.argv) > 4:
        eprint('Wrong amount of input arguments(%d), 3-4 expected', len(sys.argv))

    input_files = sys.argv[-2:]
    running_dir = split(sys.argv[0])[0]
    for idx, f in enumerate(input_files):
        if not f.startswith(sep):
            input_files[idx] = join(running_dir, f)
        if not isfile(input_files[idx]):
            eprint("Given path '%s' isn't a file", input_files[idx])

    
    argv = sys.argv[:-2]
    difference_option = 3
    if '--ndiff' in argv:
        difference_option = 1
    elif '--compare' in argv:
        difference_option = 2
    elif '--unified' in argv:
        difference_option = 3

    input_files = [codecs.open(f, encoding='cp1252') for f in input_files]
    difference_comparison(input_files, difference_option)