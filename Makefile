.PHONY: release tex

release: 
	cp COPYING ATR\ 72-500/
	cp README-release.txt ATR\ 72-500/
	7z a release/atr72_500-`git describe`.7z ATR\ 72-500/*acf ATR\ 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png ATR\ 72-500/objects/*

tex:
	inkscape svg/cockpit.svg -i std -j -C --export-png="ATR 72-500/objects/cockpit.png"
	inkscape svg/cockpit.svg -i LIT -j -C --export-png="ATR 72-500/objects/cockpit_LIT.png"
	inkscape svg/switches.svg -i std -j -C --export-png="ATR 72-500/objects/switches.png"
	inkscape svg/switches.svg -i LIT -j -C --export-png="ATR 72-500/objects/switches_LIT.png"

panel:
	inkscape svg/panel.svg -j -C --export-png="ATR 72-500/cockpit_3d/-PANELS-/Panel_Airliner.png"
	inkscape svg/panel_gen_trim_elev.svg -i gen_pointer -j -C --export-png="ATR 72-500/cockpit_3d/generic/trim_elev/gen_pointer.png"
	inkscape svg/panel_gen_trim_elev.svg -i gen_pointer-1 -j -C --export-png="ATR 72-500/cockpit_3d/generic/trim_elev/gen_pointer-1.png"
