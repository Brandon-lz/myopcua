#!/usr/bin/env python3

from threading import Thread
import logging
import random
import time
import sys
sys.path.insert(0, "..")

from opcua import ua, uamethod, Server

# method to be exposed through server
def set_method(parent, variant):
    print("set_method %d", variant.Value)
    temperature_thread.temperature.set_value(variant.Value)

class Temperature(Thread):

    def __init__(self, temperature, onoff):
        Thread.__init__(self)
        self._stop = False
        self.temperature = temperature
        self.onoff = onoff

    def stop(self):
        self._stop = True

    def run(self):
        count = 1
        while not self._stop:
            value = random.randint(-20, 100)
            self.temperature.set_value(value)

            value = bool(random.randint(0, 1))
            self.onoff.set_value(value)

            led_event.event.Message = ua.LocalizedText("high_temperature %d" % count)
            led_event.event.Severity = count
            led_event.event.temperature = random.randint(60, 100)
            led_event.event.onoff = bool(random.randint(0, 1))
            led_event.trigger()
            count += 1

            time.sleep(5)

if __name__ == "__main__":

    # optional: setup logging
    logging.basicConfig(level=logging.WARN)

    # now setup our server
    server = Server()
    server.set_endpoint("opc.tcp://0.0.0.0:4840/freeopcua/server/")
    server.set_server_name("FreeOpcUa Example Server")

    # set all possible endpoint policies for clients to connect through
    server.set_security_policy([
                ua.SecurityPolicyType.NoSecurity,
                ua.SecurityPolicyType.Basic128Rsa15_SignAndEncrypt,
                ua.SecurityPolicyType.Basic128Rsa15_Sign,
                ua.SecurityPolicyType.Basic256_SignAndEncrypt,
                ua.SecurityPolicyType.Basic256_Sign])

    # setup our own namespace
    uri = "http://examples.freeopcua.github.io"
    idx = server.register_namespace(uri)
    
    # create directly some objects and variables
    demo_led = server.nodes.objects.add_object(idx, "demo_led")
    led_temperature = demo_led.add_variable(idx, "temperature", 20)
    led_temperature.set_writable()    # Set MyVariable to be writable by clients
    led_onoff = demo_led.add_variable(idx, "onoff", True)
    led_onoff.set_writable()

    led_method = demo_led.add_method(idx, "led_method", set_method, [ua.VariantType.UInt32])

    # creating a default event object, the event object automatically will have members for all events properties
    led_event_type = server.create_custom_event_type(idx, 
                                                    'high_temperature', 
                                                    ua.ObjectIds.BaseEventType,
                                                    [('temperature', ua.VariantType.UInt32), ('onoff', ua.VariantType.Boolean)])

    led_event = server.get_event_generator(led_event_type, demo_led)
    led_event.event.Severity = 300

    # start opcua server
    server.start()
    print("Start opcua server...")

    temperature_thread = Temperature(led_temperature, led_onoff)  # just  a stupide class update a variable
    temperature_thread.start()

    try:

        led_event.trigger(message="This is BaseEvent")

        while True:
            time.sleep(5)
    finally:
        print("Exit opcua server...")
        temperature_thread.stop()
        server.stop()
