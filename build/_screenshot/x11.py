# Synthetic input for the screenshot window (see _screenshot.sh), on $DISPLAY:
#   x11.py find <title part>   - prints the id of the window whose title has it
#   x11.py key <window> <key>  - presses and releases the key (X keysym name)
#   x11.py click <window> x y  - left click at the point of the window
import ctypes
import sys

x = ctypes.cdll.LoadLibrary("libX11.so.6")
x.XOpenDisplay.restype = ctypes.c_void_p
x.XDefaultRootWindow.restype = ctypes.c_ulong
x.XStringToKeysym.restype = ctypes.c_ulong
d = ctypes.c_void_p(x.XOpenDisplay(None))
if not d:
    sys.exit("cannot open display")
root = x.XDefaultRootWindow(d)


class InputEvent(ctypes.Structure):
    # XKeyEvent and XButtonEvent share the layout; the last field is keycode or button
    _fields_ = [("type", ctypes.c_int), ("serial", ctypes.c_ulong), ("send_event", ctypes.c_int),
                ("display", ctypes.c_void_p), ("window", ctypes.c_ulong), ("root", ctypes.c_ulong),
                ("subwindow", ctypes.c_ulong), ("time", ctypes.c_ulong), ("x", ctypes.c_int),
                ("y", ctypes.c_int), ("x_root", ctypes.c_int), ("y_root", ctypes.c_int),
                ("state", ctypes.c_uint), ("code", ctypes.c_uint), ("same_screen", ctypes.c_int),
                ("pad", ctypes.c_byte * 64)]


def send(w, events):
    for t, mask, state, code, px, py in events:
        e = InputEvent(type=t, display=d, window=w, root=root, x=px, y=py,
                       state=state, code=code, same_screen=1)
        x.XSendEvent(d, ctypes.c_ulong(w), 1, mask, ctypes.byref(e))
    x.XFlush(d)


def find(part):
    stack = [root]
    while stack:
        w = stack.pop()
        name = ctypes.c_char_p()
        if x.XFetchName(d, ctypes.c_ulong(w), ctypes.byref(name)) and name.value and part in name.value.decode(errors="replace"):
            return w
        r, p, n = ctypes.c_ulong(), ctypes.c_ulong(), ctypes.c_uint()
        ch = ctypes.POINTER(ctypes.c_ulong)()
        if x.XQueryTree(d, ctypes.c_ulong(w), ctypes.byref(r), ctypes.byref(p), ctypes.byref(ch), ctypes.byref(n)):
            stack += [ch[i] for i in range(n.value)]
    return 0


cmd = sys.argv[1]
if cmd == "find":
    w = find(sys.argv[2])
    if not w:
        sys.exit("window not found")
    print(w)
elif cmd == "key":
    w = int(sys.argv[2])
    code = x.XKeysymToKeycode(d, x.XStringToKeysym(sys.argv[3].encode()))
    send(w, [(2, 1, 0, code, 10, 10), (3, 2, 0, code, 10, 10)])  # KeyPress, KeyRelease
elif cmd == "click":
    w, px, py = int(sys.argv[2]), int(sys.argv[3]), int(sys.argv[4])
    # MotionNotify, ButtonPress, ButtonRelease (with Button1 held)
    send(w, [(6, 64, 0, 0, px, py), (4, 4, 0, 1, px, py), (5, 8, 256, 1, px, py)])
else:
    sys.exit("unknown command " + cmd)
