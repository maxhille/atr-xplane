.PHONY: release tex

release: 
	7z a release/atr72_500-`git describe`.7z ATR\ 72-500/*acf ATR\ 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png ATR\ 72-500/objects/*

tex:
	inkscape svg/cockpit.svg -i std -j -C --export-png="ATR 72-500/objects/cockpit.png"
	inkscape svg/cockpit.svg -i LIT -j -C --export-png="ATR 72-500/objects/cockpit_LIT.png"
	inkscape svg/switches.svg -i std -j -C --export-png="ATR 72-500/objects/switches.png"
	inkscape svg/switches.svg -i LIT -j -C --export-png="ATR 72-500/objects/switches_LIT.png"
