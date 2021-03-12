import { Signalling, MESSAGE_TYPE } from './signalling';
import { Peer } from './peer';
import { StreamType } from './stream_types';
import { TurnCredentials } from '../api/definitions';

// Adapted From:
// https://w3c.github.io/webrtc-pc/#perfect-negotiation-example
export class WebRTCHandler {
  signaller: Signalling;
  credentials: TurnCredentials;
  peers: { [id: string]: Peer };
  tracks: { [id: string]: [MediaStream, StreamType] };

  // Callbacks
  ontrack?: (id: string, correlation: StreamType, event: RTCTrackEvent) => void;
  ontrackremove?: (id: string, correlation: StreamType, event: RTCTrackEvent) => void;
  ondisconnect?: (id: string) => void;
  onAddPeer: (polite: boolean) => void;

  constructor(signaller: Signalling, credentials: TurnCredentials, onAddPeer: (polite: boolean) => void) {
    this.signaller = signaller;
    this.onAddPeer = onAddPeer;
    this.credentials = credentials;
    this.peers = {};
    this.tracks = {};
  }

  addPeer(id: string, polite?: boolean): Peer {
    if (this.peers[id]) {
      console.log('Deleting Existing Peer: ' + id);
      this.peers[id].conn.close();
      delete this.peers[id];
    }
    console.log('Adding Peer');

    const peer = new Peer(
      id,
      new RTCPeerConnection({
        iceServers: [
          { urls: 'stun:stun.l.google.com:19302' },
          { urls: 'stun:stun1.l.google.com:19302' },
          {
            urls: 'turns:turn.astratutor.com:16500',
            username: this.credentials.username,
            credential: this.credentials.password,
          },
        ],
      }),
      polite || false,
    );

    const correlateMessage = (event: MessageEvent<any>) => {
      this.incomingCorrelation(id, event);
    };

    // Impolite Peer Creates a Data Channel
    if (!polite) {
      console.log('Creating Data Channel');
      peer.correlateChan = peer.conn.createDataChannel('correlate');
      peer.correlateChan.onmessage = correlateMessage;
    } else {
      // Polite Peer can receive a data channel
      peer.conn.ondatachannel = (event: RTCDataChannelEvent) => {
        peer.correlateChan = event.channel;
        console.log('Creating Data Channel - Receive');
        peer.correlateChan.send(JSON.stringify({ kind: StreamType._READY_ }));
        peer.correlateChan.onmessage = correlateMessage;
      };
    }

    peer.conn.oniceconnectionstatechange = () => {
      if (this.peers[id].conn.iceConnectionState === 'disconnected') {
        console.log('Peer Disconnected: ' + peer.id);
        if (this.ondisconnect) {
          this.ondisconnect(id);
        }
        this.peers[id].conn.close();
        delete this.peers[peer.id];
      }
    };

    peer.conn.onicecandidate = (event) => {
      if (event.candidate) {
        console.log('Sending Candidate: ' + id);
        this.signaller.send(MESSAGE_TYPE.CANDIDATE, id, event.candidate);
      }
    };

    peer.conn.onnegotiationneeded = async () => {
      console.log('Negotiaion Needed: ' + id);
      try {
        peer.makingOffer = true;
        await peer.conn.setLocalDescription({});
        this.signaller.send(MESSAGE_TYPE.SDP, id, peer.conn.localDescription);
      } catch (err) {
        console.error(err);
      } finally {
        peer.makingOffer = false;
      }
    };

    peer.conn.ontrack = (event: RTCTrackEvent) => {
      event.track.onunmute = () => {
        if (!this.ontrack) return;
        const sid = event.streams[0].id;
        const correlation = peer.streamCorrelations[sid];
        event.streams[0].addEventListener('removetrack', (_e) => {
          if (this.ontrackremove) {
            this.ontrackremove(id, correlation, event);
          }
        });
        this.ontrack(id, peer.streamCorrelations[sid], event);
        delete peer.streamCorrelations[sid];
      };
    };

    this.peers[id] = peer;
    return peer;
  }

