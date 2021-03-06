ATR Open Source Project for X-Plane
===================================

This project intends to build a good set of ATR aircraft and share it with the community as open source.

## Goals
- short term: ATR 72-500 which is flyable and has the most basic features one would expect
- medium term: Aircraft is closer to study level by having realistic perfomance characteristics and systems
- long term: Have a complete set including ATR 42 and the -600 series

## Download latest release
Go to the `releases` page of this Github repo and download from there

## Development set up
1. Clone this repo
2. link `ATR 72-500` into your X-Plane `Aircraft` folder

- X-Plane 11.41
- Blender 2.82
- Inkscape 0.92

## Building
1. Export `obj`s via XPlane2Blender
2. Export `png`s via Inkscape and `make tex`
3. Tag release (eg `0.3`)
4. `make release`

## Credits
```
(C) 2020 Maximilian Hille
```

```
Initial model, bitmap resources and sounds:
Authors: Narendran Muraleedharan, Donald Belcham, Dwayne Gable, Oliver (ot-666), camelon
https://bitbucket.org/muraleen/atr72-500-c-project
```
