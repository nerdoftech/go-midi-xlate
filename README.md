# go-midi-xlate

This project contains a series of modules and programs for translating MIDI messages from/to various input and output sources.

## Apps

### x32-delay

This program will take a sequence of MIDI "note-on" messages from a click track source such as ableton at the respective tempo to calulate the delay time. Then it will send to a Behringer X32 mixing console (via MIDI port only) a Sysex message that will set the delay time of the `stereo delay` effect. 

It takes 4 note-on messages at the begining of the song to calculate the delay time. Typically for a 4/4 time signature, you would send the note-on messages on the quarter notes for one measure. More than 4 will not help, and 8 beats will only make it send the delay time to the X32 twice.

__Note:__ This program is only available for 64-bit Intel based MacOS computers. If you would like to see more

Options:
```text
% x32-delay -help
Usage of ./build/mac/x32-delay:
  -ch uint
    	MIDI channel to listen on (default 1)
  -fxc int
    	fx channel of delay (default 1)
  -in int
    	input midi port index
  -list
    	shows MIDI ports
  -log string
    	sets the log level (default "info")
  -note uint
    	MIDI note that will trigger beat (default 43)
  -out int
    	output midi port index
```

Listing the ports to get the index number:
```text
% x32-delay -list
In ports:
  0 - Network In1
  1 - USB Midi
Out ports:
  0 - Network In1
  1 - USB Midi
```

Running the program:

```text
# Listen for 'C1' note-on message on MIDI channel 2 from network port 'In1', 
# output to network port 'USB Midi'

% x32-delay -in 0 -out 1 -ch 2 -fxc 2 -note 36
10:51PM INF Input port: Network Air1
10:51PM INF Output port: USB Midi

```