  // Correlate track IDs with content type e.g. Webcam, Screenshare
  incomingCorrelation(id: string, event: MessageEvent<any>) {
    const correlation = JSON.parse(event.data);
    console.log('Incoming Correlation: ' + id + ' Kind: ' + correlation.kind);
    const peer = this.peers[id];

    // _READY_ informs they are ready to start receiving
    if (correlation.kind === StreamType._READY_ && !peer.ready) {
      peer.ready = true;
      this.onAddPeer(peer.polite);
      // Respond that we are ready too
      peer.correlateChan!.send(JSON.stringify({ kind: StreamType._READY_ }));

      // Send them correlations for all tracks we are sharing
      Object.entries(this.tracks).forEach(([trackId, [stream, kind]]) => {
        const track = stream.getTrackById(trackId);
        peer.correlateChan!.send(JSON.stringify({ sid: stream.id, kind: kind }));
        const sender = peer.conn.addTrack(track!, stream);
        peer.tracks[trackId] = sender;
      });
    } else {
      peer.streamCorrelations[correlation.sid] = correlation.kind;
    }
  }

  async incomingSDP(id: string, sdp: RTCSessionDescription) {
    // Add a peer with politeness if this peer doesn't not exist already
    const peer = this.peers[id] ?? this.addPeer(id, true);
    try {
      const readyForOffer =
        !peer.makingOffer && (peer.conn.signalingState === 'stable' || peer.isSettingRemoteAnswerPending);
      const offerCollision = sdp.type === 'offer' && !readyForOffer;

      peer.ignoreOffer = !peer.polite && offerCollision;
      if (peer.ignoreOffer) return;

      peer.isSettingRemoteAnswerPending = sdp.type === 'answer';
      await peer.conn.setRemoteDescription(sdp);
      peer.isSettingRemoteAnswerPending = false;

      if (sdp.type === 'offer') {
        await peer.conn.setLocalDescription({});
        this.signaller.send(MESSAGE_TYPE.SDP, id, peer.conn.localDescription);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async incomingCandidate(id: string, candidate: RTCIceCandidate) {
    const peer = this.peers[id];
    try {
      try {
        await peer.conn.addIceCandidate(candidate);
      } catch (err) {
        if (!peer.ignoreOffer) throw err;
      }
    } catch (err) {
      console.error(err);
    }
  }

  addTrack(track: MediaStreamTrack, kind: StreamType, stream: MediaStream) {
    if (this.tracks[track.id]) return;

    console.log('Adding Track: ', track);
    this.tracks[track.id] = [stream, kind];
    Object.values(this.peers).forEach((peer) => {
      if (!peer.ready) return;

      peer.correlateChan!.send(JSON.stringify({ sid: stream.id, kind: kind }));
      const sender = peer.conn.addTrack(track, stream);
      peer.tracks[track.id] = sender;
    });
  }

  replaceTrack(oldTrack: MediaStreamTrack, newTrack: MediaStreamTrack) {
    if (!this.tracks[oldTrack.id] || newTrack.id === oldTrack.id) return;

    console.log('Replacing Track: ' + oldTrack.id + ' with ' + newTrack.id);
    this.tracks[newTrack.id] = this.tracks[oldTrack.id];
    delete this.tracks[oldTrack.id];

    Object.values(this.peers).forEach(async (peer) => {
      if (!peer.ready || !peer.tracks[oldTrack.id]) return;
      peer.tracks[newTrack.id] = peer.tracks[oldTrack.id];
      await peer.tracks[oldTrack.id].replaceTrack(newTrack);
      delete peer.tracks[oldTrack.id];
      console.log('Replaced Peer', peer.id);
    });
  }

  removeTrack(track: MediaStreamTrack) {
    if (!this.tracks[track.id]) return;

    console.log('Removing Track: ' + track.id);
    delete this.tracks[track.id];
    Object.values(this.peers).forEach((peer) => {
      if (!peer.ready || !peer.tracks[track.id]) return;

      peer.conn.removeTrack(peer.tracks[track.id]);
      delete peer.tracks[track.id];
    });
  }

  close() {
    Object.values(this.peers).forEach((peer) => {
      peer.conn.close();
    });
  }
}
