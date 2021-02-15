import { Signalling, MESSAGE_TYPE } from './signalling';
import { Peer } from './peer';
import { StreamType } from './stream_types';

// Adapted From:
// https://w3c.github.io/webrtc-pc/#perfect-negotiation-example
export class WebRTCHandler {
  signaller: Signalling;
  peers: { [id: string]: Peer };

  // Callbacks
  ontrack?: (id: string, correlation: StreamType, event: RTCTrackEvent) => void;

  constructor(signaller: Signalling) {
    this.signaller = signaller;
    this.peers = {};
  }

  addPeer(id: string, polite?: boolean): Peer {
    console.log('Adding Peer: ' + id);
    const peer = new Peer(id, new RTCPeerConnection(), polite ?? false);

    if (!polite) {
      peer.correlateChan = peer.conn.createDataChannel('correlate');
      peer.correlateChan.onmessage = (event) => {
        this.incomingCorrelation(id, event);
      };
    }

    peer.conn.oniceconnectionstatechange = () => {
      if (peer.conn.iceConnectionState === 'disconnected') {
        console.log('Peer Disconnected: ' + peer.id);
        delete this.peers[peer.id];
      }
    };

    peer.conn.ondatachannel = (event: RTCDataChannelEvent) => {
      peer.correlateChan = event.channel;
      peer.correlateChan.onmessage = (event) => {
        this.incomingCorrelation(id, event);
      };
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
      if (!this.ontrack) return;
      const sid = event.streams[0].id;
      this.ontrack(id, peer.streamCorrelations[sid], event);
      delete peer.streamCorrelations[sid];
    };

    this.peers[id] = peer;
    return peer;
  }

  incomingCorrelation(id: string, event: MessageEvent<any>) {
    const correlation = JSON.parse(event.data);
    this.peers[id].streamCorrelations[correlation.sid] = correlation.kind;
  }

  async incomingSDP(id: string, sdp: RTCSessionDescription) {
    const peer = this.peers[id] ?? this.addPeer(id, true);
    try {
      const readyForOffer =
        !peer.makingOffer && (peer.conn.signalingState == 'stable' || peer.isSettingRemoteAnswerPending);
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
    Object.values(this.peers).forEach((peer) => {
      peer.correlateChan!.send(JSON.stringify({ sid: stream.id, kind: kind }));
      const sender = peer.conn.addTrack(track, stream);
      peer.tracks[track.id] = sender;
    });
  }

  replaceTrack(oldTrack: MediaStreamTrack, newTrack: MediaStreamTrack) {
    Object.values(this.peers).forEach(async (peer) => {
      if (!peer.tracks[oldTrack.id]) return;
      await peer.tracks[oldTrack.id].replaceTrack(newTrack);
      peer.tracks[newTrack.id] = peer.tracks[oldTrack.id];
      delete peer.tracks[oldTrack.id];
    });
  }

  removeTrack(track: MediaStreamTrack) {
    Object.values(this.peers).forEach((peer) => {
      if (!peer.tracks[track.id]) return;
      peer.conn.removeTrack(peer.tracks[track.id]);
      delete peer.tracks[track.id];
    });
  }
}
