# Scale Manager

This is the Windows scale-manager application that is used to define new scales and upload them to
presets on the harp adapter.

## Why Go?

I set the following requirements:

* runs on 32bit Win7 (Iasos's machines are all Win7)
* zero or very few external dependencies (minimize chance of DLL Hell on Iasos's machines)
   * so ideally something that can be compiled to a statically linked .EXE
   * that rules out webview, python, java, .NET - even tcl/tk

Finding tools that support Win7 in 2020 is a challenge.  I originally thought I could just create an old school 
Win32/MFC statically linked C++ application, but I was not able to get modern MS tooling to cooperate.  
Perhaps if I was a more Microsoft-y kinda guy, I would have found a way to convince the latest Visual Studio to 
target Windows7 -- but alas this project's encounter with Microsoft tooling reminded me why I've avoided them 
for most of my career :).  I considered Qt or GTK, but frankly, working with C++ is something I swore I'd never 
do after writing the Tower Eiffel compiler (our motto: "You Deserve Better than a C+").  Eiffel never really 
caught on, but thankfully Java gave the world a decent alternative to C++.   But I digress...

In the end I found a pure Go library [tadvi/winc](https://github.com/tadvi/winc) that provides very simple GUI 
support with no external dependencies - easy to compile, no unexpected DLL conflicts on Iasos's machines.  
I would have prefered to use [andlabs/ui](https://github.com/andlabs/ui) library, but its dependency on the 
mingw compiler toolchain turns out to be too fragile (link issues caused by incompatible versions of mingw?). 
I wasted a lot of time trying to get it to build a 32bit windows exe and gave up when I found tadvi/winc.

