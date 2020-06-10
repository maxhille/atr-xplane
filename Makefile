.PHONY: release

release: 
	7z a release/atr72_500-`git describe`.7z ATR\ 72-500/*acf ATR\ 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png ATR\ 72-500/objects/*
