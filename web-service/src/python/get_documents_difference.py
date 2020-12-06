import diff_match_patch as dmp_module
import sys
from json import dumps as json_dumps

def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, end='', **kwargs)
    sys.exit(1)

def get_difference(tx1, txt2, args):
    # args - option={semantic|raw|efficiency}, html={true|false}, editcost=int, timeout=int
    option, html, editcost, timeout = args

    dmp = dmp_module.diff_match_patch()

    dmp.Diff_EditCost = int(editcost) if editcost else 8
    dmp.Diff_Timeout = int(timeout) if timeout else 2

    diff = dmp.diff_main(txt1, txt2)

    if option == "efficiency":
        dmp.diff_cleanupEfficiency(diff)
    elif option == "raw":
        pass
    else: # semantic
        dmp.diff_cleanupSemantic(diff)

    if html == "false":
        return json_dumps(diff)
    else:
        return dmp.diff_prettyHtml(diff)

if __name__ == '__main__':
    if len(sys.argv) != 7:
        eprint('Wrong number of input arguments(%d), 7 expected' % len(sys.argv))

    try:
        len1, len2 = int(sys.argv[1]), int(sys.argv[2])
        buf = sys.stdin.buffer

        txt1 = buf.read(len1).decode("utf-8")
        txt2 = buf.read(len2).decode("utf-8")

        print(get_difference(txt1, txt2, sys.argv[3:]), end='')
    except Exception as ex:
        eprint(ex)
