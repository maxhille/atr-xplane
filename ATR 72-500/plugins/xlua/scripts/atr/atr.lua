print("ATR XLua loading")

bus_volts = find_dataref("sim/cockpit2/electrical/bus_volts")
fuel_pump_on = find_dataref("sim/cockpit2/engine/actuators/fuel_pump_on")

fuel_pump_button_off = create_dataref("atr/fuel_pump_button_off", "array[2]")

function after_physics()
    if fuel_pump_on[0] == 0 and bus_volts[0] > 18.0 then
        fuel_pump_button_off[0] = 1
    else
        fuel_pump_button_off[0] = 0
    end
    if fuel_pump_on[1] == 0 and bus_volts[0] > 18.0 then
        fuel_pump_button_off[1] = 1
    else
        fuel_pump_button_off[1] = 0
    end
end
