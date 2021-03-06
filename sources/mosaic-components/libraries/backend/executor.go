

package backend


import "os"
import "strings"

import "mosaic-components/libraries/channels"
import "mosaic-components/libraries/messages"
import "vgl/transcript"


func Execute (_callbacks Callbacks, _componentIdentifier string, _channelEndpoint string) (error) {
	
	_transcript := packageTranscript
	
	var _error error
	
	_transcript.TraceDebugging ("initializing...")
	_transcript.TraceDebugging ("  * using the identifier `%s`;", _componentIdentifier)
	_transcript.TraceDebugging ("  * using the channel `%s`;", _channelEndpoint)
	
	_transcript.TraceDebugging ("creating the component backend...")
	var _backend Backend
	var _backendChannelCallbacks channels.Callbacks
	if _backend, _backendChannelCallbacks, _error = Create (_callbacks); _error != nil {
		panic (_error)
	}
	
	if true {
		transcript.SetBackend (& transcriptBackend { backend : _backend })
	}
	
	_transcript.TraceDebugging ("creating the component channel...")
	if _channelEndpoint == "stdio" {
		_transcript.TraceDebugging ("  * using the stdio endpoint;")
		_inboundStream := os.Stdin
		_outboundStream := os.Stdout
		if _, _error = channels.Create (_backendChannelCallbacks, _inboundStream, _outboundStream, nil); _error != nil {
			panic (_error)
		}
	} else if strings.HasPrefix (_channelEndpoint, "tcp:") {
		_channelTcpEndpoint := _channelEndpoint[4:]
		_transcript.TraceDebugging ("  * usig the TCP endpoint `%s`;", _channelTcpEndpoint)
		if _, _error = channels.CreateAndDial (_backendChannelCallbacks, "tcp", _channelTcpEndpoint); _error != nil {
			panic (_error)
		}
	} else {
		_transcript.TraceError ("invalid component channel endpoint; aborting!")
		panic ("failed")
	}
	
	_transcript.TraceInformation ("executing...")
	
	_transcript.TraceDebugging ("waiting for the termination of the component backend...")
	if _error := _backend.WaitTerminated (); _error != nil {
		panic (_error)
	}
	
	_transcript.TraceInformation ("terminated.")
	return nil
}


type transcriptBackend struct {
	backend Backend
}

func (_transcript *transcriptBackend) Consume (_trace *transcript.Trace) () {
	
	// FIXME: Make this configurable!
	if true {
		transcript.StdErrBackend.Consume (_trace)
	}
	
	// FIXME: Make this configurable!
	if false {
		if transcript.ShouldConsume (_trace, transcript.DumpingLevel) {
			_line := transcript.FormatTrace (_trace)
			_transcript.backend.TranscriptPush (messages.Attachment (_line))
		}
	}
}
